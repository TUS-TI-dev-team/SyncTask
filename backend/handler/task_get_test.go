package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系: 有効な task_id 指定時に 200 OK とタスク詳細JSONが返却されること", func(t *testing.T) {
		userID := "550e8400-e29b-41d4-a716-446655440000"
		taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		now := time.Now()
		dueTime := now.Add(24 * time.Hour)

		mockResp := &model.GetTaskResponse{
			Task: &model.Task{
				ID:          taskID,
				UserID:      userID,
				Title:       "課題レポート提出",
				Comment:     "第5章の要約を含むこと",
				Priority:    "high",
				Status:      "in_progress",
				DueDatetime: &dueTime,
				IsPinned:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockSvc := &mockTaskService{
			getTaskFunc: func(ctx context.Context, uID, tID string) (*model.GetTaskResponse, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, taskID, tID)
				return mockResp, nil
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := GetTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp model.GetTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.NotNil(t, resp.Task)
		assert.Equal(t, taskID, resp.Task.ID)
		assert.Equal(t, userID, resp.Task.UserID)
		assert.Equal(t, "課題レポート提出", resp.Task.Title)
		assert.Equal(t, "第5章の要約を含むこと", resp.Task.Comment)
		assert.Equal(t, "high", resp.Task.Priority)
		assert.Equal(t, "in_progress", resp.Task.Status)
		require.NotNil(t, resp.Task.DueDatetime)
		assert.True(t, dueTime.Equal(*resp.Task.DueDatetime))
		assert.True(t, resp.Task.IsPinned)
	})

	t.Run("異常系: 未ログイン（Context に userID なし）の場合に 401 UNAUTHORIZED を返すこと", func(t *testing.T) {
		taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		mockSvc := &mockTaskService{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		// userID をセットしない

		h := GetTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "UNAUTHORIZED", errResp.Error.Code)
		assert.NotEmpty(t, errResp.Error.Message)
		assert.NotNil(t, errResp.Error.Details)
		assert.Empty(t, errResp.Error.Details)
		assert.Contains(t, w.Body.String(), `"details":[]`)
	})

	t.Run("異常系: タスクが存在しない（または他者所有）場合に 404 NOT_FOUND を返すこと", func(t *testing.T) {
		userID := "550e8400-e29b-41d4-a716-446655440000"
		taskID := "non-existent-task-id"

		mockSvc := &mockTaskService{
			getTaskFunc: func(ctx context.Context, uID, tID string) (*model.GetTaskResponse, error) {
				return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := GetTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
		assert.Equal(t, "指定されたタスクが見つかりません。", errResp.Error.Message)
		assert.NotNil(t, errResp.Error.Details)
		assert.Empty(t, errResp.Error.Details)
		assert.Contains(t, w.Body.String(), `"details":[]`)
	})

	t.Run("異常系: サーバー内部エラー発生時に 500 INTERNAL_SERVER_ERROR を返すこと", func(t *testing.T) {
		userID := "550e8400-e29b-41d4-a716-446655440000"
		taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"

		mockSvc := &mockTaskService{
			getTaskFunc: func(ctx context.Context, uID, tID string) (*model.GetTaskResponse, error) {
				return nil, assert.AnError
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := GetTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", errResp.Error.Code)
		assert.NotEmpty(t, errResp.Error.Message)
		assert.NotNil(t, errResp.Error.Details)
		assert.Empty(t, errResp.Error.Details)
		assert.Contains(t, w.Body.String(), `"details":[]`)
	})
}
