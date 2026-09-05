package repository

import (
	"context"
	"database/sql"
	"errors"
	"synctask/backend/model"
)

const selectTaskByIDQuery = `
SELECT
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
FROM TASK
WHERE TASK_ID = $1 AND USER_ID = $2
`

// GetTaskByID は指定された taskID および userID に一致するタスクを取得します。
func (r *taskRepository) GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error) {
	row := r.db.QueryRowContext(ctx, selectTaskByIDQuery, taskID, userID)

	var task model.Task
	err := row.Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Priority,
		&task.DueDatetime,
		&task.Status,
		&task.IsPinned,
		&task.Comment,
		&task.SearchText,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}
