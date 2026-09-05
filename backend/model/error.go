package model

import "fmt"

// ErrorDetail は各フィールドのエラー詳細を表します。
type ErrorDetail struct {
	// Field はエラーが発生したフィールド名です
	Field string `json:"field" example:"title"`
	// Message はフィールドに関するエラーメッセージです
	Message string `json:"message" example:"タイトルは必須です。"`
}

// ErrorBody はエラーレスポンスの本体を表します。
type ErrorBody struct {
	// Code はエラーコードです（例: "BAD_REQUEST", "UNAUTHORIZED", "NOT_FOUND", "INTERNAL_SERVER_ERROR"）
	Code string `json:"code" example:"BAD_REQUEST"`
	// Message はエラーメッセージです
	Message string `json:"message" example:"入力内容に不備があります。"`
	// Details は各フィールドのエラー詳細リストです
	Details []ErrorDetail `json:"details"`
}

// ErrorResponse は API 共通のエラーレスポンス形式を表します。
type ErrorResponse struct {
	// Error はエラーの詳細情報です
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
	if details == nil {
		details = []ErrorDetail{}
	}
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
		Details:    []ErrorDetail{},
	}
}

// NewNotFoundError は 404 Not Found 用の AppError を生成します。
func NewNotFoundError(message string) *AppError {
	return &AppError{
		StatusCode: 404,
		Code:       "NOT_FOUND",
		Message:    message,
		Details:    []ErrorDetail{},
	}
}

// NewErrorResponse は API 共通形式の ErrorResponse を生成します。
func NewErrorResponse(code, message string, details []ErrorDetail) ErrorResponse {
	if details == nil {
		details = []ErrorDetail{}
	}
	return ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
