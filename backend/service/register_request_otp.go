package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"time"

	"synctask/backend/model"
	"synctask/backend/repository"
	"synctask/backend/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequestOtpService は新規登録OTP発行ユースケースのインターフェースです。
type RegisterRequestOtpService interface {
	RequestOtp(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error)
}

// RegisterRequestOtpDependencies は時間、待機、乱数、ハッシュ等のテスト可能な依存です。
type RegisterRequestOtpDependencies struct {
	Now               func() time.Time
	Sleep             func(time.Duration)
	GenerateSessionID func() (string, error)
	GenerateOTP       func() (string, error)
	HashPassword      func(string) (string, error)
	HashOTP           func(string) (string, error)
	ResponseDelay     func() time.Duration
}

type registerRequestOtpService struct {
	repo   repository.RegisterRequestOtpRepository
	mailer Mailer
	deps   RegisterRequestOtpDependencies
}

// NewRegisterRequestOtpService はRegisterRequestOtpServiceを生成します。
func NewRegisterRequestOtpService(
	repo repository.RegisterRequestOtpRepository,
	mailer Mailer,
	deps RegisterRequestOtpDependencies,
) RegisterRequestOtpService {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.Sleep == nil {
		deps.Sleep = time.Sleep
	}
	if deps.GenerateSessionID == nil {
		deps.GenerateSessionID = util.GenerateOTPSessionID
	}
	if deps.GenerateOTP == nil {
		deps.GenerateOTP = util.GenerateOTP
	}
	if deps.HashPassword == nil {
		deps.HashPassword = func(pw string) (string, error) {
			hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			return string(hash), err
		}
	}
	if deps.HashOTP == nil {
		deps.HashOTP = func(otp string) (string, error) {
			hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
			return string(hash), err
		}
	}
	if deps.ResponseDelay == nil {
		deps.ResponseDelay = func() time.Duration {
			// 1.0s ± 0.1s (900ms ~ 1100ms)
			n, err := rand.Int(rand.Reader, big.NewInt(201))
			if err != nil {
				return 1000 * time.Millisecond
			}
			return time.Duration(900+n.Int64()) * time.Millisecond
		}
	}

	return &registerRequestOtpService{
		repo:   repo,
		mailer: mailer,
		deps:   deps,
	}
}

// RequestOtp は新規登録OTP発行処理を行います。
//
// @spec 入力検証違反時は遅延なしで 400 Bad Request を返し、ACCESS_LOGを記録する。
// @spec 有効アカウントまたは有効OTPセッションが存在する場合はダミーセッションを作成し、メール送信をスキップする。
// @spec 未登録かつ排他なしの場合は実OTPを生成し、セッションを保存してメールを送信する。
// @spec 実メール送信失敗時は DELIVERY_STATUS='sendable'、SEND_FAILED_COUNT=1 に更新し 503 OTP_DELIVERY_FAILED を返す。
// @spec 正常処理・ダミー処理ともに 1.0s ± 0.1s のレスポンス遅延を適用し、同一構造の 200 OK を返す。
func (s *registerRequestOtpService) RequestOtp(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
	now := s.deps.Now()

	// 1. 入力正規化・バリデーション（違反時は遅延なし 400 Bad Request）
	if err := req.Validate(); err != nil {
		accessLog := &model.AccessLogRecord{
			LogID:     uuid.NewString(),
			UserID:    sql.NullString{},
			IPAddress: clientIP,
			Endpoint:  "POST auth/register/request-otp",
			AccessAt:  now,
		}
		_ = s.repo.RecordAccessLog(ctx, accessLog)
		return nil, err
	}

	// 2. 既存アカウント・既存有効OTPセッションの照会
	activeUser, err := s.repo.FindActiveUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	activeOtp, err := s.repo.FindActiveOtpSessionByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	isDummy := activeUser || activeOtp

	// 3. セッションIDおよびマスクメール生成
	sessionID, err := s.deps.GenerateSessionID()
	if err != nil {
		return nil, err
	}
	maskedEmail := util.MaskEmail(req.Email)

	var (
		otp     string
		otpHash string
		pwdHash string
	)

	if !isDummy {
		otp, err = s.deps.GenerateOTP()
		if err != nil {
			return nil, err
		}
		pwdHash, err = s.deps.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		otpHash, err = s.deps.HashOTP(otp)
		if err != nil {
			return nil, err
		}
	}

	session := &model.OtpSessionRecord{
		OtpSessionID:     sessionID,
		Purpose:          "SIGNUP",
		UserID:           sql.NullString{},
		MaskedEmail:      maskedEmail,
		Status:           "active",
		IsDummy:          isDummy,
		AttemptCount:     0,
		SendCount:        0,
		SendFailedCount:  0,
		DeliveryStatus:   "pending",
		LastSentAt:       now,
		OtpExpiresAt:     now.Add(5 * time.Minute),
		SessionExpiresAt: now.Add(15 * time.Minute),
		CreatedAt:        now,
	}

	if !isDummy {
		session.PendingUsername = sql.NullString{String: req.Username, Valid: true}
		session.PendingEmail = sql.NullString{String: req.Email, Valid: true}
		session.PendingPasswordHash = sql.NullString{String: pwdHash, Valid: true}
		session.OtpHash = sql.NullString{String: otpHash, Valid: true}
	}

	mailLog := &model.MailAuthLogRecord{
		LogID:     uuid.NewString(),
		UserID:    sql.NullString{},
		Email:     req.Email,
		AuthType:  "SIGNUP",
		IPAddress: clientIP,
		EventType: "ISSUED",
		IsSuccess: true,
		IsDummy:   isDummy,
		AccessAt:  now,
	}

	accessLog := &model.AccessLogRecord{
		LogID:      uuid.NewString(),
		UserID:     sql.NullString{},
		IPAddress:  clientIP,
		Endpoint:   "POST auth/register/request-otp",
		ResourceID: sql.NullString{String: sessionID, Valid: true},
		AccessAt:   now,
	}

	if err := s.repo.SaveSessionWithLogs(ctx, session, mailLog, accessLog); err != nil {
		return nil, err
	}

	// 4. 実メール送信（通常処理時のみ）
	if !isDummy {
		if err := s.mailer.SendOTP(ctx, req.Email, otp); err != nil {
			_ = s.repo.UpdateOtpSessionDeliveryStatus(ctx, sessionID, "sendable", 1)
			return nil, model.NewServiceUnavailableError("OTP_DELIVERY_FAILED", "メールの送信に失敗しました。")
		}
		_ = s.repo.UpdateOtpSessionDeliveryStatus(ctx, sessionID, "sent", 0)
	}

	// 5. タイミング攻撃対策遅延（1.0s ± 0.1s）を適用（通常成功時・ダミー時）
	s.deps.Sleep(s.deps.ResponseDelay())

	// 6. レスポンス返却
	return &model.RegisterRequestOtpResponse{
		OtpSessionID:     sessionID,
		MaskedEmail:      maskedEmail,
		ExpiresInSeconds: 300,
		CooldownSeconds:  60,
	}, nil
}
