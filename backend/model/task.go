package model

import "time"

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

// Validate はリクエストパラメータの検証を行います（Step 1 スタブ）。
func (r *CreateTaskRequest) Validate() error {
	return nil
}
