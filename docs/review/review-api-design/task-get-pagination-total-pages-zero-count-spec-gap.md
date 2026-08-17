# GET tasks の通常一覧表示において検索結果0件時の pagination.total_pages 返却値の未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` (3.3.1 節) の通常一覧表示において、指定された検索条件・フィルタにより該当するタスクが 0件であった場合（`total_count: 0`）、レスポンスの `pagination.total_pages` が `0` となるのか `1` となるのかの仕様定義が欠落しています。

## 2. 詳細な指摘内容
1. **通常一覧とカレンダー一覧での仕様記述の不均衡**:
   - `04_tasks.md` L61 のカレンダー期間取得（`start_date` / `end_date` 指定）時の注記では、「該当タスクが0件の場合は `limit: 0, total_count: 0, total_pages: 1`」と明記されています。
   - しかし、通常一覧表示（`GET tasks` ページネーション時）において、検索条件（`keyword` や `status` 等）に合致するタスクが存在しない場合（`total_count: 0`）の `pagination` オブジェクトの返却値（`page: 1, limit: 20, total_count: 0, total_pages: 0` なのか `total_pages: 1` なのか）について明記されていません。
2. **フロントエンドのページネーションコンポーネント制御への影響**:
   - `total_pages` が `0` か `1` かによって、フロントエンドの「1 / 1 ページ」等のページ表示や前へ/次へボタンの非活性化制御に表記ブレやバグが生じる懸念があります。

## 3. 推奨される修正案
`GET tasks` 3.3.1 節の `pagination` オブジェクト仕様注記に、「通常一覧取得において該当タスクが 0件（`total_count: 0`）の場合、`pagination` は `{"page": 1, "limit": <指定limit値>, "total_count": 0, "total_pages": 0}`（または `total_pages: 1` で統一）として返却する」旨を明確に規定してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:18:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 節 (`GET tasks`) の Response 注記に、通常一覧取得において該当タスクが 0件（`total_count: 0`）の場合の `pagination` オブジェクトが `{"page": 1, "limit": <指定limit値>, "total_count": 0, "total_pages": 0}` として返却される旨を明確化しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)

