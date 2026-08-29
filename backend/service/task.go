package service

import (
	"context"
	"synctask/backend/model"
	"synctask/backend/repository"
)

// TaskService はタスクに関するビジネスロジックのインターフェースです。
type TaskService interface {
	CreateTask(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error)
}

type taskService struct {
	repo repository.TaskRepository
}

// NewTaskService は TaskService の新しいインスタンスを生成します。
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// CreateTask は単一または繰り返しタスクの作成処理を実行します（Step 1 スタブ）。
func (s *taskService) CreateTask(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	return nil, nil
}
