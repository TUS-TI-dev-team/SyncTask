package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// NormalizeLoginEmail はログイン用メールアドレスの前後空白を除去して小文字化します。
//
// @spec 前後の半角・全角スペース、タブ、改行などのUnicode空白を除去する。
// @spec 文字列内部は変更せず、英字を小文字化する。
func NormalizeLoginEmail(email string) string {
	return strings.ToLower(strings.TrimFunc(email, unicode.IsSpace))
}

// GenerateSecureToken はreaderから指定バイト数を読み、URLセーフなトークンを生成します。
//
// @spec sizeは正数でなければならない。
// @spec パディングなしのBase64 URL Encodingを使用する。
func GenerateSecureToken(reader io.Reader, size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("token size must be positive")
	}
	raw := make([]byte, size)
	if _, err := io.ReadFull(reader, raw); err != nil {
		return "", fmt.Errorf("failed to read secure random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
