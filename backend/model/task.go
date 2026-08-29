package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
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

// RecurringRule は繰り返しタスク生成ルールを表します。
type RecurringRule struct {
	// StartDate は繰り返し開始日です（YYYY-MM-DD形式）
	StartDate string `json:"start_date" example:"2026-08-22"`
	// EndDate は繰り返し終了日です（YYYY-MM-DD形式、最大1年間）
	EndDate string `json:"end_date" example:"2026-10-31"`
	// DaysOfWeek は繰り返す曜日のリストです（monday, tuesday, wednesday, thursday, friday, saturday, sunday）
	DaysOfWeek []string `json:"days_of_week" example:"saturday"`
	// DueTime は締切時刻です（HH:mm形式、デフォルト: 23:59）
	DueTime string `json:"due_time" example:"18:00"`
}

// CreateTaskRequest はタスク新規作成リクエストボディを表します。
type CreateTaskRequest struct {
	// Title はタスクのタイトルです（1〜100文字、必須）
	Title string `json:"title" example:"課題レポート提出"`
	// Comment はタスクの補足コメントです（0〜1000文字）
	Comment string `json:"comment" example:"第5章の要約を含むこと"`
	// Priority は優先度です（high, medium, low、デフォルト: medium）
	Priority string `json:"priority" example:"high" enums:"high,medium,low"`
	// DueDatetime は締切日時です（ISO 8601形式またはYYYY-MM-DD）
	DueDatetime *string `json:"due_datetime" example:"2026-08-20T23:59:00+09:00"`
	// IsPinned はピン留めフラグです（デフォルト: false）
	IsPinned *bool `json:"is_pinned" example:"false"`
	// IsRecurring は繰り返し一括生成フラグです（デフォルト: false）
	IsRecurring *bool `json:"is_recurring" example:"false"`
	// RecurringRule は繰り返し生成ルールです（IsRecurringがtrueの場合に必須）
	RecurringRule *RecurringRule `json:"recurring_rule"`
}

// CreateTaskResponse はタスク新規作成成功時のレスポンスを表します。
type CreateTaskResponse struct {
	// CreatedCount は作成されたタスクの件数です
	CreatedCount int `json:"created_count" example:"1"`
	// Tasks は作成されたタスクのリストです
	Tasks []*Task `json:"tasks"`
}

var (
	validPriorities = map[string]bool{
		"high":   true,
		"medium": true,
		"low":    true,
	}

	validDaysOfWeek = map[string]bool{
		"monday":    true,
		"tuesday":   true,
		"wednesday": true,
		"thursday":  true,
		"friday":    true,
		"saturday":  true,
		"sunday":    true,
	}

	timeFormatRegex = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)
)

// Validate はリクエストパラメータの検証を行います。
func (r *CreateTaskRequest) Validate() error {
	var details []ErrorDetail

	// 1. タイトル検証
	trimmedTitle := strings.TrimSpace(r.Title)
	if trimmedTitle == "" {
		details = append(details, ErrorDetail{
			Field:   "title",
			Message: "タイトルは必須です。",
		})
	} else if strings.ContainsAny(r.Title, "\n\r\t") {
		details = append(details, ErrorDetail{
			Field:   "title",
			Message: "タイトルに改行やタブを含めることはできません。",
		})
	} else if utf8.RuneCountInString(trimmedTitle) > 100 {
		details = append(details, ErrorDetail{
			Field:   "title",
			Message: "タイトルは100文字以内で入力してください。",
		})
	}

	// 2. コメント検証
	if utf8.RuneCountInString(r.Comment) > 1000 {
		details = append(details, ErrorDetail{
			Field:   "comment",
			Message: "コメントは1000文字以内で入力してください。",
		})
	}

	// 3. 優先度検証
	if r.Priority != "" && !validPriorities[r.Priority] {
		details = append(details, ErrorDetail{
			Field:   "priority",
			Message: "優先度は low, medium, high のいずれかを指定してください。",
		})
	}

	// 4. 繰り返しタスクルールの検証
	if r.IsRecurring != nil && *r.IsRecurring {
		if r.RecurringRule == nil {
			details = append(details, ErrorDetail{
				Field:   "recurring_rule",
				Message: "繰り返しタスクの作成には recurring_rule が必須です。",
			})
		} else {
			rule := r.RecurringRule

			var startT, endT time.Time
			var startErr, endErr error

			if rule.StartDate == "" {
				details = append(details, ErrorDetail{
					Field:   "recurring_rule.start_date",
					Message: "開始日は必須です。",
				})
			} else {
				startT, startErr = time.Parse("2006-01-02", rule.StartDate)
				if startErr != nil {
					details = append(details, ErrorDetail{
						Field:   "recurring_rule.start_date",
						Message: "開始日の形式が正しくありません（YYYY-MM-DD）。",
					})
				}
			}

			if rule.EndDate == "" {
				details = append(details, ErrorDetail{
					Field:   "recurring_rule.end_date",
					Message: "終了日は必須です。",
				})
			} else {
				endT, endErr = time.Parse("2006-01-02", rule.EndDate)
				if endErr != nil {
					details = append(details, ErrorDetail{
						Field:   "recurring_rule.end_date",
						Message: "終了日の形式が正しくありません（YYYY-MM-DD）。",
					})
				}
			}

			if startErr == nil && endErr == nil && !startT.IsZero() && !endT.IsZero() {
				if startT.After(endT) {
					details = append(details, ErrorDetail{
						Field:   "recurring_rule",
						Message: "開始日は終了日以前の日付を指定してください。",
					})
				}
			}

			if len(rule.DaysOfWeek) == 0 {
				details = append(details, ErrorDetail{
					Field:   "recurring_rule.days_of_week",
					Message: "1つ以上の曜日を指定してください。",
				})
			} else {
				for _, day := range rule.DaysOfWeek {
					if !validDaysOfWeek[strings.ToLower(day)] {
						details = append(details, ErrorDetail{
							Field:   "recurring_rule.days_of_week",
							Message: fmt.Sprintf("不正な曜日です: %s", day),
						})
						break
					}
				}
			}

			if rule.DueTime != "" {
				if !timeFormatRegex.MatchString(rule.DueTime) {
					details = append(details, ErrorDetail{
						Field:   "recurring_rule.due_time",
						Message: "締切時刻の形式が不正です（HH:mm）。",
					})
				}
			}
		}
	}

	if len(details) > 0 {
		return NewBadRequestError("入力内容に不備があります。", details)
	}

	return nil
}
