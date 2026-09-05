package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"synctask/backend/handler"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter_Root(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := SetupRouter(db)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", body["message"])
}

func TestSetupRouter_HealthCheck_DevMode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	r := SetupRouter(db)

	req, _ := http.NewRequest(http.MethodGet, "/health-check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res handler.HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Status)
	assert.Equal(t, "healthy", res.Message)
	assert.Equal(t, "connected", res.Database)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetupRouter_HealthCheck_ReleaseMode(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := SetupRouter(db)

	req, _ := http.NewRequest(http.MethodGet, "/health-check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// テスト後にモードを元に戻す
	gin.SetMode(gin.TestMode)
}

func TestSetupRouter_LoginRoute(t *testing.T) {
	previousMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(previousMode) })

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO ACCESS_LOG").
		WithArgs(sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "POST auth/login", nil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := SetupRouter(db)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetupRouter_Tasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := SetupRouter(db)

	// 未認証（userIDなし）で POST /api/tasks にアクセスすると 401 が返る（ルートが正しくハンドラーに到達している確認）
	req, _ := http.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(`{"title":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSetupRouter_GetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := SetupRouter(db)

	// 未認証（userIDなし）で GET /api/tasks/:task_id にアクセスすると 401 が返る（ルートが正しくハンドラーに到達している確認）
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks/7c9e6679-7425-40de-944b-e07fc1f90ae7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSetupRouter_TrustedProxies_DirectAccess(t *testing.T) {
	previousMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(previousMode) })

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// TrustedProxies が空（直接アクセス時・デフォルト設定）
	r := SetupRouter(db)
	r.GET("/test-client-ip", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})

	req := httptest.NewRequest(http.MethodGet, "/test-client-ip", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.195")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// プロキシが信頼されていないため、X-Forwarded-For は無視され RemoteAddr の IP が採用される
	assert.Equal(t, "192.0.2.1", w.Body.String())
}

func TestSetupRouter_TrustedProxies_WithTrustedProxy(t *testing.T) {
	previousMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(previousMode) })

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 192.0.2.1 を信頼済みプロキシとして設定
	r := SetupRouter(db, Options{
		TrustedProxies: []string{"192.0.2.1"},
	})
	r.GET("/test-client-ip", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})

	// 信頼済みプロキシからのリクエスト
	req := httptest.NewRequest(http.MethodGet, "/test-client-ip", nil)
	req.RemoteAddr = "192.0.2.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.195")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 信頼済みプロキシ経由のため、X-Forwarded-For の IP が採用される
	assert.Equal(t, "203.0.113.195", w.Body.String())

	// 未信頼のプロキシ（198.51.100.1）からのリクエスト
	reqUntrusted := httptest.NewRequest(http.MethodGet, "/test-client-ip", nil)
	reqUntrusted.RemoteAddr = "198.51.100.1:12345"
	reqUntrusted.Header.Set("X-Forwarded-For", "203.0.113.195")
	wUntrusted := httptest.NewRecorder()
	r.ServeHTTP(wUntrusted, reqUntrusted)

	assert.Equal(t, http.StatusOK, wUntrusted.Code)
	// 未信頼プロキシからのリクエストなので X-Forwarded-For は無視され RemoteAddr が採用される
	assert.Equal(t, "198.51.100.1", wUntrusted.Body.String())
}

