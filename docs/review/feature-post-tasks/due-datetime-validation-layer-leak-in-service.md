# due_datetime のフォーマットバリデーション責務が Model 層から漏出している問題

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-29 14:27:00
- **Target Files**:
  - [backend/model/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go#L96-L219)
  - [backend/service/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task.go#L140-L164)

## 1. 問題の概要
`CreateTaskRequest` のバリデーション（`model.CreateTaskRequest.Validate()`）において、`title` や `comment`、`priority`、`recurring_rule` などの検証は Model 層で行われている一方、単一タスク作成時の `due_datetime` のフォーマット（ISO 8601 または YYYY-MM-DD）の妥当性検証が Model 層で行われず、Service 層（`service.task.go`）でパース・エラー生成が行われています。

## 2. 詳細な指摘内容
1. `backend/model/task.go` の `Validate()` メソッド内には `DueDatetime` に対する構文検証ロジックが記述されていません。
2. そのため、リクエスト構造体の自己検証（Model層バリデーション）のみを呼び出した時点では、不正な日時文字列（例: `"invalid-date"` や `"2026-13-45"`）が含まれていてもバリデーションを通過してしまいます。
3. `backend/service/task.go` の L142-164 で `time.Parse` を呼び、失敗時に `model.NewBadRequestError` を返却する形となっており、入力形式バリデーションの責務が Model 層と Service 層に分散しています。

## 3. 推奨される修正案
1. `backend/model/task.go` の `CreateTaskRequest.Validate()` 内で、`DueDatetime` が指定されている場合に許容フォーマット（ISO 8601 または YYYY-MM-DD）に合致しているかを検証する処理を追加します。
2. Service 層では検証済みの前提で安全にパース・タイムゾーン変換を行えるように責務を整理します。
3. `backend/model/task_test.go` に不正な `due_datetime` が指定された場合にバリデーションエラーとなるテストケースを追加します。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 14:38:00
- **Status**: Resolved

### 実施した修正内容
- `backend/model/task.go` の `CreateTaskRequest.Validate()` 内に、単一タスク作成時の `DueDatetime` 形式検証（`YYYY-MM-DD`, `time.RFC3339`, `2006-01-02T15:04:05`）を追加しました。不正な文字列が指定された場合は `field: "due_datetime"` のバリデーションエラーを返却します。
- これにより、入力形式バリデーションの責務を Model 層に集約し、Service 層は検証済み前提で安全に JST 変換・パースを実行できるよう責務を分離しました。
- `backend/model/task_test.go` に正常な各日時フォーマットおよび不正な日時文字列を検証するテストケースを追加しました。

### 変更したファイル
- [task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task_test.go)
