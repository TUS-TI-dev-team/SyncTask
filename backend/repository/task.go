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
	GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error)
}

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository は TaskRepository の新しいインスタンスを生成します。
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}
