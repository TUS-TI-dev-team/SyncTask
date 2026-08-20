# パスワードリセットAPIにおけるOTP検証フローとリクエストボディ記述の矛盾

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:30:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/password-reset/verify-otp`（3.1.7）でOTP検証を完了させてステータスを `verified` に遷移させる設計であるにもかかわらず、後続の `POST auth/password-reset/reset`（3.1.9）のリクエストボディ例に `otp` パラメータが記載されており、検証フローの二重性およびリクエスト構造に不整合が生じています。

## 2. 詳細な指摘内容
1. **OTP検証とパスワードリセットの二重定義・構造矛盾**:
   - `3.1.7 POST auth/password-reset/verify-otp` では `otp_session_id` と `otp` を送信して検証に成功すると `200 OK` が返却され、DB上の `OTP_SESSION.STATUS` が `verified` へ更新されます（`database_design.md` L86, L90）。
   - しかし `3.1.9 POST auth/password-reset/reset` の Request Body には `otp_session_id`, `otp`, `new_password` の3つが定義されています。
   - ステップ3.1.7で既にOTP検証が完了している場合、3.1.9で再度 `otp` を送信して再検証を求めるのは冗長であり、既に試行・消費されたOTPハッシュとの照合でエラーとなるか、3.1.7の存在意義が消失します。3.1.7を個別に呼び出す画面フロー（OTP検証画面 → 新パスワード入力画面）を採用する場合、3.1.9では `otp_session_id`（ステータス `verified` のもの）と `new_password` のみを要求すべきです。
2. **既存ログインセッション無効化仕様の欠落**:
   - ログイン状態でのパスワード変更（`3.2.4 PATCH users/{user_id}/password`）では、パスワード更新成功時に対象ユーザーの全ログインセッション（`LOGIN_SESSION`）を一括物理削除するセキュリティ仕様が明記されています。
   - しかし `3.1.9 POST auth/password-reset/reset` では、パスワードリセット成功時に第三者の不正セッションを破棄するための「既存ログインセッション一括物理削除」およびCookie消去に関する記述が欠落しています。

## 3. 推奨される修正案
1. `3.1.9 POST auth/password-reset/reset` の Request Body から冗長な `otp` フィールドを削除し、`otp_session_id` と `new_password` のみに整理してください。
2. `3.1.9` の処理仕様として、「`OTP_SESSION` のステータスが `verified` かつ有効期限内であることを確認し、パスワード更新後に当該OTPセッションおよび対象ユーザーの全ログインセッション（`LOGIN_SESSION`）を物理削除する」旨を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST auth/password-reset/reset` (3.1.9) の Request Body から冗長な `otp` フィールドを削除し、`otp_session_id` と `new_password` のみに整理しました。また、パスワードリセット成功時に `OTP_SESSION` および該当ユーザーの全 `LOGIN_SESSION` を物理削除し、Cookie を破棄（再ログイン要求）する仕様を明記しました。

### 変更したファイル
- [02_auth.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
