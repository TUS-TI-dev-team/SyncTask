package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"synctask/backend/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

type mockLogoutService struct {
	logoutFunc	func(context.Context, *model.LogoutInput) (error)
}

func (m *mockLogoutService) Logout(ctx context.Context, input *model.LogoutInput) (error) {
	return m.logoutFunc(ctx, input)
}

func TestLogoutHandker(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })

	t.Run("正常系: 200とsid, csrf を削除するCookie を返すこと", func(t *testing.T){
		svc := &mockLogoutService{logoutFunc: func(_ context.Context, input *model.LogoutInput) (error) {
			assert.Equal(t, "4c7bcdd6-898b-483f-bd0d-7a4191a0384b", input.SessionID)
			return nil
		}}

		w := performLogoutRequest(t, svc, LogoutHandlerOptions{CookieSecure: true}, func(req *http.Request) {
			req.AddCookie(&http.Cookie{Name: "sync_task_sid", Value: "4c7bcdd6-898b-483f-bd0d-7a4191a0384b"})
		})

		assert.Equal(t, http.StatusOK, w.Code)
		cookie := w.Result().Cookies()

		sessionRemove := cookieByName(t, cookie, "sync_task_sid")
		assert.Equal(t, "", sessionRemove.Value)
		assert.True(t, sessionRemove.HttpOnly)
		assert.True(t, sessionRemove.Secure)
		assert.Equal(t, http.SameSiteLaxMode, sessionRemove.SameSite)
		assert.Equal(t, "/", sessionRemove.Path)
		assert.Equal(t, -1, sessionRemove.MaxAge)

		csrfRemove := cookieByName(t, cookie, "XSRF-TOKEN")
		assert.Equal(t, "", csrfRemove.Value)
		assert.False(t, csrfRemove.HttpOnly)
		assert.True(t, csrfRemove.Secure)
		assert.Equal(t, http.SameSiteLaxMode, csrfRemove.SameSite)
		assert.Equal(t, "/", csrfRemove.Path)
		assert.Equal(t, -1, csrfRemove.MaxAge)
	})
}

/*
func cookieByName(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("cookie %s not found", name)
	return nil
}
*/

func performLogoutRequest(t *testing.T, svc *mockLogoutService, options LogoutHandlerOptions, customize func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	if customize != nil {
		customize(req)
	}
	c.Request = req
	LogoutHandler(svc, options)(c)
	return w
}
