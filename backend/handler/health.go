package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse はヘルスチェックエンドポイントのレスポンス構造体です。
//
// APIの稼働状態とメッセージを返却します。
type HealthResponse struct {
	// Status はサーバーの現在のステータスを示します（例: "ok"）
	Status string `json:"status" example:"ok"`
	// Message はヘルスチェックに関する補足メッセージです（例: "healthy"）
	Message string `json:"message" example:"healthy"`
}

// HealthCheckHandler は開発環境用のヘルスチェックエンドポイントのハンドラーです。
//
// サーバーが正常にリクエストを受信・処理できる状態にあるかを検証するために使用されます。
// 本番環境（GIN_MODE=release）ではルーティング自体が無効化されます。
//
// @Summary ヘルスチェック
// @Description サーバーの稼働状態を確認します（開発モード時のみ有効）。
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "サーバー稼働中"
// @Router /health-check [get]
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Message: "healthy",
	})
}
