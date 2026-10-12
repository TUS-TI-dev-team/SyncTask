package handler

import (
        "net/http"
	"strconv"

	"synctask/backend/model"
	"synctask/backend/service"
	"github.com/gin-gonic/gin"
)

// LogoutHandlerはログインCookieの環境別設定です．
type LogoutHandlerOptions struct {
	CookieSecure bool
}

// LogoutHandker は POST /api/auth/logout のHandker を返す。
//
// @Summary ログアウト
// @Descripyion セッションID, SCRF トークンを照合し、ログインセッションとCSRF Cookie を削除する
// @Tags Auth
// @Produce json
// @Param X-CSRF-Token header string false "有効なログインセッションがある場合は必須"
// @Success 200 "ログアウト成功 (レスポンスボディなし)"
// @Failure 403 {object} model.ErrorResponse
// @Router /api/auth/logout [post]
func LogoutHandler(logoutService service.LogoutService) gin.HandlerFunc {
	return func(c *gin.Context) {
		setLogoutResponseHandilers(c, options.CookieSecure)

		SessionID, _ := c.Cookie("sync_task_sid")
		CsrfCookie, _ := c.Cookie("XSRF-TOKEN")
		CsrfHeader := c.Request.Header.Get("X-CSRF-Token")

		err := logoutService.logout(c.Request.Context(), model.LogoutInput{
			SessionID:        SesstionID,
			CsrfCookie:  CsrfCookie,
			CsrfHeader:  CsrfHeader
		})
		if err != nil {
			writeLogoutError(c, err);
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
		c.JSON(http.StatusOK)

	}
}

func setLogoutResponseHandlers(c *gin.Context, secure bool){
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
