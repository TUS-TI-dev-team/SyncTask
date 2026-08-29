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

func TestLoginRepository_AttemptLogin(t *testing.T) {
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)
	attempt := &model.LoginAttempt{
		Email: "user@example.com", Password: "Password123!", IP: "203.0.113.10",
		UserAgent: "test-agent", OldSessionID: "old-session", SessionID: "new-session",
		CSRFToken: "csrf-token", Now: now, ExpiresAt: now.Add(30 * 24 * time.Hour),
	}

	t.Run("正常系: 認証成功時に失敗状態リセットと旧セッション差替えとログを同一トランザクションで行うこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewLoginRepository(db)

		mock.ExpectBegin()
		expectNoIPRateRow(mock, attempt.IP)
		mock.ExpectQuery("SELECT.+FROM LOGIN_ACCOUNT.+FOR UPDATE").
			WithArgs(attempt.Email).
			WillReturnRows(accountRows(now).AddRow("user-id", "example", attempt.Email, "$valid-hash", 2, now.Add(-time.Minute), nil, now, now))
		mock.ExpectExec("UPDATE LOGIN_ACCOUNT.+LOGIN_FAILED_COUNT = 0").
			WithArgs("user-id").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("DELETE FROM LOGIN_SESSION").
			WithArgs(attempt.OldSessionID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("INSERT INTO LOGIN_SESSION").
			WithArgs(attempt.SessionID, "user-id", attempt.ExpiresAt, attempt.UserAgent, attempt.Now).
			WillReturnResult(sqlmock.NewResult(0, 1))
		expectLoginLog(mock, "user-id", attempt.Email, attempt.IP, true, attempt.Now)
		expectAccessLog(mock, "user-id", attempt.IP, attempt.Now)
		mock.ExpectCommit()

		result, err := repo.AttemptLogin(context.Background(), attempt, func(hash string) bool {
			assert.Equal(t, "$valid-hash", hash)
			return true
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.LoginStatusSuccess, result.Status)
		require.NotNil(t, result.User)
		assert.Equal(t, "user-id", result.User.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: パスワード不一致時にアカウントとIP失敗回数およびログを更新すること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewLoginRepository(db)

		mock.ExpectBegin()
		expectNoIPRateRow(mock, attempt.IP)
		mock.ExpectQuery("SELECT.+FROM LOGIN_ACCOUNT.+FOR UPDATE").
			WithArgs(attempt.Email).
			WillReturnRows(accountRows(now).AddRow("user-id", "example", attempt.Email, "$valid-hash", 4, now.Add(-time.Minute), nil, now, now))
		mock.ExpectExec("UPDATE LOGIN_ACCOUNT.+LOGIN_FAILED_COUNT").
			WithArgs(5, attempt.Now, attempt.Now.Add(30*time.Minute), "user-id").
			WillReturnResult(sqlmock.NewResult(0, 1))
		expectIPFailureUpsert(mock, attempt.IP, attempt.Now)
		expectLoginLog(mock, "user-id", attempt.Email, attempt.IP, false, attempt.Now)
		expectAccessLog(mock, "user-id", attempt.IP, attempt.Now)
		mock.ExpectCommit()

		result, err := repo.AttemptLogin(context.Background(), attempt, func(string) bool { return false })

		require.NoError(t, err)
		assert.Equal(t, model.LoginStatusUnauthorized, result.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: 未登録メールでもダミーハッシュを照合してIP失敗回数とログを更新すること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewLoginRepository(db)

		mock.ExpectBegin()
		expectNoIPRateRow(mock, attempt.IP)
		mock.ExpectQuery("SELECT.+FROM LOGIN_ACCOUNT.+FOR UPDATE").
			WithArgs(attempt.Email).
			WillReturnError(sql.ErrNoRows)
		expectIPFailureUpsert(mock, attempt.IP, attempt.Now)
		expectLoginLog(mock, nil, attempt.Email, attempt.IP, false, attempt.Now)
		expectAccessLog(mock, nil, attempt.IP, attempt.Now)
		mock.ExpectCommit()
		checked := false

		result, err := repo.AttemptLogin(context.Background(), attempt, func(hash string) bool {
			checked = true
			assert.Equal(t, DummyPasswordHash, hash)
			return false
		})

		require.NoError(t, err)
		assert.True(t, checked)
		assert.Equal(t, model.LoginStatusUnauthorized, result.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: IP遮断中は認証照合せず429結果とログを返し遮断期限を延長しないこと", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewLoginRepository(db)
		blockedUntil := now.Add(899 * time.Second)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT FAILED_COUNT, LAST_FAILED_AT, BLOCKED_UNTIL.+FOR UPDATE").
			WithArgs(attempt.IP).
			WillReturnRows(sqlmock.NewRows([]string{"FAILED_COUNT", "LAST_FAILED_AT", "BLOCKED_UNTIL"}).
				AddRow(30, now.Add(-time.Minute), blockedUntil))
		expectLoginLog(mock, nil, attempt.Email, attempt.IP, false, attempt.Now)
		expectAccessLog(mock, nil, attempt.IP, attempt.Now)
		mock.ExpectCommit()

		result, err := repo.AttemptLogin(context.Background(), attempt, func(string) bool {
			t.Fatal("password check must not run")
			return false
		})

		require.NoError(t, err)
		assert.Equal(t, model.LoginStatusRateLimited, result.Status)
		assert.Equal(t, 899, result.RetryAfter)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー時にロールバックすること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewLoginRepository(db)
		dbErr := errors.New("query failed")

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT FAILED_COUNT, LAST_FAILED_AT, BLOCKED_UNTIL.+FOR UPDATE").
			WithArgs(attempt.IP).
			WillReturnError(dbErr)
		mock.ExpectRollback()

		result, err := repo.AttemptLogin(context.Background(), attempt, func(string) bool { return false })
		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLoginRepository_RecordInvalidRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewLoginRepository(db)
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)

	mock.ExpectExec("INSERT INTO ACCESS_LOG").
		WithArgs(sqlmock.AnyArg(), nil, "203.0.113.10", "POST auth/login", nil, now).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.RecordInvalidRequest(context.Background(), "203.0.113.10", now)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func accountRows(now time.Time) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"USER_ID", "USER_NAME", "EMAIL", "PASSWORD_HASH", "LOGIN_FAILED_COUNT",
		"LOGIN_LAST_FAILED_AT", "LOGIN_LOCK_UNTIL", "CREATED_AT", "UPDATED_AT",
	})
}

func expectNoIPRateRow(mock sqlmock.Sqlmock, ip string) {
	mock.ExpectQuery("SELECT FAILED_COUNT, LAST_FAILED_AT, BLOCKED_UNTIL.+FOR UPDATE").
		WithArgs(ip).
		WillReturnError(sql.ErrNoRows)
}

func expectIPFailureUpsert(mock sqlmock.Sqlmock, ip string, now time.Time) {
	mock.ExpectExec("INSERT INTO LOGIN_IP_RATE_LIMIT").
		WithArgs(ip, now, now).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func expectLoginLog(mock sqlmock.Sqlmock, userID any, email, ip string, success bool, now time.Time) {
	mock.ExpectExec("INSERT INTO LOGIN_LOG").
		WithArgs(sqlmock.AnyArg(), userID, email, ip, success, false, now).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func expectAccessLog(mock sqlmock.Sqlmock, userID any, ip string, now time.Time) {
	mock.ExpectExec("INSERT INTO ACCESS_LOG").
		WithArgs(sqlmock.AnyArg(), userID, ip, "POST auth/login", nil, now).
		WillReturnResult(sqlmock.NewResult(0, 1))
}
