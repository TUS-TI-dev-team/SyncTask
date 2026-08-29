package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeLoginEmail(t *testing.T) {
	t.Run("正常系: 前後の半角全角空白タブ改行を除去して小文字化すること", func(t *testing.T) {
		got := NormalizeLoginEmail(" \t\r\n　User.Name+Tag@Example.COM　\n")
		assert.Equal(t, "user.name+tag@example.com", got)
	})

	t.Run("正常系: 文字列内部の空白は変更しないこと", func(t *testing.T) {
		assert.Equal(t, "user name@example.com", NormalizeLoginEmail("User Name@Example.COM"))
	})

	t.Run("境界値: 空白のみの場合に空文字を返すこと", func(t *testing.T) {
		assert.Empty(t, NormalizeLoginEmail(" \t\r\n　"))
	})
}

func TestGenerateSecureToken(t *testing.T) {
	t.Run("正常系: 指定バイト数をURLセーフBase64で返すこと", func(t *testing.T) {
		source := bytes.Repeat([]byte{0xab}, 32)
		token, err := GenerateSecureToken(bytes.NewReader(source), 32)
		require.NoError(t, err)

		decoded, err := base64.RawURLEncoding.DecodeString(token)
		require.NoError(t, err)
		assert.Equal(t, source, decoded)
		assert.NotContains(t, token, "=")
	})

	t.Run("異常系: 乱数読み取り失敗時にエラーを返すこと", func(t *testing.T) {
		token, err := GenerateSecureToken(errorReader{}, 32)
		require.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("異常系: バイト数が0以下の場合にエラーを返すこと", func(t *testing.T) {
		token, err := GenerateSecureToken(bytes.NewReader(nil), 0)
		require.Error(t, err)
		assert.Empty(t, token)
	})
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("random source failed")
}
