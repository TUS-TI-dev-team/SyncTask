package service

import (
	"context"
	"time"

	"synctask/backend/model"
	"synctask/backend/util"
)

// PatchTask は指定された taskID のタスク情報を部分更新します。
//
// @spec
// - 既存タスクを取得し存在しない場合は 404 NOT_FOUND（IDOR/BOLA防止）
// - リクエストのバリデーションを実行（400 BAD_REQUEST）
// - 更新対象フィールドが1つも指定されていない場合は DB 更新をスキップし既存タスクを返却
// - 指定されたフィールドのみを既存タスクに上書き
// - title または comment 更新時は SearchText を再生成
// - UpdatedAt を JST の現在時刻に更新
// - リポジトリの UpdateTask を実行して更新結果を返却
func (s *taskService) PatchTask(ctx context.Context, userID, taskID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
	existing, err := s.repo.GetTaskByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !req.HasChanges() {
		return &model.PatchTaskResponse{
			Task: existing,
		}, nil
	}

	textChanged := false

	if req.TitlePresent() && req.Title != nil {
		existing.Title = *req.Title
		textChanged = true
	}

	if req.CommentPresent() {
		if req.CommentNull() {
			existing.Comment = ""
		} else if req.Comment != nil {
			existing.Comment = *req.Comment
		}
		textChanged = true
	}

	if req.PriorityPresent() && req.Priority != nil {
		existing.Priority = *req.Priority
	}

	if req.StatusPresent() && req.Status != nil {
		existing.Status = *req.Status
	}

	if req.DueDatetimePresent() {
		parsedDue, ok, dueErr := req.ParsedDueDatetime()
		if dueErr != nil {
			return nil, model.NewBadRequestError("締切日時の形式が不正です。", nil)
		}
		if ok {
			existing.DueDatetime = parsedDue
		}
	}

	if req.IsPinnedPresent() && req.IsPinned != nil {
		existing.IsPinned = *req.IsPinned
	}

	if textChanged {
		existing.SearchText = util.NormalizeSearchText(existing.Title, existing.Comment)
	}

	existing.UpdatedAt = time.Now().In(jst)

	updated, err := s.repo.UpdateTask(ctx, existing)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
	}

	return &model.PatchTaskResponse{
		Task: updated,
	}, nil
}
