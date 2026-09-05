package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchTaskRequest_UnmarshalAndValidate(t *testing.T) {
	t.Run("正常系: 全フィールド指定時に正常にデコードおよびバリデーションを通過すること", func(t *testing.T) {
		body := `{
			"title": "課題レポート提出（修正版）",
			"comment": "参考文献の追記完了",
			"priority": "high",
			"status": "completed",
			"due_datetime": "2026-08-21T23:59:00+09:00",
			"is_pinned": true
		}`

		var req PatchTaskRequest
		err := json.Unmarshal([]byte(body), &req)
		require.NoError(t, err)

		err = req.Validate()
		require.NoError(t, err)

		assert.True(t, req.HasChanges())
		assert.True(t, req.TitlePresent())
		require.NotNil(t, req.Title)
		assert.Equal(t, "課題レポート提出（修正版）", *req.Title)

		assert.True(t, req.CommentPresent())
		require.NotNil(t, req.Comment)
		assert.Equal(t, "参考文献の追記完了", *req.Comment)

		assert.True(t, req.PriorityPresent())
		require.NotNil(t, req.Priority)
		assert.Equal(t, "high", *req.Priority)

		assert.True(t, req.StatusPresent())
		require.NotNil(t, req.Status)
		assert.Equal(t, "completed", *req.Status)

		assert.True(t, req.DueDatetimePresent())
		require.NotNil(t, req.DueDatetime)
		assert.Equal(t, "2026-08-21T23:59:00+09:00", *req.DueDatetime)

		assert.True(t, req.IsPinnedPresent())
		require.NotNil(t, req.IsPinned)
		assert.True(t, *req.IsPinned)

		parsedDue, ok, dueErr := req.ParsedDueDatetime()
		require.NoError(t, dueErr)
		assert.True(t, ok)
		require.NotNil(t, parsedDue)
		assert.Equal(t, 2026, parsedDue.Year())
		assert.Equal(t, time.Month(8), parsedDue.Month())
		assert.Equal(t, 21, parsedDue.Day())
		assert.Equal(t, 23, parsedDue.Hour())
		assert.Equal(t, 59, parsedDue.Minute())
	})

	t.Run("正常系: 単一フィールド（statusのみ、is_pinnedのみ等）指定時に正常にデコードされること", func(t *testing.T) {
		// status のみ
		var reqStatus PatchTaskRequest
		err := json.Unmarshal([]byte(`{"status": "in_progress"}`), &reqStatus)
		require.NoError(t, err)
		require.NoError(t, reqStatus.Validate())
		assert.True(t, reqStatus.HasChanges())
		assert.True(t, reqStatus.StatusPresent())
		assert.False(t, reqStatus.TitlePresent())
		assert.False(t, reqStatus.CommentPresent())
		assert.False(t, reqStatus.PriorityPresent())
		assert.False(t, reqStatus.DueDatetimePresent())
		assert.False(t, reqStatus.IsPinnedPresent())
		require.NotNil(t, reqStatus.Status)
		assert.Equal(t, "in_progress", *reqStatus.Status)

		// is_pinned のみ
		var reqPinned PatchTaskRequest
		err = json.Unmarshal([]byte(`{"is_pinned": false}`), &reqPinned)
		require.NoError(t, err)
		require.NoError(t, reqPinned.Validate())
		assert.True(t, reqPinned.HasChanges())
		assert.True(t, reqPinned.IsPinnedPresent())
		assert.False(t, reqPinned.StatusPresent())
		require.NotNil(t, reqPinned.IsPinned)
		assert.False(t, *reqPinned.IsPinned)
	})

	t.Run("正常系: 空ボディ {} 指定時に HasChanges が false となりバリデーションを通過すること", func(t *testing.T) {
		var req PatchTaskRequest
		err := json.Unmarshal([]byte(`{}`), &req)
		require.NoError(t, err)
		require.NoError(t, req.Validate())
		assert.False(t, req.HasChanges())
	})

	t.Run("正常系: 読み取り専用フィールド（id, user_id 等）が含まれていても無視されて通過すること", func(t *testing.T) {
		body := `{
			"id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
			"user_id": "550e8400-e29b-41d4-a716-446655440000",
			"created_at": "2026-08-17T10:00:00+09:00",
			"updated_at": "2026-08-17T11:00:00+09:00",
			"title": "更新後タイトル"
		}`
		var req PatchTaskRequest
		err := json.Unmarshal([]byte(body), &req)
		require.NoError(t, err)
		require.NoError(t, req.Validate())
		assert.True(t, req.HasChanges())
		assert.True(t, req.TitlePresent())
		require.NotNil(t, req.Title)
		assert.Equal(t, "更新後タイトル", *req.Title)
	})

	t.Run("正常系: comment に null または空文字が指定された場合クリア対象となること", func(t *testing.T) {
		// comment: null
		var reqNull PatchTaskRequest
		err := json.Unmarshal([]byte(`{"comment": null}`), &reqNull)
		require.NoError(t, err)
		require.NoError(t, reqNull.Validate())
		assert.True(t, reqNull.HasChanges())
		assert.True(t, reqNull.CommentPresent())
		assert.True(t, reqNull.CommentNull())
		require.NotNil(t, reqNull.Comment)
		assert.Equal(t, "", *reqNull.Comment)

		// comment: ""
		var reqEmpty PatchTaskRequest
		err = json.Unmarshal([]byte(`{"comment": ""}`), &reqEmpty)
		require.NoError(t, err)
		require.NoError(t, reqEmpty.Validate())
		assert.True(t, reqEmpty.HasChanges())
		assert.True(t, reqEmpty.CommentPresent())
		assert.False(t, reqEmpty.CommentNull())
		require.NotNil(t, reqEmpty.Comment)
		assert.Equal(t, "", *reqEmpty.Comment)
	})

	t.Run("正常系: due_datetime に null が指定された場合クリア対象となること", func(t *testing.T) {
		var req PatchTaskRequest
		err := json.Unmarshal([]byte(`{"due_datetime": null}`), &req)
		require.NoError(t, err)
		require.NoError(t, req.Validate())
		assert.True(t, req.HasChanges())
		assert.True(t, req.DueDatetimePresent())
		assert.True(t, req.DueDatetimeNull())
		assert.Nil(t, req.DueDatetime)

		parsedDue, ok, dueErr := req.ParsedDueDatetime()
		require.NoError(t, dueErr)
		assert.True(t, ok)
		assert.Nil(t, parsedDue)
	})

	t.Run("正常系: due_datetime に YYYY-MM-DD 形式が指定された場合 JST 23:59:00 として解釈されること", func(t *testing.T) {
		var req PatchTaskRequest
		err := json.Unmarshal([]byte(`{"due_datetime": "2026-08-20"}`), &req)
		require.NoError(t, err)
		require.NoError(t, req.Validate())

		parsedDue, ok, dueErr := req.ParsedDueDatetime()
		require.NoError(t, dueErr)
		assert.True(t, ok)
		require.NotNil(t, parsedDue)
		assert.Equal(t, 2026, parsedDue.Year())
		assert.Equal(t, time.Month(8), parsedDue.Month())
		assert.Equal(t, 20, parsedDue.Day())
		assert.Equal(t, 23, parsedDue.Hour())
		assert.Equal(t, 59, parsedDue.Minute())
		assert.Equal(t, 0, parsedDue.Second())
		_, offset := parsedDue.Zone()
		assert.Equal(t, 9*60*60, offset)
	})

	t.Run("異常系: title に null が指定された場合に 400 エラーを返すこと", func(t *testing.T) {
		var req PatchTaskRequest
		err := json.Unmarshal([]byte(`{"title": null}`), &req)
		require.NoError(t, err)

		err = req.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		require.NotEmpty(t, appErr.Details)
		assert.Equal(t, "title", appErr.Details[0].Field)
	})

	t.Run("異常系: title が空文字または空白のみの場合に 400 エラーを返すこと", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
		}{
			{name: "空文字", title: `""`},
			{name: "半角スペースのみ", title: `"   "`},
			{name: "全角スペースのみ", title: `"　　"`},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var req PatchTaskRequest
				err := json.Unmarshal([]byte(`{"title": `+tc.title+`}`), &req)
				require.NoError(t, err)

				err = req.Validate()
				require.Error(t, err)
				appErr, ok := err.(*AppError)
				require.True(t, ok)
				assert.Equal(t, 400, appErr.StatusCode)
				require.NotEmpty(t, appErr.Details)
				assert.Equal(t, "title", appErr.Details[0].Field)
			})
		}
	})

	t.Run("異常系: title に改行やタブが含まれる場合に 400 エラーを返すこと", func(t *testing.T) {
		testCases := []struct {
			name  string
			title string
		}{
			{name: "改行を含む", title: `"タスク\nタイトル"`},
			{name: "タブを含む", title: `"タスク\tタイトル"`},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var req PatchTaskRequest
				err := json.Unmarshal([]byte(`{"title": `+tc.title+`}`), &req)
				require.NoError(t, err)

				err = req.Validate()
				require.Error(t, err)
				appErr, ok := err.(*AppError)
				require.True(t, ok)
				assert.Equal(t, 400, appErr.StatusCode)
				require.NotEmpty(t, appErr.Details)
				assert.Equal(t, "title", appErr.Details[0].Field)
			})
		}
	})

	t.Run("境界値: title が 100 文字の場合に通過し、101 文字で 400 エラーを返すこと", func(t *testing.T) {
		// 100文字: 正常
		title100 := strings.Repeat("あ", 100)
		bodyValid, _ := json.Marshal(map[string]string{"title": title100})
		var reqValid PatchTaskRequest
		err := json.Unmarshal(bodyValid, &reqValid)
		require.NoError(t, err)
		assert.NoError(t, reqValid.Validate())

		// 101文字: エラー
		title101 := strings.Repeat("あ", 101)
		bodyInvalid, _ := json.Marshal(map[string]string{"title": title101})
		var reqInvalid PatchTaskRequest
		err = json.Unmarshal(bodyInvalid, &reqInvalid)
		require.NoError(t, err)
		err = reqInvalid.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		require.NotEmpty(t, appErr.Details)
		assert.Equal(t, "title", appErr.Details[0].Field)
	})

	t.Run("境界値: comment が 1000 文字の場合に通過し、1001 文字で 400 エラーを返すこと", func(t *testing.T) {
		// 1000文字: 正常
		comment1000 := strings.Repeat("あ", 1000)
		bodyValid, _ := json.Marshal(map[string]string{"comment": comment1000})
		var reqValid PatchTaskRequest
		err := json.Unmarshal(bodyValid, &reqValid)
		require.NoError(t, err)
		assert.NoError(t, reqValid.Validate())

		// 1001文字: エラー
		comment1001 := strings.Repeat("あ", 1001)
		bodyInvalid, _ := json.Marshal(map[string]string{"comment": comment1001})
		var reqInvalid PatchTaskRequest
		err = json.Unmarshal(bodyInvalid, &reqInvalid)
		require.NoError(t, err)
		err = reqInvalid.Validate()
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		require.NotEmpty(t, appErr.Details)
		assert.Equal(t, "comment", appErr.Details[0].Field)
	})

	t.Run("異常系: priority に null または無効な値が指定された場合に 400 エラーを返すこと", func(t *testing.T) {
		invalidCases := []string{
			`{"priority": null}`,
			`{"priority": "invalid"}`,
			`{"priority": "HIGH"}`,
			`{"priority": 123}`,
		}

		for _, b := range invalidCases {
			var req PatchTaskRequest
			err := json.Unmarshal([]byte(b), &req)
			require.NoError(t, err)

			err = req.Validate()
			require.Error(t, err, "input: %s", b)
			appErr, ok := err.(*AppError)
			require.True(t, ok)
			assert.Equal(t, 400, appErr.StatusCode)
			require.NotEmpty(t, appErr.Details)
			assert.Equal(t, "priority", appErr.Details[0].Field)
		}
	})

	t.Run("異常系: status に null または無効な値が指定された場合に 400 エラーを返すこと", func(t *testing.T) {
		invalidCases := []string{
			`{"status": null}`,
			`{"status": "pending"}`,
			`{"status": "COMPLETED"}`,
			`{"status": 99}`,
		}

		for _, b := range invalidCases {
			var req PatchTaskRequest
			err := json.Unmarshal([]byte(b), &req)
			require.NoError(t, err)

			err = req.Validate()
			require.Error(t, err, "input: %s", b)
			appErr, ok := err.(*AppError)
			require.True(t, ok)
			assert.Equal(t, 400, appErr.StatusCode)
			require.NotEmpty(t, appErr.Details)
			assert.Equal(t, "status", appErr.Details[0].Field)
		}
	})

	t.Run("異常系: is_pinned に null または非 boolean 値が指定された場合に 400 エラーを返すこと", func(t *testing.T) {
		invalidCases := []string{
			`{"is_pinned": null}`,
			`{"is_pinned": "true"}`,
			`{"is_pinned": 1}`,
			`{"is_pinned": 0}`,
		}

		for _, b := range invalidCases {
			var req PatchTaskRequest
			err := json.Unmarshal([]byte(b), &req)
			require.NoError(t, err)

			err = req.Validate()
			require.Error(t, err, "input: %s", b)
			appErr, ok := err.(*AppError)
			require.True(t, ok)
			assert.Equal(t, 400, appErr.StatusCode)
			require.NotEmpty(t, appErr.Details)
			assert.Equal(t, "is_pinned", appErr.Details[0].Field)
		}
	})

	t.Run("異常系: due_datetime の形式が不正な場合に 400 エラーを返すこと", func(t *testing.T) {
		invalidDates := []string{
			"invalid-date",
			"2026/08/20",
			"2026-13-45",
			"2026-08-20 23:59:00",
		}

		for _, dt := range invalidDates {
			var req PatchTaskRequest
			err := json.Unmarshal([]byte(`{"due_datetime": "`+dt+`"}`), &req)
			require.NoError(t, err)

			err = req.Validate()
			require.Error(t, err, "date: %s", dt)
			appErr, ok := err.(*AppError)
			require.True(t, ok)
			assert.Equal(t, 400, appErr.StatusCode)
			require.NotEmpty(t, appErr.Details)
			assert.Equal(t, "due_datetime", appErr.Details[0].Field)
		}
	})

	t.Run("異常系: JSON の形式自体が不正な場合に 400 エラーを返すこと", func(t *testing.T) {
		var req PatchTaskRequest
		err := req.UnmarshalJSON([]byte(`{invalid-json`))
		require.Error(t, err)
		appErr, ok := err.(*AppError)
		require.True(t, ok)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		assert.Equal(t, "不正なJSONフォーマットです。", appErr.Message)
	})
}
