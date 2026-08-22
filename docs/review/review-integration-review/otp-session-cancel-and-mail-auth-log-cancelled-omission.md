# OTPセッションキャンセルAPIおよびMAIL_AUTH_LOG.EVENT_TYPEの設計書間不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 13:58:00
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)

## 1. 問題の概要

ユーザー操作（OTP入力画面やパスワード入力画面での「戻る」ボタン押下、画面離脱等）により進行中のOTPセッションを即時無効化するAPI `POST auth/otp-session/cancel`、およびその際のメール認証ログ（`MAIL_AUTH_LOG.EVENT_TYPE='CANCELLED'`）について、API設計書（`02_auth.md` 3.1.13）で定義されている一方で、データベース設計書（`database_design.md`）、画面設計書（`screen_design.md`）、および各処理設計書（`01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md`）への反映が不完全であり、ドキュメント間で不整合が発生しています。

## 2. 詳細な指摘内容

1. **`database_design.md` における `CANCELLED` イベント値の未記載**:
   - `02_auth.md` 3.1.13節（672行目）には「メール認証ログテーブル（`MAIL_AUTH_LOG`）に `EVENT_TYPE='CANCELLED'` を記録した上で、該当レコードを DB から直ちに物理削除（`DELETE FROM OTP_SESSION`）します」と規定されています。
   - しかし、`docs/design/database_design.md` 6.3節（179行目）の `MAIL_AUTH_LOG.EVENT_TYPE` カラム定義の備考には：
     `ISSUED (発行), VERIFY_SUCCESS (検証成功), VERIFY_FAILED (検証失敗), RESEND_REQUESTED (手動再送), AUTO_RESEND (5回失敗時自動処理), EXPIRED (有効期限切れ)`
     と記載されており、`CANCELLED`（ユーザーキャンセル）が追加されていません。

2. **`screen_design.md` におけるキャンセル時動作記述の不一致**:
   - 行62（アカウント関連/OTP入力画面）の遷移先欄には「`戻る・キャンセル：プロフィール編集画面（クライアント側 otp_session_id を破棄し、POST auth/otp-session/cancel を呼び出してサーバー側セッションを即時無効化）`」と明記されています。
   - しかし、行7（アカウント作成/OTP入力画面）、行9（パスワードリセット/OTP入力画面）、行10（新パスワード入力画面）の遷移先欄には「`戻る・キャンセル：...`」としか書かれておらず、クライアント側セッション破棄および `POST auth/otp-session/cancel` 呼び出しの規定が欠落しています。

3. **処理設計書（`01_account_creation.md` / `02_account_edit.md` / `06_password_reset.md`）におけるキャンセル処理の未記載**:
   - `01_account_creation.md`（1.1節 表および1.5節 シーケンス図）、`02_account_edit.md`（2.5節および2.6節 シーケンス図）、`06_password_reset.md`（6.1節 表および6.6節 シーケンス図）において、対象API一覧やフローに `POST auth/otp-session/cancel` の呼び出しおよびログ記録（`EVENT_TYPE='CANCELLED'`）に関する記述が欠落しています。

## 3. 推奨される修正案

1. **`database_design.md`**: 6.3節 `MAIL_AUTH_LOG.EVENT_TYPE` の備考に `CANCELLED (ユーザーキャンセル)` を追加する。
2. **`screen_design.md`**: 行7（アカウント作成/OTP入力画面）、行9（パスワードリセット/OTP入力画面）、行10（新パスワード入力画面）の戻る・キャンセル時動作について、行62と同様に「クライアント側 `otp_session_id` を破棄し、`POST auth/otp-session/cancel` を呼び出してサーバー側セッションを即時無効化」と統一して追記する。
3. **`01_account_creation.md`, `02_account_edit.md`, `06_password_reset.md`**: 対象API表に `POST auth/otp-session/cancel` を追加し、戻る・画面離脱時のキャンセル処理およびログ記録（`EVENT_TYPE='CANCELLED'`）の記述・シーケンス図への反映を行う。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` の `MAIL_AUTH_LOG.EVENT_TYPE` 定義に `CANCELLED` を追加しました。
- `docs/design/screen_design.md` の各OTP画面の戻る・キャンセル時動作に `POST auth/otp-session/cancel` 呼び出しを統一明記しました。
- `docs/design/process_design/01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md` の対象API表・本文・シーケンス図にキャンセルAPI呼び出しおよび `MAIL_AUTH_LOG(CANCELLED)` 記録・セッション物理削除仕様を反映しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
- [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
- [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
- [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
