# 処理設計書におけるOTP発行・再送レスポンス例のcooldown_seconds欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:58:00
- **Target Files**:
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)

## 1. 問題の概要

API設計書（`02_auth.md`）において、新規登録、パスワードリセット、メールアドレス変更の全 OTP 発行（`request-otp`）および再送（`resend-otp`）API のレスポンススキーマに `cooldown_seconds: 60` が共通フィールドとして追加・統一されました。

しかし、処理設計書（`01_account_creation.md` および `06_password_reset.md`）のレスポンスJSON例に `cooldown_seconds: 60` が反映されておらず、設計書間でスキーマ記述の不整合が生じています。

## 2. 詳細な指摘内容

1. **`docs/design/process_design/01_account_creation.md` 1.2.2節**:
   - 40-46行目のレスポンス例：
     ```json
     {
       "otp_session_id": "otp_sess_a1b2c3d4e5",
       "masked_email": "user**********@example.com",
       "expires_in_seconds": 300
     }
     ```
     → `cooldown_seconds: 60` が欠落しています（`02_auth.md` 3.1.1節では `cooldown_seconds: 60` を定義）。

2. **`docs/design/process_design/06_password_reset.md` 6.2.2節**:
   - 34-40行目のレスポンス例：
     ```json
     {
       "otp_session_id": "otp_sess_reset_12345",
       "masked_email": "user**********@example.com",
       "expires_in_seconds": 300
     }
     ```
     → `cooldown_seconds: 60` が欠落しています（`02_auth.md` 3.1.6節では `cooldown_seconds: 60` を定義）。

## 3. 推奨される修正案

- `docs/design/process_design/01_account_creation.md`（40-46行目）および `docs/design/process_design/06_password_reset.md`（34-40行目）のレスポンスJSON例に `"cooldown_seconds": 60` を追加し、`02_auth.md` のスキーマ定義と完全同期させる。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/process_design/01_account_creation.md` および `docs/design/process_design/06_password_reset.md` のレスポンスJSON例に `"cooldown_seconds": 60` を追記し、API設計書のレスポンススキーマと完全に整合させました。

### 変更したファイル
- [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
- [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
