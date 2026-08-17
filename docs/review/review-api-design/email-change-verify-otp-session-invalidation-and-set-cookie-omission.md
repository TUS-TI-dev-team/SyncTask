# メールアドレス変更確定API (`POST auth/change-email/verify-otp`) におけるセッション破棄・Cookie消去ヘッダー記述の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`POST auth/change-email/verify-otp`（3.1.11）において、OTP検証成功によりメールアドレスの変更が確定した際、要件定義書で定められた「該当ユーザーの全ログインセッションの一括物理削除」および「Cookieの消去（再ログイン要求）」に対応するレスポンスヘッダー（`Set-Cookie`）と処理記述が欠落しています。また、当該エンドポイントにおけるエラー仕様（`##### Errors`）セクションが存在しません。

## 2. 詳細な指摘内容
1. **全ログインセッション物理削除および Cookie 消去仕様の欠落**:
   - 要件定義書（`requirements.md` L44）には「変更確定と同時に、該当ユーザーのすべての既存ログインセッション（他端末・他ブラウザを含む）を一括無効化（物理削除）し、新メールアドレスでの再度ログインを要求する」と明記されています。
   - しかし `02_auth.md` 3.1.11 節（L278-L299）のレスポンス仕様（Response 200 OK）には、`Set-Cookie: sync_task_sid=; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Max-Age=0` のヘッダー記述が存在せず、`LOGIN_SESSION` レコードの全削除処理についての記載も欠落しています（`3.1.5 logout`, `3.2.3 DELETE user`, `3.2.4 PATCH password` では明記済み）。
2. **エラー仕様 (`##### Errors`) セクションの完全欠落**:
   - `3.1.11` に `##### Errors` セクションが存在しないため、未ログイン（`401 Unauthorized`）、CSRF検証失敗（`403 Forbidden`）、OTP不一致（`400 Bad Request`）、有効期限切れ（`410 Gone`）、5回失敗時自動再送（`422 Unprocessable Entity` `"code": "OTP_REISSUED_DUE_TO_FAILURES"`）のエラー定義が不透明になっています。

## 3. 推奨される修正案
1. `3.1.11 POST auth/change-email/verify-otp` のレスポンス仕様（200 OK）に以下の Cookie 消去ヘッダーおよび処理注記を追記してください:
   ```markdown
   ##### Response (200 OK)
   - **Set-Cookie**: `sync_task_sid=; Max-Age=0`
   - **Set-Cookie**: `XSRF-TOKEN=; Max-Age=0`

   ```json
   {
     "message": "Email address has been updated successfully. Please log in again with your new email address."
   }
   ```
   ※検証成功後、DB上の当該ユーザーの全ログインセッション（`LOGIN_SESSION`）を物理削除し、Cookieを消去して新メールアドレスでの再ログインを要求します。
   ```
2. 欠落している `##### Errors` セクションを追加し、`400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `410 Gone`, `422 Unprocessable Entity` を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST auth/change-email/verify-otp` (3.1.11) のレスポンスヘッダーに `Set-Cookie: sync_task_sid=; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Max-Age=0` を追加し、メール変更確定時の全ログインセッション一括物理削除・再ログイン要求仕様および `##### Errors` セクションを追記しました。

### 変更したファイル
- [02_auth.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
