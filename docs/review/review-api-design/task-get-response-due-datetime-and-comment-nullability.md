# GET tasksレスポンスにおける due_datetime および comment の null 許容性・初期値定義の不足

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` および `GET tasks/{task_id}` のレスポンススキーマにおいて、締切日時が設定されていないタスクにおける `due_datetime` の表現型（`null` 返却）、およびコメントが未入力である場合の `comment` の表現型（空文字 `""` または `null`）の定義が明確に記載されていません。

## 2. 詳細な指摘内容
1. **`due_datetime` の null 許容性の型定義不足**:
   - レスポンス例 (L27-40, L142-157) では `"due_datetime": "2026-08-20T23:59:00+09:00"` と値が存在するパターンのみが掲載されていますが、締切日時が未設定（NULL）のタスクの場合に JSON レスポンス上で `"due_datetime": null` として返却される旨の明記がありません。

2. **`comment` の未入力時レスポンス表現の未定義**:
   - タスク作成・更新時にコメントが入力されなかった場合、DBの `COMMENT` カラムは `NULL` または空文字となりますが、APIレスポンスの JSON オブジェクトにおいて `"comment": ""` なのか `"comment": null` なのかの統一ルールが記載されていません。
   - クライアント側（TypeScript等）で型安全なデータ層を実装する際、`null` / `undefined` / 空文字のハンドリングに揺れが生じるリスクがあります。

## 3. 推奨される修正案
1. `GET tasks` および `GET tasks/{task_id}` のレスポンススキーマ定義に以下を明確に追記してください:
   - `due_datetime`: 型 `string / null`（締切日時未設定時は `null` を返却）。
   - `comment`: 型 `string`（コメント未入力・クリア時は空文字 `""` を返却、または `string / null` と表現を統一する規約を明記）。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`GET tasks` および `GET tasks/{task_id}` のレスポンス仕様注記に、締切日時未設定時は `"due_datetime": null` となる型 `string / null`、およびコメント未入力時は `"comment": ""` となる型 `string` を明確に規定しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
