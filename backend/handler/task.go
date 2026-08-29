package handler

import (
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// CreateTaskHandler は POST /api/tasks のハンドラーを返します（Step 1 スタブ）。
func CreateTaskHandler(service service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
