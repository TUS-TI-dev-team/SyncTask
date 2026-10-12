package model

import (
	"unicode/utf8"
)

type LogoutInput struct {
	SessionID  string
	CsrfCookie string
	CsrfHeader string
}
