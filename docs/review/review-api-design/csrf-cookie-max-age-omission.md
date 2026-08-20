# CSRFトークンCookieのMax-Age属性未指定による永続セッション時のトークン消失問題

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1.1 節においてログインセッションCookie（`sync_task_sid`）は有効期限 1ヶ月（`Max-Age=2592000`）の永続Cookieとして規定されているのに対し、1.2 節の CSRF トークンCookie（`XSRF-TOKEN`）の発行仕様 `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/` には `Max-Age` 属性が指定されていない。これにより `XSRF-TOKEN` はブラウザ終了時に削除されるセッションCookieとなり、ブラウザ再起動後にログインセッションが有効であるにもかかわらず CSRF トークンが消失して状態変更 API（`POST`, `PUT`, `PATCH`, `DELETE`）が全滅する問題が生じる。

## 2. 詳細な指摘内容
1. **Cookie有効期限（Max-Age）の不一致**:
   - **セッションCookie（`sync_task_sid`）**: 1.1 節および 3.1.4 節において `Max-Age=2592000`（30日間）が指定され、ブラウザを閉じても維持される。
   - **CSRFトークンCookie（`XSRF-TOKEN`）**: 1.2 節（L27）および 3.1.2/3.1.4 節において `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/` と規定されており、`Max-Age`（または `Expires`）が未指定であるため、ブラウザ終了時に破棄されるセッションCookieとなる。

2. **ブラウザ再起動後のデッドロック現象**:
   ユーザーがブラウザを閉じて再起動した後、`sync_task_sid` Cookie によりログイン状態は維持されるが、`XSRF-TOKEN` Cookie は破棄されている。この状態で JavaScript が Cookie から `XSRF-TOKEN` を取得しようとしても存在しないため、`X-CSRF-Token` ヘッダーを付与できず、タスク作成・更新・削除やプロフ変更等の状態変更 API リクエストがすべて `403 Forbidden` (`FORBIDDEN`) エラーとなる。
   また、専用の CSRF トークン取得 API や GET リクエスト時の自動再発行仕様が未定義のため、ユーザーはログアウト・再ログインを行わない限り状態変更操作が不可能になる。

## 3. 推奨される修正案
`01_overview.md` の 1.2 節（L27）の `Set-Cookie` 仕様に `Max-Age=2592000` を追加し、セッションCookie（`sync_task_sid`）と同一の保持期間を設定してください。

```markdown
- CSRFトークンは **Double Submit Cookie 方式** にて管理します。ログイン成功（`auth/login`）およびアカウント新規登録完了（`auth/register/verify-otp`）時に、レスポンスヘッダーで `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000`（JavaScriptから読み取り可能な `HttpOnly=false`）を発行します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
CSRFトークンCookie（`XSRF-TOKEN`）の発行仕様に `Max-Age=2592000`（30日間）を追加し、ログインセッションCookie（`sync_task_sid`）と有効期限を統一して、ブラウザ再起動後のデッドロック現象を防止しました。

### 変更したファイル
- [01_overview.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/01_overview.md)
