package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRequestOtpRepository_FindActiveUserByEmail(t *testing.T) {
	t.Run("正常系: 有効アカウントが存在する場合にtrueを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectQuery("SELECT.+FROM LOGIN_ACCOUNT.+WHERE EMAIL = .+ AND IS_DELETED = FALSE").
			WithArgs("user@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

		found, err := repo.FindActiveUserByEmail(context.Background(), "user@example.com")
		require.NoError(t, err)
		assert.True(t, found)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("正常系: アカウントが存在しない場合にfalseを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectQuery("SELECT.+FROM LOGIN_ACCOUNT.+WHERE EMAIL = .+ AND IS_DELETED = FALSE").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		found, err := repo.FindActiveUserByEmail(context.Background(), "notfound@example.com")
		require.NoError(t, err)
		assert.False(t, found)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー時にエラーを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)
		dbErr := errors.New("db error")

		mock.ExpectQuery("SELECT.+FROM LOGIN_ACCOUNT.+WHERE EMAIL = .+ AND IS_DELETED = FALSE").
			WithArgs("user@example.com").
			WillReturnError(dbErr)

		found, err := repo.FindActiveUserByEmail(context.Background(), "user@example.com")
		assert.ErrorIs(t, err, dbErr)
		assert.False(t, found)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegisterRequestOtpRepository_FindActiveOtpSessionByEmail(t *testing.T) {
	t.Run("正常系: 有効期限内のactiveまたはverifiedなOTPセッションが存在する場合にtrueを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectQuery("SELECT.+FROM OTP_SESSION.+WHERE PENDING_EMAIL = .+").
			WithArgs("user@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

		found, err := repo.FindActiveOtpSessionByEmail(context.Background(), "user@example.com")
		require.NoError(t, err)
		assert.True(t, found)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("正常系: セッションが存在しない場合にfalseを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectQuery("SELECT.+FROM OTP_SESSION.+WHERE PENDING_EMAIL = .+").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		found, err := repo.FindActiveOtpSessionByEmail(context.Background(), "notfound@example.com")
		require.NoError(t, err)
		assert.False(t, found)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー時にエラーを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)
		dbErr := errors.New("db error")

		mock.ExpectQuery("SELECT.+FROM OTP_SESSION.+WHERE PENDING_EMAIL = .+").
			WithArgs("user@example.com").
			WillReturnError(dbErr)

		found, err := repo.FindActiveOtpSessionByEmail(context.Background(), "user@example.com")
		assert.ErrorIs(t, err, dbErr)
		assert.False(t, found)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegisterRequestOtpRepository_CreateOtpSession(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	session := &model.OtpSessionRecord{
		OtpSessionID:        "otp_sess_12345",
		Purpose:             "SIGNUP",
		PendingUsername:     sql.NullString{String: "testUser", Valid: true},
		PendingEmail:        sql.NullString{String: "user@example.com", Valid: true},
		MaskedEmail:         "user**********@example.com",
		PendingPasswordHash: sql.NullString{String: "pwdhash", Valid: true},
		OtpHash:             sql.NullString{String: "otphash", Valid: true},
		Status:              "active",
		IsDummy:             false,
		AttemptCount:        0,
		SendCount:           0,
		SendFailedCount:     0,
		DeliveryStatus:      "pending",
		LastSentAt:          now,
		OtpExpiresAt:        now.Add(5 * time.Minute),
		SessionExpiresAt:    now.Add(15 * time.Minute),
		CreatedAt:           now,
	}

	t.Run("正常系: 通常セッションレコードを正しくINSERTできること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectExec("INSERT INTO OTP_SESSION").
			WithArgs(
				session.OtpSessionID,
				session.Purpose,
				session.UserID,
				session.PendingUsername,
				session.PendingEmail,
				session.MaskedEmail,
				session.PendingPasswordHash,
				session.OtpHash,
				session.Status,
				session.IsDummy,
				session.AttemptCount,
				session.SendCount,
				session.SendFailedCount,
				session.DeliveryStatus,
				session.LastSentAt,
				session.OtpExpiresAt,
				session.SessionExpiresAt,
				session.CreatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateOtpSession(context.Background(), session)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー時にエラーを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)
		dbErr := errors.New("insert failed")

		mock.ExpectExec("INSERT INTO OTP_SESSION").
			WillReturnError(dbErr)

		err = repo.CreateOtpSession(context.Background(), session)
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegisterRequestOtpRepository_UpdateOtpSessionDeliveryStatus(t *testing.T) {
	t.Run("正常系: 配信ステータスと送信失敗回数を正しくUPDATEできること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectExec("UPDATE OTP_SESSION SET DELIVERY_STATUS = .+, SEND_FAILED_COUNT = .+ WHERE OTP_SESSION_ID = .+").
			WithArgs("sent", 0, "otp_sess_12345").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = repo.UpdateOtpSessionDeliveryStatus(context.Background(), "otp_sess_12345", "sent", 0)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー時にエラーを返すこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)
		dbErr := errors.New("update failed")

		mock.ExpectExec("UPDATE OTP_SESSION SET DELIVERY_STATUS = .+, SEND_FAILED_COUNT = .+ WHERE OTP_SESSION_ID = .+").
			WithArgs("sendable", 1, "otp_sess_12345").
			WillReturnError(dbErr)

		err = repo.UpdateOtpSessionDeliveryStatus(context.Background(), "otp_sess_12345", "sendable", 1)
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegisterRequestOtpRepository_RecordMailAuthLog(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	log := &model.MailAuthLogRecord{
		LogID:     "log-123",
		Email:     "user@example.com",
		AuthType:  "SIGNUP",
		IPAddress: "192.0.2.1",
		EventType: "ISSUED",
		IsSuccess: true,
		IsDummy:   false,
		AccessAt:  now,
	}

	t.Run("正常系: メール認証ログを正しくINSERTできること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectExec("INSERT INTO MAIL_AUTH_LOG").
			WithArgs(log.LogID, log.UserID, log.Email, log.AuthType, log.IPAddress, log.EventType, log.IsSuccess, log.IsDummy, log.AccessAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.RecordMailAuthLog(context.Background(), log)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegisterRequestOtpRepository_RecordAccessLog(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	log := &model.AccessLogRecord{
		LogID:     "log-456",
		IPAddress: "192.0.2.1",
		Endpoint:  "POST auth/register/request-otp",
		AccessAt:  now,
	}

	t.Run("正常系: アクセスログを正しくINSERTできること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectExec("INSERT INTO ACCESS_LOG").
			WithArgs(log.LogID, log.UserID, log.IPAddress, log.Endpoint, log.ResourceID, log.AccessAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.RecordAccessLog(context.Background(), log)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegisterRequestOtpRepository_SaveSessionWithLogs(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	session := &model.OtpSessionRecord{
		OtpSessionID:     "otp_sess_12345",
		Purpose:          "SIGNUP",
		MaskedEmail:      "user**********@example.com",
		Status:           "active",
		LastSentAt:       now,
		OtpExpiresAt:     now.Add(5 * time.Minute),
		SessionExpiresAt: now.Add(15 * time.Minute),
		CreatedAt:        now,
	}
	mailLog := &model.MailAuthLogRecord{
		LogID:     "mail-log-1",
		Email:     "user@example.com",
		AuthType:  "SIGNUP",
		IPAddress: "192.0.2.1",
		EventType: "ISSUED",
		IsSuccess: true,
		IsDummy:   false,
		AccessAt:  now,
	}
	accessLog := &model.AccessLogRecord{
		LogID:     "access-log-1",
		IPAddress: "192.0.2.1",
		Endpoint:  "POST auth/register/request-otp",
		AccessAt:  now,
	}

	t.Run("正常系: 単一トランザクションでセッション・メール認証ログ・アクセスログをコミットすること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO OTP_SESSION").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO MAIL_AUTH_LOG").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO ACCESS_LOG").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.SaveSessionWithLogs(context.Background(), session, mailLog, accessLog)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: トランザクション途中のエラーでロールバックすること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewRegisterRequestOtpRepository(db)
		dbErr := errors.New("insert session failed")

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO OTP_SESSION").WillReturnError(dbErr)
		mock.ExpectRollback()

		err = repo.SaveSessionWithLogs(context.Background(), session, mailLog, accessLog)
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
