package model

import (
	"github.com/google/uuid"
)

// LogoutInput はログアウトに必要な情報をまとめる.
type LogoutInput struct {
	SessionID string `json:"sid"`
}

const reqVersion uuid.Version = 4

func (r *LogoutInput) Validate() bool {
	//func (r* LogoutInput) Validate() uuid.Version {
	parsedUUID, err := uuid.Parse(r.SessionID)
	if err != nil {
		return false
	}
	//return parsedUUID.Version()
	return parsedUUID.Version() == reqVersion
}

type LogoutStatus string

const (
	LogoutStatusNotFound LogoutStatus = "notfound"
	LogoutStatusSuccess  LogoutStatus = "success"
)

type LogoutAttempt struct {
	SessionID string
}

type LogoutAttemptResult struct {
	Status LogoutStatus
}
