package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLoginService struct {
	loginFunc         func(context.Context, *model.LoginRequest, model.LoginMetadata) (*model.LoginServiceResult, error)
	recordInvalidFunc func(context.Context, string) error
}

func (m *mockLoginService) Login(ctx context.Context, req *model.LoginRequest, meta model.LoginMetadata) (*model.LoginServiceResult, error) {
	return m.loginFunc(ctx, req, meta)
}

func (m *mockLoginService) RecordInvalidRequest(ctx context.Context, ip string) error {
	if m.recordInvalidFunc == nil {
		return nil
	}
	return m.recordInvalidFunc(ctx, ip)
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })

	t.Run("正常系: 200とユーザーおよび仕様どおりのCookieを返すこと", func(t *testing.T) {
		now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))
		svc := &mockLoginService{loginFunc: func(_ context.Context, req *model.LoginRequest, meta model.LoginMetadata) (*model.LoginServiceResult, error) {
			assert.Equal(t, "User@Example.COM", req.Email)
			assert.Equal(t, "Password123!", req.Password)
			assert.Equal(t, "192.0.2.1", meta.IP)
			assert.Equal(t, "test-agent", meta.UserAgent)
			assert.Equal(t, "old-session", meta.OldSessionID)
			return &model.LoginServiceResult{
				Response:  model.LoginResponse{User: model.LoginUser{ID: "user-id", Username: "example", Email: "user@example.com", CreatedAt: now, UpdatedAt: now}},
				SessionID: "new-session", CSRFToken: "csrf-token", MaxAge: 30 * 24 * time.Hour,
			}, nil
		}}

		w := performLoginRequest(t, svc, LoginHandlerOptions{CookieSecure: true}, "{\"email\":\"User@Example.COM\",\"password\":\"Password123!\"}", func(req *http.Request) {
			req.RemoteAddr = "192.0.2.1:54321"
			req.Header.Set("User-Agent", "test-agent")
			req.AddCookie(&http.Cookie{Name: "sync_task_sid", Value: "old-session"})
		})

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.LoginResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "user-id", response.User.ID)
		cookies := w.Result().Cookies()
		require.Len(t, cookies, 2)

		session := cookieByName(t, cookies, "sync_task_sid")
		assert.Equal(t, "new-session", session.Value)
		assert.True(t, session.HttpOnly)
		assert.True(t, session.Secure)
		assert.Equal(t, http.SameSiteLaxMode, session.SameSite)
		assert.Equal(t, "/", session.Path)
		assert.Equal(t, 2592000, session.MaxAge)

		csrf := cookieByName(t, cookies, "XSRF-TOKEN")
		assert.Equal(t, "csrf-token", csrf.Value)
		assert.False(t, csrf.HttpOnly)
		assert.True(t, csrf.Secure)
		assert.Equal(t, 2592000, csrf.MaxAge)
	})

	t.Run("正常系: 未知のJSONフィールドを無視すること", func(t *testing.T) {
		called := false
		svc := successLoginMock(func(req *model.LoginRequest) {
			called = true
			assert.Equal(t, "user@example.com", req.Email)
		})
		w := performLoginRequest(t, svc, LoginHandlerOptions{}, "{\"email\":\"user@example.com\",\"password\":\"Password123!\",\"unknown\":\"ignored\"}", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, called)
	})

	t.Run("異常系: 不正JSONは400と空detailsを返してACCESS_LOG記録を依頼すること", func(t *testing.T) {
		invalidCalls := 0
		svc := &mockLoginService{
			loginFunc: func(context.Context, *model.LoginRequest, model.LoginMetadata) (*model.LoginServiceResult, error) {
				t.Fatal("Login must not be called")
				return nil, nil
			},
			recordInvalidFunc: func(_ context.Context, ip string) error {
				invalidCalls++
				assert.Equal(t, "192.0.2.1", ip)
				return nil
			},
		}
		w := performLoginRequest(t, svc, LoginHandlerOptions{}, "{\"email\":", func(req *http.Request) {
			req.RemoteAddr = "192.0.2.1:1234"
		})
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, 1, invalidCalls)
		assertErrorResponse(t, w, "BAD_REQUEST")
	})

	t.Run("異常系: 401ではCookieを発行しないこと", func(t *testing.T) {
		svc := errorLoginMock(model.NewUnauthorizedError("メールアドレスまたはパスワードが正しくありません。"))
		w := performLoginRequest(t, svc, LoginHandlerOptions{}, validLoginJSON(), nil)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Empty(t, w.Result().Cookies())
		assertErrorResponse(t, w, "UNAUTHORIZED")
	})

	t.Run("異常系: 429でRetry-Afterを返すこと", func(t *testing.T) {
		svc := errorLoginMock(model.NewRateLimitError("ログイン試行回数が上限に達しました。", 899))
		w := performLoginRequest(t, svc, LoginHandlerOptions{}, validLoginJSON(), nil)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "899", w.Header().Get("Retry-After"))
		assertErrorResponse(t, w, "RATE_LIMIT_EXCEEDED")
	})

	t.Run("異常系: 予期しないServiceエラーは500でCookieを発行しないこと", func(t *testing.T) {
		svc := errorLoginMock(errors.New("unexpected"))
		w := performLoginRequest(t, svc, LoginHandlerOptions{}, validLoginJSON(), nil)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Empty(t, w.Result().Cookies())
		assertErrorResponse(t, w, "INTERNAL_SERVER_ERROR")
	})

	t.Run("準正常系: 全レスポンスにキャッシュ抑止ヘッダーを付与すること", func(t *testing.T) {
		svc := errorLoginMock(model.NewUnauthorizedError("メールアドレスまたはパスワードが正しくありません。"))
		w := performLoginRequest(t, svc, LoginHandlerOptions{}, validLoginJSON(), nil)
		assert.Equal(t, "no-store, no-cache, must-revalidate", w.Header().Get("Cache-Control"))
		assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	})
}

func successLoginMock(inspect func(*model.LoginRequest)) *mockLoginService {
	return &mockLoginService{loginFunc: func(_ context.Context, req *model.LoginRequest, _ model.LoginMetadata) (*model.LoginServiceResult, error) {
		inspect(req)
		return &model.LoginServiceResult{
			Response:  model.LoginResponse{User: model.LoginUser{ID: "user-id"}},
			SessionID: "session", CSRFToken: "csrf", MaxAge: 30 * 24 * time.Hour,
		}, nil
	}}
}

func errorLoginMock(err error) *mockLoginService {
	return &mockLoginService{loginFunc: func(context.Context, *model.LoginRequest, model.LoginMetadata) (*model.LoginServiceResult, error) {
		return nil, err
	}}
}

func validLoginJSON() string {
	return "{\"email\":\"user@example.com\",\"password\":\"Password123!\"}"
}

func performLoginRequest(t *testing.T, svc *mockLoginService, options LoginHandlerOptions, body string, customize func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if customize != nil {
		customize(req)
	}
	c.Request = req
	LoginHandler(svc, options)(c)
	return w
}

func cookieByName(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("cookie %s not found", name)
	return nil
}

func assertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, code string) {
	t.Helper()
	var response model.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, code, response.Error.Code)
	assert.NotNil(t, response.Error.Details)
}
