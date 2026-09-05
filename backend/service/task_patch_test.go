package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskService_PatchTask(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	userID := "550e8400-e29b-41d4-a716-446655440000"
	taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
	now := time.Date(2026, 8, 17, 10, 0, 0, 0, jst)
	existingDue := time.Date(2026, 8, 20, 23, 59, 0, 0, jst)

	baseTask := &model.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       "課題レポート提出",
		Comment:     "第5章の要約を含むこと",
		Priority:    "medium",
		Status:      "not_started",
		DueDatetime: &existingDue,
		IsPinned:    false,
		SearchText:  "課題レポート提出 第5章の要約を含むこと",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("正常系: 指定フィールドが正常に更新され PatchTaskResponse が返却されること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{
			"title": "課題レポート提出（修正版）",
			"status": "completed",
			"is_pinned": true
		}`), &req)
		require.NoError(t, err)

		taskCopy := *baseTask
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, taskID, tID)
				return &taskCopy, nil
			},
			updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
				assert.Equal(t, "課題レポート提出（修正版）", task.Title)
				assert.Equal(t, "completed", task.Status)
				assert.True(t, task.IsPinned)
				assert.Equal(t, "medium", task.Priority) // 変更なし
				return task, nil
			},
		}

		svc := NewTaskService(repo)
		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)

		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Task)
		assert.Equal(t, "課題レポート提出（修正版）", res.Task.Title)
		assert.Equal(t, "completed", res.Task.Status)
		assert.True(t, res.Task.IsPinned)
		assert.Equal(t, 1, repo.getTaskByIDCalls)
		assert.Equal(t, 1, repo.updateTaskCalls)
	})

	t.Run("正常系: title または comment 更新時に SearchText が再生成されること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{
			"title": "タスク名更新",
			"comment": "コメント更新"
		}`), &req)
		require.NoError(t, err)

		taskCopy := *baseTask
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				return &taskCopy, nil
			},
			updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
				// util.NormalizeSearchText("タスク名更新", "コメント更新") の結果がセットされていること
				assert.Equal(t, "タスク名更新 コメント更新", task.SearchText)
				return task, nil
			},
		}

		svc := NewTaskService(repo)
		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "タスク名更新 コメント更新", res.Task.SearchText)
	})

	t.Run("正常系: 空リクエストボディ {} の場合に DB 更新をスキップして既存タスクがそのまま返却されること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{}`), &req)
		require.NoError(t, err)

		taskCopy := *baseTask
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				return &taskCopy, nil
			},
			updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
				t.Fatal("DB更新が呼ばれてはならない")
				return nil, nil
			},
		}

		svc := NewTaskService(repo)
		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, baseTask.Title, res.Task.Title)
		assert.Equal(t, 1, repo.getTaskByIDCalls)
		assert.Equal(t, 0, repo.updateTaskCalls)
	})

	t.Run("正常系: comment に null 指定時にコメントが空文字にクリアされること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{"comment": null}`), &req)
		require.NoError(t, err)

		taskCopy := *baseTask
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				return &taskCopy, nil
			},
			updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
				assert.Equal(t, "", task.Comment)
				return task, nil
			},
		}

		svc := NewTaskService(repo)
		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "", res.Task.Comment)
		assert.Equal(t, 1, repo.updateTaskCalls)
	})

	t.Run("正常系: due_datetime に null 指定時に締切日時が nil にクリアされること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{"due_datetime": null}`), &req)
		require.NoError(t, err)

		taskCopy := *baseTask
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				return &taskCopy, nil
			},
			updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
				assert.Nil(t, task.DueDatetime)
				return task, nil
			},
		}

		svc := NewTaskService(repo)
		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Nil(t, res.Task.DueDatetime)
		assert.Equal(t, 1, repo.updateTaskCalls)
	})

	t.Run("異常系: リクエストバリデーションエラー時に 400 BAD_REQUEST が返却されること", func(t *testing.T) {
		var req model.PatchTaskRequest
		// title に null
		err := json.Unmarshal([]byte(`{"title": null}`), &req)
		require.NoError(t, err)

		repo := &mockTaskRepository{}
		svc := NewTaskService(repo)

		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)
		require.Error(t, err)
		assert.Nil(t, res)

		var appErr *model.AppError
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		assert.Equal(t, 0, repo.getTaskByIDCalls)
	})

	t.Run("異常系: 対象タスクが存在しない（または他者所有）場合に 404 NOT_FOUND が返却されること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{"title": "更新後タイトル"}`), &req)
		require.NoError(t, err)

		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				return nil, nil // レコードなし
			},
		}
		svc := NewTaskService(repo)

		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)
		require.Error(t, err)
		assert.Nil(t, res)

		var appErr *model.AppError
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, 404, appErr.StatusCode)
		assert.Equal(t, "NOT_FOUND", appErr.Code)
		assert.Equal(t, "指定されたタスクが見つかりません。", appErr.Message)
		assert.Equal(t, 0, repo.updateTaskCalls)
	})

	t.Run("異常系: リポジトリ層でエラーが発生した場合にエラーがそのまま返却されること", func(t *testing.T) {
		var req model.PatchTaskRequest
		err := json.Unmarshal([]byte(`{"title": "更新後タイトル"}`), &req)
		require.NoError(t, err)

		dbErr := errors.New("db error")
		repo := &mockTaskRepository{
			getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
				return nil, dbErr
			},
		}
		svc := NewTaskService(repo)

		res, err := svc.PatchTask(context.Background(), userID, taskID, &req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, dbErr, err)
	})
}
