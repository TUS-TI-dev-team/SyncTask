package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"synctask/backend/model"

	"github.com/google/uuid"
)

// DummyPasswordHash は未登録アカウントでも実ハッシュ相当の照合を行うための固定bcryptハッシュです。
const DummyPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

// PasswordCheck は取得したハッシュに対するパスワード照合処理です。
type PasswordCheck func(passwordHash string) bool

// LoginRepository はログイン試行の永続化境界です。
type LoginRepository interface {
	AttemptLogin(context.Context, *model.LoginAttempt, PasswordCheck) (*model.LoginAttemptResult, error)
	RecordInvalidRequest(context.Context, string, time.Time) error
}

type loginRepository struct {
	db *sql.DB
}

// NewLoginRepository はLoginRepositoryを生成します。
func NewLoginRepository(db *sql.DB) LoginRepository {
	return &loginRepository{db: db}
}

type loginAccount struct {
	user       model.LoginUser
	hash       string
	failed     int
	lastFailed sql.NullTime
	lockUntil  sql.NullTime
}

// AttemptLogin は制限確認、照合、カウンター、セッション、監査ログを単一トランザクションで処理します。
//
// @spec IP行とアカウント行をロックし、同一IP・メールの更新を直列化する。
// @spec 成功時の失敗状態リセット、旧セッション削除、新セッション登録、ログを原子的に行う。
// @spec 認証失敗時のアカウント/IPカウンターとログを原子的に行う。
func (r *loginRepository) AttemptLogin(ctx context.Context, attempt *model.LoginAttempt, check PasswordCheck) (_ *model.LoginAttemptResult, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	blockedUntil, err := selectIPLimit(ctx, tx, attempt.IP)
	if err != nil {
		return nil, err
	}
	if blockedUntil.Valid && blockedUntil.Time.After(attempt.Now) {
		if err = insertAttemptLogs(ctx, tx, nil, attempt, false); err != nil {
			return nil, err
		}
		if err = tx.Commit(); err != nil {
			return nil, err
		}
		retry := int(math.Ceil(blockedUntil.Time.Sub(attempt.Now).Seconds()))
		return &model.LoginAttemptResult{Status: model.LoginStatusRateLimited, RetryAfter: retry}, nil
	}

	account, err := selectLoginAccount(ctx, tx, attempt.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	accountExists := err == nil
	locked := accountExists && account.lockUntil.Valid && account.lockUntil.Time.After(attempt.Now)

	hash := DummyPasswordHash
	if accountExists {
		hash = account.hash
	}
	passwordMatches := false
	if !locked {
		passwordMatches = check(hash)
	}

	if !accountExists || locked || !passwordMatches {
		var userID any
		if accountExists {
			userID = account.user.ID
		}
		if accountExists && !locked {
			count := nextAccountFailureCount(account, attempt.Now)
			var lockUntil any
			if count >= 5 {
				lockUntil = attempt.Now.Add(30 * time.Minute)
			}
			if _, err = tx.ExecContext(ctx, updateAccountFailureQuery, count, attempt.Now, lockUntil, account.user.ID); err != nil {
				return nil, err
			}
		}
		if err = upsertIPFailure(ctx, tx, attempt); err != nil {
			return nil, err
		}
		if err = insertAttemptLogs(ctx, tx, userID, attempt, false); err != nil {
			return nil, err
		}
		if err = tx.Commit(); err != nil {
			return nil, err
		}
		return &model.LoginAttemptResult{Status: model.LoginStatusUnauthorized}, nil
	}

	if _, err = tx.ExecContext(ctx, resetAccountFailureQuery, account.user.ID); err != nil {
		return nil, err
	}
	if attempt.OldSessionID != "" {
		if _, err = tx.ExecContext(ctx, deleteOldSessionQuery, attempt.OldSessionID); err != nil {
			return nil, err
		}
	}
	if _, err = tx.ExecContext(ctx, insertSessionQuery,
		attempt.SessionID, account.user.ID, attempt.ExpiresAt, attempt.UserAgent, attempt.Now,
	); err != nil {
		return nil, err
	}
	if err = insertAttemptLogs(ctx, tx, account.user.ID, attempt, true); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &model.LoginAttemptResult{Status: model.LoginStatusSuccess, User: &account.user}, nil
}

// RecordInvalidRequest は入力を評価できない要求をACCESS_LOGだけへ記録します。
func (r *loginRepository) RecordInvalidRequest(ctx context.Context, ip string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, insertAccessLogQuery, uuid.NewString(), nil, ip, "POST auth/login", nil, now)
	return err
}

func selectIPLimit(ctx context.Context, tx *sql.Tx, ip string) (sql.NullTime, error) {
	var failed int
	var lastFailed time.Time
	var blockedUntil sql.NullTime
	err := tx.QueryRowContext(ctx, selectIPLimitQuery, ip).Scan(&failed, &lastFailed, &blockedUntil)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullTime{}, nil
	}
	return blockedUntil, err
}

