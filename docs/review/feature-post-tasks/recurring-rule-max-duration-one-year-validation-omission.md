# 繰り返しタスク（recurring_rule）の期間制約（最大1年間 / 52週以内）バリデーションの欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-29 14:27:00
- **Target Files**:
  - [backend/model/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go#L177-L184)
  - [docs/plans/backend/post-tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/plans/backend/post-tasks.md#L57)

## 1. 問題の概要
仕様書（`docs/design/api_design/04_tasks.md` L129, L181）および開発計画書（`docs/plans/backend/post-tasks.md` L57）では、繰り返しタスクの終了日（`end_date`）について「最大1年間（52週以内）」「日付範囲不整合（`start_date > end_date` または 1年超）」と規定されています。しかし、現在の `model.CreateTaskRequest.Validate()` では `start_date > end_date` のチェックのみが存在し、期間が1年間（366日 / 52週）を超えているかどうかの検証が欠落しています。

## 2. 詳細な指摘内容
1. `backend/model/task.go` の L177-184 において、`startT.After(endT)` のチェックは行われていますが、`startT` から `endT` までの期間が1年を超えている場合の判定がありません。
2. 生成件数が100件以内の制約はあるものの、例えば特定月日のみや隔週指定などの指定パターンによって「2年間にわたり合計50件作成する」といったリクエストが送信された場合、期間制約（最大1年）のバリデーションをすり抜けてしまいます。

## 3. 推奨される修正案
1. `model.CreateTaskRequest.Validate()` において、`startT.AddDate(1, 0, 0).Before(endT)`（または `endT.Sub(startT) > 366*24*time.Hour`）の判定を追加し、1年を超える期間が指定された場合は `recurring_rule`（または `recurring_rule.end_date`）に対して `400 BAD_REQUEST` エラー（例: `"終了日は開始日から1年以内の日付を指定してください。"`）を返却するように実装します。
2. `backend/model/task_test.go` に、期間が1年を超える（367日以上）場合にバリデーションエラーとなるテストケースを追加します。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 14:38:00
- **Status**: Resolved

### 実施した修正内容
- `backend/model/task.go` の `CreateTaskRequest.Validate()` 内で、`startT.AddDate(1, 0, 0).Before(endT)` の期間超過判定を追加しました。
- 1年（366日 / 52週）を超える期間が指定された場合は `recurring_rule.end_date` フィールドに対して `終了日は開始日から1年以内の日付を指定してください。` の 400 Bad Request バリデーションエラーを返却するように実装しました。
- `backend/model/task_test.go` に1年以内の境界値（同年月日の1年後）および1年超過（1年後 + 1日）のテストケースを追加し、検証を完了しました。

### 変更したファイル
- [task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task_test.go)
