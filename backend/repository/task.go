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

// CreateTask は単一タスクを DB に登録します（Step 1 スタブ）。
func (r *taskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	return nil
}

// CreateTasks は複数タスクをトランザクション内で一括登録します（Step 1 スタブ）。
func (r *taskRepository) CreateTasks(ctx context.Context, tasks []*model.Task) error {
	return nil
}
