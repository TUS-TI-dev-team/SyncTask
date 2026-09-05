package handler

import (
	"errors"
	"net/http"

	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// RegisterRequestOtpHandler は POST /api/auth/register/request-otp のHandlerを返します。
//
// @Summary 新規登録OTP発行
// @Description 新規登録情報のバリデーション・OTP発行・メール送信を行います。
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequestOtpRequest true "新規登録リクエスト"
// @Success 200 {object} model.RegisterRequestOtpResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 503 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/auth/register/request-otp [post]
func RegisterRequestOtpHandler(svc service.RegisterRequestOtpService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")

		var req model.RegisterRequestOtpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse(
				"BAD_REQUEST",
				"不正なリクエスト形式です。",
				nil,
			))
			return
		}

		res, err := svc.RequestOtp(c.Request.Context(), &req, c.ClientIP())
		if err != nil {
			var appErr *model.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
				return
			}
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
				"INTERNAL_SERVER_ERROR",
				"サーバー内部でエラーが発生しました。",
				nil,
			))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
