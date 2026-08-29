package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"synctask/backend/model"
	"synctask/backend/util"
)

// CreateTask は単一または繰り返しタスクの作成処理を実行します。
func (s *taskService) CreateTask(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	// 1. バリデーション実行
	if err := req.Validate(); err != nil {
		return nil, err
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	isPinned := false
	if req.IsPinned != nil {
		isPinned = *req.IsPinned
	}

	isRecurring := req.IsRecurring != nil && *req.IsRecurring
	searchText := util.NormalizeSearchText(req.Title, req.Comment)
	now := time.Now().In(jst)

	// 2. 繰り返しタスク作成
	if isRecurring {
		rule := req.RecurringRule
		startT, err := time.ParseInLocation("2006-01-02", rule.StartDate, jst)
		if err != nil {
			return nil, model.NewBadRequestError("開始日の形式が正しくありません（YYYY-MM-DD）。", []model.ErrorDetail{
				{Field: "recurring_rule.start_date", Message: "開始日の形式が正しくありません。"},
			})
		}
		endT, err := time.ParseInLocation("2006-01-02", rule.EndDate, jst)
		if err != nil {
			return nil, model.NewBadRequestError("終了日の形式が正しくありません（YYYY-MM-DD）。", []model.ErrorDetail{
				{Field: "recurring_rule.end_date", Message: "終了日の形式が正しくありません。"},
			})
		}

		dueHour, dueMin := 23, 59
		if rule.DueTime != "" {
			_, _ = fmt.Sscanf(rule.DueTime, "%d:%d", &dueHour, &dueMin)
		}

		targetDays := make(map[string]bool)
		for _, day := range rule.DaysOfWeek {
			targetDays[strings.ToLower(day)] = true
		}

		var tasks []*model.Task
		for cur := startT; !cur.After(endT); cur = cur.AddDate(0, 0, 1) {
			dayName := weekdayMap[cur.Weekday()]
			if targetDays[dayName] {
				dueDt := time.Date(cur.Year(), cur.Month(), cur.Day(), dueHour, dueMin, 0, 0, jst)
				task := &model.Task{
					ID:          uuid.NewString(),
					UserID:      userID,
					Title:       req.Title,
					Comment:     req.Comment,
					Priority:    priority,
					Status:      "not_started",
					DueDatetime: &dueDt,
					IsPinned:    isPinned,
					SearchText:  searchText,
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				tasks = append(tasks, task)
			}
		}

		if len(tasks) == 0 {
			return nil, model.NewBadRequestError("入力内容に不備があります。", []model.ErrorDetail{
				{
					Field:   "recurring_rule",
					Message: "指定された期間内に該当する曜日が存在しません",
				},
			})
		}

		if len(tasks) > 100 {
			return nil, model.NewBadRequestError("入力内容に不備があります。", []model.ErrorDetail{
				{
					Field:   "recurring_rule",
					Message: "生成件数が上限（100件）を超えています",
				},
			})
		}

		if err := s.repo.CreateTasks(ctx, tasks); err != nil {
			return nil, err
		}

		return &model.CreateTaskResponse{
			CreatedCount: len(tasks),
			Tasks:        tasks,
		}, nil
	}

	// 3. 単一タスク作成
	var dueDatetime *time.Time
	if req.DueDatetime != nil && *req.DueDatetime != "" {
		dueStr := *req.DueDatetime
		if len(dueStr) == 10 {
			t, err := time.ParseInLocation("2006-01-02", dueStr, jst)
			if err != nil {
				return nil, model.NewBadRequestError("締切日時の形式が不正です。", []model.ErrorDetail{
					{Field: "due_datetime", Message: "YYYY-MM-DD 形式で指定してください。"},
				})
			}
			jstTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, jst)
			dueDatetime = &jstTime
		} else {
			t, err := time.Parse(time.RFC3339, dueStr)
			if err != nil {
				t, err = time.ParseInLocation("2006-01-02T15:04:05", dueStr, jst)
				if err != nil {
					return nil, model.NewBadRequestError("締切日時の形式が不正です。", []model.ErrorDetail{
						{Field: "due_datetime", Message: "ISO 8601 形式で指定してください。"},
					})
				}
			}
			jstTime := t.In(jst)
			dueDatetime = &jstTime
		}
	}

	task := &model.Task{
		ID:          uuid.NewString(),
		UserID:      userID,
		Title:       req.Title,
		Comment:     req.Comment,
		Priority:    priority,
		Status:      "not_started",
		DueDatetime: dueDatetime,
		IsPinned:    isPinned,
		SearchText:  searchText,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	return &model.CreateTaskResponse{
		CreatedCount: 1,
		Tasks:        []*model.Task{task},
	}, nil
}
