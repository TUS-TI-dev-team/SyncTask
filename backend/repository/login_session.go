package repository

import (
	"context"
	"database/sql"

	"synctask/backend/model"
)

type LoginSessionRepository interface {
	FindBySessionID(context.Context, string) (*model.LoginSession, error)
}

type loginSessionRepository struct {
	db *sql.DB
}

func NewLoginSessionRepository(db *sql.DB) LoginSessionRepository {
	return &loginSessionRepository{db: db}
}

func (r *loginSessionRepository) FindBySessionID(ctx context.Context, sessionID string) (_ *model.FindBySessionIDResult, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// sessionID, err := selectLoginSession(ctx, tx, sessionID)
	// TODO: ここからの処理(値を取り出してResultとして返す)を実装する
}

func selectLoginSession(ctx context.Context, tx *sql.Tx, sessionID string) (*model.LoginSession, error) {
	session := &model.LoginSession{}
	err := tx.QueryRowContext(ctx, selectLoginSessionQuery, sessionID).Scan(
		&session.SessionID,
		&session.UserID,
		&session.ExpiresAt,
	)
	// TODO: selectLoginSession用のSQL文を書く
	// TODO: ここからの処理を書く
}
