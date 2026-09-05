package util

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// MaskEmail はメールアドレスをマスク処理します。
//
// @spec 先頭4文字（ローカル部が4文字未満の場合は先頭1文字）とドメイン以外を固定10文字のアスタリスク（'**********'）でマスクする。
// @spec 不正形式や空文字の場合は加工せず返すか安全に処理する。
func MaskEmail(email string) string {
	atIdx := strings.LastIndex(email, "@")
	if atIdx == -1 {
		return email
	}

	localPart := email[:atIdx]
	domainPart := email[atIdx:]

	localRunes := []rune(localPart)
	localLen := len(localRunes)
	if localLen == 0 {
		return email
	}

	prefixLen := 4
	if localLen < 4 {
		prefixLen = 1
	}

	prefix := string(localRunes[:prefixLen])
	return prefix + "**********" + domainPart
}

// GenerateOTPSessionID はOTPセッション用の推測困難なIDを生成します。
//
// @spec 接頭辞 'otp_sess_' とURL-safeなランダム文字列で構成される。
func GenerateOTPSessionID() (string, error) {
	token, err := GenerateSecureToken(rand.Reader, 24)
	if err != nil {
		return "", err
	}
	return "otp_sess_" + token, nil
}

const otpChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateOTP は8桁のOTPコードを生成します。
//
// @spec 暗号論的擬似乱数生成器を使用し、大文字英字および数字（全36種）からなる8桁の文字列を生成する。
func GenerateOTP() (string, error) {
	const otpLength = 8
	charsLen := big.NewInt(int64(len(otpChars)))
	result := make([]byte, otpLength)

	for i := 0; i < otpLength; i++ {
		n, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", err
		}
		result[i] = otpChars[n.Int64()]
	}

	return string(result), nil
}
