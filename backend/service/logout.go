package service

import (
	"context"
	"fmt"

	"synctask/backend/model"
	"synctask/backend/repository"
)

type LogoutService interface {
	Logout(context.Context, *model.LogoutInput) error
}

type logoutService struct {
	repo repository.LogoutRepository
}

func NewLogoutService(repo repository.LogoutRepository) LogoutService {
	return &logoutService{repo: repo}
}

func (s *logoutService) Logout(ctx context.Context, input *model.LogoutInput) error {
	if input == nil || !input.Validate() {
		return nil
	}

	attempt := &model.LogoutAttempt{
		SessionID: input.SessionID,
	}
	result, err := s.repo.AttemptLogout(ctx, attempt)
	if err != nil {
		return err
	}

	switch result.Status {
	case model.LogoutStatusSuccess:
		return nil
	case model.LogoutStatusNotFound:
		return nil
	default:
		return fmt.Errorf("unknown logout status: %s", result.Status)
	}
}
