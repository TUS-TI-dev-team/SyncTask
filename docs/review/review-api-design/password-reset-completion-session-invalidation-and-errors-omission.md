# パスワードリセット完了API (`POST auth/password-reset/reset`) におけるセッション物理削除・Cookie消去・エラー仕様の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`POST auth/password-reset/reset`（3.1.9）において、パスワードリセット完了時に要件定義書で義務付けられている「対象ユーザーの既存ログインセッション一括物理削除」および「Cookie/仮セッションの消去」に対応する仕様記述・レスポンスヘッダーが欠落しています。また、`3.1.9` にエラー仕様（`##### Errors`）セクションが存在せず、未検証・期限切れセッションやパスワード要件違反時の応答が不透明です。

## 2. 詳細な指摘内容
1. **全ログインセッション一括物理削除および Cookie 消去仕様の欠落**:
   - 要件定義書（`requirements.md` L75）には「新パスワードを設定完了すると同時に当該パスワードリセット用OTP/仮セッションを直ちに失効・DBから物理削除し（クライアント側の仮セッションCookieも消去）、該当ユーザーのすべての既存ログインセッション（他端末・他ブラウザを含む）を一括無効化（物理削除）して再度ログインを要求する」と規定されています。
   - しかし `02_auth.md` 3.1.9 節（L232-L253）では `Response (200 OK)` に `"message"` のみが返却され、`LOGIN_SESSION` および `OTP_SESSION` の物理削除処理についての記載や `Set-Cookie: sync_task_sid=; Max-Age=0` ヘッダーの定義が存在しません。
2. **エラー仕様 (`##### Errors`) セクションの完全欠落**:
   - `3.1.9` に `##### Errors` セクションが存在しないため、以下のエラー定義が欠落しています:
     - `400 Bad Request`: 新パスワード要件違反（文字数、文字種、ユーザー名/ローカル部包含違反）
     - `403 Forbidden`: `otp_session_id` のステータスが `VERIFIED` でない場合（OTP未検証状態でのアクセス）
     - `410 Gone`: OTP検証成功後の仮セッション有効期限切れ（検証後15分経過）

## 3. 推奨される修正案
1. `3.1.9 POST auth/password-reset/reset` のレスポンス仕様（200 OK）を以下のように更新し、物理削除および Cookie 消去ヘッダーを明記してください:
   ```markdown
   ##### Response (200 OK)
   - **Set-Cookie**: `sync_task_sid=; Max-Age=0`
   - **Set-Cookie**: `XSRF-TOKEN=; Max-Age=0`

   ```json
   {
     "message": "Password has been reset successfully. Please log in with your new password."
   }
   ```
   ※処理成功後、当該OTPセッション（`OTP_SESSION`）および該当ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）をDBから直ちに物理削除し、Cookieを消去して再ログインを要求します。
   ```
2. `3.1.9` に `##### Errors` セクションを追加し、`400 Bad Request`（入力検証違反）、`403 Forbidden`（未検証OTPセッション指定）、`410 Gone`（仮セッション期限切れ）を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST auth/password-reset/reset` (3.1.9) に `Set-Cookie: sync_task_sid=; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Max-Age=0` を追加し、パスワードリセット成功時のOTPセッションおよび全ログインセッションの物理削除仕様を明記するとともに、`##### Errors` セクション（400, 403, 410）を追加しました。

### 変更したファイル
- [02_auth.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
