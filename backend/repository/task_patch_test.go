package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"synctask/backend/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_UpdateTask(t *testing.T) {
	columns := []string{
		"TASK_ID",
		"USER_ID",
		"TITLE",
		"PRIORITY",
		"DUE_DATETIME",
		"STATUS",
		"IS_PINNED",
		"COMMENT",
		"SEARCH_TEXT",
		"CREATED_AT",
		"UPDATED_AT",
	}

	taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
	userID := "550e8400-e29b-41d4-a716-446655440000"
	due := time.Date(2026, 8, 21, 23, 59, 0, 0, time.UTC)
	created := time.Date(2026, 8, 17, 10, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC)

	t.Run("正常系: タスクが正常に更新され、更新後のモデルが返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		inputTask := &model.Task{
			ID:          taskID,
			UserID:      userID,
			Title:       "課題レポート提出（修正版）",
			Comment:     "参考文献の追記完了",
			Priority:    "high",
			Status:      "completed",
			DueDatetime: &due,
			IsPinned:    true,
			SearchText:  "課題レポート提出（修正版） 参考文献の追記完了",
			CreatedAt:   created,
			UpdatedAt:   updated,
		}

		rows := sqlmock.NewRows(columns).AddRow(
			taskID,
			userID,
			inputTask.Title,
			inputTask.Priority,
			inputTask.DueDatetime,
			inputTask.Status,
			inputTask.IsPinned,
			inputTask.Comment,
			inputTask.SearchText,
			created,
			updated,
		)

		mock.ExpectQuery("UPDATE TASK SET").
			WithArgs(
				inputTask.Title,
				inputTask.Priority,
				inputTask.DueDatetime,
				inputTask.Status,
				inputTask.IsPinned,
				inputTask.Comment,
				inputTask.SearchText,
				inputTask.UpdatedAt,
				inputTask.ID,
				inputTask.UserID,
			).
			WillReturnRows(rows)

		res, err := repo.UpdateTask(ctx, inputTask)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, taskID, res.ID)
		assert.Equal(t, userID, res.UserID)
		assert.Equal(t, "課題レポート提出（修正版）", res.Title)
		assert.Equal(t, "参考文献の追記完了", res.Comment)
		assert.Equal(t, "high", res.Priority)
		assert.Equal(t, "completed", res.Status)
		require.NotNil(t, res.DueDatetime)
		assert.True(t, due.Equal(*res.DueDatetime))
		assert.True(t, res.IsPinned)
		assert.Equal(t, "課題レポート提出（修正版） 参考文献の追記完了", res.SearchText)
		assert.Equal(t, created, res.CreatedAt)
		assert.Equal(t, updated, res.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("準正常系: 更新対象レコードが存在しない（または他者所有）場合に nil が返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		inputTask := &model.Task{
			ID:        "non-existent-task-id",
			UserID:    userID,
			Title:     "存在しないタスク",
			Priority:  "low",
			Status:    "not_started",
			UpdatedAt: updated,
		}

		mock.ExpectQuery("UPDATE TASK SET").
			WithArgs(
				inputTask.Title,
				inputTask.Priority,
				inputTask.DueDatetime,
				inputTask.Status,
				inputTask.IsPinned,
				inputTask.Comment,
				inputTask.SearchText,
				inputTask.UpdatedAt,
				inputTask.ID,
				inputTask.UserID,
			).
			WillReturnError(sql.ErrNoRows)

		res, err := repo.UpdateTask(ctx, inputTask)
		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBクエリエラー発生時にエラーが返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		inputTask := &model.Task{
			ID:        taskID,
			UserID:    userID,
			Title:     "エラー発生タスク",
			Priority:  "medium",
			Status:    "in_progress",
			UpdatedAt: updated,
		}

		dbErr := errors.New("database connection failed")
		mock.ExpectQuery("UPDATE TASK SET").
			WithArgs(
				inputTask.Title,
				inputTask.Priority,
				inputTask.DueDatetime,
				inputTask.Status,
				inputTask.IsPinned,
				inputTask.Comment,
				inputTask.SearchText,
				inputTask.UpdatedAt,
				inputTask.ID,
				inputTask.UserID,
			).
			WillReturnError(dbErr)

		res, err := repo.UpdateTask(ctx, inputTask)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
