package middleware

import (
	"errors"
	"log"
	"net/http"
	"synctask/backend/model"
	"synctask/backend/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const sessionCookieKey = "sync_task_sid"

const userIDKey = "user_id"
const expiresAtKey = "login_session_expires_at"

func Auth(s service.LoginSessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(sessionCookieKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewUnauthorizedError("not logged in"))
			return
		}

		if _, err := uuid.Parse(sid); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewUnauthorizedError("not logged in"))
			return
		}

		res, err := s.ValidateSession(c, sid)
		if err != nil {
			var appErr *model.AppError
			if errors.As(err, &appErr) {
				c.AbortWithStatusJSON(appErr.StatusCode, model.NewErrorResponse(appErr.Code, appErr.Message, appErr.Details))
			} else {
				log.Fatal(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_SERVER_ERROR", "サーバー内部でエラーが発生しました。", nil))
			}
		}

		SetUserID(c, res.UserID)
		SetExpiresAt(c, res.ExpiresAt)

		c.Next()
	}
}

func SetUserID(c *gin.Context, userID string) {
	c.Set(userIDKey, userID)
}

func GetUserID(c *gin.Context) (string, bool) {
	val, exists := c.Get(userIDKey)
	if !exists {
		return "", false
	}
	id, ok := val.(string)
	return id, ok
}

func SetExpiresAt(c *gin.Context, expiresAt time.Time) {
	c.Set(expiresAtKey, expiresAt)
}

func GetExpiresAt(c *gin.Context) (time.Time, bool) {
	val, exists := c.Get(expiresAtKey)
	if !exists {
		return time.Time{}, false
	}
	e, ok := val.(time.Time)
	return e, ok
}
