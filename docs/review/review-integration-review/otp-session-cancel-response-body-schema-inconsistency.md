# OTPセッションキャンセルAPIのレスポンスボディスキーマ不整合

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:15:00
- **Target Files**:
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)

## 1. 問題の概要

`POST auth/otp-session/cancel` APIの正常終了レスポンスボディについて、API設計書（`docs/design/api_design/02_auth.md`）では `{"message": "OTP session cancelled successfully."}` と定義されているのに対し、各処理設計書（`01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md`）では `{"status": "cancelled"}` と記載されており、ドキュメント間でレスポンスJSONスキーマの不整合が発生しています。また、`06_password_reset.md` 6.7節の `MAIL_AUTH_LOG` イベント一覧から `CANCELLED` が漏れています。

## 2. 詳細な指摘内容

1. **レスポンスボディスキーマの不一致**:
   - `docs/design/api_design/02_auth.md` 3.1.13節（673〜678行目）:
     ```json
     {
       "message": "OTP session cancelled successfully."
     }
     ```
     フィールド名: `message`（string）
   - `docs/design/process_design/01_account_creation.md` 1.4.2節（102行目）:
     `処理完了後は遅延なしで 200 OK（{"status": "cancelled"}）を返却する。`
   - `docs/design/process_design/02_account_edit.md` 2.5節 セッション破棄・キャンセル（127行目）:
     `処理完了後は遅延なしで 200 OK（{"status": "cancelled"}）を返却する。`
   - `docs/design/process_design/06_password_reset.md` 6.4.2節（77行目）:
     `処理完了後は遅延なしで 200 OK（{"status": "cancelled"}）を返却する。`

   処理設計書側で `{"status": "cancelled"}` と記載されているため、API実装時やフロントエンド連携時にレスポンスパースで不具合が生じるリスクがあります。

2. **`06_password_reset.md` 6.7節のイベント一覧における `CANCELLED` 記述漏れ**:
   - `06_password_reset.md` 6.4.2節では `MAIL_AUTH_LOG` に `EVENT_TYPE='CANCELLED'` を記録すると規定されていますが、6.7節（187行目）のイベント一覧列挙（`ISSUED / VERIFY_SUCCESS / VERIFY_FAILED / RESEND_REQUESTED / AUTO_RESEND / EXPIRED`）に `CANCELLED` が含まれていません。

## 3. 推奨される修正案

1. **レスポンスボディの統一**:
   - `docs/design/process_design/01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md` における記述を、API設計書（`02_auth.md` 3.1.13節）に合わせて `200 OK（{"message": "OTP session cancelled successfully."}）` に統一する。
2. **`06_password_reset.md` 6.7節のイベント列挙の更新**:
   - 187行目のイベント一覧に `CANCELLED` を追加する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:16:00
- **Status**: Resolved

### 実施した修正内容
1. `docs/design/process_design/01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md` における `POST auth/otp-session/cancel` 完了時のレスポンス記述を、API設計書（`02_auth.md`）の定義に合わせて `200 OK（{"message": "OTP session cancelled successfully."}）` に修正・統一しました。
2. `docs/design/process_design/06_password_reset.md` 6.7節の `MAIL_AUTH_LOG` イベント一覧に `CANCELLED` を追加しました。

### 変更したファイル
- [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
- [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
- [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
