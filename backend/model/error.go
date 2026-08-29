package model

import "fmt"

// ErrorDetail は各フィールドのエラー詳細を表します。
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorBody はエラーレスポンスの本体を表します。
type ErrorBody struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ErrorResponse は API 共通のエラーレスポンス形式を表します。
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// AppError はアプリケーション内部で利用するエラー型です。
type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Details    []ErrorDetail
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewBadRequestError は 400 Bad Request 用の AppError を生成します。
func NewBadRequestError(message string, details []ErrorDetail) *AppError {
	return &AppError{
		StatusCode: 400,
		Code:       "BAD_REQUEST",
		Message:    message,
		Details:    details,
	}
}

// NewUnauthorizedError は 401 Unauthorized 用の AppError を生成します。
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		StatusCode: 401,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}
