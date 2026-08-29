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

// mockTaskService はテスト用の TaskService モックです。
type mockTaskService struct {
	createTaskFunc func(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error)
}

func (m *mockTaskService) CreateTask(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, userID, req)
	}
	return nil, nil
}

func TestCreateTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系: 正しいリクエストで 201 Created とレスポンスJSONが返却されること", func(t *testing.T) {
		now := time.Now()
		dueTime := now.Add(24 * time.Hour)
		mockResp := &model.CreateTaskResponse{
			CreatedCount: 1,
			Tasks: []*model.Task{
				{
					ID:          "7c9e6679-7425-40de-944b-e07fc1f90ae7",
					UserID:      "550e8400-e29b-41d4-a716-446655440000",
					Title:       "課題レポート提出",
					Comment:     "第5章の要約を含むこと",
					Priority:    "high",
					Status:      "not_started",
					DueDatetime: &dueTime,
					IsPinned:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		}

		mockSvc := &mockTaskService{
			createTaskFunc: func(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", userID)
				assert.Equal(t, "課題レポート提出", req.Title)
				assert.Equal(t, "第5章の要約を含むこと", req.Comment)
				assert.Equal(t, "high", req.Priority)
				return mockResp, nil
			},
		}

		reqBody := `{
			"title": "課題レポート提出",
			"comment": "第5章の要約を含むこと",
			"priority": "high",
			"due_datetime": "2026-08-20T23:59:00+09:00",
			"is_pinned": false
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		h := CreateTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp model.CreateTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, 1, resp.CreatedCount)
		require.Len(t, resp.Tasks, 1)
		assert.Equal(t, "7c9e6679-7425-40de-944b-e07fc1f90ae7", resp.Tasks[0].ID)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.Tasks[0].UserID)
		assert.Equal(t, "課題レポート提出", resp.Tasks[0].Title)
		assert.Equal(t, "第5章の要約を含むこと", resp.Tasks[0].Comment)
		assert.Equal(t, "high", resp.Tasks[0].Priority)
		assert.Equal(t, "not_started", resp.Tasks[0].Status)
		assert.False(t, resp.Tasks[0].IsPinned)
	})

	t.Run("正常系: 繰り返し作成で 201 Created と created_count, tasks 配列が返却されること", func(t *testing.T) {
		now := time.Now()
		due1 := now.Add(24 * time.Hour)
		due2 := now.Add(48 * time.Hour)
		mockResp := &model.CreateTaskResponse{
			CreatedCount: 2,
			Tasks: []*model.Task{
				{
					ID:          "task-1",
					UserID:      "550e8400-e29b-41d4-a716-446655440000",
					Title:       "週次ゼミ発表準備",
					Comment:     "進捗スライド作成",
					Priority:    "medium",
					Status:      "not_started",
					DueDatetime: &due1,
					IsPinned:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          "task-2",
					UserID:      "550e8400-e29b-41d4-a716-446655440000",
					Title:       "週次ゼミ発表準備",
					Comment:     "進捗スライド作成",
					Priority:    "medium",
					Status:      "not_started",
					DueDatetime: &due2,
					IsPinned:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		}

		mockSvc := &mockTaskService{
			createTaskFunc: func(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", userID)
				require.NotNil(t, req.IsRecurring)
				assert.True(t, *req.IsRecurring)
				require.NotNil(t, req.RecurringRule)
				assert.Equal(t, "2026-08-22", req.RecurringRule.StartDate)
				return mockResp, nil
			},
		}

		reqBody := `{
			"title": "週次ゼミ発表準備",
			"comment": "進捗スライド作成",
			"priority": "medium",
			"is_pinned": false,
			"is_recurring": true,
			"recurring_rule": {
				"start_date": "2026-08-22",
				"end_date": "2026-10-31",
				"days_of_week": ["saturday"],
				"due_time": "18:00"
			}
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		h := CreateTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp model.CreateTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, 2, resp.CreatedCount)
		require.Len(t, resp.Tasks, 2)
		assert.Equal(t, "task-1", resp.Tasks[0].ID)
		assert.Equal(t, "task-2", resp.Tasks[1].ID)
	})

	t.Run("正常系: タイトルの空白トリムおよびコメントの改行正規化が適用されたタスクがレスポンスに含まれること", func(t *testing.T) {
		now := time.Now()
		mockResp := &model.CreateTaskResponse{
			CreatedCount: 1,
			Tasks: []*model.Task{
				{
					ID:        "7c9e6679-7425-40de-944b-e07fc1f90ae7",
					UserID:    "550e8400-e29b-41d4-a716-446655440000",
					Title:     "課題レポート提出",
					Comment:   "第1行\n第2行",
					Priority:  "medium",
					Status:    "not_started",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		}

		mockSvc := &mockTaskService{
			createTaskFunc: func(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
				return mockResp, nil
			},
		}

		reqBody := `{
			"title": "  課題レポート提出  ",
			"comment": "第1行\r\n第2行",
			"priority": "medium"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		h := CreateTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp model.CreateTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "課題レポート提出", resp.Tasks[0].Title)
		assert.Equal(t, "第1行\n第2行", resp.Tasks[0].Comment)
	})

	t.Run("異常系: 未ログイン（Context に userID なし）の場合に 401 UNAUTHORIZED を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{}

		reqBody := `{
			"title": "課題レポート提出"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		// userID をセットしない

		h := CreateTaskHandler(mockSvc)
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

	t.Run("異常系: リクエストバリデーション違反時に 400 BAD_REQUEST と詳細 details を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{
			createTaskFunc: func(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
				return nil, model.NewBadRequestError("入力内容に不備があります。", []model.ErrorDetail{
					{
						Field:   "title",
						Message: "タイトルは必須です。",
					},
				})
			},
		}

		reqBody := `{
			"title": "",
			"priority": "medium"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		h := CreateTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "BAD_REQUEST", errResp.Error.Code)
		assert.Equal(t, "入力内容に不備があります。", errResp.Error.Message)
		require.Len(t, errResp.Error.Details, 1)
		assert.Equal(t, "title", errResp.Error.Details[0].Field)
		assert.Equal(t, "タイトルは必須です。", errResp.Error.Details[0].Message)
	})

	t.Run("異常系: 不正な JSON ボディの場合に 400 BAD_REQUEST と空の details 配列を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{}

		invalidJSON := `{"title": "不正なJSON`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		h := CreateTaskHandler(mockSvc)
		h(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResp model.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "BAD_REQUEST", errResp.Error.Code)
		assert.NotEmpty(t, errResp.Error.Message)
		assert.NotNil(t, errResp.Error.Details)
		assert.Empty(t, errResp.Error.Details)
		assert.Contains(t, w.Body.String(), `"details":[]`)
	})

	t.Run("異常系: サーバー内部エラー発生時に 500 INTERNAL_SERVER_ERROR と空の details 配列を返すこと", func(t *testing.T) {
		mockSvc := &mockTaskService{
			createTaskFunc: func(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
				return nil, assert.AnError
			},
		}

		reqBody := `{
			"title": "有効なタイトル",
			"priority": "medium"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		h := CreateTaskHandler(mockSvc)
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
