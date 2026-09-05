package handler

import (
	"bytes"
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

func TestPatchTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
	now := time.Now()
	due := now.Add(24 * time.Hour)

	t.Run("正常系: 有効なリクエストで 200 OK と更新後のタスク詳細が返却されること", func(t *testing.T) {
		mockResp := &model.PatchTaskResponse{
			Task: &model.Task{
				ID:          taskID,
				UserID:      userID,
				Title:       "課題レポート提出（修正版）",
				Comment:     "参考文献の追記完了",
				Priority:    "high",
				Status:      "completed",
				DueDatetime: &due,
				IsPinned:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockSvc := &mockTaskService{
			patchTaskFunc: func(ctx context.Context, uID, tID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, taskID, tID)
				require.NotNil(t, req.Title)
				assert.Equal(t, "課題レポート提出（修正版）", *req.Title)
				return mockResp, nil
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"title": "課題レポート提出（修正版）", "status": "completed"}`
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := PatchTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp model.PatchTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.NotNil(t, resp.Task)
		assert.Equal(t, taskID, resp.Task.ID)
		assert.Equal(t, "課題レポート提出（修正版）", resp.Task.Title)
		assert.Equal(t, "completed", resp.Task.Status)
	})

	t.Run("正常系: 空ボディ {} でも 200 OK と現在のタスク詳細が返却されること", func(t *testing.T) {
		mockResp := &model.PatchTaskResponse{
			Task: &model.Task{
				ID:          taskID,
				UserID:      userID,
				Title:       "既存タスク",
				Comment:     "既存コメント",
				Priority:    "medium",
				Status:      "not_started",
				DueDatetime: nil,
				IsPinned:    false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockSvc := &mockTaskService{
			patchTaskFunc: func(ctx context.Context, uID, tID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
				assert.False(t, req.HasChanges())
				return mockResp, nil
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(`{}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := PatchTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp model.PatchTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.NotNil(t, resp.Task)
		assert.Equal(t, "既存タスク", resp.Task.Title)
	})

	t.Run("異常系: 未ログイン（Context に userID なし）の場合に 401 UNAUTHORIZED を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(`{"title": "test"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		// userID をセットしない

		h := PatchTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "UNAUTHORIZED", errResp.Error.Code)
		assert.NotEmpty(t, errResp.Error.Message)
		assert.NotNil(t, errResp.Error.Details)
		assert.Empty(t, errResp.Error.Details)
	})

	t.Run("異常系: JSON 構文不正またはバリデーションエラー時に 400 BAD_REQUEST を返すこと", func(t *testing.T) {
		// 1. JSON構文不正
		mockSvc := &mockTaskService{}
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(`{invalid json`))
		c1.Request.Header.Set("Content-Type", "application/json")
		c1.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c1.Set("userID", userID)

		h := PatchTaskHandler(mockSvc)
		h(c1)

		assert.Equal(t, http.StatusBadRequest, w1.Code)
		var errResp1 model.ErrorResponse
		err := json.Unmarshal(w1.Body.Bytes(), &errResp1)
		require.NoError(t, err)
		assert.Equal(t, "BAD_REQUEST", errResp1.Error.Code)

		// 2. バリデーションエラー（Service が 400 を返却）
		mockSvcValidate := &mockTaskService{
			patchTaskFunc: func(ctx context.Context, uID, tID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
				return nil, model.NewBadRequestError("入力内容に不備があります。", []model.ErrorDetail{
					{Field: "title", Message: "タイトルは1文字以上100文字以内で入力してください。"},
				})
			},
		}
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(`{"title": ""}`))
		c2.Request.Header.Set("Content-Type", "application/json")
		c2.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c2.Set("userID", userID)

		h2 := PatchTaskHandler(mockSvcValidate)
		h2(c2)

		assert.Equal(t, http.StatusBadRequest, w2.Code)
		var errResp2 model.ErrorResponse
		err = json.Unmarshal(w2.Body.Bytes(), &errResp2)
		require.NoError(t, err)
		assert.Equal(t, "BAD_REQUEST", errResp2.Error.Code)
		require.Len(t, errResp2.Error.Details, 1)
		assert.Equal(t, "title", errResp2.Error.Details[0].Field)
	})

	t.Run("異常系: 対象タスクが存在しない（または他者所有）場合に 404 NOT_FOUND を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{
			patchTaskFunc: func(ctx context.Context, uID, tID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
				return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(`{"title": "更新"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := PatchTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
		assert.Equal(t, "指定されたタスクが見つかりません。", errResp.Error.Message)
		assert.Empty(t, errResp.Error.Details)
	})

	t.Run("異常系: サーバー内部エラー発生時に 500 INTERNAL_SERVER_ERROR を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{
			patchTaskFunc: func(ctx context.Context, uID, tID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
				return nil, assert.AnError
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/tasks/"+taskID, bytes.NewBufferString(`{"title": "更新"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "task_id", Value: taskID}}
		c.Set("userID", userID)

		h := PatchTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", errResp.Error.Code)
		assert.NotEmpty(t, errResp.Error.Message)
	})
}
