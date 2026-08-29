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
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Comment     string     `json:"comment"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	DueDatetime *time.Time `json:"due_datetime"`
	IsPinned    bool       `json:"is_pinned"`
	SearchText  string     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RecurringRule は繰り返しタスク生成ルールを表します。
type RecurringRule struct {
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
	DaysOfWeek []string `json:"days_of_week"`
	DueTime    string   `json:"due_time"`
}

// CreateTaskRequest はタスク新規作成リクエストボディを表します。
type CreateTaskRequest struct {
	Title         string         `json:"title"`
	Comment       string         `json:"comment"`
	Priority      string         `json:"priority"`
	DueDatetime   *string        `json:"due_datetime"`
	IsPinned      *bool          `json:"is_pinned"`
	IsRecurring   *bool          `json:"is_recurring"`
	RecurringRule *RecurringRule `json:"recurring_rule"`
}

// CreateTaskResponse はタスク新規作成成功時のレスポンスを表します。
type CreateTaskResponse struct {
	CreatedCount int     `json:"created_count"`
	Tasks        []*Task `json:"tasks"`
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
