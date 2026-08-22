# 処理設計書におけるOTPセッションキャンセルAPIの記載欠落および画面設計書との不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [README.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/README.md)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
  - [01_overview.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/01_overview.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)

## 1. 問題の概要
API共通仕様書（`01_overview.md`）および認証API設計書（`02_auth.md` 3.1.13）において、ユーザーの「戻る」操作や画面離脱時にサーバー側OTPセッションを即時物理削除するエンドポイント `POST auth/otp-session/cancel` が定義されています。
しかし、処理設計書（`process_design/` 配下の各ファイル）において本APIに関する記載やシーケンスが完全に欠落しています。
さらに、画面設計書（`screen_design.md`）ではメールアドレス変更時のOTP画面にのみ `/api/` プレフィックス付きでキャンセルAPI呼び出しが記載されている一方、新規登録やパスワードリセットのOTP画面ではキャンセルAPI呼び出しの記載が漏れており、ドキュメント間で設計の不整合が発生しています。

## 2. 詳細な指摘内容

1. **処理設計書における `POST auth/otp-session/cancel` の欠落**:
   - `docs/design/process_design/01_account_creation.md` の対象API表（1.1節）およびシーケンス図（1.5節）に、画面の「戻る」ボタン押下時や離脱時のキャンセル処理（`POST auth/otp-session/cancel` の呼び出しおよび `OTP_SESSION` 物理削除）が定義されていません。
   - `docs/design/process_design/02_account_edit.md` および `docs/design/process_design/06_password_reset.md` においても同様に、キャンセルAPIの対象API一覧、処理詳細、シーケンス図への記載がありません。
   - `docs/design/process_design/README.md` の設計上の共通原則にも、OTPキャンセル時のサーバー側・クライアント側後処理に関する記述がありません。

2. **画面設計書における記述の不統一・漏れ**:
   - `docs/design/screen_design.md` の行62（アカウント関連/OTP入力画面）では「`戻る・キャンセル：プロフィール編集画面（クライアント側 otp_session_id を破棄し、POST /api/auth/otp-session/cancel を呼び出してサーバー側セッションを即時無効化）`」と記載されていますが、APIパスに他と異なる `/api/` プレフィックスが付与されています。
   - 一方、行7（アカウント作成/OTP入力画面）および行9（パスワードリセット/OTP入力画面）の遷移先欄には「`戻る・キャンセル：...`」としか書かれておらず、クライアント側のセッション破棄や `POST auth/otp-session/cancel` の呼び出し仕様が記載されていません。また、行10（新パスワード入力画面）のキャンセル時のセッション破棄についても未記載です。

## 3. 推奨される修正案

1. **処理設計書の改訂**:
   - `docs/design/process_design/01_account_creation.md`、`02_account_edit.md`、`06_password_reset.md` の対象API一覧に `POST auth/otp-session/cancel` を追加し、ユーザーの「戻る」ボタン押下や画面離脱時に本APIを呼び出してサーバー側 `OTP_SESSION` を物理削除（ダミー時も区別せず200応答）し、クライアント側の `otp_session_id` を破棄するフローを明記・シーケンス図に反映する。
   - `docs/design/process_design/README.md` の「設計上の共通原則」に、ユーザー操作によるOTP中断・キャンセル時の処理方針を追加する。

2. **画面設計書の統一**:
   - `docs/design/screen_design.md` の行7（アカウント作成/OTP入力画面）、行9（パスワードリセット/OTP入力画面）、行10（新パスワード入力画面）、行62（アカウント関連/OTP入力画面）の遷移先・動作説明について、すべて「クライアント側 `otp_session_id` を破棄し、`POST auth/otp-session/cancel` を呼び出してサーバー側セッションを即時無効化」と統一的に明記する（プレフィックス表記も `POST auth/otp-session/cancel` に統一する）。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
処理設計書（`01_account_creation.md`, `02_account_edit.md`, `06_password_reset.md`, `README.md`）に対象API `POST auth/otp-session/cancel` およびキャンセル時フローを追加し、画面設計書の全OTP・パスワード入力画面のキャンセル記述を統一しました。

### 変更したファイル
- [README.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\README.md)
- [01_account_creation.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\01_account_creation.md)
- [02_account_edit.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\02_account_edit.md)
- [06_password_reset.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\06_password_reset.md)
- [screen_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\screen_design.md)
