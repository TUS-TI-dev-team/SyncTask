# view_type と明示的な sort_by パラメータ併用時のソート順序優先関係の未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` において、特定ビューを指定する `view_type`（`high_priority`, `near_deadline`, `pinned`）と、ソート条件を指定する `sort_by`（例: `due_date_asc`, `created_at_desc` 等）が同時にリクエストされた際の優先適用ルールが明記されていません。

## 2. 詳細な指摘内容
`04_tasks.md` (L14) では `view_type` パラメータが定義されており、`pinned` 等のビューでは要件定義書によりデフォルトの並び順ルール（ピン留め優先→締切昇順→作成日時降順）が規定されています。
一方、L22 では `sort_by` パラメータのデフォルト値が `default` であると記載されています。

クライアントが `view_type=pinned&sort_by=created_at_desc` または `view_type=high_priority&sort_by=due_date_asc` のようにビュー切り替えと明示的なソート指定を同時に行った場合に、`sort_by` の指定がビューの並び順を上書き（優先適用）するのか、それとも `view_type` 固有のソート順が固定適用されて `sort_by` が無視されるのかが仕様上未規定です。

## 3. 推奨される修正案
`04_tasks.md` (L22) の `sort_by` パラメータ説明欄に、以下のルールを明記してください。

> ※ `view_type` 指定時に `sort_by` が明示的に指定された場合は、指定された `sort_by` のソート条件が最優先で適用されます（`sort_by` パラメータ省略時は `view_type` ごとのデフォルトソートが適用されます）。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 (`GET tasks`) の `sort_by` パラメータ説明欄に、「`view_type` 指定時に `sort_by` が明示的に指定された場合は、指定された `sort_by` のソート条件が最優先で適用されます」という優先適用ルールを明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
