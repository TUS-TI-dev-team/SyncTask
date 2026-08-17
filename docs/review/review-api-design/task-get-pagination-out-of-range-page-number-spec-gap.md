# GET tasks における総ページ数超過ページ番号（page > total_pages）指定時のレスポンス仕様欠落

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` (3.3.1 節) において、該当タスクが存在する状態（`total_count > 0`）で、総ページ数を超えるページ番号（`page > total_pages`、例: 全45件・limit 20で `total_pages: 3` の際に `page: 5` を指定）がリクエストされた場合のレスポンス挙動（200 OK + 空配列 `items: []` なのか 400 Bad Request なのか）が明記されていません。

## 2. 詳細な指摘内容
1. **境界値リクエスト時のレスポンス未定義**:
   - `04_tasks.md` L12 では `page < 1` の場合に `400 Bad Request` を返却する旨が記載されており、L78 では検索結果が0件（`total_count: 0`）の場合の `pagination` オブジェクトの構造（`page: 1, limit: limit, total_count: 0, total_pages: 0`）が規定されています。
   - しかし、検索条件に該当するタスクが存在し `total_pages = 3` である場合において、リクエストパラメータとして `page = 5` など `total_pages` を超えるページ番号が送信された場合の振る舞いが未定義です。
2. **フロントエンド・バックエンド間の認識齟齬リスク**:
   - REST API の標準的な振る舞いとして `200 OK` と空配列 `items: []` および `pagination` オブジェクト (`{"page": 5, "limit": 20, "total_count": 45, "total_pages": 3}`) を返却するのか、あるいはパラメータ不正として `400 Bad Request` を返却するのかについて明示がないため、フロントエンドの画面描画制御やページ切替ロジックで予期せぬエラーが発生する懸念があります。

## 3. 推奨される修正案
`04_tasks.md` 3.3.1 節 (`GET tasks`) の Response 注記に、総ページ数を超えるページ番号（`page > total_pages`）が指定された場合の挙動として以下の仕様を明記してください：

```markdown
※ `total_count > 0` の状態で総ページ数を超えるページ番号（`page > total_pages`）が指定された場合、エラーとはせず `200 OK` レスポンスとして空配列 `items: []` と指定された `page` 番号を含む `pagination` オブジェクト（`{"page": <指定page>, "limit": <指定limit>, "total_count": <該当総件数>, "total_pages": <全総ページ数>}`）を返却します。
```
