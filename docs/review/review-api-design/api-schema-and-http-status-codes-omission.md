# API詳細スキーマおよびHTTPステータスコード定義の欠如

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
現在の `api_design.md` は表形式による日本語の概略記載にとどまっており、各エンドポイントで送受信される具体的な JSON フィールド名、データ型、必須/任意制約、バリデーションルール、および成功・エラー時の HTTP ステータスコードが定義されていません。このため、フロントエンド・バックエンドの実装者が迷わず開発を進められない状態です。

## 2. 詳細な指摘内容
1. **JSON スキーマの未定義**:
   - 各リクエスト/レスポンスのフィールド名（例: `task_name` vs `title` vs `name`、`due_date` のフォーマット等）や型が明記されていません。
   - 要件定義書で指定されている制約（タスク名 1〜100文字、コメント 0〜1000文字、優先度 `high` | `medium` | `low`、ステータス `not_started` | `in_progress` | `completed` 等）が API スキーマとして定義されていません。
2. **HTTP ステータスコードおよびエラーパターンの未定義**:
   - 成功時（`200 OK`, `201 Created`, `204 No Content` 等）および想定されるエラー時（`400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `409 Conflict`, `422 Unprocessable Entity`, `429 Too Many Requests`）のステータスコード対応が未定義です。
   - `requirements.md` L189 に記載の「認可エラー時の `404 Not Found` 返却（存在秘匿）」や、L195 の「ログインレートリミット時の `429 Too Many Requests`」などの非機能要件・セキュリティ要件が API 設計に反映されていません。

## 3. 推奨される修正案
1. 各エンドポイントについて、具体的な Request Body / Query Parameters / Response Body の JSON スキーマ（フィールド名・型・必須区分・制約）を定義してください。
2. エンドポイントごとに返却される主な HTTP ステータスコードと、共通エラーレスポンス構造におけるエラーコード（`code` 値: 例: `INVALID_CREDENTIALS`, `TASK_NOT_FOUND`, `RATE_LIMIT_EXCEEDED` 等）の対応一覧を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` において、全エンドポイント（認証、ユーザー、タスク）に対する JSON リクエスト/レスポンススキーマ、データ型、必須/任意制約、バリデーション制約を網羅的に定義しました。
- 共通エラーレスポンス構造および HTTP ステータスコード（200, 201, 400, 401, 403, 404, 409, 422, 429, 500）と対応するエラーコード一覧を明文化しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
