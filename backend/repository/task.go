package repository

import (
	"context"
	"database/sql"
	"synctask/backend/model"
)

// TaskRepository はタスクのデータ永続化インターフェースです。
type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	CreateTasks(ctx context.Context, tasks []*model.Task) error
}

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository は TaskRepository の新しいインスタンスを生成します。
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

const insertTaskQuery = `
INSERT INTO TASK (
	TASK_ID,
	USER_ID,
	TITLE,
	PRIORITY,
	DUE_DATETIME,
	STATUS,
	IS_PINNED,
	COMMENT,
	SEARCH_TEXT,
	CREATED_AT,
	UPDATED_AT
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

// CreateTask は単一タスクを DB に登録します。
func (r *taskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	_, err := r.db.ExecContext(
		ctx,
		insertTaskQuery,
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
	)
	return err
}

// CreateTasks は複数タスクをトランザクション内で一括登録します。
func (r *taskRepository) CreateTasks(ctx context.Context, tasks []*model.Task) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, task := range tasks {
		_, err := tx.ExecContext(
			ctx,
			insertTaskQuery,
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
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