func selectLoginAccount(ctx context.Context, tx *sql.Tx, email string) (*loginAccount, error) {
	account := &loginAccount{}
	err := tx.QueryRowContext(ctx, selectLoginAccountQuery, email).Scan(
		&account.user.ID,
		&account.user.Username,
		&account.user.Email,
		&account.hash,
		&account.failed,
		&account.lastFailed,
		&account.lockUntil,
		&account.user.CreatedAt,
		&account.user.UpdatedAt,
	)
	return account, err
}

func nextAccountFailureCount(account *loginAccount, now time.Time) int {
	if !account.lastFailed.Valid || now.Sub(account.lastFailed.Time) > 15*time.Minute {
		return 1
	}
	return account.failed + 1
}

func upsertIPFailure(ctx context.Context, tx *sql.Tx, attempt *model.LoginAttempt) error {
	_, err := tx.ExecContext(ctx, upsertIPFailureQuery, attempt.IP, attempt.Now, attempt.Now)
	return err
}

func insertAttemptLogs(ctx context.Context, tx *sql.Tx, userID any, attempt *model.LoginAttempt, success bool) error {
	if _, err := tx.ExecContext(ctx, insertLoginLogQuery,
		uuid.NewString(), userID, attempt.Email, attempt.IP, success, false, attempt.Now,
	); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, insertAccessLogQuery,
		uuid.NewString(), userID, attempt.IP, "POST auth/login", nil, attempt.Now,
	)
	return err
}

const selectIPLimitQuery = `
SELECT FAILED_COUNT, LAST_FAILED_AT, BLOCKED_UNTIL
FROM LOGIN_IP_RATE_LIMIT
WHERE IP_ADDRESS = $1
FOR UPDATE
`

const selectLoginAccountQuery = `
SELECT USER_ID, USER_NAME, EMAIL, PASSWORD_HASH, LOGIN_FAILED_COUNT,
       LOGIN_LAST_FAILED_AT, LOGIN_LOCK_UNTIL, CREATED_AT, UPDATED_AT
FROM LOGIN_ACCOUNT
WHERE EMAIL = $1 AND IS_DELETED = FALSE
FOR UPDATE
`

const updateAccountFailureQuery = `
UPDATE LOGIN_ACCOUNT
SET LOGIN_FAILED_COUNT = $1,
    LOGIN_LAST_FAILED_AT = $2,
    LOGIN_LOCK_UNTIL = $3,
    UPDATED_AT = $2
WHERE USER_ID = $4
`

const resetAccountFailureQuery = `
UPDATE LOGIN_ACCOUNT
SET LOGIN_FAILED_COUNT = 0,
    LOGIN_LAST_FAILED_AT = NULL,
    LOGIN_LOCK_UNTIL = NULL,
    UPDATED_AT = CURRENT_TIMESTAMP
WHERE USER_ID = $1
`

const upsertIPFailureQuery = `
INSERT INTO LOGIN_IP_RATE_LIMIT (IP_ADDRESS, FAILED_COUNT, LAST_FAILED_AT, BLOCKED_UNTIL, UPDATED_AT)
VALUES ($1, 1, $2, NULL, $3)
ON CONFLICT (IP_ADDRESS) DO UPDATE SET
    FAILED_COUNT = CASE
        WHEN EXCLUDED.LAST_FAILED_AT - LOGIN_IP_RATE_LIMIT.LAST_FAILED_AT > INTERVAL '5 minutes' THEN 1
        ELSE LOGIN_IP_RATE_LIMIT.FAILED_COUNT + 1
    END,
    LAST_FAILED_AT = EXCLUDED.LAST_FAILED_AT,
    BLOCKED_UNTIL = CASE
        WHEN (
            CASE
                WHEN EXCLUDED.LAST_FAILED_AT - LOGIN_IP_RATE_LIMIT.LAST_FAILED_AT > INTERVAL '5 minutes' THEN 1
                ELSE LOGIN_IP_RATE_LIMIT.FAILED_COUNT + 1
            END
        ) >= 30 THEN EXCLUDED.LAST_FAILED_AT + INTERVAL '15 minutes'
        ELSE LOGIN_IP_RATE_LIMIT.BLOCKED_UNTIL
    END,
    UPDATED_AT = EXCLUDED.UPDATED_AT
`

const deleteOldSessionQuery = `DELETE FROM LOGIN_SESSION WHERE SESSION_ID = $1`

const insertSessionQuery = `
INSERT INTO LOGIN_SESSION (SESSION_ID, USER_ID, EXPIRES_AT, USER_AGENT, CREATED_AT)
VALUES ($1, $2, $3, $4, $5)
`

const insertLoginLogQuery = `
INSERT INTO LOGIN_LOG (LOG_ID, USER_ID, EMAIL, IP_ADDRESS, IS_SUCCESS, IS_SESSION_USED, ACCESS_AT)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

const insertAccessLogQuery = `
INSERT INTO ACCESS_LOG (LOG_ID, USER_ID, IP_ADDRESS, ENDPOINT, RESOURCE_ID, ACCESS_AT)
VALUES ($1, $2, $3, $4, $5, $6)
`
