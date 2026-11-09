package service

import (
	"context"
	"database/sql"
	"errors"
	"synctask/backend/model"
	"synctask/backend/repository"
	"time"
)

type LoginSessionService interface {
	ValidateSession(ctx context.Context, ID string) (*model.LoginSessionServiceResult, error)
}

type LoginSessionDependencies struct {
	Now func() time.Time
}

type loginSessionService struct {
	repo repository.LoginSessionRepository
	deps LoginSessionDependencies
}

func NewLoginSessionService(repo repository.LoginSessionRepository, deps LoginSessionDependencies) LoginSessionService {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	return &loginSessionService{repo: repo, deps: deps}
}

func (s *loginSessionService) ValidateSession(ctx context.Context, ID string) (*model.LoginSessionServiceResult, error) {
	res, err := s.repo.FindBySessionID(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.NewUnauthorizedError("session not found")
		}
		return nil, err
	}

	if res.ExpiresAt.Before(s.deps.Now()) {
		return nil, model.NewUnauthorizedError("session is already expired")
	}

	return &model.LoginSessionServiceResult{UserID: res.UserID, ExpiresAt: res.ExpiresAt}, nil
}
