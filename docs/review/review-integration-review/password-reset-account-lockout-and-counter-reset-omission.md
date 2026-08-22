# パスワードリセット完了時におけるアカウントロック解除および各種失敗カウンター初期化の規定欠落

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-22 14:25:00
- **Target Files**:
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

## 1. 問題の概要
通常ログインのパスワード連続誤入力によってアカウントロックアウト（`LOGIN_LOCK_UNTIL` 設定、`LOGIN_FAILED_COUNT >= 5`）となったユーザーが、パスワードリセット手続き（`POST auth/password-reset/reset`）を正常に完了した場合、新パスワードでのログインが要求されます。
しかし、パスワードリセット完了処理（`06_password_reset.md` 6.5.2節、`02_auth.md` 3.1.9節）の確定トランザクションにおいて、`LOGIN_LOCK_UNTIL = NULL`、`LOGIN_FAILED_COUNT = 0`、`LOGIN_LAST_FAILED_AT = NULL`、および再認証失敗カウンター（`REAUTH_FAILED_COUNT = 0`、`REAUTH_LAST_FAILED_AT = NULL`）の初期化リセット処理が規定されていません。
このため、ロックアウト解除を目的としてパスワードリセットを完了しても、直後の新パスワードでのログイン時に `LOGIN_LOCK_UNTIL` によるロックアウト判定で弾かれてしまい（401 `UNAUTHORIZED`）、ログイン不能となる重大なデッドロック・仕様不整合が発生します。

## 2. 詳細な指摘内容

1. **パスワードリセット確定トランザクションにおけるロック解除・カウンターリセットの欠落**:
   - `docs/design/process_design/06_password_reset.md`（6.5.2節 更新トランザクション）では以下の4点のみが記載されています：
     1. `LOGIN_ACCOUNT.PASSWORD_HASH` と `UPDATED_AT` の更新
     2. 対象ユーザーの全 `LOGIN_SESSION` 物理削除
     3. 検証済み `OTP_SESSION` およびアクティブな `OTP_SESSION` の一括物理削除
     4. `ACCESS_LOG` 記録
   - `docs/design/api_design/02_auth.md`（3.1.9節 リクエスト評価順序ステップ4）でも同様にパスワード更新とセッション・OTP削除のみとなっています。
   - ここに `LOGIN_ACCOUNT` テーブルのログイン失敗管理カラム（`LOGIN_FAILED_COUNT = 0`, `LOGIN_LAST_FAILED_AT = NULL`, `LOGIN_LOCK_UNTIL = NULL`）および再認証失敗管理カラム（`REAUTH_FAILED_COUNT = 0`, `REAUTH_LAST_FAILED_AT = NULL`）を初期化する更新処理が含まれていません。

2. **データベース設計書の備考欄との不整合**:
   - `docs/design/database_design.md` の `LOGIN_ACCOUNT` テーブル定義において：
     - `LOGIN_FAILED_COUNT`: 「最後の失敗から15分経過またはログイン成功時に0にリセット」
     - `LOGIN_LOCK_UNTIL`: 「5回連続失敗時にロックアウト発生時刻から30分後を設定」
     - `REAUTH_FAILED_COUNT`: 「成功時、パスワード変更時、ログアウト時、5回失敗セッション破棄時に0リセット」
   - このように、いずれのカラムにおいても「パスワードリセット完了時（新パスワード設定成功時）」における初期化リセット契機の記載が漏れています。

## 3. 推奨される修正案

1. `docs/design/process_design/06_password_reset.md`（6.5.2節）および `docs/design/api_design/02_auth.md`（3.1.9節ステップ4）のパスワードリセット完了トランザクションに以下を明記する：
   - `LOGIN_ACCOUNT` の更新処理において、`PASSWORD_HASH` と `UPDATED_AT` に加え、`LOGIN_FAILED_COUNT = 0`, `LOGIN_LAST_FAILED_AT = NULL`, `LOGIN_LOCK_UNTIL = NULL`, `REAUTH_FAILED_COUNT = 0`, `REAUTH_LAST_FAILED_AT = NULL` を同時に初期化リセットする。
2. `docs/design/database_design.md` の `LOGIN_ACCOUNT` テーブル定義（`LOGIN_FAILED_COUNT`, `LOGIN_LOCK_UNTIL`, `REAUTH_FAILED_COUNT`）の備考欄に、パスワードリセット完了時にもこれらが 0 / NULL にリセット・解除される旨を明記する。
