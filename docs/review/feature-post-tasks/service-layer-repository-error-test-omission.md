# Service 層における Repository エラー伝播テストの欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-29 14:55:00
- **Target Files**:
  - [backend/service/task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task_test.go)

## 1. 問題の概要
`backend/service/task_test.go` において、単一タスク作成（`CreateTask`）および繰り返しタスク作成（`CreateTasks`）時に Repository がエラーを返却した場合の異常系テストケースが実装されていません。

## 2. 詳細な指摘内容
1. `TestTaskService_CreateTask` では、正常系、境界値（1件、100件）、バリデーションエラー（0件、101件、日付フォーマット不正等）は網羅されていますが、`mockTaskRepository.CreateTask` または `mockTaskRepository.CreateTasks` が DB エラー等のエラーを返した際に Service がそのエラーを呼び出し元へ正しく伝播するかを検証するテストが存在しません。
2. `backend/TESTING_GUIDE.md` の設計方針に基づき、各層のモック連携において依存先からのエラーが正しくハンドリングされることを検証することは重要な品質保証要件です。

## 3. 推奨される修正案
`backend/service/task_test.go` に以下の2つの異常系サブテストケースを追加します：
1. 単一タスク作成時に `repo.CreateTask` がエラーを返した場合に、Service がそのエラーをそのまま返却すること。
2. 繰り返しタスク作成時に `repo.CreateTasks` がエラーを返した場合に、Service がそのエラーをそのまま返却すること。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 15:13:00
- **Status**: Resolved

### 実施した修正内容
- `backend/service/task_test.go` の `TestTaskService_CreateTask` 内に、以下の2つの異常系テストケースを追加しました：
  1. 「異常系: 単一タスク作成時にリポジトリがエラーを返却した場合、Serviceがそのエラーを伝播すること」
  2. 「異常系: 繰り返しタスク作成時にリポジトリがエラーを返却した場合、Serviceがそのエラーを伝播すること」
- それぞれモックリポジトリからエラーを返却させ、Service が nil レスポンスおよび該当のエラーを正しく呼び出し元へ伝播することを検証しました。

### 変更したファイル
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task_test.go)
