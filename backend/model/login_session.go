package model

import (
	"time"
)

type LoginSessionServiceResult struct {
	UserID    string
	ExpiresAt time.Time
}

type FindBySessionIDResult struct {
	UserID    string
	ExpiresAt time.Time
}
