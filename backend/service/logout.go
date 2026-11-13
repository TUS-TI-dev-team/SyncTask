package service

import (
	"context"
	"fmt"
	"time"

	"synctask/backend/model"
	"synctask/backend/repository"
)

type LogoutService interface {
	Logout(context.Context, *model.LogoutInput) error
}

type LogoutDependencies struct {
	Now func() time.Time
}

type logoutService struct {
	repo repository.LogoutRepository
	deps LogoutDependencies
}

func NewLogoutService(repo repository.LogoutRepository, deps LogoutDependencies) LogoutService {
	if deps.Now == nil {
		deps.Now = time.Now
	}

	return &logoutService{repo: repo, deps: deps}
}

// Logout は入力検証後にRepositoryへ原初的なログアウト試行を依頼します。
//
// @spec 入力不正はreturn でnilを返す
// @spec セッション削除成功時・失敗時および対象セッション不在時はいずれもnilを返す
// @spec サーバ内の処理で未知のエラーが発生した場合はそれを受け取りそれをそのままhandlerへ返す
func (s *logoutService) Logout(ctx context.Context, input *model.LogoutInput) error {
	startedAt := s.deps.Now()
	if input == nil || !input.Validate() {
		return nil
	}

	attempt := &model.LogoutAttempt{
		SessionID: input.SessionID,
		Now:	   startedAt,
	}

	result, err := s.repo.AttemptLogout(ctx, attempt)
	if result == nil {
		return fmt.Errorf("logout result is nil")
	}
	if err != nil && result.Status != model.LogoutStatusDbFailed{
		return err
	}

	switch result.Status {
	case model.LogoutStatusSuccess:
		return nil
	case model.LogoutStatusNotFound:
		return nil
	case model.LogoutStatusDbFailed:
		return nil
	default:
		return fmt.Errorf("unknown logout status: %s", result.Status)
	}
}
