package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_GetTaskByID(t *testing.T) {
	t.Run("正常系: 指定した task_id と user_id に一致するタスクが取得できること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		userID := "550e8400-e29b-41d4-a716-446655440000"
		due := time.Date(2026, 8, 20, 23, 59, 0, 0, time.UTC)
		now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

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

		rows := sqlmock.NewRows(columns).AddRow(
			taskID,
			userID,
			"課題レポート提出",
			"high",
			due,
			"not_started",
			true,
			"第5章の要約を含むこと",
			"課題レポート提出 第5章の要約を含むこと",
			now,
			now,
		)

		mock.ExpectQuery("SELECT TASK_ID, USER_ID, TITLE, PRIORITY, DUE_DATETIME, STATUS, IS_PINNED, COMMENT, SEARCH_TEXT, CREATED_AT, UPDATED_AT FROM TASK WHERE TASK_ID = \\$1 AND USER_ID = \\$2").
			WithArgs(taskID, userID).
			WillReturnRows(rows)

		task, err := repo.GetTaskByID(ctx, userID, taskID)
		require.NoError(t, err)
		require.NotNil(t, task)
		assert.Equal(t, taskID, task.ID)
		assert.Equal(t, userID, task.UserID)
		assert.Equal(t, "課題レポート提出", task.Title)
		assert.Equal(t, "第5章の要約を含むこと", task.Comment)
		assert.Equal(t, "high", task.Priority)
		assert.Equal(t, "not_started", task.Status)
		assert.Equal(t, &due, task.DueDatetime)
		assert.True(t, task.IsPinned)
		assert.Equal(t, "課題レポート提出 第5章の要約を含むこと", task.SearchText)
		assert.Equal(t, now, task.CreatedAt)
		assert.Equal(t, now, task.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("正常系: 締切日時が未設定（NULL）かつコメントが空文字のタスクが正常に取得できること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		userID := "550e8400-e29b-41d4-a716-446655440000"
		now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

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

		rows := sqlmock.NewRows(columns).AddRow(
			taskID,
			userID,
			"買い物リスト作成",
			"medium",
			nil, // DUE_DATETIME is NULL
			"not_started",
			false,
			"", // COMMENT is empty string
			"買い物リスト作成",
			now,
			now,
		)

		mock.ExpectQuery("SELECT TASK_ID, USER_ID, TITLE, PRIORITY, DUE_DATETIME, STATUS, IS_PINNED, COMMENT, SEARCH_TEXT, CREATED_AT, UPDATED_AT FROM TASK WHERE TASK_ID = \\$1 AND USER_ID = \\$2").
			WithArgs(taskID, userID).
			WillReturnRows(rows)

		task, err := repo.GetTaskByID(ctx, userID, taskID)
		require.NoError(t, err)
		require.NotNil(t, task)
		assert.Equal(t, taskID, task.ID)
		assert.Equal(t, userID, task.UserID)
		assert.Equal(t, "買い物リスト作成", task.Title)
		assert.Equal(t, "", task.Comment)
		assert.Equal(t, "medium", task.Priority)
		assert.Equal(t, "not_started", task.Status)
		assert.Nil(t, task.DueDatetime)
		assert.False(t, task.IsPinned)
		assert.Equal(t, "買い物リスト作成", task.SearchText)
		assert.Equal(t, now, task.CreatedAt)
		assert.Equal(t, now, task.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("準正常系: 該当レコードが存在しない場合（sql.ErrNoRows）に nil が返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		taskID := "non-existent-task-id"
		userID := "550e8400-e29b-41d4-a716-446655440000"

		mock.ExpectQuery("SELECT TASK_ID, USER_ID, TITLE, PRIORITY, DUE_DATETIME, STATUS, IS_PINNED, COMMENT, SEARCH_TEXT, CREATED_AT, UPDATED_AT FROM TASK WHERE TASK_ID = \\$1 AND USER_ID = \\$2").
			WithArgs(taskID, userID).
			WillReturnError(sql.ErrNoRows)

		task, err := repo.GetTaskByID(ctx, userID, taskID)
		assert.NoError(t, err)
		assert.Nil(t, task)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("異常系: DBクエリエラー発生時にエラーが返却されること", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewTaskRepository(db)
		ctx := context.Background()

		taskID := "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		userID := "550e8400-e29b-41d4-a716-446655440000"
		dbErr := errors.New("database error")

		mock.ExpectQuery("SELECT TASK_ID, USER_ID, TITLE, PRIORITY, DUE_DATETIME, STATUS, IS_PINNED, COMMENT, SEARCH_TEXT, CREATED_AT, UPDATED_AT FROM TASK WHERE TASK_ID = \\$1 AND USER_ID = \\$2").
			WithArgs(taskID, userID).
			WillReturnError(dbErr)

		task, err := repo.GetTaskByID(ctx, userID, taskID)
		require.Error(t, err)
		assert.Nil(t, task)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
