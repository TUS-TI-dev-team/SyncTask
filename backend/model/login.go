package model

import (
	"time"
	"unicode/utf8"

	"synctask/backend/util"
)

// LoginRequest はログインリクエストを表します。
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"Password123!"`
}

// Validate はログイン入力を検証し、メールアドレスだけを正規化します。
//
// @spec email と password は必須文字列である。
// @spec email は前後空白を除去して小文字化し、password は加工しない。
// @spec password は8文字以上128文字以下である。
func (r *LoginRequest) Validate() error {
	r.Email = util.NormalizeLoginEmail(r.Email)
	details := make([]ErrorDetail, 0, 2)

	if r.Email == "" {
		details = append(details, ErrorDetail{Field: "email", Message: "メールアドレスは必須です。"})
	}

	passwordLength := utf8.RuneCountInString(r.Password)
	if passwordLength == 0 {
		details = append(details, ErrorDetail{Field: "password", Message: "パスワードは必須です。"})
	} else if passwordLength < 8 || passwordLength > 128 {
		details = append(details, ErrorDetail{Field: "password", Message: "パスワードは8文字以上128文字以下で入力してください。"})
	}

	if len(details) > 0 {
		return NewBadRequestError("入力内容に不備があります。", details)
	}
	return nil
}

// LoginMetadata はHTTPリクエスト由来のログイン情報を表します。
type LoginMetadata struct {
	IP           string
	UserAgent    string
	OldSessionID string
}

// LoginUser はログイン成功時に返すユーザー情報です。
type LoginUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse はログイン成功レスポンスです。
type LoginResponse struct {
	User LoginUser `json:"user"`
}

// LoginServiceResult はHandlerへ渡すログイン結果とCookie情報です。
type LoginServiceResult struct {
	Response  LoginResponse
	SessionID string
	CSRFToken string
	MaxAge    time.Duration
}

// LoginAttempt はRepositoryで原子的に評価するログイン試行です。
type LoginAttempt struct {
	Email        string
	Password     string
	IP           string
	UserAgent    string
	OldSessionID string
	SessionID    string
	CSRFToken    string
	Now          time.Time
	ExpiresAt    time.Time
}

// LoginStatus はRepositoryで判定したログイン結果です。
type LoginStatus string

const (
	LoginStatusSuccess      LoginStatus = "success"
	LoginStatusUnauthorized LoginStatus = "unauthorized"
	LoginStatusRateLimited  LoginStatus = "rate_limited"
)

// LoginAttemptResult はRepositoryで確定したログイン試行結果です。
type LoginAttemptResult struct {
	Status     LoginStatus
	User       *LoginUser
	RetryAfter int
}
