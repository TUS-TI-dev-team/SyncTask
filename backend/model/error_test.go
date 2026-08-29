package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponse_JSONSerialization(t *testing.T) {
	t.Run("正常系: Details が空スライスの場合に JSON に details: [] が出力されること", func(t *testing.T) {
		resp := NewErrorResponse("UNAUTHORIZED", "認証が必要です。", nil)
		bytes, err := json.Marshal(resp)
		require.NoError(t, err)

		jsonStr := string(bytes)
		assert.Contains(t, jsonStr, `"details":[]`)

		expected := `{"error":{"code":"UNAUTHORIZED","message":"認証が必要です。","details":[]}}`
		assert.JSONEq(t, expected, jsonStr)
	})

	t.Run("正常系: Details に要素が含まれる場合に JSON に詳細リストが出力されること", func(t *testing.T) {
		details := []ErrorDetail{
			{Field: "title", Message: "タイトルは必須です。"},
		}
		resp := NewErrorResponse("BAD_REQUEST", "入力内容に不備があります。", details)
		bytes, err := json.Marshal(resp)
		require.NoError(t, err)

		jsonStr := string(bytes)
		expected := `{"error":{"code":"BAD_REQUEST","message":"入力内容に不備があります。","details":[{"field":"title","message":"タイトルは必須です。"}]}}`
		assert.JSONEq(t, expected, jsonStr)
	})

	t.Run("正常系: NewBadRequestError で Details が nil の場合でも空スライスに初期化されること", func(t *testing.T) {
		appErr := NewBadRequestError("入力内容に不備があります。", nil)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "BAD_REQUEST", appErr.Code)
		assert.NotNil(t, appErr.Details)
		assert.Empty(t, appErr.Details)

		resp := NewErrorResponse(appErr.Code, appErr.Message, appErr.Details)
		bytes, err := json.Marshal(resp)
		require.NoError(t, err)
		assert.Contains(t, string(bytes), `"details":[]`)
	})

	t.Run("正常系: NewUnauthorizedError で Details が空スライスで初期化されること", func(t *testing.T) {
		appErr := NewUnauthorizedError("認証が必要です。")
		assert.Equal(t, 401, appErr.StatusCode)
		assert.Equal(t, "UNAUTHORIZED", appErr.Code)
		assert.NotNil(t, appErr.Details)
		assert.Empty(t, appErr.Details)

		resp := NewErrorResponse(appErr.Code, appErr.Message, appErr.Details)
		bytes, err := json.Marshal(resp)
		require.NoError(t, err)
		assert.Contains(t, string(bytes), `"details":[]`)
	})
}
