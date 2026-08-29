# エラーレスポンスにおける details キー省略（omitempty）がAPI共通仕様に違反する問題

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-29 14:27:00
- **Target Files**:
  - [backend/model/error.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/error.go#L20)

## 1. 問題の概要
API共通仕様書（`docs/design/api_design/01_overview.md` L87）では、エラーレスポンスの `error.details` フィールドについて「フィールド単位のバリデーション詳細情報リスト。対象フィールドが存在しないエラー応答の場合は空配列 `[]` を返却（`null` やキー省略は不可）」と明記されています。しかし、現在の `ErrorBody` 構造体の定義では `Details []ErrorDetail json:"details,omitempty"` と `omitempty` タグが付与されているため、詳細情報がない場合にキーごと JSON から除外されてしまいます。

## 2. 詳細な指摘内容
1. `backend/model/error.go` の L20 にて `Details []ErrorDetail json:"details,omitempty"` と定義されています。
2. これにより、401 Unauthorized や 400（不正なJSON形式）、500（内部エラー）など、`Details` が空または `nil` のエラーレスポンスにおいて、JSON 出力結果が `{"error":{"code":"UNAUTHORIZED","message":"認証が必要です。"}}` となり、`details` キーが含まれません。
3. フロントエンドが共通エラーハンドリングで `res.error.details` を配列前提として操作（例: `res.error.details.find(...)`）した場合に JavaScript の `TypeError`（Cannot read properties of undefined）を誘発するリスクがあります。

## 3. 推奨される修正案
1. `backend/model/error.go` の `Details` フィールドのタグから `omitempty` を削除し、`Details []ErrorDetail json:"details"` とします。
2. レスポンス生成時やエラーインスタンス作成時に `details` が `nil` の場合は `details = []ErrorDetail{}`（空スライス）で初期化されるように保証し、常に `[]` が出力されるようにします。
3. ハンドラーのテスト等で、401 エラーや単一エラー時にも `details: []` が JSON に含まれていることを検証するアサーションを追加します。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 14:38:00
- **Status**: Resolved

### 実施した修正内容
- `backend/model/error.go` の `ErrorBody.Details` タグから `omitempty` を削除し、`Details []ErrorDetail json:"details"` としました。
- `NewBadRequestError`、`NewUnauthorizedError`、および `NewErrorResponse` 関数内で `Details` が `nil` の場合に `[]ErrorDetail{}` で初期化されるように保証しました。
- `backend/handler/task.go` で共通の `model.NewErrorResponse` を使用するように統一し、401/400/500 のすべてのエラーレスポンスで `details: []` が常に JSON 出力されることを `backend/handler/task_test.go` のテストで検証しました。

### 変更したファイル
- [error.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/error.go)
- [task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/handler/task.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/handler/task_test.go)
