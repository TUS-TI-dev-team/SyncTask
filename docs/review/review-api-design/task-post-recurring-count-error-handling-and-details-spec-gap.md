# POST tasks 繰り返し一括作成時の件数オーバー・0件エラーのメッセージおよびエラーコード仕様の未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` API で毎週繰り返し一括作成（`is_recurring: true`）を実行した際、計算される生成件数が0件または100件超となった場合に返却されるエラーのメッセージ文字列、レスポンス形式（`error.details`）についての明確な定義が不足しています。

## 2. 詳細な指摘内容
要件定義書（`docs/req-def/requirements.md` L88）では、繰り返しタスク一括作成時のエラーメッセージが以下のように規定されています：
- **生成件数 0件の場合**: 「指定された期間内に該当する曜日が存在しません」
- **生成件数 100件超の場合**: 「生成件数が上限（100件）を超えています」

しかし、`04_tasks.md` 3.3.2 の Errors セクション（L130-L134）には以下のように記載されています：
> - `400 Bad Request`: タイトル文字数違反、期間・曜日不整合、生成件数超過（0件または101件以上）等（code: `"BAD_REQUEST"`）  
> - `422 Unprocessable Entity`: ビジネスルール違反（code: `"UNPROCESSABLE_ENTITY"`）

共通エラー仕様（`01_overview.md`）では `error.details` フィールドにエラーの詳細を含める形式となっていますが、件数上限違反や 0件エラー時の `field` 名（例: `field: "recurring_rule"` または `field: "recurring_rule.days_of_week"`）および上記の特定エラーメッセージを提示する定義が漏れています。

また、`422 Unprocessable Entity` が並記されているため、件数オーバーが 400 なのか 422 なのか、どちらが返却されるのか実装者の解釈が分かれる可能性があります。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 の Errors セクション（L130 付近）を以下のように更新してください。

```markdown
##### Errors
- `400 Bad Request`: リクエスト形式またはバリデーション不正（code: `"BAD_REQUEST"`）
  - 生成件数0件時: `error.details: [{ "field": "recurring_rule", "message": "指定された期間内に該当する曜日が存在しません" }]`
  - 生成件数100件超過時: `error.details: [{ "field": "recurring_rule", "message": "生成件数が上限（100件）を超えています" }]`
  - タイトル文字数違反、日付範囲不整合（`start_date > end_date` または 1年超）、無効な曜日指定等
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
```
