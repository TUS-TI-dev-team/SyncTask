# MAIL_AUTH_LOGにおけるOTPセッションキャンセル時イベント種別の未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)

## 1. 問題の概要
`POST auth/otp-session/cancel` によるユーザー操作でのOTPセッション破棄時について、データベース設計書（`database_design.md`）のメール認証ログテーブル `MAIL_AUTH_LOG.EVENT_TYPE` カラムの定義にキャンセル時のイベント値（例: `CANCELLED`）が含まれておらず、キャンセル時の監査ログ記録仕様が未定義となっています。

## 2. 詳細な指摘内容

1. **DB設計における `EVENT_TYPE` のEnum値**:
   - `docs/design/database_design.md` 6.3節 `MAIL_AUTH_LOG`:
     > `EVENT_TYPE`: `ISSUED` (発行), `VERIFY_SUCCESS` (検証成功), `VERIFY_FAILED` (検証失敗), `RESEND_REQUESTED` (手動再送), `AUTO_RESEND` (5回失敗時自動処理), `EXPIRED` (有効期限切れ)
   - キャンセル操作を表す `CANCELLED` などの値が定義されていません。

2. **API設計および処理設計におけるログ仕様の欠落**:
   - `docs/design/api_design/02_auth.md` 3.1.13節（`POST auth/otp-session/cancel`）および各処理設計書において、キャンセル実行時に `MAIL_AUTH_LOG` や `ACCESS_LOG` へどのように記録を行うのか（あるいは `MAIL_AUTH_LOG` には記録せず `ACCESS_LOG` のみ記録するのか）が明記されていません。

## 3. 推奨される修正案

- `docs/design/database_design.md` の `MAIL_AUTH_LOG.EVENT_TYPE` の定義に `CANCELLED`（ユーザーキャンセル）を追加し、`02_auth.md` 3.1.13節および関連する処理設計書（`01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md`）のキャンセル処理詳細にログ記録仕様を明記する。
- または、セッションキャンセル時は `ACCESS_LOG` のみへの記録とする方針を明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
`database_design.md` の `MAIL_AUTH_LOG.EVENT_TYPE` 定義に `CANCELLED`（ユーザー操作によるキャンセル）を追加し、各認証APIおよび処理設計書にキャンセル時の監査ログ記録仕様を明記しました。

### 変更したファイル
- [database_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\database_design.md)
- [02_auth.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\02_auth.md)
- [README.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\README.md)
- [01_account_creation.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\01_account_creation.md)
- [02_account_edit.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\02_account_edit.md)
- [06_password_reset.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\06_password_reset.md)
