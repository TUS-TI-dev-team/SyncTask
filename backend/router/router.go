package router

import (
	"database/sql"
	"net/http"

	_ "synctask/backend/docs"
	"synctask/backend/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter は Gin ルーターを初期化し、各ルートを登録します。
//
// 開発モード（gin.Mode() != gin.ReleaseMode）の場合にのみ、
// /health-check エンドポイントおよび /swagger/*any (Swagger UI) が有効化されます。
func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	// ルートエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// 開発時のみ有効なルート設定
	if gin.Mode() != gin.ReleaseMode {
		// ヘルスチェックエンドポイント
		r.GET("/health-check", handler.HealthCheckHandler(db))

		// Swagger UI エンドポイント
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
