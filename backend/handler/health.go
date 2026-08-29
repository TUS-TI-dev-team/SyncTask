package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse はヘルスチェックエンドポイントのレスポンス構造体です。
//
// APIの稼働状態とDB接続状態を返却します。
type HealthResponse struct {
	// Status はサーバーの現在のステータスを示します（例: "ok"）
	Status string `json:"status" example:"ok"`
	// Message はヘルスチェックに関する補足メッセージです（例: "healthy"）
	Message string `json:"message" example:"healthy"`
	// Database はデータベースの接続状態を示します（例: "connected", "disconnected"）
	Database string `json:"database" example:"connected"`
}

// HealthCheckHandler は開発環境用のヘルスチェックエンドポイントのハンドラーを返します。
//
// サーバーが正常にリクエストを受信・処理できる状態にあるかを検証し、DBへのPing結果を返却します。
// 本番環境（GIN_MODE=release）ではルーティング自体が無効化されます。
//
// @Summary ヘルスチェック
// @Description サーバーおよびデータベースの稼働状態を確認します（開発モード時のみ有効）。
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "サーバー稼働中"
// @Router /health-check [get]
func HealthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbStatus := "connected"
		if db == nil || db.Ping() != nil {
			dbStatus = "disconnected"
		}
		c.JSON(http.StatusOK, HealthResponse{
			Status:   "ok",
			Message:  "healthy",
			Database: dbStatus,
		})
	}
}
