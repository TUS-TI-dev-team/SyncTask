package model

import (
	"database/sql"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`)
)

// RegisterRequestOtpRequest は新規登録OTP発行リクエストを表します。
type RegisterRequestOtpRequest struct {
	Username string `json:"username" example:"exampleUser"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"Password123!"`
}

// Validate は新規登録OTP発行リクエストを検証します。
//
// @spec username は前後の空白をトリムし、2〜20文字の半角英数字であること。
// @spec email は前後の空白をトリムし小文字化して、RFC 5322準拠形式かつ255文字以下であること。
// @spec password はトリムせず、8〜128文字で、英大文字・英小文字・数字・記号の4種中3種以上を含むこと。
// @spec password は、4文字以上のユーザー名またはメールローカル部（大文字小文字不問）を含まないこと。
// @spec 不備がある場合は model.NewBadRequestError を返すこと。
func (r *RegisterRequestOtpRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	var details []ErrorDetail

	// 1. username のバリデーション
	usernameLen := utf8.RuneCountInString(r.Username)
	if usernameLen == 0 {
		details = append(details, ErrorDetail{Field: "username", Message: "ユーザー名は必須です。"})
	} else if usernameLen < 2 || usernameLen > 20 {
		details = append(details, ErrorDetail{Field: "username", Message: "ユーザー名は2文字以上20文字以下で入力してください。"})
	} else if !usernameRegex.MatchString(r.Username) {
		details = append(details, ErrorDetail{Field: "username", Message: "ユーザー名は半角英数字のみ使用できます。"})
	}

	// 2. email のバリデーション
	emailLen := utf8.RuneCountInString(r.Email)
	if emailLen == 0 {
		details = append(details, ErrorDetail{Field: "email", Message: "メールアドレスは必須です。"})
	} else if len(r.Email) > 255 {
		details = append(details, ErrorDetail{Field: "email", Message: "メールアドレスは255文字以下で入力してください。"})
	} else if !emailRegex.MatchString(r.Email) {
		details = append(details, ErrorDetail{Field: "email", Message: "メールアドレスの形式が正しくありません。"})
	}

	// 3. password のバリデーション
	passwordLen := utf8.RuneCountInString(r.Password)
	if passwordLen == 0 {
		details = append(details, ErrorDetail{Field: "password", Message: "パスワードは必須です。"})
	} else if passwordLen < 8 || passwordLen > 128 {
		details = append(details, ErrorDetail{Field: "password", Message: "パスワードは8文字以上128文字以下で入力してください。"})
	} else {
		// 文字種チェック（大文字、小文字、数字、記号のうち3種以上）
		var hasUpper, hasLower, hasDigit, hasSymbol bool
		for _, ch := range r.Password {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsLower(ch):
				hasLower = true
			case unicode.IsDigit(ch):
				hasDigit = true
			case (ch >= 0x21 && ch <= 0x2f) || (ch >= 0x3a && ch <= 0x40) || (ch >= 0x5b && ch <= 0x60) || (ch >= 0x7b && ch <= 0x7e):
				hasSymbol = true
			}
		}

		typeCount := 0
		if hasUpper {
			typeCount++
		}
		if hasLower {
			typeCount++
		}
		if hasDigit {
			typeCount++
		}
		if hasSymbol {
			typeCount++
		}

		if typeCount < 3 {
			details = append(details, ErrorDetail{Field: "password", Message: "パスワードは英大文字、英小文字、数字、記号のうち3種類以上を含める必要があります。"})
		} else {
			// 禁止パターンチェック: 4文字以上のユーザー名またはメールローカル部（大文字小文字不問）
			pwLower := strings.ToLower(r.Password)

			// ユーザー名チェック
			if usernameLen >= 4 && strings.Contains(pwLower, strings.ToLower(r.Username)) {
				details = append(details, ErrorDetail{Field: "password", Message: "パスワードにユーザー名を含めることはできません。"})
			} else {
				// メールローカル部チェック
				localPart := r.Email
				if atIdx := strings.Index(r.Email, "@"); atIdx != -1 {
					localPart = r.Email[:atIdx]
				}
				if utf8.RuneCountInString(localPart) >= 4 && strings.Contains(pwLower, strings.ToLower(localPart)) {
					details = append(details, ErrorDetail{Field: "password", Message: "パスワードにメールアドレスのローカル部を含めることはできません。"})
				}
			}
		}
	}

	if len(details) > 0 {
		return NewBadRequestError("入力内容に不備があります。", details)
	}

	return nil
}

// RegisterRequestOtpResponse は新規登録OTP発行成功レスポンスを表します。
type RegisterRequestOtpResponse struct {
	OtpSessionID     string `json:"otp_session_id" example:"otp_sess_a1b2c3d4e5"`
	MaskedEmail      string `json:"masked_email" example:"user**********@example.com"`
	ExpiresInSeconds int    `json:"expires_in_seconds" example:"300"`
	CooldownSeconds  int    `json:"cooldown_seconds" example:"60"`
}

// OtpSessionRecord は OTP_SESSION テーブルのレコード構造を表します。
type OtpSessionRecord struct {
	OtpSessionID        string
	Purpose             string
	UserID              sql.NullString
	PendingUsername     sql.NullString
	PendingEmail        sql.NullString
	MaskedEmail         string
	PendingPasswordHash sql.NullString
	OtpHash             sql.NullString
	Status              string
	IsDummy             bool
	AttemptCount        int
	SendCount           int
	SendFailedCount     int
	DeliveryStatus      string
	LastSentAt          time.Time
	OtpExpiresAt        time.Time
	SessionExpiresAt    time.Time
	CreatedAt           time.Time
}

// MailAuthLogRecord は MAIL_AUTH_LOG テーブルのレコード構造を表します。
type MailAuthLogRecord struct {
	LogID     string
	UserID    sql.NullString
	Email     string
	AuthType  string
	IPAddress string
	EventType string
	IsSuccess bool
	IsDummy   bool
	AccessAt  time.Time
}

// AccessLogRecord は ACCESS_LOG テーブルのレコード構造を表します。
type AccessLogRecord struct {
	LogID      string
	UserID     sql.NullString
	IPAddress  string
	Endpoint   string
	ResourceID sql.NullString
	AccessAt   time.Time
}
