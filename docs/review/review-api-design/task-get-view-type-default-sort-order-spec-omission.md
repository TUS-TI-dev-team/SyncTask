# GET tasks における各 view_type 指定時のデフォルトソート順序の定義欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` API において、`view_type` パラメータ（`high_priority`, `near_deadline`, `pinned`）を指定し、`sort_by` パラメータを省略した場合に適用される「ビューごとのデフォルトソート順序」の具体例・定義が文章中に欠落しています。

## 2. 詳細な指摘内容
`04_tasks.md` L22 には以下の注記があります：
> なお `view_type` 指定時に `sort_by` が明示的に指定された場合は、指定された `sort_by` のソート条件が最優先で適用されます（`sort_by` 省略時は `view_type` ごとのデフォルトソートが適用されます）。

しかし、`view_type`（`high_priority`, `near_deadline`, `pinned`）ごとに具体的にどのようなデフォルトソート条件が適用されるかが `04_tasks.md` 内で明記されていません。

要件定義書（`docs/req-def/requirements.md` L103-L116）における定義：
1. **優先タスク (`high_priority`)**: ピン留め優先（`is_pinned DESC`） → 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）
2. **締切間近タスク (`near_deadline`)**: 締切日時昇順（`due_datetime ASC`） → 作成日時降順（`created_at DESC`）
3. **ピン留めタスク (`pinned`)**: 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）

これが明確に記述されていない場合、バックエンド開発者が独自のソート順で実装してしまい、UI画面で期待される順序でタスクが表示されない原因となります。

## 3. 推奨される修正案
`04_tasks.md` L22 または `view_type` パラメータ説明（L14）付近に、以下の注記を追記してください。

```markdown
※ `sort_by` 省略時に `view_type` ごとに適用されるデフォルトソート順序は以下の通りです：
- `high_priority`: ピン留め優先（`is_pinned DESC`） → 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）
- `near_deadline`: 締切日時昇順（`due_datetime ASC`） → 作成日時降順（`created_at DESC`）
- `pinned`: 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）
```
