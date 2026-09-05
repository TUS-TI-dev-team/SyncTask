package model

import (
	"time"
)

// Task はタスクエンティティを表します。
type Task struct {
	// ID はタスクの一意識別子（UUID）です
	ID string `json:"id" example:"7c9e6679-7425-40de-944b-e07fc1f90ae7"`
	// UserID は所有ユーザーの一意識別子（UUID）です
	UserID string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Title はタスクのタイトルです（1〜100文字）
	Title string `json:"title" example:"課題レポート提出"`
	// Comment はタスクの補足コメント・詳細説明です（0〜1000文字）
	Comment string `json:"comment" example:"第5章の要約を含むこと"`
	// Priority は優先度です（high, medium, low）
	Priority string `json:"priority" example:"high" enums:"high,medium,low"`
	// Status はタスクのステータスです（not_started, in_progress, completed）
	Status string `json:"status" example:"not_started" enums:"not_started,in_progress,completed"`
	// DueDatetime は締切日時です（ISO 8601形式）
	DueDatetime *time.Time `json:"due_datetime" example:"2026-08-20T23:59:00+09:00"`
	// IsPinned はピン留めフラグです
	IsPinned bool `json:"is_pinned" example:"false"`
	// SearchText は検索用の正規化文字列です（JSON出力対象外）
	SearchText string `json:"-"`
	// CreatedAt は作成日時です
	CreatedAt time.Time `json:"created_at" example:"2026-08-17T12:00:00+09:00"`
	// UpdatedAt は更新日時です
	UpdatedAt time.Time `json:"updated_at" example:"2026-08-17T12:00:00+09:00"`
}
