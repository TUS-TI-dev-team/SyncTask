package handler

import (
	"net/http"

	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// CreateTaskHandler は POST /api/tasks のハンドラーを返します。
//
// 新規タスクを作成します。単一タスクの登録、および繰り返しルールに基づくタスクの一括即時生成（最大100件）に対応します。
//
// @Summary 新規タスク作成
// @Description 新規タスクを作成します。単一タスクの登録、および繰り返しルールに基づくタスクの一括即時生成（最大100件）に対応します。
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body model.CreateTaskRequest true "タスク作成リクエスト"
// @Success 201 {object} model.CreateTaskResponse "タスク作成成功"
// @Failure 400 {object} model.ErrorResponse "不正なリクエスト形式またはバリデーションエラー"
// @Failure 401 {object} model.ErrorResponse "認証エラー（未ログイン）"
// @Failure 500 {object} model.ErrorResponse "サーバー内部エラー"
// @Router /api/tasks [post]
func CreateTaskHandler(service service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "認証が必要です。", nil))
			return
		}

		var req model.CreateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse("BAD_REQUEST", "不正なリクエスト形式です。", nil))
			return
		}

		res, err := service.CreateTask(c.Request.Context(), userID, &req)
		if err != nil {
			if appErr, ok := err.(*model.AppError); ok {
				c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
				return
			}
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_SERVER_ERROR", "サーバー内部でエラーが発生しました。", nil))
			return
		}

		c.JSON(http.StatusCreated, res)
	}
}
