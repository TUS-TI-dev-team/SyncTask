package repository

import (
	"testing"
	"time"
	"context"
	"database/sql"

	"synctask/backend/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutRepository_AttemptLogout(t *testing.T) {
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)
	attempt := &model.LogoutAttempt	{
		SessionID: "e81d3567-07fd-43b0-9f90-d465fbeed493",
		Now: now,
	}

	t.Run("正常系: 正常なセッションID を取得際にuid の取得,セッションの削除,再認証のリセットがすべて正常に commit され完了ステータスを返すこと。", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()
		repo := NewLogoutRepository(db)

		mock.ExpectBegin()
	})

	t.Run("正常系: 正常なセッションID を取得際にuid の取得,セッションの削除,再認証のリセットがすべて正常に commit され完了ステータスを返すこと。", func(t *testing.T) {
	})

	t.Run("準正常系: 期限が切れたセッションID 取得時に uid の取得,セッションの削除を行い,再認証のリセットを行わないこ", func(t *testing.T) {
	})

	t.Run("準正常系: 存在しないセッションIDが送られてきた際にな削除処理を行わずにnotfound ステータスを返すこと。", func(t *testing.T) {
	})

	t.Run("準正常系: uid の取得にて不正な値が帰ってきた際にnotfoundステータスを返すこと。", func(t *testing.T) {
	})

	t.Run("異常系: 正常なセッションIDを取得時にuid の取得にて未知のエラーが発生した際にdb エラーステータスを返しロールバックされること。", func(t *testing.T) {
	})

	t.Run("異常系: 正常なセッションIDを取得時にuid の取得を行い,セッションの削除に未知のエラーが発生した際にdb エラーステータスを返すこと。", func(t *testing.T) {
	})

	t.Run("異常系: 正常なセッションIDを取得時に再認証のリセット時に未知のエラーが発生した際にdb エラーステータスを返しロールバックを行うこと。", func(t *testing.T) {
	})

	t.Run("異常系: トランザクションの生成失敗時にdb エラーを返しロールバック処理が行われないこと。", func(t *testing.T) {
	})

	t.Run("異常系: 正常なセッションID 取得時に commit にて未知のエラーが発生した際にdb エラーを返しロールバックを行うこと。", func(t *testing.T) {
	})

	t.Run("異常系: キャンセル済みのコンテキストを渡した際に適切なエラー及びdbfaild を返すこと。", func(t *testing.T) {
	})

	t.Run("境界値: EXPIRES_AT.Valid がfalse の際は再認証のリセットは行わずにcommit し完了ステータスを返すこと。", func(t *testing.T) {
	})

	t.Run("境界値: 正常なセッションID を取得時にEXPIRES_AT = NOW の際に再認証のリセットを行わず,正常にcommit され完了ステータスを返すこと。", func(t *testing.T) {
	})
}
