package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	// テストモードに設定
	gin.SetMode(gin.TestMode)

	// テスト用レスポンスレコーダーとコンテキスト作成
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// ハンドラー実行
	HealthCheckHandler(c)

	// ステータスコード検証 (200 OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディ検証
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "healthy", response.Message)
}
