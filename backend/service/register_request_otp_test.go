package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"synctask/backend/model"
	"synctask/backend/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRegisterRequestOtpRepository struct {
	findActiveUserFunc       func(ctx context.Context, email string) (bool, error)
	findActiveOtpSessionFunc func(ctx context.Context, email string) (bool, error)
	createOtpSessionFunc     func(ctx context.Context, session *model.OtpSessionRecord) error
	updateDeliveryStatusFunc func(ctx context.Context, sessionID, status string, sendFailedCount int) error
	recordMailAuthLogFunc    func(ctx context.Context, log *model.MailAuthLogRecord) error
	recordAccessLogFunc      func(ctx context.Context, log *model.AccessLogRecord) error
	saveSessionWithLogsFunc  func(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error
}

func (m *mockRegisterRequestOtpRepository) FindActiveUserByEmail(ctx context.Context, email string) (bool, error) {
	if m.findActiveUserFunc != nil {
		return m.findActiveUserFunc(ctx, email)
	}
	return false, nil
}

func (m *mockRegisterRequestOtpRepository) FindActiveOtpSessionByEmail(ctx context.Context, email string) (bool, error) {
	if m.findActiveOtpSessionFunc != nil {
		return m.findActiveOtpSessionFunc(ctx, email)
	}
	return false, nil
}

func (m *mockRegisterRequestOtpRepository) CreateOtpSession(ctx context.Context, session *model.OtpSessionRecord) error {
	if m.createOtpSessionFunc != nil {
		return m.createOtpSessionFunc(ctx, session)
	}
	return nil
}

func (m *mockRegisterRequestOtpRepository) UpdateOtpSessionDeliveryStatus(ctx context.Context, sessionID, status string, sendFailedCount int) error {
	if m.updateDeliveryStatusFunc != nil {
		return m.updateDeliveryStatusFunc(ctx, sessionID, status, sendFailedCount)
	}
	return nil
}

func (m *mockRegisterRequestOtpRepository) RecordMailAuthLog(ctx context.Context, log *model.MailAuthLogRecord) error {
	if m.recordMailAuthLogFunc != nil {
		return m.recordMailAuthLogFunc(ctx, log)
	}
	return nil
}

func (m *mockRegisterRequestOtpRepository) RecordAccessLog(ctx context.Context, log *model.AccessLogRecord) error {
	if m.recordAccessLogFunc != nil {
		return m.recordAccessLogFunc(ctx, log)
	}
	return nil
}

