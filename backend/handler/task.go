package handler

import (
	"net/http"

	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// CreateTaskHandler は POST /api/tasks のハンドラーを返します。
func CreateTaskHandler(service service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "認証が必要です。",
				},
			})
			return
		}

		var req model.CreateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    "BAD_REQUEST",
					Message: "不正なリクエスト形式です。",
				},
			})
			return
		}

		res, err := service.CreateTask(c.Request.Context(), userID, &req)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				c.JSON(appErr.StatusCode, model.ErrorResponse{
					Error: model.ErrorBody{
						Code:    appErr.Code,
						Message: appErr.Message,
						Details: appErr.Details,
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "サーバー内部でエラーが発生しました。",
				},
			})
			return
		}

		c.JSON(http.StatusCreated, res)
	}
}
