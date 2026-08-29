package service

import (
	"context"
	"synctask/backend/model"
)

// GetTask は指定された taskID のタスク詳細を取得します。
func (s *taskService) GetTask(ctx context.Context, userID, taskID string) (*model.GetTaskResponse, error) {
	task, err := s.repo.GetTaskByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
	}

	return &model.GetTaskResponse{
		Task: task,
	}, nil
}
