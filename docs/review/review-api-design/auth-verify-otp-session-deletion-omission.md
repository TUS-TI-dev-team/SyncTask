# アカウント登録完了およびメールアドレス変更確定時における OTP_SESSION 物理削除仕様の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`02_auth.md` の `POST auth/register/verify-otp` (3.1.2) および `POST auth/change-email/verify-otp` (3.1.11) において、検証成功・手続き確定時の DB 処理説明として既存ログインセッション（`LOGIN_SESSION`）の物理削除や Cookie 発行/消去は記述されているものの、手続き完了に伴う使用済み `OTP_SESSION` レコード自身の直ちの物理削除（`DELETE`）仕様が記述から欠落しています。

## 2. 詳細な指摘内容
1. **データベース設計書および要件定義書との不整合**:
   - `database_design.md` L95 には「新パスワード設定完了時やアカウント作成確定時等に直ちにDBから物理削除されます」と明確に規定されています。
   - `requirements.md` L264 にも「有効期限切れまたは無効化・使用済みとなったOTPセッションレコードは、新パスワード設定完了時等に直ちに物理削除する」と記載されています。
   - `02_auth.md` 3.1.9 (`POST auth/password-reset/reset`) では「当該OTPセッション（`OTP_SESSION`）および該当ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）をDBから直ちに物理削除し...」と正しく明記されています。
   - しかし、3.1.2 (`register/verify-otp`) L69 および 3.1.11 (`change-email/verify-otp`) L362 では、`LOGIN_SESSION` の削除や Cookie 操作のみが言及されており、`OTP_SESSION` の即時物理削除についての記述が省略されています。

2. **運用・セキュリティ上の影響**:
   - `OTP_SESSION` が手続き完了時に即時削除されない場合、定期クリーンアップ Cron ジョブ（15分間隔）が実行されるまでの間、DB 上で使用済み session ID レコードが `active` や `verified` のまま残留します。
   - `database_design.md` L204 で定義されている部分一意インデックス `uq_otp_session_active_pending_email ON OTP_SESSION (PENDING_EMAIL) WHERE STATUS IN ('active', 'verified')` の制約により、登録やメール変更が完了した直後に同一メールアドレスで再リクエストや他手続きを行おうとした際、残留した `OTP_SESSION` が衝突の原因となります。

## 3. 推奨される修正案
`02_auth.md` の 3.1.2 (`POST auth/register/verify-otp`) および 3.1.11 (`POST auth/change-email/verify-otp`) の説明注記（`※...`）を更新し、検証成功時に当該 `OTP_SESSION` レコードを DB から直ちに物理削除する旨を明記してください。

- **3.1.2 (`register/verify-otp`) 修正案**:
  ```markdown
  ※検証成功後、自動ログイン処理を行います。本登録完了と同時に、使用済みの手続き用OTPセッション（`OTP_SESSION`）をDBから直ちに物理削除します。なおリクエスト時に既存のログインセッションCookie（`sync_task_sid`）が送信された場合は、複数アカウントへの同時重複ログインを防止するため、その旧セッション（`LOGIN_SESSION`）もDBから物理削除した上で新しいログインセッションを発行します。
  ```

- **3.1.11 (`change-email/verify-otp`) 修正案**:
  ```markdown
  ※指定された `otp_session_id` に紐づくユーザーID（`OTP_SESSION.USER_ID`）が現在認証中のログインユーザーIDと一致していることを検証します。検証成功後、アカウントのメールアドレスを更新し、旧メールアドレス宛てに変更完了通知メールを送信（非同期処理）します。同時に、使用済みの手続き用OTPセッション（`OTP_SESSION`）および当該ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）をDBから直ちに物理削除してCookieを消去し、新メールアドレスでの再ログインを要求します。
  ```
