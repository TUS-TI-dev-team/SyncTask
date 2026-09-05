package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginRequest_Validate(t *testing.T) {
	tests := []struct {
		name       string
		req        LoginRequest
		wantEmail  string
		wantFields []string
	}{
		{
			name:      "正常系: メールをトリムして小文字化しパスワードは加工しないこと",
			req:       LoginRequest{Email: " \t\n　User@Example.COM　\r\n", Password: " Password123! "},
			wantEmail: "user@example.com",
		},
		{
			name:      "正常系: メールの未知形式は入力検証では拒否しないこと",
			req:       LoginRequest{Email: "not-an-email", Password: "Password123!"},
			wantEmail: "not-an-email",
		},
		{
			name:       "異常系: メール欠落時にemail詳細を返すこと",
			req:        LoginRequest{Password: "Password123!"},
			wantFields: []string{"email"},
		},
		{
			name:       "異常系: パスワード欠落時にpassword詳細を返すこと",
			req:        LoginRequest{Email: "user@example.com"},
			wantEmail:  "user@example.com",
			wantFields: []string{"password"},
		},
		{
			name:       "境界値: パスワード7文字を拒否すること",
			req:        LoginRequest{Email: "user@example.com", Password: "Abc12!x"},
			wantEmail:  "user@example.com",
			wantFields: []string{"password"},
		},
		{
			name:       "境界値: パスワード129文字を拒否すること",
			req:        LoginRequest{Email: "user@example.com", Password: strings.Repeat("a", 129)},
			wantEmail:  "user@example.com",
			wantFields: []string{"password"},
		},
		{
			name:      "境界値: パスワード8文字を受理すること",
			req:       LoginRequest{Email: "user@example.com", Password: "Abc123!x"},
			wantEmail: "user@example.com",
		},
		{
			name:      "境界値: パスワード128文字を受理すること",
			req:       LoginRequest{Email: "user@example.com", Password: strings.Repeat("a", 128)},
			wantEmail: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPassword := tt.req.Password
			err := tt.req.Validate()
			assert.Equal(t, tt.wantEmail, tt.req.Email)
			assert.Equal(t, originalPassword, tt.req.Password)

			if len(tt.wantFields) == 0 {
				require.NoError(t, err)
				return
			}
			appErr, ok := err.(*AppError)
			require.True(t, ok)
			assert.Equal(t, 400, appErr.StatusCode)
			assert.Equal(t, "BAD_REQUEST", appErr.Code)
			require.Len(t, appErr.Details, len(tt.wantFields))
			for i, field := range tt.wantFields {
				assert.Equal(t, field, appErr.Details[i].Field)
			}
		})
	}
}
