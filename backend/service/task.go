package service

import (
	"context"
	"time"

	"synctask/backend/model"
	"synctask/backend/repository"
)

var jst = time.FixedZone("JST", 9*60*60)

// TaskService はタスクに関するビジネスロジックのインターフェースです。
type TaskService interface {
	CreateTask(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error)
	GetTask(ctx context.Context, userID, taskID string) (*model.GetTaskResponse, error)
}

type taskService struct {
	repo repository.TaskRepository
}

// NewTaskService は TaskService の新しいインスタンスを生成します。
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

var weekdayMap = map[time.Weekday]string{
	time.Sunday:    "sunday",
	time.Monday:    "monday",
	time.Tuesday:   "tuesday",
	time.Wednesday: "wednesday",
	time.Thursday:  "thursday",
	time.Friday:    "friday",
	time.Saturday:  "saturday",
}
