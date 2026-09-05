package handler

import (
	"errors"
	"net/http"

	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// PatchTaskHandler は PATCH /api/tasks/:task_id のハンドラーを返します。
//
// 指定されたタスクIDの情報を部分更新します。
//
// @Summary タスク部分更新
// @Description 指定されたタスクIDの情報を部分更新します。リクエストボディに含まれるフィールドのみが更新されます。
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "タスクID（UUID）"
// @Param task body model.PatchTaskRequest true "タスク部分更新リクエスト"
// @Success 200 {object} model.PatchTaskResponse "タスク部分更新成功"
// @Failure 400 {object} model.ErrorResponse "バリデーションエラー・JSON形式不正"
// @Failure 401 {object} model.ErrorResponse "認証エラー（未ログイン）"
// @Failure 404 {object} model.ErrorResponse "タスクが存在しない、または他ユーザー所有"
// @Failure 500 {object} model.ErrorResponse "サーバー内部エラー"
// @Router /api/tasks/{task_id} [patch]
func PatchTaskHandler(service service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です。", nil))
			return
		}

		taskID := c.Param("task_id")

		var req model.PatchTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			var appErr *model.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
				return
			}
			c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_REQUEST", "リクエストボディのJSON形式が不正です。", nil))
			return
		}

		res, err := service.PatchTask(c.Request.Context(), userID, taskID, &req)
		if err != nil {
			var appErr *model.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
				return
			}
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_SERVER_ERROR", "サーバー内部でエラーが発生しました。", nil))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
