package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	validStatuses = map[string]bool{
		"not_started": true,
		"in_progress": true,
		"completed":   true,
	}
	jst = time.FixedZone("JST", 9*60*60)
)

// PatchTaskRequest はタスク部分更新リクエストを表します。
//
// @spec PATCH /api/tasks/{task_id}
// - 指定されたフィールドのみを更新対象とする部分更新リクエスト
// - 未指定のフィールドは更新せず既存の値を保持する
// - 空ボディ {} の場合は変更を行わずに現在のタスク情報を返却する
// - システムの読み取り専用フィールド（id, user_id, created_at, updated_at）は無視する
type PatchTaskRequest struct {
	Title       *string
	Comment     *string
	Priority    *string
	Status      *string
	DueDatetime *string
	IsPinned    *bool

	titlePresent       bool
	commentPresent     bool
	priorityPresent    bool
	statusPresent      bool
	dueDatetimePresent bool
	isPinnedPresent    bool

	commentNull     bool
	dueDatetimeNull bool

	typeErrors []ErrorDetail
}

// TitlePresent は title フィールドがリクエストに含まれていたかを返します。
func (r *PatchTaskRequest) TitlePresent() bool { return r.titlePresent }

// CommentPresent は comment フィールドがリクエストに含まれていたかを返します。
func (r *PatchTaskRequest) CommentPresent() bool { return r.commentPresent }

// PriorityPresent は priority フィールドがリクエストに含まれていたかを返します。
func (r *PatchTaskRequest) PriorityPresent() bool { return r.priorityPresent }

// StatusPresent は status フィールドがリクエストに含まれていたかを返します。
func (r *PatchTaskRequest) StatusPresent() bool { return r.statusPresent }

// DueDatetimePresent は due_datetime フィールドがリクエストに含まれていたかを返します。
func (r *PatchTaskRequest) DueDatetimePresent() bool { return r.dueDatetimePresent }

// IsPinnedPresent は is_pinned フィールドがリクエストに含まれていたかを返します。
func (r *PatchTaskRequest) IsPinnedPresent() bool { return r.isPinnedPresent }

// CommentNull は comment フィールドに明示的に null が指定されたかを返します。
func (r *PatchTaskRequest) CommentNull() bool { return r.commentNull }

// DueDatetimeNull は due_datetime フィールドに明示的に null が指定されたかを返します。
func (r *PatchTaskRequest) DueDatetimeNull() bool { return r.dueDatetimeNull }

// HasChanges は更新対象フィールドが1つ以上指定されているかを判定します。
func (r *PatchTaskRequest) HasChanges() bool {
	return r.titlePresent || r.commentPresent || r.priorityPresent ||
		r.statusPresent || r.dueDatetimePresent || r.isPinnedPresent
}

// UnmarshalJSON は JSON から PatchTaskRequest へのカスタムデコードを行います。
// 各フィールドの未指定・値指定・明示的 null を正確に識別します。
func (r *PatchTaskRequest) UnmarshalJSON(data []byte) error {
	// JSON 構文自体のバリデーション
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var rawMap map[string]json.RawMessage
	if err := dec.Decode(&rawMap); err != nil {
		return NewBadRequestError("不正なJSONフォーマットです。", nil)
	}

	r.typeErrors = nil

	// 1. title (非Null許容 string)
	if raw, ok := rawMap["title"]; ok {
		r.titlePresent = true
		if string(raw) == "null" {
			r.typeErrors = append(r.typeErrors, ErrorDetail{
				Field:   "title",
				Message: "タイトルにnullを指定することはできません。",
			})
		} else {
			var val string
			if err := json.Unmarshal(raw, &val); err != nil {
				r.typeErrors = append(r.typeErrors, ErrorDetail{
					Field:   "title",
					Message: "タイトルの型が不正です。",
				})
			} else {
				r.Title = &val
			}
		}
	}

	// 2. comment (Null許容 string)
	if raw, ok := rawMap["comment"]; ok {
		r.commentPresent = true
		if string(raw) == "null" {
			r.commentNull = true
			empty := ""
			r.Comment = &empty
		} else {
			var val string
			if err := json.Unmarshal(raw, &val); err != nil {
				r.typeErrors = append(r.typeErrors, ErrorDetail{
					Field:   "comment",
					Message: "コメントの型が不正です。",
				})
			} else {
				r.Comment = &val
			}
		}
	}

	// 3. priority (非Null許容 string enum)
	if raw, ok := rawMap["priority"]; ok {
		r.priorityPresent = true
		if string(raw) == "null" {
			r.typeErrors = append(r.typeErrors, ErrorDetail{
				Field:   "priority",
				Message: "優先度にnullを指定することはできません。",
			})
		} else {
			var val string
			if err := json.Unmarshal(raw, &val); err != nil {
				r.typeErrors = append(r.typeErrors, ErrorDetail{
					Field:   "priority",
					Message: "優先度の型が不正です。",
				})
			} else {
				r.Priority = &val
			}
		}
	}

	// 4. status (非Null許容 string enum)
	if raw, ok := rawMap["status"]; ok {
		r.statusPresent = true
		if string(raw) == "null" {
			r.typeErrors = append(r.typeErrors, ErrorDetail{
				Field:   "status",
				Message: "ステータスにnullを指定することはできません。",
			})
		} else {
			var val string
			if err := json.Unmarshal(raw, &val); err != nil {
				r.typeErrors = append(r.typeErrors, ErrorDetail{
					Field:   "status",
					Message: "ステータスの型が不正です。",
				})
			} else {
				r.Status = &val
			}
		}
	}

	// 5. due_datetime (Null許容 string)
	if raw, ok := rawMap["due_datetime"]; ok {
		r.dueDatetimePresent = true
		if string(raw) == "null" {
			r.dueDatetimeNull = true
			r.DueDatetime = nil
		} else {
			var val string
			if err := json.Unmarshal(raw, &val); err != nil {
				r.typeErrors = append(r.typeErrors, ErrorDetail{
					Field:   "due_datetime",
					Message: "締切日時の型が不正です。",
				})
			} else {
				r.DueDatetime = &val
			}
		}
	}

	// 6. is_pinned (非Null許容 bool)
	if raw, ok := rawMap["is_pinned"]; ok {
		r.isPinnedPresent = true
		if string(raw) == "null" {
			r.typeErrors = append(r.typeErrors, ErrorDetail{
				Field:   "is_pinned",
				Message: "ピン留めフラグにnullを指定することはできません。",
			})
		} else {
			var val bool
			if err := json.Unmarshal(raw, &val); err != nil {
				r.typeErrors = append(r.typeErrors, ErrorDetail{
					Field:   "is_pinned",
					Message: "ピン留めフラグの型が不正です（booleanを指定してください）。",
				})
			} else {
				r.IsPinned = &val
			}
		}
	}

	return nil
}

