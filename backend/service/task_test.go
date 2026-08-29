package service

import (
	"context"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTaskRepository は TaskRepository インターフェースの手動モックです。
type mockTaskRepository struct {
	createTaskFunc   func(ctx context.Context, task *model.Task) error
	createTasksFunc  func(ctx context.Context, tasks []*model.Task) error
	createTaskCalls  int
	createTasksCalls int
}

func (m *mockTaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	m.createTaskCalls++
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, task)
	}
	return nil
}

func (m *mockTaskRepository) CreateTasks(ctx context.Context, tasks []*model.Task) error {
	m.createTasksCalls++
	if m.createTasksFunc != nil {
		return m.createTasksFunc(ctx, tasks)
	}
	return nil
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func TestTaskService_CreateTask(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	userID := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("正常系: 単一タスクの作成処理が成功し、生成されたタスクが返却されること", func(t *testing.T) {
		repo := &mockTaskRepository{
			createTaskFunc: func(ctx context.Context, task *model.Task) error {
				assert.Equal(t, "課題レポート提出", task.Title)
				assert.Equal(t, "第5章の要約を含むこと", task.Comment)
				assert.Equal(t, "high", task.Priority)
				assert.Equal(t, userID, task.UserID)
				assert.Equal(t, "not_started", task.Status)
				assert.False(t, task.IsPinned)
				require.NotNil(t, task.DueDatetime)
				return nil
			},
		}
		svc := NewTaskService(repo)

		dueStr := "2026-08-20T23:59:00+09:00"
		req := &model.CreateTaskRequest{
			Title:       "課題レポート提出",
			Comment:     "第5章の要約を含むこと",
			Priority:    "high",
			DueDatetime: &dueStr,
			IsPinned:    boolPtr(false),
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.CreatedCount)
		require.Len(t, res.Tasks, 1)

		task := res.Tasks[0]
		assert.NotEmpty(t, task.ID)
		_, parseErr := uuid.Parse(task.ID)
		assert.NoError(t, parseErr)
		assert.Equal(t, userID, task.UserID)
		assert.Equal(t, "課題レポート提出", task.Title)
		assert.Equal(t, "第5章の要約を含むこと", task.Comment)
		assert.Equal(t, "high", task.Priority)
		assert.Equal(t, "not_started", task.Status)
		assert.False(t, task.IsPinned)
		require.NotNil(t, task.DueDatetime)
		expectedTime := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)
		assert.True(t, expectedTime.Equal(*task.DueDatetime))
		assert.False(t, task.CreatedAt.IsZero())
		assert.False(t, task.UpdatedAt.IsZero())

		assert.Equal(t, 1, repo.createTaskCalls)
		assert.Equal(t, 0, repo.createTasksCalls)
	})

	t.Run("正常系: 繰り返しタスクの作成処理で期間内の該当日タスクが昇順で正しく生成されること", func(t *testing.T) {
		repo := &mockTaskRepository{
			createTasksFunc: func(ctx context.Context, tasks []*model.Task) error {
				require.Len(t, tasks, 3)
				assert.Equal(t, "週次ゼミ発表準備", tasks[0].Title)
				assert.Equal(t, userID, tasks[0].UserID)
				return nil
			},
		}
		svc := NewTaskService(repo)

		req := &model.CreateTaskRequest{
			Title:       "週次ゼミ発表準備",
			Comment:     "進捗スライド作成",
			Priority:    "medium",
			IsPinned:    boolPtr(false),
			IsRecurring: boolPtr(true),
			RecurringRule: &model.RecurringRule{
				StartDate:  "2026-08-22", // Saturday
				EndDate:    "2026-09-05", // Saturday
				DaysOfWeek: []string{"saturday"},
				DueTime:    "18:00",
			},
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 3, res.CreatedCount)
		require.Len(t, res.Tasks, 3)

		expectedDates := []time.Time{
			time.Date(2026, 8, 22, 18, 0, 0, 0, jst),
			time.Date(2026, 8, 29, 18, 0, 0, 0, jst),
			time.Date(2026, 9, 5, 18, 0, 0, 0, jst),
		}

		for i, task := range res.Tasks {
			assert.NotEmpty(t, task.ID)
			_, parseErr := uuid.Parse(task.ID)
			assert.NoError(t, parseErr)
			assert.Equal(t, userID, task.UserID)
			assert.Equal(t, "週次ゼミ発表準備", task.Title)
			assert.Equal(t, "進捗スライド作成", task.Comment)
			assert.Equal(t, "medium", task.Priority)
			assert.Equal(t, "not_started", task.Status)
			assert.False(t, task.IsPinned)
			require.NotNil(t, task.DueDatetime)
			assert.True(t, expectedDates[i].Equal(*task.DueDatetime))
			if i > 0 {
				assert.True(t, res.Tasks[i].DueDatetime.After(*res.Tasks[i-1].DueDatetime))
			}
		}

		assert.Equal(t, 0, repo.createTaskCalls)
		assert.Equal(t, 1, repo.createTasksCalls)
	})

	t.Run("境界値: 繰り返し作成で生成件数がちょうど1件の場合に成功すること", func(t *testing.T) {
		repo := &mockTaskRepository{
			createTasksFunc: func(ctx context.Context, tasks []*model.Task) error {
				require.Len(t, tasks, 1)
				return nil
			},
		}
		svc := NewTaskService(repo)

		req := &model.CreateTaskRequest{
			Title:       "単発繰り返しタスク",
			Priority:    "low",
			IsRecurring: boolPtr(true),
			RecurringRule: &model.RecurringRule{
				StartDate:  "2026-08-22", // Saturday
				EndDate:    "2026-08-22", // Saturday
				DaysOfWeek: []string{"saturday"},
				DueTime:    "12:00",
			},
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.CreatedCount)
		require.Len(t, res.Tasks, 1)

		expectedDate := time.Date(2026, 8, 22, 12, 0, 0, 0, jst)
		require.NotNil(t, res.Tasks[0].DueDatetime)
		assert.True(t, expectedDate.Equal(*res.Tasks[0].DueDatetime))

		assert.Equal(t, 0, repo.createTaskCalls)
		assert.Equal(t, 1, repo.createTasksCalls)
	})

	t.Run("境界値: 繰り返し作成で生成件数がちょうど100件の場合に成功すること", func(t *testing.T) {
		repo := &mockTaskRepository{
			createTasksFunc: func(ctx context.Context, tasks []*model.Task) error {
				require.Len(t, tasks, 100)
				return nil
			},
		}
		svc := NewTaskService(repo)

		// 2026-01-01 から 100日間 (2026-01-01 〜 2026-04-10)
		req := &model.CreateTaskRequest{
			Title:       "毎日タスク100日分",
			Priority:    "medium",
			IsRecurring: boolPtr(true),
			RecurringRule: &model.RecurringRule{
				StartDate:  "2026-01-01",
				EndDate:    "2026-04-10",
				DaysOfWeek: []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
				DueTime:    "09:00",
			},
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 100, res.CreatedCount)
		require.Len(t, res.Tasks, 100)

		assert.Equal(t, 0, repo.createTaskCalls)
		assert.Equal(t, 1, repo.createTasksCalls)
	})

	t.Run("異常系: 繰り返し作成で生成件数が0件の場合に指定のエラーメッセージを返すこと", func(t *testing.T) {
		repo := &mockTaskRepository{}
		svc := NewTaskService(repo)

		// 2026-08-24 (月) 〜 2026-08-25 (火) の期間で日曜日を指定 -> 0件
		req := &model.CreateTaskRequest{
			Title:       "該当なし繰り返しタスク",
			Priority:    "medium",
			IsRecurring: boolPtr(true),
			RecurringRule: &model.RecurringRule{
				StartDate:  "2026-08-24", // Monday
				EndDate:    "2026-08-25", // Tuesday
				DaysOfWeek: []string{"sunday"},
				DueTime:    "10:00",
			},
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.Error(t, err)
		assert.Nil(t, res)

		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		require.NotEmpty(t, appErr.Details)
		assert.Equal(t, "recurring_rule", appErr.Details[0].Field)
		assert.Equal(t, "指定された期間内に該当する曜日が存在しません", appErr.Details[0].Message)

		assert.Equal(t, 0, repo.createTaskCalls)
		assert.Equal(t, 0, repo.createTasksCalls)
	})

	t.Run("異常系: 繰り返し作成で生成件数が101件以上の場合に指定のエラーメッセージを返すこと", func(t *testing.T) {
		repo := &mockTaskRepository{}
		svc := NewTaskService(repo)

		// 2026-01-01 から 101日間 (2026-01-01 〜 2026-04-11)
		req := &model.CreateTaskRequest{
			Title:       "上限超過タスク",
			Priority:    "medium",
			IsRecurring: boolPtr(true),
			RecurringRule: &model.RecurringRule{
				StartDate:  "2026-01-01",
				EndDate:    "2026-04-11",
				DaysOfWeek: []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
				DueTime:    "09:00",
			},
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.Error(t, err)
		assert.Nil(t, res)

		appErr, ok := err.(*model.AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		require.NotEmpty(t, appErr.Details)
		assert.Equal(t, "recurring_rule", appErr.Details[0].Field)
		assert.Equal(t, "生成件数が上限（100件）を超えています", appErr.Details[0].Message)

		assert.Equal(t, 0, repo.createTaskCalls)
		assert.Equal(t, 0, repo.createTasksCalls)
	})

	t.Run("準正常系: due_datetime に日付のみ（YYYY-MM-DD）が指定された場合、23:59:00+09:00 が補完されること", func(t *testing.T) {
		repo := &mockTaskRepository{
			createTaskFunc: func(ctx context.Context, task *model.Task) error {
				expectedTime := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)
				require.NotNil(t, task.DueDatetime)
				assert.True(t, expectedTime.Equal(*task.DueDatetime))
				return nil
			},
		}
		svc := NewTaskService(repo)

		req := &model.CreateTaskRequest{
			Title:       "日付のみ締切タスク",
			Priority:    "medium",
			DueDatetime: strPtr("2026-08-20"),
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Tasks, 1)

		expectedTime := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)
		require.NotNil(t, res.Tasks[0].DueDatetime)
		assert.True(t, expectedTime.Equal(*res.Tasks[0].DueDatetime))

		assert.Equal(t, 1, repo.createTaskCalls)
		assert.Equal(t, 0, repo.createTasksCalls)
	})

	t.Run("準正常系: UTCや他タイムゾーンの due_datetime が JST に変換・正規化されること", func(t *testing.T) {
		repo := &mockTaskRepository{
			createTaskFunc: func(ctx context.Context, task *model.Task) error {
				expectedTime := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)
				require.NotNil(t, task.DueDatetime)
				assert.True(t, expectedTime.Equal(*task.DueDatetime))
				return nil
			},
		}
		svc := NewTaskService(repo)

		// UTC 14:59 は JST (UTC+9) で 23:59
		req := &model.CreateTaskRequest{
			Title:       "UTC指定タスク",
			Priority:    "medium",
			DueDatetime: strPtr("2026-08-20T14:59:00Z"),
		}

		res, err := svc.CreateTask(context.Background(), userID, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Tasks, 1)

		expectedTime := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)
		require.NotNil(t, res.Tasks[0].DueDatetime)
		assert.True(t, expectedTime.Equal(*res.Tasks[0].DueDatetime))

		assert.Equal(t, 1, repo.createTaskCalls)
		assert.Equal(t, 0, repo.createTasksCalls)
	})
}
