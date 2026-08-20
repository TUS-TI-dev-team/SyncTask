# 未認証POSTリクエストとCSRF検証対象範囲の記述矛盾

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1.2 セキュリティ仕様において「状態を変更するすべてのHTTPメソッド（`POST`, `PUT`, `PATCH`, `DELETE`）において CSRFトークンの検証を必須とします」と記述されているが、未ログイン状態で行われる `POST auth/login` や `POST auth/register/request-otp` などのリクエストでは、クライアントはまだ `XSRF-TOKEN` Cookie を保持していないため、すべての `POST` メソッドで CSRF 検証を必須とすると未ログインユーザーがログイン・新規登録を実行できなくなる矛盾が生じている。

## 2. 詳細な指摘内容
1. **CSRF検証の適用範囲の誤解**:
   `01_overview.md` L22-23 にて「状態を変更するすべてのHTTPメソッド（`POST`, `PUT`, `PATCH`, `DELETE`）において CSRFトークンの検証を必須とします。」と一括指定されている。
   しかし、L24 にて Double Submit Cookie 方式の CSRF トークン（`XSRF-TOKEN` Cookie）は **ログイン成功（`auth/login`）** および **新規登録完了（`auth/register/verify-otp`）** 時に初めてレスポンスヘッダーで発行されると定義されている。
   そのため、ログイン前または新規登録開始時点の未認証 `POST` リクエスト（`POST auth/login`, `POST auth/register/request-otp`, `POST auth/password-reset/*` 等）においてはクライアント側に CSRF トークンが存在しない。

2. **詳細設計書（`02_auth.md`）との記述不一致**:
   `02_auth.md` においては、`POST auth/logout` や `POST auth/change-email/*` などの認証必須エンドポイントにのみ `Headers: X-CSRF-Token: <token>` が記載されており、未認証の `POST auth/login` や `POST auth/register/request-otp` には同ヘッダーの記載がない。概要書 1.2 の「すべてのHTTPメソッド」という記述が詳細設計と食い違っている。

## 3. 推奨される修正案
`01_overview.md` の 1.2 セキュリティ仕様（L22-23）の記述を以下のように修正し、CSRF検証が「認証済みの状態変更リクエスト」を対象とすることを明記してください。

```markdown
- **CSRF対策**:
  - Cookieベースの認証を行うため、**認証を必要とするすべての状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）**において CSRFトークン（`X-CSRF-Token` ヘッダー）の検証を必須とします（未認証のログイン・会員登録リクエスト等を除く）。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.2 節の CSRF 検証仕様を「認証を必要とするすべての状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）において CSRFトークンの検証を必須とする（未認証のログイン・会員登録・パスワードリセット等のリクエストを除く）」に修正しました。

### 変更したファイル
- [01_overview.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/01_overview.md)
