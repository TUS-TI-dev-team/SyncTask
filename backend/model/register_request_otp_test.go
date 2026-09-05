package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRequestOtpRequest_Validate(t *testing.T) {
	t.Run("正常系: 有効な入力の場合に検証を通過すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, "exampleUser", req.Username)
		assert.Equal(t, "user@example.com", req.Email)
		assert.Equal(t, "Password123!", req.Password)
	})

	t.Run("正常系: ユーザー名の前後空白をトリムすること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "  testUser  ",
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, "testUser", req.Username)
	})

	t.Run("正常系: メールアドレスの前後空白をトリムして小文字化すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "  User.Name+Tag@Example.COM  ",
			Password: "Password123!",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, "user.name+tag@example.com", req.Email)
	})

	t.Run("正常系: パスワードの前後空白はトリムせず保持すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: " Password123! ",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, " Password123! ", req.Password)
	})

	t.Run("正常系: ユーザー名またはメールローカル部が3文字以下の場合はパスワードに含まれていても通過すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "usr",
			Email:    "abc@example.com",
			Password: "usrAbc123!", // "usr" と "abc" (各3文字) を含む
		}
		err := req.Validate()
		require.NoError(t, err)
	})

	t.Run("境界値: ユーザー名最小長2文字の場合に通過すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "ab",
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, "ab", req.Username)
	})

	t.Run("境界値: ユーザー名最大長20文字の場合に通過すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "abcdefghijklmnopqrst", // 20 chars
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, "abcdefghijklmnopqrst", req.Username)
	})

	t.Run("境界値: パスワード最小長8文字の場合に通過すること", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: "Pass12!x", // 8 chars (upper, lower, digit, symbol)
		}
		err := req.Validate()
		require.NoError(t, err)
	})

	t.Run("境界値: パスワード最大長128文字の場合に通過すること", func(t *testing.T) {
		// 128文字 (大文字, 小文字, 数字, 記号を含む)
		pw := "Pass12!" + strings.Repeat("a", 121)
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: pw,
		}
		err := req.Validate()
		require.NoError(t, err)
	})

	t.Run("境界値: メールアドレス最大長255文字の場合に通過すること", func(t *testing.T) {
		localPart := strings.Repeat("a", 64)
		domainPart := strings.Repeat("b", 186) + ".com" // 64 + 1 + 190 = 255
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    localPart + "@" + domainPart,
			Password: "Password123!",
		}
		err := req.Validate()
		require.NoError(t, err)
	})

	t.Run("異常系: ユーザー名が未指定の場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "   ",
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		require.NotEmpty(t, appErr.Details)
		assert.Equal(t, "username", appErr.Details[0].Field)
	})

	t.Run("異常系: ユーザー名が1文字の場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "a",
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "username", appErr.Details[0].Field)
	})

	t.Run("異常系: ユーザー名が21文字以上の場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "abcdefghijklmnopqrstu", // 21 chars
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "username", appErr.Details[0].Field)
	})

	t.Run("異常系: ユーザー名に記号や不正な文字が含まれる場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "user_name!",
			Email:    "user@example.com",
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "username", appErr.Details[0].Field)
	})

	t.Run("異常系: メールアドレスが未指定の場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "   ",
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "email", appErr.Details[0].Field)
	})

	t.Run("異常系: メールアドレスの形式が不正な場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "invalid-email-format",
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "email", appErr.Details[0].Field)
	})

	t.Run("境界値: メールアドレスが256文字以上の場合はエラーを返すこと", func(t *testing.T) {
		localPart := strings.Repeat("a", 64)
		domainPart := strings.Repeat("b", 187) + ".com" // 64 + 1 + 191 = 256
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    localPart + "@" + domainPart,
			Password: "Password123!",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "email", appErr.Details[0].Field)
	})

	t.Run("異常系: パスワードが未指定の場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: "",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "password", appErr.Details[0].Field)
	})

	t.Run("境界値: パスワードが7文字以下の場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: "Pass12!", // 7 chars
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "password", appErr.Details[0].Field)
	})

	t.Run("境界値: パスワードが129文字以上の場合はエラーを返すこと", func(t *testing.T) {
		pw := "Pass12!" + strings.Repeat("a", 122) // 129 chars
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: pw,
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "password", appErr.Details[0].Field)
	})

	t.Run("異常系: パスワードが文字種要件（4種中3種以上）を満たさない場合はエラーを返すこと", func(t *testing.T) {
		tests := []struct {
			name     string
			password string
		}{
			{"小文字+数字のみ（2種）", "password12345"},
			{"大文字+小文字のみ（2種）", "PasswordPassword"},
			{"数字+記号のみ（2種）", "12345678!@#$%"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := RegisterRequestOtpRequest{
					Username: "testUser",
					Email:    "user@example.com",
					Password: tc.password,
				}
				err := req.Validate()
				require.Error(t, err)
				appErr, ok := err.(*AppError)
				require.True(t, ok)
				assert.Equal(t, 400, appErr.StatusCode)
				assert.Equal(t, "password", appErr.Details[0].Field)
			})
		}
	})

	t.Run("異常系: パスワードに4文字以上のユーザー名（大文字小文字不問）を含む場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "user@example.com",
			Password: "MyTESTUSER123!", // "TESTUSER" を含む
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "password", appErr.Details[0].Field)
	})

	t.Run("異常系: パスワードに4文字以上のメールローカル部（大文字小文字不問）を含む場合はエラーを返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "testUser",
			Email:    "myEmail@example.com",
			Password: "MYEMAIL12345!", // "MYEMAIL" を含む
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "password", appErr.Details[0].Field)
	})

	t.Run("異常系: 複数フィールドに違反がある場合に全エラー詳細を返すこと", func(t *testing.T) {
		req := RegisterRequestOtpRequest{
			Username: "",
			Email:    "",
			Password: "",
		}
		err := req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		require.Len(t, appErr.Details, 3)
		fields := []string{appErr.Details[0].Field, appErr.Details[1].Field, appErr.Details[2].Field}
		assert.Contains(t, fields, "username")
		assert.Contains(t, fields, "email")
		assert.Contains(t, fields, "password")
	})
}
