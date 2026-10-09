package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("sync_task_sid")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Print("sync_task_sid not specified!")
			return 
		}

		log.Printf("sync_task_sid: %v", sid)
	}
}