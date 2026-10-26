package util

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskEmail(t *testing.T) {
	t.Run("正常系: ローカル部が4文字の場合に先頭4文字と固定10文字アスタリスクとドメインにマスクすること", func(t *testing.T) {
		got := MaskEmail("user@example.com")
		assert.Equal(t, "user**********@example.com", got)
	})

	t.Run("正常系: ローカル部が4文字より長い場合に先頭4文字と固定10文字アスタリスクとドメインにマスクすること", func(t *testing.T) {
		got1 := MaskEmail("username@example.com")
		assert.Equal(t, "user**********@example.com", got1)

		got2 := MaskEmail("john.doe@example.com")
		assert.Equal(t, "john**********@example.com", got2)
	})

	t.Run("境界値: ローカル部が3文字の場合に先頭1文字と固定10文字アスタリスクとドメインにマスクすること", func(t *testing.T) {
		got := MaskEmail("abc@example.com")
		assert.Equal(t, "a**********@example.com", got)
	})

	t.Run("境界値: ローカル部が1文字の場合に先頭1文字と固定10文字アスタリスクとドメインにマスクすること", func(t *testing.T) {
		got := MaskEmail("a@example.com")
		assert.Equal(t, "a**********@example.com", got)
	})

	t.Run("異常系: @が含まれない不正なメールアドレスの場合はそのまま返すこと", func(t *testing.T) {
		got := MaskEmail("invalidemail")
		assert.Equal(t, "invalidemail", got)
	})
}

func TestGenerateOTPSessionID(t *testing.T) {
	t.Run("正常系: @spec UUIDv4形式のOTPセッションIDを生成すること", func(t *testing.T) {
		id, err := GenerateOTPSessionID()
		require.NoError(t, err)
		parsed, parseErr := uuid.Parse(id)
		require.NoError(t, parseErr)
		assert.Equal(t, uuid.Version(4), parsed.Version())
		// URL-safe checks
		assert.NotContains(t, id, "/")
		assert.NotContains(t, id, "+")
		assert.NotContains(t, id, "=")
	})

	t.Run("正常系: 複数回生成した際に重複しないこと", func(t *testing.T) {
		set := make(map[string]struct{})
		for i := 0; i < 50; i++ {
			id, err := GenerateOTPSessionID()
			require.NoError(t, err)
			assert.NotContains(t, set, id)
			set[id] = struct{}{}
		}
	})
}

func TestGenerateOTP(t *testing.T) {
	t.Run("正常系: 8桁の大文字英数字（大文字英字および数字のみ）を生成すること", func(t *testing.T) {
		otpRegex := regexp.MustCompile(`^[0-9A-Z]{8}$`)
		for i := 0; i < 20; i++ {
			otp, err := GenerateOTP()
			require.NoError(t, err)
			assert.Len(t, otp, 8)
			assert.True(t, otpRegex.MatchString(otp), "OTP %s should match ^[0-9A-Z]{8}$", otp)
		}
	})

	t.Run("正常系: 複数回生成した際に重複しないこと", func(t *testing.T) {
		set := make(map[string]struct{})
		for i := 0; i < 50; i++ {
			otp, err := GenerateOTP()
			require.NoError(t, err)
			assert.NotContains(t, set, otp)
			set[otp] = struct{}{}
		}
	})
}
