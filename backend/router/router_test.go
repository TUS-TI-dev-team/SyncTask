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
