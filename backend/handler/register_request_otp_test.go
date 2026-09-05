package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"synctask/backend/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRegisterRequestOtpService struct {
	requestOtpFunc func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error)
}

func (m *mockRegisterRequestOtpService) RequestOtp(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
	if m.requestOtpFunc != nil {
		return m.requestOtpFunc(ctx, req, clientIP)
	}
	return nil, nil
}

func performRegisterRequestOtp(t *testing.T, svc *mockRegisterRequestOtpService, body string, customize func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register/request-otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if customize != nil {
		customize(req)
	}
	c.Request = req
	RegisterRequestOtpHandler(svc)(c)
	return w
}

func validRegisterJSON() string {
	return `{"username":"exampleUser","email":"user@example.com","password":"Password123!"}`
}

func TestRegisterRequestOtpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })

	t.Run("正常系: 200とレスポンスボディおよびキャッシュ抑止ヘッダーを返すこと", func(t *testing.T) {
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				assert.Equal(t, "exampleUser", req.Username)
				assert.Equal(t, "user@example.com", req.Email)
				assert.Equal(t, "Password123!", req.Password)
				assert.Equal(t, "192.0.2.1", clientIP)
				return &model.RegisterRequestOtpResponse{
					OtpSessionID:     "otp_sess_test_12345",
					MaskedEmail:      "user**********@example.com",
					ExpiresInSeconds: 300,
					CooldownSeconds:  60,
				}, nil
			},
		}

		w := performRegisterRequestOtp(t, svc, validRegisterJSON(), func(req *http.Request) {
			req.RemoteAddr = "192.0.2.1:54321"
		})

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.RegisterRequestOtpResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "otp_sess_test_12345", response.OtpSessionID)
		assert.Equal(t, "user**********@example.com", response.MaskedEmail)
		assert.Equal(t, 300, response.ExpiresInSeconds)
		assert.Equal(t, 60, response.CooldownSeconds)

		assert.Equal(t, "no-store, no-cache, must-revalidate", w.Header().Get("Cache-Control"))
		assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	})

	t.Run("正常系: 未知のJSONフィールドを無視すること", func(t *testing.T) {
		called := false
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				called = true
				return &model.RegisterRequestOtpResponse{
					OtpSessionID:     "otp_sess_123",
					MaskedEmail:      "user**********@example.com",
					ExpiresInSeconds: 300,
					CooldownSeconds:  60,
				}, nil
			},
		}

		w := performRegisterRequestOtp(t, svc, `{"username":"exampleUser","email":"user@example.com","password":"Password123!","extra":"field"}`, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, called)
	})

	t.Run("異常系: 不正なJSON形式の場合は400を返すこと", func(t *testing.T) {
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				t.Fatal("RequestOtp must not be called")
				return nil, nil
			},
		}

		w := performRegisterRequestOtp(t, svc, `{"username":`, nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response model.ErrorResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "BAD_REQUEST", response.Error.Code)
	})

	t.Run("異常系: 入力バリデーションエラー時は400とエラー詳細を返すこと", func(t *testing.T) {
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				return nil, model.NewBadRequestError("入力内容に不備があります。", []model.ErrorDetail{
					{Field: "username", Message: "ユーザー名は必須です。"},
				})
			},
		}

		w := performRegisterRequestOtp(t, svc, validRegisterJSON(), nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response model.ErrorResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "BAD_REQUEST", response.Error.Code)
		require.Len(t, response.Error.Details, 1)
		assert.Equal(t, "username", response.Error.Details[0].Field)
	})

	t.Run("異常系: メール送信失敗時は503を返すこと", func(t *testing.T) {
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				return nil, model.NewServiceUnavailableError("OTP_DELIVERY_FAILED", "メールの送信に失敗しました。")
			},
		}

		w := performRegisterRequestOtp(t, svc, validRegisterJSON(), nil)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		var response model.ErrorResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "OTP_DELIVERY_FAILED", response.Error.Code)
	})

	t.Run("異常系: 予期しないServiceエラーは500を返すこと", func(t *testing.T) {
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				return nil, errors.New("db crash")
			},
		}

		w := performRegisterRequestOtp(t, svc, validRegisterJSON(), nil)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response model.ErrorResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Error.Code)
	})

	t.Run("準正常系: 全レスポンスにキャッシュ抑止ヘッダーを付与すること", func(t *testing.T) {
		svc := &mockRegisterRequestOtpService{
			requestOtpFunc: func(ctx context.Context, req *model.RegisterRequestOtpRequest, clientIP string) (*model.RegisterRequestOtpResponse, error) {
				return nil, model.NewServiceUnavailableError("OTP_DELIVERY_FAILED", "メールの送信に失敗しました。")
			},
		}

		w := performRegisterRequestOtp(t, svc, validRegisterJSON(), nil)
		assert.Equal(t, "no-store, no-cache, must-revalidate", w.Header().Get("Cache-Control"))
		assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	})
}
