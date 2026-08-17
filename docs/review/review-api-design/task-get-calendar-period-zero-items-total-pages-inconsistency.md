# `GET tasks` カレンダー表示用期間取得（`start_date` / `end_date` 指定）における該当件数 0件時の `total_pages` 返却値不整合

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-17 18:05:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`04_tasks.md` の `GET tasks` (L78, L84) において、検索・絞り込み結果が0件（`total_count: 0`）の際の `pagination` オブジェクトの挙動について、通常一覧取得時とカレンダー期間取得時（`start_date` / `end_date` 指定時）の間で `total_pages` の値に不整合が存在します。

## 2. 詳細な指摘内容
1. **通常一覧取得時とカレンダー期間取得時の `total_pages` 不整合**:
   - `04_tasks.md` L78 では、通常一覧取得において該当タスクが 0件（`total_count: 0`）の場合、`pagination` オブジェクトは `{"page": 1, "limit": <指定limit値>, "total_count": 0, "total_pages": 0}` を返却すると明確に定義されています。
   - 一方、L84 では、`start_date` / `end_date` を指定したカレンダー期間取得において該当タスクが 0件の場合、`limit: 0, total_count: 0, total_pages: 1` と規定されており、データが 0件であるにもかかわらず `total_pages: 1` が返却される仕様となっています。

2. **フロントエンドの判定処理および数学的ページネーション理論との矛盾**:
   - データ件数が0件（`total_count: 0`）である場合、ページ数は 0 ページ（`total_pages: 0`）となるのが一般的かつ数学的計算規則に整合します。
   - `total_count: 0` に対し `total_pages: 1` を返却すると、フロントエンド側で `total_pages === 0` を判定条件としてデータ非存在メッセージやUIの制御を行う際に誤動作を引き起こす原因となります。

## 3. 推奨される修正案
`04_tasks.md` L84 の `start_date` / `end_date` 指定時のレスポンス補足説明を以下のように修正し、該当タスクが 0件の場合は通常一覧取得と同様に `total_pages: 0` を返却する仕様に統一してください。

```markdown
また、ページネーションが無効化されて期間内の全タスクが一括返却されるため、`pagination` オブジェクトは `page: 1`, `limit: total_count`（取得件数と同値）, `total_pages: 1` として返却されます（該当タスクが0件の場合は `limit: 0, total_count: 0, total_pages: 0`）。
```
