package service

import(
	"context"

	"synctask/backend/model"
)

type LogoutService interface {
	Logout(context.Context, *model.LogoutInput) (error)
}
