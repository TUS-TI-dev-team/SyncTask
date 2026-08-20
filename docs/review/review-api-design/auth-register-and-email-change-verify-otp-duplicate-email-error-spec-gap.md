# OTP検証本登録・メールアドレス変更確定時のDB一意制約競合エラー定義漏れ

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/register/verify-otp`（新規登録OTP検証）および `POST auth/change-email/verify-otp`（メールアドレス変更確定）において、OTP検証成功後にDBへアカウント登録またはメールアドレス更新を行う際、メールアドレスやユーザー名が他ユーザーと競合した場合のエラーハンドリングおよびエラーコード定義が欠落しています。

## 2. 詳細な指摘内容
1. `3.1.1 POST auth/register/request-otp` および `3.1.10 POST auth/change-email/request-otp` では、アカウント列挙（User Enumeration）防止のため、指定されたメールアドレスが既に他ユーザーに登録されていてもエラーとせず `200 OK`（ダミーセッション発行）を返却します。
2. しかし、`request-otp` 呼び出し時点では未登録であったメールアドレスやユーザー名が、`request-otp` と `verify-otp` の間のタイムラグ中に別のユーザーによって登録されるレースコンディションが発生する可能性があります（また、`request-otp` でユーザー名重複チェックを行わない設計の場合も発生します）。
3. この状態でユーザーが `verify-otp`（`3.1.2` または `3.1.11`）を呼び出し正しいOTPを入力した場合、DBのユーザー登録・更新処理（`LOGIN_ACCOUNT.EMAIL` や `USERS.USERNAME` の UNIQUE 制約）で例外が発生します。
4. しかし、`3.1.2` および `3.1.11` の Errors 節には、重複時のエラー（例: `422 Unprocessable Entity` `EMAIL_ALREADY_EXISTS` / `USERNAME_ALREADY_EXISTS` や `400 BAD_REQUEST`）が定義されておらず、DBエラーに起因する未ハンドルの `500 Internal Server Error` が発生する懸念があります。

## 3. 推奨される修正案
`3.1.2` および `3.1.11` のリクエスト評価順序および Errors 節に、本登録・更新実行時のDB一意制約違反（メールアドレス・ユーザー名の重複）に対するエラーハンドリングを追記し、適切なHTTPステータスコード（`422 Unprocessable Entity`）およびエラーコード（`EMAIL_ALREADY_EXISTS` / `USERNAME_ALREADY_EXISTS`）を規定してください。
