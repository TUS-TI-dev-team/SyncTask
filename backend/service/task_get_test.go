package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskService_GetTask(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	userID := "550e8400-e29b-41d4-a716-446655440000"
	taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"

	t.Run("正常系: 指定した task_id のタスクが正常に取得され GetTaskResponse が返却されること", func(t *testing.T) {
		dueTime := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)
		createdAt := time.Date(2026, 8, 17, 10, 0, 0, 0, jst)
		updatedAt := time.Date(2026, 8, 17, 11, 30, 0, 0, jst)

		expectedTask := &model.Task{
			ID:          taskID,
			UserID:      userID,
			Title:       "課題レポート提出",
			Comment:     "第5章の要約を含むこと",
			Priority:    "high",
			Status:      "in_progress",
			DueDatetime: &dueTime,
			IsPinned:    true,
			SearchText:  "課題レポート提出 第5章の要約を含むこと",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, taskID, tID)
				return expectedTask, nil
			},
		}
		svc := NewTaskService(repo)

		res, err := svc.GetTask(context.Background(), userID, taskID)

		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Task)
		assert.Equal(t, taskID, res.Task.ID)
		assert.Equal(t, userID, res.Task.UserID)
		assert.Equal(t, "課題レポート提出", res.Task.Title)
		assert.Equal(t, "第5章の要約を含むこと", res.Task.Comment)
		assert.Equal(t, "high", res.Task.Priority)
		assert.Equal(t, "in_progress", res.Task.Status)
		require.NotNil(t, res.Task.DueDatetime)
		assert.True(t, dueTime.Equal(*res.Task.DueDatetime))
		assert.True(t, res.Task.IsPinned)
		assert.True(t, createdAt.Equal(res.Task.CreatedAt))
		assert.True(t, updatedAt.Equal(res.Task.UpdatedAt))

		assert.Equal(t, 1, repo.getTaskByIDCalls)
	})

	t.Run("異常系: リポジトリが nil を返した場合（該当タスクなし）に 404 NOT_FOUND エラーが返却されること", func(t *testing.T) {
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, taskID, tID)
				return nil, nil
			},
		}
		svc := NewTaskService(repo)

		res, err := svc.GetTask(context.Background(), userID, taskID)

		require.Error(t, err)
		assert.Nil(t, res)

		var appErr *model.AppError
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, 404, appErr.StatusCode)
		assert.Equal(t, "NOT_FOUND", appErr.Code)
		assert.Equal(t, "指定されたタスクが見つかりません。", appErr.Message)
		assert.Empty(t, appErr.Details)

		assert.Equal(t, 1, repo.getTaskByIDCalls)
	})

	t.Run("異常系: リポジトリ層で予期せぬエラーが発生した場合にエラーがそのまま返却されること", func(t *testing.T) {
		expectedErr := errors.New("db query error")
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, taskID, tID)
				return nil, expectedErr
			},
		}
		svc := NewTaskService(repo)

		res, err := svc.GetTask(context.Background(), userID, taskID)

		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expectedErr, err)

		assert.Equal(t, 1, repo.getTaskByIDCalls)
	})
}
