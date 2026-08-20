# ユーザー管理APIにおけるCookie無効化・削除ヘッダーのセキュリティ属性（HttpOnly, Secure, SameSite）の記載漏れ

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-17 18:00:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`03_users.md` の `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) のレスポンス仕様（200 OK および 401 `SESSION_DESTROYED`）において、セッションCookie（`sync_task_sid`）および CSRF トークンCookie（`XSRF-TOKEN`）を削除・無効化する `Set-Cookie` ヘッダーからセキュリティ属性（`HttpOnly`, `Secure`, `SameSite=Lax`）の記述が漏れている。

## 2. 詳細な指摘内容
1. **Cookie発行時と削除時における属性不一致**:
   - `01_overview.md` 1.1 節および `02_auth.md` にて規定されている通り、セッションCookie（`sync_task_sid`）は `HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000` 属性で発行され、CSRF トークンCookie（`XSRF-TOKEN`）は `Secure; SameSite=Lax; Path=/; Max-Age=2592000` 属性で発行される。
   - `01_overview.md` 1.1 節（L24-L25）では、セッション削除・Cookie消去時のレスポンスヘッダーとして以下の完全な属性が明記されている：
     - `Set-Cookie: sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0`
     - `Set-Cookie: XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`
   - しかし、`03_users.md` の 3.2.3 (`DELETE`) および 3.2.4 (`PATCH .../password`) の Response (200 OK) および 401 (`SESSION_DESTROYED`) エラーにおける `Set-Cookie` の記述は以下のようになっており、`HttpOnly`, `Secure`, `SameSite=Lax` 属性が省略されている：
     - `Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`
     - `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0`

2. **ブラウザ側でのCookie削除失敗リスク**:
   - RFC 6265 および主要なモダンブラウザ（Chrome, Firefox, Safari等）のセキュリティポリシーでは、`HttpOnly` や `Secure`, `SameSite` 属性が付与されて保存された Cookie を `Max-Age=0` で破棄する際、元々の Cookie 属性と一致しない `Set-Cookie` ヘッダー（例: `HttpOnly` や `Secure` が欠落したヘッダー）を受信した場合、削除要求を拒否または無視することがある。
   - その結果、アカウント削除後やパスワード変更後にもかかわらず、クライアントブラウザ上に旧セッションCookieやCSRFトークンCookieが削除されずに残留し、セッション管理上の脆弱性やフロントエンドの動作不整合を引き起こす危険性がある。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 (`DELETE users/{user_id}`) および 3.2.4 (`PATCH users/{user_id}/password`) の `Response (200 OK)` および `Errors (401 SESSION_DESTROYED)` における `Set-Cookie` ヘッダーの記述を `01_overview.md` 1.1 節の規定通り完全なセキュリティ属性を含んだ記述へ修正してください。

```markdown
- **Set-Cookie**: `sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`
```
