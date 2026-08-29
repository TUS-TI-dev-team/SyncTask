package handler

import (
	"net/http"
	"strconv"

	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
)

// LoginHandlerOptions はログインCookieの環境別設定です。
type LoginHandlerOptions struct {
	CookieSecure bool
}

// LoginHandler は POST /api/auth/login のHandlerを返します。
//
// @Summary ログイン
// @Description メールアドレスとパスワードを照合し、ログインセッションとCSRF Cookieを発行します。
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "ログインリクエスト"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 429 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/auth/login [post]
func LoginHandler(loginService service.LoginService, options LoginHandlerOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		setLoginResponseHeaders(c, options.CookieSecure)

		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if logErr := loginService.RecordInvalidRequest(c.Request.Context(), c.ClientIP()); logErr != nil {
				writeLoginError(c, logErr)
				return
			}
			c.JSON(http.StatusBadRequest, model.NewErrorResponse(
				"BAD_REQUEST",
				"不正なリクエスト形式です。",
				nil,
			))
			return
		}

		oldSessionID, _ := c.Cookie("sync_task_sid")
		result, err := loginService.Login(c.Request.Context(), &req, model.LoginMetadata{
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			OldSessionID: oldSessionID,
		})
		if err != nil {
			writeLoginError(c, err)
			return
		}

		maxAge := int(result.MaxAge.Seconds())
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "sync_task_sid",
			Value:    result.SessionID,
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
			Secure:   options.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "XSRF-TOKEN",
			Value:    result.CSRFToken,
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: false,
			Secure:   options.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
		c.JSON(http.StatusOK, result.Response)
	}
}

func setLoginResponseHeaders(c *gin.Context, secure bool) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
	if secure {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}
}

func writeLoginError(c *gin.Context, err error) {
	if appErr, ok := err.(*model.AppError); ok {
		if appErr.RetryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(appErr.RetryAfter))
		}
		c.JSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
		return
	}
	c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
		"INTERNAL_SERVER_ERROR",
		"サーバー内部でエラーが発生しました。",
		nil,
	))
}
