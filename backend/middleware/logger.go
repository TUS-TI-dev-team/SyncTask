package middleware

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Printf("request: %v", c.Request)
		// c.Abort()
		c.Next()
		status := c.Writer.Status()
		log.Printf("Status: %v", status)
		elapsed := time.Since(start)
		log.Printf("Took: %v", elapsed)
	}
}