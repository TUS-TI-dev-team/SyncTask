# コメントの改行コード正規化（\n）およびトリム処理の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-29 14:27:00
- **Target Files**:
  - [backend/model/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go#L118-L124)
  - [backend/service/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task.go#L97)
  - [backend/service/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task.go#L170)

## 1. 問題の概要
仕様書（`docs/design/api_design/04_tasks.md`）では、コメント（`comment`）について「0〜1000文字（トリム後）。改行は `\n` に正規化。未入力時は空文字 `""` として登録」と明記されています。しかし、現在の実装ではトリム後の文字数カウントが行われておらず、CRLF（`\r\n`）や CR（`\r`）から LF（`\n`）への改行コード正規化処理が実装されていません。

## 2. 詳細な指摘内容
1. `backend/model/task.go` の L119 で `utf8.RuneCountInString(r.Comment) > 1000` としており、トリム前の文字数でチェックしています。仕様上は「0〜1000文字（トリム後）」であるため、前後に空白がある場合に本来許容される入力が弾かれるか、逆にトリム・正規化後の文字数制約と不一致になる可能性があります。
2. 改行コード（`\r\n`, `\r`）を `\n` に変換する処理が存在しないため、Windows 環境や一部の HTTP クライアントから送信された改行コードがそのまま DB に永続化されてしまいます。
3. `comment` 未入力時（省略時または空文字時）のハンドリングにおいて、トリムした結果空文字となる場合の正規化が明示されていません。

## 3. 推奨される修正案
1. `comment` に対しても `strings.TrimSpace` でトリムを行い、`strings.ReplaceAll(comment, "\r\n", "\n")` および `strings.ReplaceAll(comment, "\r", "\n")` 等により改行コードを `\n` に正規化する処理を追加します。
2. バリデーション時は正規化後の文字数（0〜1000文字）をチェックします。
3. 改行を含むコメント（`\r\n`, `\r`）が `\n` に正規化されて DB 保存・返却されることを検証する単体テストを追加します。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 14:38:00
- **Status**: Resolved

### 実施した修正内容
- `backend/model/task.go` の `CreateTaskRequest.Validate()` 内で、コメントの前後の空白文字トリム（`strings.TrimSpace`）および改行コード正規化（`\r\n`, `\r` → `\n`）を行うように実装しました。
- 正規化・トリム後の文字列で 0〜1000文字以内の制約を検証し、未入力や空白のみの場合は空文字 `""` に正規化して永続化・返却されるようにしました。
- `backend/model/task_test.go` および `backend/service/task_test.go` にコメントのトリムと改行コード正規化の単体テストを追加しました。

### 変更したファイル
- [task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/model/task_test.go)
- [task_test.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/service/task_test.go)
