package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"synctask/backend/model"
	"synctask/backend/repository"
	"synctask/backend/util"

	"golang.org/x/crypto/bcrypt"
)

const loginSessionDuration = 30 * 24 * time.Hour

// LoginService はログインユースケースのインターフェースです。
type LoginService interface {
	Login(context.Context, *model.LoginRequest, model.LoginMetadata) (*model.LoginServiceResult, error)
	RecordInvalidRequest(context.Context, string) error
}

// LoginDependencies は時間、待機、乱数、パスワード照合のテスト可能な依存です。
type LoginDependencies struct {
	Now             func() time.Time
	Sleep           func(time.Duration)
	GenerateToken   func() (string, error)
	ComparePassword func(string, string) bool
	FailureDelay    func() time.Duration
}

type loginService struct {
	repo repository.LoginRepository
	deps LoginDependencies
}

// NewLoginService はLoginServiceを生成します。
func NewLoginService(repo repository.LoginRepository, deps LoginDependencies) LoginService {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.Sleep == nil {
		deps.Sleep = time.Sleep
	}
	if deps.GenerateToken == nil {
		deps.GenerateToken = func() (string, error) {
			return util.GenerateSecureToken(rand.Reader, 32)
		}
	}
	if deps.ComparePassword == nil {
		deps.ComparePassword = func(hash, password string) bool {
			return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
		}
	}
	if deps.FailureDelay == nil {
		deps.FailureDelay = randomLoginFailureDelay
	}
	return &loginService{repo: repo, deps: deps}
}

// Login は入力検証後にRepositoryへ原子的なログイン試行を依頼します。
//
// @spec 入力不正は遅延と失敗カウンター加算を行わずACCESS_LOGだけを記録する。
// @spec 401と429は処理開始から1.0秒±0.1秒となるよう待機する。
// @spec 成功時は新規セッション/CSRFトークンを返し人工遅延しない。
func (s *loginService) Login(ctx context.Context, req *model.LoginRequest, meta model.LoginMetadata) (*model.LoginServiceResult, error) {
	startedAt := s.deps.Now()
	if err := req.Validate(); err != nil {
		if logErr := s.repo.RecordInvalidRequest(ctx, meta.IP, startedAt); logErr != nil {
			return nil, fmt.Errorf("failed to record invalid login request: %w", logErr)
		}
		return nil, err
	}

	sessionID, err := s.deps.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}
	csrfToken, err := s.deps.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate csrf token: %w", err)
	}

	attempt := &model.LoginAttempt{
		Email:        req.Email,
		Password:     req.Password,
		IP:           meta.IP,
		UserAgent:    meta.UserAgent,
		OldSessionID: meta.OldSessionID,
		SessionID:    sessionID,
		CSRFToken:    csrfToken,
		Now:          startedAt,
		ExpiresAt:    startedAt.Add(loginSessionDuration),
	}
	result, err := s.repo.AttemptLogin(ctx, attempt, func(hash string) bool {
		return s.deps.ComparePassword(hash, req.Password)
	})
	if err != nil {
		return nil, err
	}

	switch result.Status {
	case model.LoginStatusSuccess:
		if result.User == nil {
			return nil, fmt.Errorf("login repository returned success without user")
		}
		return &model.LoginServiceResult{
			Response:  model.LoginResponse{User: *result.User},
			SessionID: sessionID,
			CSRFToken: csrfToken,
			MaxAge:    loginSessionDuration,
		}, nil
	case model.LoginStatusUnauthorized:
		s.waitForFailureDelay(startedAt)
		return nil, model.NewUnauthorizedError("メールアドレスまたはパスワードが正しくありません。")
	case model.LoginStatusRateLimited:
		s.waitForFailureDelay(startedAt)
		return nil, model.NewRateLimitError("ログイン試行回数が上限に達しました。しばらくしてから再試行してください。", result.RetryAfter)
	default:
		return nil, fmt.Errorf("unknown login status: %s", result.Status)
	}
}

// RecordInvalidRequest はJSON構文等でLoginRequestを生成できない要求を記録します。
func (s *loginService) RecordInvalidRequest(ctx context.Context, ip string) error {
	return s.repo.RecordInvalidRequest(ctx, ip, s.deps.Now())
}

func (s *loginService) waitForFailureDelay(startedAt time.Time) {
	target := s.deps.FailureDelay()
	elapsed := s.deps.Now().Sub(startedAt)
	if remaining := target - elapsed; remaining > 0 {
		s.deps.Sleep(remaining)
	}
}

func randomLoginFailureDelay() time.Duration {
	value, err := rand.Int(rand.Reader, big.NewInt(200_000_001))
	if err != nil {
		return time.Second
	}
	return 900*time.Millisecond + time.Duration(value.Int64())
}
