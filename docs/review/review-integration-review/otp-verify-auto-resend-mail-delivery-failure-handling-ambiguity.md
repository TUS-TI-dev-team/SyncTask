# OTP検証5回連続失敗時の自動再送におけるメール送信失敗時エラーハンドリング仕様の未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 14:15:00
- **Target Files**:
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md#L111-L120)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md#L347-L356)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md#L570-L583)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md#L58)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md#L129)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md#L50)

## 1. 問題の概要
API設計書（`02_auth.md` 3.1.2, 3.1.7, 3.1.11）および各処理設計書（`01_account_creation.md`, `02_account_edit.md`, `06_password_reset.md`）において、OTP照合が5回連続で失敗した際に、サーバー側で新しい8桁OTPを自動再発行してメールを自動送信し、`422 Unprocessable Entity`（code: `OTP_REISSUED_DUE_TO_FAILURES`）を返却する仕様が定められています。
しかし、**この自動再送処理においてメールサーバー障害等によりメール送信に失敗した場合（1〜4回目の送信失敗時）**のAPIレスポンスステータスおよびエラーハンドリングが未定義となっています。

## 2. 詳細な指摘内容
1. **API設計および処理設計の現状記述**:
   - `02_auth.md` の各 `verify-otp` API リクエスト評価順序:
     - 「不一致（試行5回達成）: 失敗回数をリセットし、OTP自動再発行通知 `422 Unprocessable Entity`（code: `"OTP_REISSUED_DUE_TO_FAILURES"`、遅延 1.0s ± 0.1s）を返却します。」
   - `02_auth.md` 4節「OTP API共通規則」7-8行目:
     - 「実メール送信に失敗した場合は `DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT+=1` とし、`503 Service Unavailable`（code: `OTP_DELIVERY_FAILED`）の `error.otp_session_id` に再送用 `otp_session_id` を返します。」
     - 「初回送信、手動再送、5回照合失敗による自動再送を通じて5回連続で送信に失敗した場合、対象OTPセッションを物理削除し、`410 Gone`（code: `OTP_SESSION_INVALIDATED`）を返します。」
2. **問題点・曖昧さ**:
   - `verify-otp` 実行時に5回目の不一致によって自動再送が実行され、そのメール送信が失敗した場合：
     - パターンA: メール送信失敗を優先して `503 Service Unavailable`（code: `OTP_DELIVERY_FAILED`）を返却し、画面側で「新しい認証コードの送信に失敗しました。再送信ボタンから再試行してください」と案内するのか。
     - パターンB: 通常通り `422 Unprocessable Entity`（code: `OTP_REISSUED_DUE_TO_FAILURES`）を返却し、DB上の `DELIVERY_STATUS` を `sendable`（`SEND_FAILED_COUNT+=1`）にしておき、ユーザーがメールを受信できない場合に手動再送ボタンを押下させるのか。
   - パターンBの場合、ユーザーには「新しい認証コードを再送しました」と表示されるが実際にはメールが届いておらず、UX上の重大な混乱を招きます。また、パターンAの場合の `verify-otp` の Errors 定義に `503` が記載されていません。

## 3. 推奨される修正案
1. `02_auth.md` の各 `verify-otp` API（3.1.2, 3.1.7, 3.1.11）の評価順序および Errors 定義に、自動再送時のメール送信失敗ハンドリングを明記する：
   - 5回目不一致時の自動再送においてメール送信に失敗した場合（1〜4回目の送信失敗）、`DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT+=1` とし、`503 Service Unavailable`（code: `"OTP_DELIVERY_FAILED"`）を返却する。
   - 自動再送を含めて5回連続送信失敗となった場合は、対象セッションを物理削除し、`410 Gone`（code: `"OTP_SESSION_INVALIDATED"`）を返却する。
   - `verify-otp` の Errors 一覧に `503 Service Unavailable`（code: `"OTP_DELIVERY_FAILED"`）を追加する。
2. 処理設計書（`01_account_creation.md`, `02_account_edit.md`, `06_password_reset.md`）のシーケンスおよび評価順序にも同様の分岐を追記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design/02_auth.md` の各 `verify-otp` API（3.1.2, 3.1.7, 3.1.11）において、5回目不一致による自動再送メール送信失敗時の `503 Service Unavailable`（code: `OTP_DELIVERY_FAILED`）および5回連続失敗時 `410 Gone`（code: `OTP_SESSION_INVALIDATED`）のエラーハンドリングを評価順序・Errors一覧に明記しました。
- `docs/design/process_design/01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md` のシーケンス図および評価順序に同様のメール送信失敗ハンドリング分岐を追記しました。

### 変更したファイル
- [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
- [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
- [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
- [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
