package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func TestCreateTaskRequest_Validate(t *testing.T) {
	t.Run("正常系: 有効な単一タスクリクエストがバリデーションを通過すること", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:       "課題レポート提出",
			Comment:     "第5章の要約を含むこと",
			Priority:    "high",
			DueDatetime: strPtr("2026-08-20T23:59:00+09:00"),
			IsPinned:    boolPtr(false),
			IsRecurring: boolPtr(false),
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("正常系: 有効な繰り返しタスクリクエストがバリデーションを通過すること", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:       "週次ゼミ発表準備",
			Comment:     "進捗スライド作成",
			Priority:    "medium",
			IsPinned:    boolPtr(false),
			IsRecurring: boolPtr(true),
			RecurringRule: &RecurringRule{
				StartDate:  "2026-08-22",
				EndDate:    "2026-10-31",
				DaysOfWeek: []string{"saturday"},
				DueTime:    "18:00",
			},
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("異常系: タイトルが空文字または空白文字のみの場合にエラーとなること", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
		}{
			{name: "空文字", title: ""},
			{name: "半角スペースのみ", title: "   "},
			{name: "全角スペースのみ", title: "　　"},
			{name: "混在スペースのみ", title: " 　 \t "},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := CreateTaskRequest{
					Title:    tc.title,
					Priority: "medium",
				}
				err := req.Validate()
				require.Error(t, err)
			})
		}
	})

	t.Run("異常系: タイトルに改行やタブ等の制御文字が含まれる場合にエラーとなること", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
		}{
			{name: "改行(LF)を含む", title: "タスク\nタイトル"},
			{name: "改行(CR)を含む", title: "タスク\rタイトル"},
			{name: "改行(CRLF)を含む", title: "タスク\r\nタイトル"},
			{name: "タブを含む", title: "タスク\tタイトル"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := CreateTaskRequest{
					Title:    tc.title,
					Priority: "medium",
				}
				err := req.Validate()
				require.Error(t, err)
			})
		}
	})

	t.Run("境界値: タイトルが100文字の場合に通過し、101文字でエラーとなること", func(t *testing.T) {
		// 100文字 (正常)
		reqValid := CreateTaskRequest{
			Title:    strings.Repeat("あ", 100),
			Priority: "medium",
		}
		assert.NoError(t, reqValid.Validate())

		// 101文字 (エラー)
		reqInvalid := CreateTaskRequest{
			Title:    strings.Repeat("あ", 101),
			Priority: "medium",
		}
		assert.Error(t, reqInvalid.Validate())
	})

	t.Run("境界値: コメントが1000文字の場合に通過し、1001文字でエラーとなること", func(t *testing.T) {
		// 1000文字 (正常)
		reqValid := CreateTaskRequest{
			Title:    "有効なタイトル",
			Comment:  strings.Repeat("あ", 1000),
			Priority: "medium",
		}
		assert.NoError(t, reqValid.Validate())

		// 1001文字 (エラー)
		reqInvalid := CreateTaskRequest{
			Title:    "有効なタイトル",
			Comment:  strings.Repeat("あ", 1001),
			Priority: "medium",
		}
		assert.Error(t, reqInvalid.Validate())
	})

	t.Run("異常系: priority に無効な値が指定された場合にエラーとなること", func(t *testing.T) {
		invalidPriorities := []string{"invalid", "urgent", "HIGH", "123", "none"}
		for _, p := range invalidPriorities {
			req := CreateTaskRequest{
				Title:    "有効なタイトル",
				Priority: p,
			}
			err := req.Validate()
			require.Error(t, err)
		}
	})

	t.Run("異常系: is_recurring=true で recurring_rule が nil の場合にエラーとなること", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:         "有効なタイトル",
			Priority:      "medium",
			IsRecurring:   boolPtr(true),
			RecurringRule: nil,
		}
		err := req.Validate()
		require.Error(t, err)
	})

	t.Run("異常系: start_date > end_date の場合にエラーとなること", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:       "有効なタイトル",
			Priority:    "medium",
			IsRecurring: boolPtr(true),
			RecurringRule: &RecurringRule{
				StartDate:  "2026-10-31",
				EndDate:    "2026-08-22",
				DaysOfWeek: []string{"saturday"},
				DueTime:    "18:00",
			},
		}
		err := req.Validate()
		require.Error(t, err)
	})

	t.Run("異常系: days_of_week が空配列または不正な曜日の場合にエラーとなること", func(t *testing.T) {
		testCases := []struct {
			name       string
			daysOfWeek []string
		}{
			{name: "空配列", daysOfWeek: []string{}},
			{name: "不正な曜日", daysOfWeek: []string{"funday"}},
			{name: "一部不正な曜日", daysOfWeek: []string{"monday", "invalid"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := CreateTaskRequest{
					Title:       "有効なタイトル",
					Priority:    "medium",
					IsRecurring: boolPtr(true),
					RecurringRule: &RecurringRule{
						StartDate:  "2026-08-22",
						EndDate:    "2026-10-31",
						DaysOfWeek: tc.daysOfWeek,
						DueTime:    "18:00",
					},
				}
				err := req.Validate()
				require.Error(t, err)
			})
		}
	})

	t.Run("異常系: due_time のフォーマットが不正（秒付きや範囲外）な場合にエラーとなること", func(t *testing.T) {
		invalidDueTimes := []struct {
			name    string
			dueTime string
		}{
			{name: "秒付き形式 (HH:mm:ss)", dueTime: "18:00:00"},
			{name: "時間範囲外 (24時以上)", dueTime: "25:00"},
			{name: "分範囲外 (60分以上)", dueTime: "12:60"},
			{name: "不正文字列", dueTime: "invalid"},
			{name: "区切りなし形式", dueTime: "1800"},
		}

		for _, tc := range invalidDueTimes {
			t.Run(tc.name, func(t *testing.T) {
				req := CreateTaskRequest{
					Title:       "有効なタイトル",
					Priority:    "medium",
					IsRecurring: boolPtr(true),
					RecurringRule: &RecurringRule{
						StartDate:  "2026-08-22",
						EndDate:    "2026-10-31",
						DaysOfWeek: []string{"saturday"},
						DueTime:    tc.dueTime,
					},
				}
				err := req.Validate()
				require.Error(t, err)
			})
		}
	})

	t.Run("正常系: タイトルとコメントの前後空白がトリムされ、改行コードが正規化されること", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:    "  課題レポート提出　　",
			Comment:  "\r\n第1章\r\n第2章\r第3章\n\r\n",
			Priority: "high",
		}
		err := req.Validate()
		require.NoError(t, err)
		assert.Equal(t, "課題レポート提出", req.Title)
		assert.Equal(t, "第1章\n第2章\n第3章", req.Comment)
	})

	t.Run("正常系/異常系: 単一タスクの due_datetime 形式バリデーション", func(t *testing.T) {
		validDueDatetimes := []string{
			"2026-08-20",
			"2026-08-20T23:59:00+09:00",
			"2026-08-20T14:59:00Z",
			"2026-08-20T23:59:00",
		}
		for _, dt := range validDueDatetimes {
			req := CreateTaskRequest{
				Title:       "有効なタイトル",
				DueDatetime: strPtr(dt),
			}
			assert.NoError(t, req.Validate(), "valid datetime: %s", dt)
		}

		invalidDueDatetimes := []string{
			"invalid-date",
			"2026-13-45",
			"2026/08/20",
			"2026-08-20 23:59:00",
		}
		for _, dt := range invalidDueDatetimes {
			req := CreateTaskRequest{
				Title:       "有効なタイトル",
				DueDatetime: strPtr(dt),
			}
			err := req.Validate()
			require.Error(t, err, "invalid datetime: %s", dt)
			appErr, ok := err.(*AppError)
			require.True(t, ok)
			require.Len(t, appErr.Details, 1)
			assert.Equal(t, "due_datetime", appErr.Details[0].Field)
		}
	})

	t.Run("境界値/異常系: 繰り返しルールの期間が1年以内の場合は通過し、1年を超える場合はエラーとなること", func(t *testing.T) {
		// 1年以内（同年月日の1年後）: 正常
		reqValid := CreateTaskRequest{
			Title:       "有効なタイトル",
			IsRecurring: boolPtr(true),
			RecurringRule: &RecurringRule{
				StartDate:  "2026-08-22",
				EndDate:    "2027-08-22",
				DaysOfWeek: []string{"saturday"},
			},
		}
		assert.NoError(t, reqValid.Validate())

		// 1年超（1年後 + 1日）: エラー
		reqInvalid := CreateTaskRequest{
			Title:       "有効なタイトル",
			IsRecurring: boolPtr(true),
			RecurringRule: &RecurringRule{
				StartDate:  "2026-08-22",
				EndDate:    "2027-08-23",
				DaysOfWeek: []string{"saturday"},
			},
		}
		err := reqInvalid.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		require.Len(t, appErr.Details, 1)
		assert.Equal(t, "recurring_rule.end_date", appErr.Details[0].Field)
		assert.Equal(t, "終了日は開始日から1年以内の日付を指定してください。", appErr.Details[0].Message)
	})

	t.Run("正常系: コメントが空白文字のみの場合は空文字に正規化されること", func(t *testing.T) {
		testCases := []struct {
			name    string
			comment string
		}{
			{name: "半角スペースのみ", comment: "   "},
			{name: "全角スペースのみ", comment: "　　"},
			{name: "改行とタブ混在", comment: "\r\n \t \n \r"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := CreateTaskRequest{
					Title:    "タイトル",
					Comment:  tc.comment,
					Priority: "medium",
				}
				err := req.Validate()
				require.NoError(t, err)
				assert.Equal(t, "", req.Comment)
			})
		}
	})

	t.Run("正常系: is_recurring=true の場合は due_datetime が不正な値でも無視されてバリデーションを通過すること", func(t *testing.T) {
		req := CreateTaskRequest{
			Title:       "繰り返しタスク",
			Priority:    "medium",
			DueDatetime: strPtr("invalid-datetime-format"),
			IsRecurring: boolPtr(true),
			RecurringRule: &RecurringRule{
				StartDate:  "2026-08-22",
				EndDate:    "2026-10-31",
				DaysOfWeek: []string{"saturday"},
				DueTime:    "18:00",
			},
		}
		err := req.Validate()
		assert.NoError(t, err)
	})
}
