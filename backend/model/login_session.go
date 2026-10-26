package model

import (
	"synctask/backend/util"
	"time"
)

type LoginSession struct {
	SessionID string
	UserID    string
	ExpiresAt time.Time
}

// TODO: LoginAttemptResult を参考に、FindBySessionIDのResultを定義する。
