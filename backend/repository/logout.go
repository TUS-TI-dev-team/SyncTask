package repository

import (
	"context"
	"database/sql"
	"errors"

	"synctask/backend/model"
	"github.com/google/uuid"
)


// LogoutRepository はログアウト試行の永続化境界です。
type LogoutRepository interface {
	AttemptLogout(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error)
}

type logoutRepository struct {
	db *sql.DB
}

func NewLogoutRepository(db *sql.DB) LogoutRepository {
	return &logoutRepository{db: db}
}

// AttemptLogout はログインセッションの有無確認、対象セッションの削除、ログイン後の再認証状況のリセットを単一トランザクションで処理します。
//
// @spec 受け取ったセッションIDからUserID を検索
// @spec 該当セッションタプルの削除、再認証カウントのリセットを原子的に行う
func (r *logoutRepository) AttemptLogout(ctx context.Context, attempt *model.LogoutAttempt) (_ *model.LogoutAttemptResult, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Get UId
	uid, expires_at, err := selectLoginUserID(ctx, tx, attempt.SessionID)
	if err != nil {
		return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed },err
	}

	if _, parseErr := uuid.Parse(uid.String); !uid.Valid || parseErr != nil {
		if err = tx.Commit(); err != nil {
			return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
		}
		return &model.LogoutAttemptResult{ Status: model.LogoutStatusNotFound }, nil
	}

	// delete Session
	if _, err = tx.ExecContext(ctx, deleteLoginSessionQuery, attempt.SessionID); err != nil {
		return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
	}

	// reset reauth count
	if expires_at.Valid && expires_at.Time.After(attempt.Now) {
		if _, err = tx.ExecContext(ctx, resetReauthFailedCountQuery, attempt.Now, uid.String); err != nil {
			return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
		}
	}

	if err = tx.Commit(); err != nil {
		return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
	}

	return &model.LogoutAttemptResult{ Status: model.LogoutStatusSuccess }, nil
}

func selectLoginUserID(ctx context.Context, tx *sql.Tx, sid string) (uid sql.NullString, expires_at sql.NullTime, err error) {
	err = tx.QueryRowContext(ctx, selectTargetUserIDQuery, sid).Scan(&uid, &expires_at)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullString{}, sql.NullTime{}, nil
	}
	return uid, expires_at, err
}

const selectTargetUserIDQuery = `
SELECT USER_ID, EXPIRES_AT
FROM LOGIN_SESSION
WHERE SESSION_ID = $1
`

const deleteLoginSessionQuery = `
DELETE FROM LOGIN_SESSION
WHERE SESSION_ID = $1
`

const resetReauthFailedCountQuery = `
UPDATE LOGIN_ACCOUNT
SET REAUTH_FAILED_COUNT = 0,
    REAUTH_LAST_FAILED_AT = NULL,
    UPDATED_AT = $1
WHERE USER_ID = $2
`
