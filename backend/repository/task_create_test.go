package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_CreateTask(t *testing.T) {
	t.Run("正常系: 単一タスクの INSERT が正常に実行されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		due := time.Date(2026, 8, 20, 23, 59, 0, 0, time.UTC)
		now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
		task := &model.Task{
			ID:          "7c9e6679-7425-40de-944b-e07fc1f90ae7",
			UserID:      "550e8400-e29b-41d4-a716-446655440000",
			Title:       "課題レポート提出",
			Comment:     "第5章の要約を含むこと",
			Priority:    "high",
			Status:      "not_started",
			DueDatetime: &due,
			IsPinned:    false,
			SearchText:  "課題レポート提出 第5章の要約を含むこと",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO TASK").
			WithArgs(
				task.ID,
				task.UserID,
				task.Title,
				task.Priority,
				task.DueDatetime,
				task.Status,
				task.IsPinned,
				task.Comment,
				task.SearchText,
				task.CreatedAt,
				task.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateTask(ctx, task)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー発生時にエラーが返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		task := &model.Task{
			ID:         "7c9e6679-7425-40de-944b-e07fc1f90ae7",
			UserID:     "550e8400-e29b-41d4-a716-446655440000",
			Title:      "課題レポート提出",
			Priority:   "high",
			Status:     "not_started",
			SearchText: "課題レポート提出",
		}

		dbErr := errors.New("database error")
		mock.ExpectExec("INSERT INTO TASK").
			WithArgs(
				task.ID,
				task.UserID,
				task.Title,
				task.Priority,
				task.DueDatetime,
				task.Status,
				task.IsPinned,
				task.Comment,
				task.SearchText,
				task.CreatedAt,
				task.UpdatedAt,
			).
			WillReturnError(dbErr)

		err = repo.CreateTask(ctx, task)
		require.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTaskRepository_CreateTasks(t *testing.T) {
	t.Run("正常系: 繰り返しタスクの複数件 INSERT（トランザクション）が正常に実行されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		due1 := time.Date(2026, 8, 22, 18, 0, 0, 0, time.UTC)
		due2 := time.Date(2026, 8, 29, 18, 0, 0, 0, time.UTC)
		now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

		tasks := []*model.Task{
			{
				ID:          "task-id-1",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				Title:       "週次ゼミ発表準備 1",
				Comment:     "進捗スライド作成",
				Priority:    "medium",
				Status:      "not_started",
				DueDatetime: &due1,
				IsPinned:    false,
				SearchText:  "週次ゼミ発表準備 1 進捗スライド作成",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "task-id-2",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
				Title:       "週次ゼミ発表準備 2",
				Comment:     "進捗スライド作成",
				Priority:    "medium",
				Status:      "not_started",
				DueDatetime: &due2,
				IsPinned:    false,
				SearchText:  "週次ゼミ発表準備 2 進捗スライド作成",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mock.ExpectBegin()
		for _, task := range tasks {
			mock.ExpectExec("INSERT INTO TASK").
				WithArgs(
					task.ID,
					task.UserID,
					task.Title,
					task.Priority,
					task.DueDatetime,
					task.Status,
					task.IsPinned,
					task.Comment,
					task.SearchText,
					task.CreatedAt,
					task.UpdatedAt,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		err = repo.CreateTasks(ctx, tasks)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBエラー発生時にロールバックされエラーが返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		tasks := []*model.Task{
			{
				ID:         "task-id-1",
				UserID:     "550e8400-e29b-41d4-a716-446655440000",
				Title:      "週次ゼミ発表準備 1",
				Priority:   "medium",
				Status:     "not_started",
				SearchText: "週次ゼミ発表準備 1",
			},
			{
				ID:         "task-id-2",
				UserID:     "550e8400-e29b-41d4-a716-446655440000",
				Title:      "週次ゼミ発表準備 2",
				Priority:   "medium",
				Status:     "not_started",
				SearchText: "週次ゼミ発表準備 2",
			},
		}

		dbErr := errors.New("insert execution failed")

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO TASK").
			WithArgs(
				tasks[0].ID,
				tasks[0].UserID,
				tasks[0].Title,
				tasks[0].Priority,
				tasks[0].DueDatetime,
				tasks[0].Status,
				tasks[0].IsPinned,
				tasks[0].Comment,
				tasks[0].SearchText,
				tasks[0].CreatedAt,
				tasks[0].UpdatedAt,
			).
			WillReturnError(dbErr)
		mock.ExpectRollback()

		err = repo.CreateTasks(ctx, tasks)
		require.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
