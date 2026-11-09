package repository

import (
	"context"

	"synctask/backend/model"
)

type LogoutRepository interface {
	AttemptLogout(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error)
}
