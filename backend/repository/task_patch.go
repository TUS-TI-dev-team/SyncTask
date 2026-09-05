package repository

import (
	"context"
	"database/sql"
	"errors"

	"synctask/backend/model"
)

const updateTaskQuery = `
UPDATE TASK
SET
	TITLE = $1,
	PRIORITY = $2,
	DUE_DATETIME = $3,
	STATUS = $4,
	IS_PINNED = $5,
	COMMENT = $6,
	SEARCH_TEXT = $7,
	UPDATED_AT = $8
WHERE TASK_ID = $9 AND USER_ID = $10
RETURNING
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
`

// UpdateTask は指定されたタスク情報を DB 上で更新し、更新後のタスクエンティティを返します。
//
// @spec
// - 指定された taskID かつ userID に一致するレコードを更新
// - 対象レコードが存在しない（または他者所有）場合は nil, nil を返却
// - DBエラー発生時は nil, error を返却
func (r *taskRepository) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	row := r.db.QueryRowContext(
		ctx,
		updateTaskQuery,
		task.Title,
		task.Priority,
		task.DueDatetime,
		task.Status,
		task.IsPinned,
		task.Comment,
		task.SearchText,
		task.UpdatedAt,
		task.ID,
		task.UserID,
	)

	var updated model.Task
	err := row.Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Title,
		&updated.Priority,
		&updated.DueDatetime,
		&updated.Status,
		&updated.IsPinned,
		&updated.Comment,
		&updated.SearchText,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &updated, nil
}
