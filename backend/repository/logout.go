package repository

import (
	"context"
	"database/sql"

	"synctask/backend/model"
)

type LogoutRepository interface {
	AttemptLogout(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error)
}

type logoutRepository struct {
	db *sql.DB
}

func NewLogoutRepository(db *sql.DB) LogoutRepository {
	return &logoutRepository{db: db}
}

func (r *logoutRepository) AttemptLogout(ctx context.Context, attempt *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
}
