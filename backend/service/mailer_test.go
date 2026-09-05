package service

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// @spec LogMailer.SendOTP はログ出力時に平文OTPを出力せず、マスクすること。
func TestLogMailer_SendOTP(t *testing.T) {
	t.Run("正常系: 平文OTPを出力せずマスクしてログ出力すること", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		mailer := NewLogMailer()
		err := mailer.SendOTP(context.Background(), "test@example.com", "A1B2C3D4")
		require.NoError(t, err)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "test@example.com")
		assert.Contains(t, logOutput, "[REDACTED]")
		assert.NotContains(t, logOutput, "A1B2C3D4")
	})
}
