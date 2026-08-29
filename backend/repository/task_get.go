package repository

import (
	"context"
	"synctask/backend/model"
)

// GetTaskByID は指定された taskID および userID に一致するタスクを取得します。
func (r *taskRepository) GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error) {
	// Step 2 で実装
	return nil, nil
}
