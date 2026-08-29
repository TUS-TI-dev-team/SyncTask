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
	"golang.org/x/crypto/bcrypt"
)

type mockLoginRepository struct {
	attemptFunc       func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error)
	recordInvalidFunc func(context.Context, string, time.Time) error
}

func (m *mockLoginRepository) AttemptLogin(ctx context.Context, attempt *model.LoginAttempt, check repository.PasswordCheck) (*model.LoginAttemptResult, error) {
	return m.attemptFunc(ctx, attempt, check)
}

func (m *mockLoginRepository) RecordInvalidRequest(ctx context.Context, ip string, now time.Time) error {
	if m.recordInvalidFunc == nil {
		return nil
	}
	return m.recordInvalidFunc(ctx, ip, now)
}

func TestLoginService_Login(t *testing.T) {
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	meta := model.LoginMetadata{IP: "203.0.113.10", UserAgent: "test-agent", OldSessionID: "old-session"}

	t.Run("正常系: 認証成功時に新規トークンとユーザーを返し人工遅延しないこと", func(t *testing.T) {
		tokenCalls := 0
		sleepCalls := 0
		repo := &mockLoginRepository{attemptFunc: func(_ context.Context, attempt *model.LoginAttempt, check repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			assert.Equal(t, "user@example.com", attempt.Email)
			assert.Equal(t, meta.IP, attempt.IP)
			assert.Equal(t, meta.UserAgent, attempt.UserAgent)
			assert.Equal(t, meta.OldSessionID, attempt.OldSessionID)
			assert.Equal(t, "session-token", attempt.SessionID)
			assert.Equal(t, "csrf-token", attempt.CSRFToken)
			assert.True(t, check("$hash"))
			return &model.LoginAttemptResult{
				Status: model.LoginStatusSuccess,
				User:   &model.LoginUser{ID: "user-id", Username: "example", Email: "user@example.com", CreatedAt: now, UpdatedAt: now},
			}, nil
		}}

		svc := NewLoginService(repo, LoginDependencies{
			Now:   func() time.Time { return now },
			Sleep: func(time.Duration) { sleepCalls++ },
			GenerateToken: func() (string, error) {
				tokenCalls++
				if tokenCalls == 1 {
					return "session-token", nil
				}
				return "csrf-token", nil
			},
			ComparePassword: func(hash, password string) bool {
				assert.Equal(t, "$hash", hash)
				assert.Equal(t, "Password123!", password)
				return true
			},
			FailureDelay: func() time.Duration { return time.Second },
		})

		got, err := svc.Login(context.Background(), &model.LoginRequest{Email: " User@Example.COM ", Password: "Password123!"}, meta)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "session-token", got.SessionID)
		assert.Equal(t, "csrf-token", got.CSRFToken)
		assert.Equal(t, 30*24*time.Hour, got.MaxAge)
		assert.Equal(t, "user-id", got.Response.User.ID)
		assert.Equal(t, 0, sleepCalls)
	})

	t.Run("異常系: 認証失敗時に401へ統一し処理時間が目標未満なら残時間を待機すること", func(t *testing.T) {
		clockValues := []time.Time{now, now.Add(200 * time.Millisecond)}
		clockIndex := 0
		var slept time.Duration
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			return &model.LoginAttemptResult{Status: model.LoginStatusUnauthorized}, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now: func() time.Time {
				value := clockValues[clockIndex]
				if clockIndex < len(clockValues)-1 {
					clockIndex++
				}
				return value
			},
			Sleep:           func(d time.Duration) { slept = d },
			GenerateToken:   sequentialTokens("session-token", "csrf-token"),
			ComparePassword: func(string, string) bool { return false },
			FailureDelay:    func() time.Duration { return time.Second },
		})

		got, err := svc.Login(context.Background(), &model.LoginRequest{Email: "missing@example.com", Password: "Password123!"}, meta)
		assert.Nil(t, got)
		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 401, appErr.StatusCode)
		assert.Equal(t, "UNAUTHORIZED", appErr.Code)
		assert.Empty(t, appErr.Details)
		assert.Equal(t, 800*time.Millisecond, slept)
	})

	t.Run("異常系: IP遮断中は429とRetry-Afterを返し同じ遅延を適用すること", func(t *testing.T) {
		var slept time.Duration
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			return &model.LoginAttemptResult{Status: model.LoginStatusRateLimited, RetryAfter: 899}, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now:             fixedAdvancingClock(now, 100*time.Millisecond),
			Sleep:           func(d time.Duration) { slept = d },
			GenerateToken:   sequentialTokens("session-token", "csrf-token"),
			ComparePassword: func(string, string) bool { return false },
			FailureDelay:    func() time.Duration { return time.Second },
		})

		got, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)
		assert.Nil(t, got)
		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 429, appErr.StatusCode)
		assert.Equal(t, "RATE_LIMIT_EXCEEDED", appErr.Code)
		assert.Equal(t, 899, appErr.RetryAfter)
		assert.Equal(t, 900*time.Millisecond, slept)
	})

	t.Run("異常系: 入力不正はACCESS_LOGだけを記録して待機も認証試行もしないこと", func(t *testing.T) {
		attemptCalls := 0
		invalidCalls := 0
		sleepCalls := 0
		repo := &mockLoginRepository{
			attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
				attemptCalls++
				return nil, nil
			},
			recordInvalidFunc: func(_ context.Context, ip string, at time.Time) error {
				invalidCalls++
				assert.Equal(t, meta.IP, ip)
				assert.Equal(t, now, at)
				return nil
			},
		}
		svc := NewLoginService(repo, LoginDependencies{
			Now:             func() time.Time { return now },
			Sleep:           func(time.Duration) { sleepCalls++ },
			GenerateToken:   sequentialTokens("unused", "unused"),
			ComparePassword: func(string, string) bool { return false },
			FailureDelay:    func() time.Duration { return time.Second },
		})

		got, err := svc.Login(context.Background(), &model.LoginRequest{Email: "", Password: "short"}, meta)
		assert.Nil(t, got)
		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, 0, attemptCalls)
		assert.Equal(t, 1, invalidCalls)
		assert.Equal(t, 0, sleepCalls)
	})

	t.Run("異常系: Repositoryエラー時はエラーを伝播して認証結果を返さないこと", func(t *testing.T) {
		dbErr := errors.New("database failed")
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			return nil, dbErr
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now:             func() time.Time { return now },
			Sleep:           func(time.Duration) {},
			GenerateToken:   sequentialTokens("session-token", "csrf-token"),
			ComparePassword: func(string, string) bool { return false },
			FailureDelay:    func() time.Duration { return time.Second },
		})

		got, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestLoginService_AdditionalBranches(t *testing.T) {
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)
	meta := model.LoginMetadata{IP: "203.0.113.10"}

	t.Run("正常系: デフォルト依存関係でbcrypt照合とトークン生成が動作すること", func(t *testing.T) {
		hash, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.MinCost)
		require.NoError(t, err)
		repo := &mockLoginRepository{attemptFunc: func(_ context.Context, attempt *model.LoginAttempt, check repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			assert.NotEmpty(t, attempt.SessionID)
			assert.NotEmpty(t, attempt.CSRFToken)
			assert.NotEqual(t, attempt.SessionID, attempt.CSRFToken)
			assert.True(t, check(string(hash)))
			return &model.LoginAttemptResult{
				Status: model.LoginStatusSuccess,
				User:   &model.LoginUser{ID: "user-id", Email: "user@example.com"},
			}, nil
		}}

		result, err := NewLoginService(repo, LoginDependencies{}).Login(
			context.Background(),
			&model.LoginRequest{Email: "user@example.com", Password: "Password123!"},
			meta,
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "user-id", result.Response.User.ID)
	})

	t.Run("異常系: セッショントークン生成失敗時にRepositoryを呼ばないこと", func(t *testing.T) {
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			t.Fatal("repository must not be called")
			return nil, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now:           func() time.Time { return now },
			GenerateToken: func() (string, error) { return "", errors.New("token failed") },
		})

		result, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "session token")
	})

	t.Run("異常系: CSRFトークン生成失敗時にRepositoryを呼ばないこと", func(t *testing.T) {
		calls := 0
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			t.Fatal("repository must not be called")
			return nil, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now: func() time.Time { return now },
			GenerateToken: func() (string, error) {
				calls++
				if calls == 1 {
					return "session-token", nil
				}
				return "", errors.New("token failed")
			},
		})

		result, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "csrf token")
	})

	t.Run("正常系: RecordInvalidRequestを現在時刻付きでRepositoryへ委譲すること", func(t *testing.T) {
		called := false
		repo := &mockLoginRepository{
			attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
				return nil, nil
			},
			recordInvalidFunc: func(_ context.Context, ip string, at time.Time) error {
				called = true
				assert.Equal(t, meta.IP, ip)
				assert.Equal(t, now, at)
				return nil
			},
		}
		svc := NewLoginService(repo, LoginDependencies{Now: func() time.Time { return now }})

		err := svc.RecordInvalidRequest(context.Background(), meta.IP)

		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("準正常系: 処理時間が目標遅延を超えた場合に追加待機しないこと", func(t *testing.T) {
		sleepCalls := 0
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			return &model.LoginAttemptResult{Status: model.LoginStatusUnauthorized}, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now:             fixedAdvancingClock(now, 2*time.Second),
			Sleep:           func(time.Duration) { sleepCalls++ },
			GenerateToken:   sequentialTokens("session", "csrf"),
			ComparePassword: func(string, string) bool { return false },
			FailureDelay:    func() time.Duration { return time.Second },
		})

		_, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)

		require.Error(t, err)
		assert.Equal(t, 0, sleepCalls)
	})

	t.Run("異常系: 成功結果にユーザーがない場合に内部エラーを返すこと", func(t *testing.T) {
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			return &model.LoginAttemptResult{Status: model.LoginStatusSuccess}, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now:             func() time.Time { return now },
			GenerateToken:   sequentialTokens("session", "csrf"),
			ComparePassword: func(string, string) bool { return true },
		})

		result, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "without user")
	})

	t.Run("異常系: 未知のRepositoryステータスに内部エラーを返すこと", func(t *testing.T) {
		repo := &mockLoginRepository{attemptFunc: func(context.Context, *model.LoginAttempt, repository.PasswordCheck) (*model.LoginAttemptResult, error) {
			return &model.LoginAttemptResult{Status: model.LoginStatus("unknown")}, nil
		}}
		svc := NewLoginService(repo, LoginDependencies{
			Now:             func() time.Time { return now },
			GenerateToken:   sequentialTokens("session", "csrf"),
			ComparePassword: func(string, string) bool { return false },
		})

		result, err := svc.Login(context.Background(), &model.LoginRequest{Email: "user@example.com", Password: "Password123!"}, meta)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "unknown login status")
	})
}

func sequentialTokens(tokens ...string) func() (string, error) {
	index := 0
	return func() (string, error) {
		token := tokens[index]
		if index < len(tokens)-1 {
			index++
		}
		return token, nil
	}
}

func fixedAdvancingClock(start time.Time, elapsed time.Duration) func() time.Time {
	calls := 0
	return func() time.Time {
		calls++
		if calls == 1 {
			return start
		}
		return start.Add(elapsed)
	}
}
