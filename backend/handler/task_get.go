package handler

import (
	"net/http"

	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// GetTaskHandler は GET /api/tasks/:task_id のハンドラーを返します。
//
// 指定されたタスクIDの詳細情報を取得します。
//
// @Summary タスク詳細取得
// @Description 指定されたタスクIDの詳細情報を取得します。
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "タスクID（UUID）"
// @Success 200 {object} model.GetTaskResponse "タスク詳細取得成功"
// @Failure 401 {object} model.ErrorResponse "認証エラー（未ログイン）"
// @Failure 404 {object} model.ErrorResponse "タスクが存在しない、または他ユーザー所有"
// @Failure 500 {object} model.ErrorResponse "サーバー内部エラー"
// @Router /api/tasks/{task_id} [get]
func GetTaskHandler(service service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です。", nil))
			return
		}

		taskID := c.Param("task_id")
		res, err := service.GetTask(c.Request.Context(), userID, taskID)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
				return
			}
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_SERVER_ERROR", "サーバー内部でエラーが発生しました。", nil))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
