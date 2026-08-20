package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"synctask/backend/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter_Root(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", body["message"])
}

func TestSetupRouter_HealthCheck_DevMode(t *testing.T) {
	// TestMode は ReleaseMode ではないため、/health-check が登録される
	gin.SetMode(gin.TestMode)
	r := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/health-check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res handler.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Status)
	assert.Equal(t, "healthy", res.Message)
}

func TestSetupRouter_HealthCheck_ReleaseMode(t *testing.T) {
	// ReleaseMode では /health-check が登録されないため 404 Not Found になる
	gin.SetMode(gin.ReleaseMode)
	r := SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/health-check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// テスト後にモードを元に戻す
	gin.SetMode(gin.TestMode)
}
