package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"synctask/backend/model"
	"synctask/backend/service"
)

// LogoutHandlerOptionsはログアウトCookieの環境別設定です。
type LogoutHandlerOptions struct {
	CookieSecure bool
}

// LogoutHandler は POST /api/auth/logout のHandler を返す。
//
// @Summary ログアウト
// @Description セッションID, CSRF トークンを照合し、ログインセッションとCSRF Cookie を削除する
// @Tags Auth
// @Produce json
// @Param X-CSRF-Token header string false "有効なログインセッションがある場合は必須"
// @Success 200 "ログアウト成功 (レスポンスボディなし)"
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/auth/logout [post]
func LogoutHandler(logoutService service.LogoutService, options LogoutHandlerOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		setLogoutResponseHeaders(c, options.CookieSecure)

		sessionID, _ := c.Cookie("sync_task_sid")

		err := logoutService.Logout(c.Request.Context(), &model.LogoutInput{
			SessionID: sessionID,
		})
		if err != nil {
			writeLogoutError(c, err)
			return
		}

		maxAge := -1
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "sync_task_sid",
			Value:    "",
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
			Secure:   options.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "XSRF-TOKEN",
			Value:    "",
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: false,
			Secure:   options.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
		c.Status(http.StatusOK)

	}
}

func setLogoutResponseHeaders(c *gin.Context, secure bool) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
	if secure {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}
}

func writeLogoutError(c *gin.Context, err error) {
	if appErr, ok := err.(*model.AppError); ok {
		c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}
	c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
		"INTERNAL_SERVER_ERROR",
		"サーバー内部でエラーが発生しました。",
		nil,
	))
}
