package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"synctask/backend/model"
)

type LoginSessionRepository interface {
	FindBySessionID(context.Context, string) (*model.FindBySessionIDResult, error)
}

type loginSessionRepository struct {
	db *sql.DB
}

func NewLoginSessionRepository(db *sql.DB) LoginSessionRepository {
	return &loginSessionRepository{db: db}
}

type loginSession struct {
	SessionID string
	UserID    string
	ExpiresAt time.Time
}

func (r *loginSessionRepository) FindBySessionID(ctx context.Context, ID string) (*model.FindBySessionIDResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	s, err := selectLoginSession(ctx, tx, ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return &model.FindBySessionIDResult{UserID: s.UserID, ExpiresAt: s.ExpiresAt}, nil
}

func selectLoginSession(ctx context.Context, tx *sql.Tx, ID string) (*loginSession, error) {
	s := &loginSession{}

	err := tx.QueryRowContext(ctx, selectLoginSessionQuery, ID).Scan(
		&s.SessionID,
		&s.UserID,
		&s.ExpiresAt,
	)

	return s, err
}

const selectLoginSessionQuery = `
SELECT SESSION_ID, USER_ID, EXPIRES_AT
FROM LOGIN_SESSION
WHERE SESSION_ID = $1
FOR UPDATE
`
