package service

import (
	"context"
	"log"
)

// Mailer はメール送信インターフェースです。
type Mailer interface {
	// SendOTP は認証コード(OTP)をメール送信します。
	SendOTP(ctx context.Context, toEmail, otp string) error
}

// LogMailer は開発・テスト用のログ出力Mailerです。
type LogMailer struct{}

// NewLogMailer はLogMailerを生成します。
func NewLogMailer() Mailer {
	return &LogMailer{}
}

// SendOTP はメール送信をログ出力としてシミュレートします（平文OTPは秘匿します）。
func (m *LogMailer) SendOTP(ctx context.Context, toEmail, otp string) error {
	log.Printf("[LogMailer] Send OTP to %s (length: %d, [REDACTED])", toEmail, len(otp))
	return nil
}
