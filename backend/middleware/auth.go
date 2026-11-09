package middleware

import (
	"errors"
	"log"
	"net/http"
	"synctask/backend/model"
	"synctask/backend/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Auth(s service.LoginSessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("sync_task_sid")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "not logged in", nil))
			log.Print("sync_task_sid not specified!")
			return
		}

		if _, err := uuid.Parse(sid); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("UNAUTHORIZED", "not logged in", nil))
			log.Print("sync_task_sid is not uuid format!")
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

		c.Set("userID", res.UserID)
		c.Set("loginSessionExpiresAt", res.ExpiresAt)

		c.Next()
	}
}
