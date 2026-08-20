# `DELETE` および `PATCH` API レスポンスにおける Cookie 削除ヘッダーの `Path=/` 属性の記載漏れ

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` の `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) の Response (200 OK) や Errors (401 SESSION_DESTROYED) の `Set-Cookie` 削除ヘッダーに `Path=/` 属性を明記し、Cookie 削除処理が RFC 6265 に準拠して正しくルートパスの Cookie を無効化するよう修正しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
- **Severity**: Medium
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`03_users.md` の `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) の Response (200 OK) および 401 (`SESSION_DESTROYED`) エラーにおける Cookie 削除ヘッダーの記述で、`Path=/` 属性が指定されていない。

## 2. 詳細な指摘内容
1. **Cookie 削除仕様（RFC 6265）との不整合**:
   - 本システムのログインセッションCookie（`sync_task_sid`）および CSRF トークンCookie（`XSRF-TOKEN`）は、発行時（`01_overview.md` 1.1 節, 1.2 節, `02_auth.md` 3.1.2/3.1.4）にすべて `Path=/` 属性が付与されてクライアントへ送信される。
   - RFC 6265 仕様上、発行時に `Path=/` が指定された Cookie を `Max-Age=0` で無効化・削除する場合、削除用レスポンスヘッダー（`Set-Cookie`）においても一致する `Path=/` 属性を明記しなければならない。
   - `03_users.md` の 3.2.3 および 3.2.4 では以下のように記述されている：
     - `- **Set-Cookie**: sync_task_sid=; Max-Age=0`
     - `- **Set-Cookie**: XSRF-TOKEN=; Max-Age=0`
   - `Path=/` 属性が省略された場合、ブラウザはリクエストパス（`/api/users/...`）にスコープされた Cookie の無効化と解釈するため、ルートパス (`Path=/`) に存在する元のセッションCookieが消去されず、ブラウザ上に無効なトークンが残留するリスクがある。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 (`DELETE users/{user_id}`) および 3.2.4 (`PATCH users/{user_id}/password`) の `Response (200 OK)` および `Errors (401 SESSION_DESTROYED)` における `Set-Cookie` ヘッダー記述に `Path=/` 属性を明記してください。

**修正案の例**:
```markdown
- **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`
```
