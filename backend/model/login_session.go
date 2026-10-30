package model

import (
	"time"
)

type FindBySessionIDResult struct {
	UserID    string
	ExpiresAt time.Time
}
