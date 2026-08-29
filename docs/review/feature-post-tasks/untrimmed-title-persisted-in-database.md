# タイトルの前後空白トリム結果がDB永続化・レスポンスに反映されない問題

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-29 14:27:00
- **Target Files**:
  - [backend/model/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go#L100-L116)
  - [backend/service/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task.go#L96)
  - [backend/service/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task.go#L169)

## 1. 問題の概要
仕様書（`docs/design/api_design/04_tasks.md`）では「前後の空白文字（半角・全角スペース、タブ、改行）を除去（トリム）した上で1〜100文字必須」と定義されています。しかし、現在の実装ではバリデーション時にトリム後の文字数判定を行っているものの、トリム結果が構造体に反映されず、前後に空白を含んだ元の文字列（`req.Title`）のまま DB に保存およびレスポンス返却されています。

## 2. 詳細な指摘内容
1. `backend/model/task.go` の `CreateTaskRequest.Validate()` 内（L100）で `trimmedTitle := strings.TrimSpace(r.Title)` を用いて文字数チェックを行っていますが、`r.Title` 自体は更新されません。
2. `backend/service/task.go` の `CreateTask()` 内（L96, L169）で `Task` エンティティを生成する際、`Title: req.Title` と元の未トリム文字列をそのまま代入しています。
3. これにより、例えば `"  課題レポート提出  "` というリクエストが送信された場合、バリデーションは通過するものの、DB およびレスポンスには前後に空白が付与された状態で永続化・返却されます。検索文字列（`SEARCH_TEXT`）生成時のみ `util.NormalizeSearchText` 側でトリムされているため、表示用タイトルと検索用テキストの間でも不整合が生じます。

## 3. 推奨される修正案
1. `backend/service/task.go` 内でタスク生成を行う前に、タイトルおよびコメントをトリム・正規化した変数を用意し、それを `Task` エンティティの `Title` に設定するように修正します。
2. または、リクエストの正規化処理を統一的に行うヘルパーを用意するか、`req.Title = strings.TrimSpace(req.Title)` を適用してから処理を行うようにします。
3. 単体テスト（`backend/service/task_test.go`）に、前後に半角・全角空白を含むタイトルが指定された場合にトリムされたタイトルで永続化・返却されることを検証するテストケースを追加します。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 14:38:00
- **Status**: Resolved

### 実施した修正内容
- `backend/model/task.go` の `CreateTaskRequest.Validate()` 内で `r.Title = strings.TrimSpace(r.Title)` を実行し、リクエスト構造体自身のタイトルを正規化するように修正しました。
- これにより、Service 層および DB 永続化・レスポンス返却において一貫して前後の空白（半角・全角空白、タブ等）がトリムされたタイトルが反映されるようになりました。
- `backend/model/task_test.go` および `backend/service/task_test.go` にトリム処理と永続化・返却結果を検証するテストを追加しました。

### 変更したファイル
- [task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task_test.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task_test.go)
