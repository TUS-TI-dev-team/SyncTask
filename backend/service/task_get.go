package service

import (
	"context"
	"synctask/backend/model"
)

// GetTask は指定された taskID のタスク詳細を取得します。
func (s *taskService) GetTask(ctx context.Context, userID, taskID string) (*model.GetTaskResponse, error) {
	// Step 2 で実装
	return nil, nil
}