// Validate はリクエストパラメータの業務バリデーションを行います。
//
// @spec
// - title: 1〜100文字、トリム後空文字不可、制御文字（\n\r\t）不可
// - comment: 0〜1000文字、改行は \n に正規化
// - priority: low, medium, high のいずれか
// - status: not_started, in_progress, completed のいずれか
// - due_datetime: ISO 8601 または YYYY-MM-DD 形式
func (r *PatchTaskRequest) Validate() error {
	var details []ErrorDetail
	details = append(details, r.typeErrors...)

	// 1. title 検証
	if r.titlePresent && r.Title != nil {
		trimmed := strings.TrimSpace(*r.Title)
		r.Title = &trimmed
		if *r.Title == "" {
			details = append(details, ErrorDetail{
				Field:   "title",
				Message: "タイトルは1文字以上100文字以内で入力してください。",
			})
		} else if strings.ContainsAny(*r.Title, "\n\r\t") {
			details = append(details, ErrorDetail{
				Field:   "title",
				Message: "タイトルに改行やタブを含めることはできません。",
			})
		} else if utf8.RuneCountInString(*r.Title) > 100 {
			details = append(details, ErrorDetail{
				Field:   "title",
				Message: "タイトルは100文字以内で入力してください。",
			})
		}
	}

	// 2. comment 検証
	if r.commentPresent && r.Comment != nil && !r.commentNull {
		comment := strings.TrimSpace(*r.Comment)
		comment = strings.ReplaceAll(comment, "\r\n", "\n")
		comment = strings.ReplaceAll(comment, "\r", "\n")
		r.Comment = &comment
		if utf8.RuneCountInString(*r.Comment) > 1000 {
			details = append(details, ErrorDetail{
				Field:   "comment",
				Message: "コメントは1000文字以内で入力してください。",
			})
		}
	}

	// 3. priority 検証
	if r.priorityPresent && r.Priority != nil {
		if !validPriorities[*r.Priority] {
			details = append(details, ErrorDetail{
				Field:   "priority",
				Message: "優先度は low, medium, high のいずれかを指定してください。",
			})
		}
	}

	// 4. status 検証
	if r.statusPresent && r.Status != nil {
		if !validStatuses[*r.Status] {
			details = append(details, ErrorDetail{
				Field:   "status",
				Message: "ステータスは not_started, in_progress, completed のいずれかを指定してください。",
			})
		}
	}

	// 5. due_datetime 検証
	if r.dueDatetimePresent && r.DueDatetime != nil && !r.dueDatetimeNull {
		dueStr := *r.DueDatetime
		var valid bool
		if len(dueStr) == 10 {
			if _, err := time.Parse("2006-01-02", dueStr); err == nil {
				valid = true
			}
		} else {
			if _, err := time.Parse(time.RFC3339, dueStr); err == nil {
				valid = true
			} else if _, err := time.Parse("2006-01-02T15:04:05", dueStr); err == nil {
				valid = true
			}
		}
		if !valid {
			details = append(details, ErrorDetail{
				Field:   "due_datetime",
				Message: "締切日時の形式が不正です（ISO 8601 または YYYY-MM-DD）。",
			})
		}
	}

	if len(details) > 0 {
		return NewBadRequestError("入力内容に不備があります。", details)
	}

	return nil
}

// ParsedDueDatetime は due_datetime 文字列を time.Time（JST）としてパースして返します。
// 未指定の場合は (nil, false, nil)、明示的 null の場合は (nil, true, nil) を返します。
func (r *PatchTaskRequest) ParsedDueDatetime() (*time.Time, bool, error) {
	if !r.dueDatetimePresent {
		return nil, false, nil
	}
	if r.dueDatetimeNull || r.DueDatetime == nil {
		return nil, true, nil
	}

	dueStr := *r.DueDatetime
	if len(dueStr) == 10 {
		t, err := time.ParseInLocation("2006-01-02", dueStr, jst)
		if err != nil {
			return nil, false, fmt.Errorf("invalid date format: %w", err)
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, jst)
		return &t, true, nil
	}

	t, err := time.Parse(time.RFC3339, dueStr)
	if err == nil {
		tJST := t.In(jst)
		return &tJST, true, nil
	}

	t, err = time.ParseInLocation("2006-01-02T15:04:05", dueStr, jst)
	if err == nil {
		return &t, true, nil
	}

	return nil, false, fmt.Errorf("invalid datetime format: %s", dueStr)
}

// PatchTaskResponse はタスク部分更新成功時のレスポンスを表します。
type PatchTaskResponse struct {
	// Task は更新後のタスクオブジェクトです
	Task *Task `json:"task"`
}