func (m *mockRegisterRequestOtpRepository) SaveSessionWithLogs(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error {
	if m.saveSessionWithLogsFunc != nil {
		return m.saveSessionWithLogsFunc(ctx, session, mailLog, accessLog)
	}
	return nil
}

type mockMailer struct {
	sendOTPFunc func(ctx context.Context, toEmail, otp string) error
}

func (m *mockMailer) SendOTP(ctx context.Context, toEmail, otp string) error {
	if m.sendOTPFunc != nil {
		return m.sendOTPFunc(ctx, toEmail, otp)
	}
	return nil
}

func defaultTestDeps(now time.Time) RegisterRequestOtpDependencies {
	return RegisterRequestOtpDependencies{
		Now:   func() time.Time { return now },
		Sleep: func(d time.Duration) {},
		GenerateSessionID: func() (string, error) {
			return "otp_sess_test_12345", nil
		},
		GenerateOTP: func() (string, error) {
			return "ABC12345", nil
		},
		HashPassword: func(pw string) (string, error) {
			return "hashed_" + pw, nil
		},
		HashOTP: func(otp string) (string, error) {
			return "hashed_" + otp, nil
		},
		ResponseDelay: func() time.Duration {
			return 1000 * time.Millisecond
		},
	}
}

func TestRegisterRequestOtpService_RequestOtp(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)

	t.Run("正常系: 未登録かつ排他なしの場合に実OTP生成・DB登録・メール送信し遅延後に200結果を返すこと", func(t *testing.T) {
		sleepCalled := false
		mailerCalled := false
		updateDeliveryCalled := false
		savedSession := (*model.OtpSessionRecord)(nil)
		savedMailLog := (*model.MailAuthLogRecord)(nil)

		repo := &mockRegisterRequestOtpRepository{
			findActiveUserFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			findActiveOtpSessionFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			saveSessionWithLogsFunc: func(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error {
				savedSession = session
				savedMailLog = mailLog
				return nil
			},
			updateDeliveryStatusFunc: func(ctx context.Context, sessionID, status string, count int) error {
				updateDeliveryCalled = true
				assert.Equal(t, "otp_sess_test_12345", sessionID)
				assert.Equal(t, "sent", status)
				assert.Equal(t, 0, count)
				return nil
			},
		}

		mailer := &mockMailer{
			sendOTPFunc: func(ctx context.Context, toEmail, otp string) error {
				mailerCalled = true
				assert.Equal(t, "user@example.com", toEmail)
				assert.Equal(t, "ABC12345", otp)
				return nil
			},
		}

		deps := defaultTestDeps(now)
		deps.Sleep = func(d time.Duration) {
			sleepCalled = true
		}

		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "user@example.com",
			Password: "Password123!",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "otp_sess_test_12345", res.OtpSessionID)
		assert.Equal(t, "user**********@example.com", res.MaskedEmail)
		assert.Equal(t, 300, res.ExpiresInSeconds)
		assert.Equal(t, 60, res.CooldownSeconds)

		assert.True(t, mailerCalled)
		assert.True(t, updateDeliveryCalled)
		assert.True(t, sleepCalled)

		require.NotNil(t, savedSession)
		assert.False(t, savedSession.IsDummy)
		assert.Equal(t, "exampleUser", savedSession.PendingUsername.String)
		assert.Equal(t, "user@example.com", savedSession.PendingEmail.String)
		assert.Equal(t, "hashed_Password123!", savedSession.PendingPasswordHash.String)
		assert.Equal(t, "hashed_ABC12345", savedSession.OtpHash.String)

		require.NotNil(t, savedMailLog)
		assert.False(t, savedMailLog.IsDummy)
		assert.Equal(t, "user@example.com", savedMailLog.Email)
	})

	t.Run("正常系: 既に登録済みメールアドレスの場合はダミーセッションを作成し実メール送信をスキップして遅延後に200結果を返すこと", func(t *testing.T) {
		sleepCalled := false
		mailerCalled := false
		savedSession := (*model.OtpSessionRecord)(nil)
		savedMailLog := (*model.MailAuthLogRecord)(nil)

		repo := &mockRegisterRequestOtpRepository{
			findActiveUserFunc: func(ctx context.Context, email string) (bool, error) {
				return true, nil // 登録済み
			},
			saveSessionWithLogsFunc: func(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error {
				savedSession = session
				savedMailLog = mailLog
				return nil
			},
		}

		mailer := &mockMailer{
			sendOTPFunc: func(ctx context.Context, toEmail, otp string) error {
				mailerCalled = true
				return nil
			},
		}

		deps := defaultTestDeps(now)
		deps.Sleep = func(d time.Duration) {
			sleepCalled = true
		}

		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "registered@example.com",
			Password: "Password123!",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "otp_sess_test_12345", res.OtpSessionID)
		assert.Equal(t, "regi**********@example.com", res.MaskedEmail)

		assert.False(t, mailerCalled, "実メール送信はスキップされること")
		assert.True(t, sleepCalled, "遅延が適用されること")

		require.NotNil(t, savedSession)
		assert.True(t, savedSession.IsDummy)
		assert.False(t, savedSession.PendingUsername.Valid)
		assert.False(t, savedSession.PendingEmail.Valid)
		assert.False(t, savedSession.PendingPasswordHash.Valid)
		assert.False(t, savedSession.OtpHash.Valid)

		require.NotNil(t, savedMailLog)
		assert.True(t, savedMailLog.IsDummy)
		assert.Equal(t, "registered@example.com", savedMailLog.Email)
	})

	t.Run("正常系: 既に有効なOTPセッションが存在する場合はダミーセッションを作成し実メール送信をスキップして遅延後に200結果を返すこと", func(t *testing.T) {
		mailerCalled := false
		savedSession := (*model.OtpSessionRecord)(nil)

		repo := &mockRegisterRequestOtpRepository{
			findActiveUserFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			findActiveOtpSessionFunc: func(ctx context.Context, email string) (bool, error) {
				return true, nil // 有効OTPセッションあり
			},
			saveSessionWithLogsFunc: func(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error {
				savedSession = session
				return nil
			},
		}

		mailer := &mockMailer{
			sendOTPFunc: func(ctx context.Context, toEmail, otp string) error {
				mailerCalled = true
				return nil
			},
		}

		deps := defaultTestDeps(now)
		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "active-otp@example.com",
			Password: "Password123!",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.False(t, mailerCalled)
		require.NotNil(t, savedSession)
		assert.True(t, savedSession.IsDummy)
	})

	t.Run("異常系: 入力バリデーションエラーの場合は遅延なしでエラーを返しメール送信しないこと", func(t *testing.T) {
		accessLogRecorded := false
		sleepCalled := false
		mailerCalled := false

		repo := &mockRegisterRequestOtpRepository{
			recordAccessLogFunc: func(ctx context.Context, log *model.AccessLogRecord) error {
				accessLogRecorded = true
				return nil
			},
		}
		mailer := &mockMailer{
			sendOTPFunc: func(ctx context.Context, toEmail, otp string) error {
				mailerCalled = true
				return nil
			},
		}
		deps := defaultTestDeps(now)
		deps.Sleep = func(d time.Duration) {
			sleepCalled = true
		}

		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "", // 不正
			Email:    "invalid",
			Password: "",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		require.Error(t, err)
		assert.Nil(t, res)
		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.False(t, sleepCalled, "入力不正時は遅延なし")
		assert.False(t, mailerCalled, "メール送信なし")
		assert.True(t, accessLogRecorded, "ACCESS_LOG記録")
	})

	t.Run("異常系: 実メール送信失敗時にステータスをsendableと失敗回数1に更新し503エラーを返すこと", func(t *testing.T) {
		updateDeliveryCalled := false
		repo := &mockRegisterRequestOtpRepository{
			updateDeliveryStatusFunc: func(ctx context.Context, sessionID, status string, count int) error {
				updateDeliveryCalled = true
				assert.Equal(t, "otp_sess_test_12345", sessionID)
				assert.Equal(t, "sendable", status)
				assert.Equal(t, 1, count)
				return nil
			},
		}
		mailer := &mockMailer{
			sendOTPFunc: func(ctx context.Context, toEmail, otp string) error {
				return errors.New("smtp connection refused")
			},
		}
		deps := defaultTestDeps(now)
		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "user@example.com",
			Password: "Password123!",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		require.Error(t, err)
		assert.Nil(t, res)
		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 503, appErr.StatusCode)
		assert.Equal(t, "OTP_DELIVERY_FAILED", appErr.Code)
		assert.True(t, updateDeliveryCalled)
	})

	t.Run("異常系: リポジトリエラー発生時にエラーを返すこと", func(t *testing.T) {
		dbErr := errors.New("db query error")
		repo := &mockRegisterRequestOtpRepository{
			findActiveUserFunc: func(ctx context.Context, email string) (bool, error) {
				return false, dbErr
			},
		}
		mailer := &mockMailer{}
		deps := defaultTestDeps(now)
		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "user@example.com",
			Password: "Password123!",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, res)
	})

	t.Run("準正常系: 一意制約競合時にロールバックしてダミーセッションを作成し実メール送信をスキップして遅延後に200結果を返すこと", func(t *testing.T) {
		mailerCalled := false
		sleepCalled := false
		savedSessions := []*model.OtpSessionRecord{}

		saveCallCount := 0
		repo := &mockRegisterRequestOtpRepository{
			findActiveUserFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			findActiveOtpSessionFunc: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
			saveSessionWithLogsFunc: func(ctx context.Context, session *model.OtpSessionRecord, mailLog *model.MailAuthLogRecord, accessLog *model.AccessLogRecord) error {
				saveCallCount++
				savedSessions = append(savedSessions, session)
				if saveCallCount == 1 {
					// 1回目の通常セッション保存で一意制約競合
					return repository.ErrConflict
				}
				// 2回目のダミーセッション保存は成功
				return nil
			},
		}

		mailer := &mockMailer{
			sendOTPFunc: func(ctx context.Context, toEmail, otp string) error {
				mailerCalled = true
				return nil
			},
		}

		deps := defaultTestDeps(now)
		deps.Sleep = func(d time.Duration) {
			sleepCalled = true
		}

		svc := NewRegisterRequestOtpService(repo, mailer, deps)
		req := &model.RegisterRequestOtpRequest{
			Username: "exampleUser",
			Email:    "conflict@example.com",
			Password: "Password123!",
		}

		res, err := svc.RequestOtp(context.Background(), req, "192.0.2.1")
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "otp_sess_test_12345", res.OtpSessionID)
		assert.Equal(t, "conf**********@example.com", res.MaskedEmail)
		assert.False(t, mailerCalled, "実メール送信はスキップされること")
		assert.True(t, sleepCalled, "遅延が適用されること")

		require.Len(t, savedSessions, 2)
		// 2回目に保存されたセッションはダミーセッションであること
		dummySession := savedSessions[1]
		assert.True(t, dummySession.IsDummy)
		assert.False(t, dummySession.PendingUsername.Valid)
		assert.False(t, dummySession.PendingEmail.Valid)
	})
}
