package repository

import (
	"context"
	"database/sql"
	"errors"

	"synctask/backend/model"
)

// RegisterRequestOtpRepository は新規登録OTP発行の永続化インターフェースです。
type RegisterRequestOtpRepository interface {
	FindActiveUserByEmail(ctx context.Context, email string) (bool, error)
	FindActiveOtpSessionByEmail(ctx context.Context, email string) (bool, error)
	CreateOtpSession(ctx context.Context, session *model.OtpSessionRecord) error
	UpdateOtpSessionDeliveryStatus(ctx context.Context, sessionID, status string, sendFailedCount int) error
	RecordMailAuthLog(ctx context.Context, log *model.MailAuthLogRecord) error
	RecordAccessLog(ctx context.Context, log *model.AccessLogRecord) error
	SaveSessionWithLogs(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error
}

type registerRequestOtpRepository struct {
	db *sql.DB
}

// NewRegisterRequestOtpRepository はRegisterRequestOtpRepositoryを生成します。
func NewRegisterRequestOtpRepository(db *sql.DB) RegisterRequestOtpRepository {
	return &registerRequestOtpRepository{db: db}
}

// FindActiveUserByEmail は有効アカウントの存在を確認します。
//
// @spec IS_DELETED = FALSE のアカウントが存在する場合に true を返す。
func (r *registerRequestOtpRepository) FindActiveUserByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM LOGIN_ACCOUNT WHERE EMAIL = $1 AND IS_DELETED = FALSE LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FindActiveOtpSessionByEmail は有効期限内のアクティブなOTPセッションの存在を確認します。
//
// @spec STATUS IN ('active', 'verified') かつ SESSION_EXPIRES_AT > NOW() のセッションが存在する場合に true を返す。
func (r *registerRequestOtpRepository) FindActiveOtpSessionByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM OTP_SESSION WHERE PENDING_EMAIL = $1 AND STATUS IN ('active', 'verified') AND SESSION_EXPIRES_AT > NOW() LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateOtpSession はOTP_SESSIONテーブルにレコードを挿入します。
//
// @spec セッション情報をOTP_SESSIONに永続化する。
func (r *registerRequestOtpRepository) CreateOtpSession(ctx context.Context, session *model.OtpSessionRecord) error {
	query := `
		INSERT INTO OTP_SESSION (
			OTP_SESSION_ID,
			PURPOSE,
			USER_ID,
			PENDING_USERNAME,
			PENDING_EMAIL,
			MASKED_EMAIL,
			PENDING_PASSWORD_HASH,
			OTP_HASH,
			STATUS,
			IS_DUMMY,
			ATTEMPT_COUNT,
			SEND_COUNT,
			SEND_FAILED_COUNT,
			DELIVERY_STATUS,
			LAST_SENT_AT,
			OTP_EXPIRES_AT,
			SESSION_EXPIRES_AT,
			CREATED_AT
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		session.OtpSessionID,
		session.Purpose,
		session.UserID,
		session.PendingUsername,
		session.PendingEmail,
		session.MaskedEmail,
		session.PendingPasswordHash,
		session.OtpHash,
		session.Status,
		session.IsDummy,
		session.AttemptCount,
		session.SendCount,
		session.SendFailedCount,
		session.DeliveryStatus,
		session.LastSentAt,
		session.OtpExpiresAt,
		session.SessionExpiresAt,
		session.CreatedAt,
	)
	return err
}

// UpdateOtpSessionDeliveryStatus は送信状態と失敗回数を更新します。
//
// @spec DELIVERY_STATUS と SEND_FAILED_COUNT を指定セッションIDに対して更新する。
func (r *registerRequestOtpRepository) UpdateOtpSessionDeliveryStatus(ctx context.Context, sessionID, status string, sendFailedCount int) error {
	query := `UPDATE OTP_SESSION SET DELIVERY_STATUS = $1, SEND_FAILED_COUNT = $2 WHERE OTP_SESSION_ID = $3`
	_, err := r.db.ExecContext(ctx, query, status, sendFailedCount, sessionID)
	return err
}

// RecordMailAuthLog はMAIL_AUTH_LOGテーブルにログを記録します。
//
// @spec メール認証ログを永続化する。
func (r *registerRequestOtpRepository) RecordMailAuthLog(ctx context.Context, log *model.MailAuthLogRecord) error {
	query := `
		INSERT INTO MAIL_AUTH_LOG (
			LOG_ID,
			USER_ID,
			EMAIL,
			AUTH_TYPE,
			IP_ADDRESS,
			EVENT_TYPE,
			IS_SUCCESS,
			IS_DUMMY,
			ACCESS_AT
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		log.LogID,
		log.UserID,
		log.Email,
		log.AuthType,
		log.IPAddress,
		log.EventType,
		log.IsSuccess,
		log.IsDummy,
		log.AccessAt,
	)
	return err
}

// RecordAccessLog はACCESS_LOGテーブルにアクセスログを記録します。
//
// @spec APIアクセスログを永続化する。
func (r *registerRequestOtpRepository) RecordAccessLog(ctx context.Context, log *model.AccessLogRecord) error {
	query := `
		INSERT INTO ACCESS_LOG (
			LOG_ID,
			USER_ID,
			IP_ADDRESS,
			ENDPOINT,
			RESOURCE_ID,
			ACCESS_AT
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		log.LogID,
		log.UserID,
		log.IPAddress,
		log.Endpoint,
		log.ResourceID,
		log.AccessAt,
	)
	return err
}

// SaveSessionWithLogs は単一トランザクションでセッション・メール認証ログ・アクセスログを原子的保存します。
//
// @spec 単一トランザクション内でセッション登録、メール認証ログ、アクセスログを記録しコミットする。
// @spec いずれかでエラーが発生した場合はロールバックする。
func (r *registerRequestOtpRepository) SaveSessionWithLogs(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. OTP_SESSION
	sessionQuery := `
		INSERT INTO OTP_SESSION (
			OTP_SESSION_ID,
			PURPOSE,
			USER_ID,
			PENDING_USERNAME,
			PENDING_EMAIL,
			MASKED_EMAIL,
			PENDING_PASSWORD_HASH,
			OTP_HASH,
			STATUS,
			IS_DUMMY,
			ATTEMPT_COUNT,
			SEND_COUNT,
			SEND_FAILED_COUNT,
			DELIVERY_STATUS,
			LAST_SENT_AT,
			OTP_EXPIRES_AT,
			SESSION_EXPIRES_AT,
			CREATED_AT
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err = tx.ExecContext(
		ctx,
		sessionQuery,
		session.OtpSessionID,
		session.Purpose,
		session.UserID,
		session.PendingUsername,
		session.PendingEmail,
		session.MaskedEmail,
		session.PendingPasswordHash,
		session.OtpHash,
		session.Status,
		session.IsDummy,
		session.AttemptCount,
		session.SendCount,
		session.SendFailedCount,
		session.DeliveryStatus,
		session.LastSentAt,
		session.OtpExpiresAt,
		session.SessionExpiresAt,
		session.CreatedAt,
	)
	if err != nil {
		return err
	}

	// 2. MAIL_AUTH_LOG
	mailLogQuery := `
		INSERT INTO MAIL_AUTH_LOG (
			LOG_ID,
			USER_ID,
			EMAIL,
			AUTH_TYPE,
			IP_ADDRESS,
			EVENT_TYPE,
			IS_SUCCESS,
			IS_DUMMY,
			ACCESS_AT
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.ExecContext(
		ctx,
		mailLogQuery,
		mailLog.LogID,
		mailLog.UserID,
		mailLog.Email,
		mailLog.AuthType,
		mailLog.IPAddress,
		mailLog.EventType,
		mailLog.IsSuccess,
		mailLog.IsDummy,
		mailLog.AccessAt,
	)
	if err != nil {
		return err
	}

	// 3. ACCESS_LOG
	accessLogQuery := `
		INSERT INTO ACCESS_LOG (
			LOG_ID,
			USER_ID,
			IP_ADDRESS,
			ENDPOINT,
			RESOURCE_ID,
			ACCESS_AT
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.ExecContext(
		ctx,
		accessLogQuery,
		accessLog.LogID,
		accessLog.UserID,
		accessLog.IPAddress,
		accessLog.Endpoint,
		accessLog.ResourceID,
		accessLog.AccessAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
