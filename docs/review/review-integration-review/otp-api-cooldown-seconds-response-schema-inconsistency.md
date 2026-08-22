# OTP発行・再送APIレスポンススキーマにおけるcooldown_secondsの不整合

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)

## 1. 問題の概要
認証API設計書（`02_auth.md`）において、メールアドレス変更用API（`POST auth/change-email/request-otp` および `POST auth/change-email/resend-otp`）のレスポンスボディには `cooldown_seconds: 60` フィールドが定義されているのに対し、新規登録およびパスワードリセットの対応エンドポイント（`auth/register/request-otp`、`auth/register/resend-otp`、`auth/password-reset/request-otp`、`auth/password-reset/resend-otp`）のレスポンスには `cooldown_seconds` フィールドが含まれていません。
画面設計書上、いずれの画面も全く同様に「再送信ボタン（60秒クールダウン・カウントダウン表示あり）」を実装することになっているため、レスポンススキーマの設計が一貫していません。

## 2. 詳細な指摘内容

1. **APIレスポンススキーマの比較**:
   - `auth/change-email/request-otp` (3.1.10) および `resend-otp` (3.1.12):
     ```json
     {
       "otp_session_id": "otp_sess_chg_998877",
       "masked_email": "new_**********@example.com",
       "expires_in_seconds": 300,
       "cooldown_seconds": 60
     }
     ```
   - `auth/register/request-otp` (3.1.1) および `auth/password-reset/request-otp` (3.1.6):
     ```json
     {
       "otp_session_id": "otp_sess_a1b2c3d4e5",
       "masked_email": "user**********@example.com",
       "expires_in_seconds": 300
     }
     ```
   - `auth/register/resend-otp` (3.1.3) および `auth/password-reset/resend-otp` (3.1.8):
     ```json
     {
       "message": "OTP has been resent successfully.",
       "masked_email": "user**********@example.com",
       "expires_in_seconds": 300
     }
     ```

2. **画面要件との不整合**:
   - 画面設計書（`docs/design/screen_design.md`）では、新規登録OTP画面（行7）、パスワードリセットOTP画面（行9）、アカウント関連OTP画面（行62）のいずれも共通して「`「再送信」ボタン（60秒クールダウン・カウントダウン表示あり）`」と定義されています。
   - レスポンスからクールダウン秒数を動的に取得してカウントダウンタイマーを制御する場合、メール変更API以外では `cooldown_seconds` が取得できず、フロントエンド実装で用途ごとの分岐が必要になります。

## 3. 推奨される修正案

- OTP発行・再送レスポンスのスキーマ方針を統一する：
  - **方針A (推奨)**: 全用途の `request-otp` および `resend-otp` のレスポンスに `cooldown_seconds`（integer, デフォルト: 60）を含めるよう `02_auth.md` および関連する処理設計書（`01_account_creation.md`、`06_password_reset.md`）のレスポンス定義を統一する。
  - **方針B**: クライアント側で固定60秒タイマーを管理する方針とし、`auth/change-email/` 側のレスポンスから `cooldown_seconds` を削除して全エンドポイントを同一構造に揃える。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
新規登録、パスワードリセット、メールアドレス変更の全 OTP 発行（`request-otp`）および再送（`resend-otp`）API レスポンススキーマに `cooldown_seconds: 60` を共通で追加・統一しました。

### 変更したファイル
- [02_auth.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\02_auth.md)
- [01_account_creation.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\01_account_creation.md)
- [06_password_reset.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\06_password_reset.md)
