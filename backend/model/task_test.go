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
}
