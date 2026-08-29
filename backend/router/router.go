package router

import (
	"database/sql"
	"net/http"

	_ "synctask/backend/docs"
	"synctask/backend/handler"
	"synctask/backend/repository"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter は Gin ルーターを初期化し、各ルートを登録します。
//
// 開発モード（gin.Mode() != gin.ReleaseMode）の場合にのみ、
// /health-check エンドポイントおよび /swagger/*any (Swagger UI) が有効化されます。
type Options struct {
	CookieSecure   bool
	TrustedProxies []string
}

func SetupRouter(db *sql.DB, configured ...Options) *gin.Engine {
	r := gin.Default()
	options := Options{}
	if len(configured) > 0 {
		options = configured[0]
	}
	if err := r.SetTrustedProxies(options.TrustedProxies); err != nil {
		panic("invalid trusted proxy configuration: " + err.Error())
	}

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

	// 依存性の注入 (DI)
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	loginRepo := repository.NewLoginRepository(db)
	loginService := service.NewLoginService(loginRepo, service.LoginDependencies{})
	// API ルーティング
	api := r.Group("/api")
	{
		api.POST("/tasks", handler.CreateTaskHandler(taskService))
		api.GET("/tasks/:task_id", handler.GetTaskHandler(taskService))
	}

	api.POST("/auth/login", handler.LoginHandler(loginService, handler.LoginHandlerOptions{CookieSecure: options.CookieSecure}))
	return r
}
