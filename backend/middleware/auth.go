package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("sync_task_sid")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Print("sync_task_sid not specified!")
			return
		}

		if _, err := uuid.Parse(sid); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Print("sync_task_sid not specified!")
			return
		}

		// TODO: service/login_session.go を作り、それを呼んで認証を実現する。

		c.Next()
	}
}
