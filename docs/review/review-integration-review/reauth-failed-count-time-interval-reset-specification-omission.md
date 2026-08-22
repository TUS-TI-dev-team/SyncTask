# 再認証失敗カウンターにおける時間インターバル経過リセット仕様の未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md)
  - [03_account_delete.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/03_account_delete.md)
  - [07_password_change.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/07_password_change.md)
  - [01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md)

## 1. 問題の概要
データベース設計書（`database_design.md`）には、`LOGIN_ACCOUNT` テーブルに再認証失敗日時を記録するカラム `REAUTH_LAST_FAILED_AT` が定義されており、各処理設計書でも失敗時に更新、成功・強制破棄時にNULL初期化する処理が記載されています。
しかし、通常ログイン失敗カウンター（`LOGIN_FAILED_COUNT`）では「15分のインターバルを挟まずに5回連続失敗でロック / 最後の失敗から15分経過で0にリセット」という時間経過による失効判定が明記されているのに対し、再認証失敗カウンター（`REAUTH_FAILED_COUNT`）については「前回の失敗から何分経過でリセットされるのか（または時間経過リセットは行わないのか）」という評価ロジックがどのドキュメントにも定義されていません。

## 2. 詳細な指摘内容

1. **時間経過によるリセット判定ロジックの欠落**:
   - `docs/design/database_design.md` の第1節では、`LOGIN_FAILED_COUNT` の備考に「`15分間のインターバルを挟まずに5回連続失敗で30分間ロック / 最後の失敗から15分経過またはログイン成功時に0にリセット`」とインターバルが明記されています。
   - 一方、`REAUTH_FAILED_COUNT` の備考は「`パスワード変更・アカウント削除時の再認証失敗回数。5回連続失敗でログインセッション物理削除（成功時、パスワード変更時、ログアウト時、5回失敗セッション破棄時に0リセット）`」となっており、`REAUTH_LAST_FAILED_AT` からの経過時間判定が書かれていません。
   - `docs/design/process_design/07_password_change.md`（7.4.2節）および `03_account_delete.md`（3.4節）のバックエンド処理詳細においても、再認証失敗時に `REAUTH_FAILED_COUNT += 1`、`REAUTH_LAST_FAILED_AT = NOW()` と記録する手順はありますが、照合前に `REAUTH_LAST_FAILED_AT` を確認して一定時間（例: 15分）経過している場合にカウンターを0に戻す判定ステップが存在しません。

2. **実装・運用の曖昧さ**:
   - このままでは、数日前に1回パスワード再認証に失敗したユーザーが、数日後にログインセッション内で4回間違えただけで即座にセッション強制破棄（5回連続達成扱い）されてしまう恐れがあり、「連続失敗」の定義が時間軸で曖昧になっています。

## 3. 推奨される修正案

- 再認証失敗（`REAUTH_FAILED_COUNT`）のインターバル仕様を決定し、ドキュメント間で整合させる：
  - 例: `LOGIN_FAILED_COUNT` と同様に「直前の失敗（`REAUTH_LAST_FAILED_AT`）から15分以上経過している場合は、照合前に `REAUTH_FAILED_COUNT` を 0 にリセットしてから加算判定を行う」方針を `database_design.md`、`03_users.md`、`03_account_delete.md`、`07_password_change.md` の評価順序・処理詳細に明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
再認証失敗カウンター（`REAUTH_FAILED_COUNT`）について、直前の失敗（`REAUTH_LAST_FAILED_AT`）から15分以上経過している場合は照合前にカウンターを 0 にリセットしてから加算判定を行う仕様を、DB設計書・API設計書・処理設計書に明記しました。

### 変更したファイル
- [database_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\database_design.md)
- [03_users.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\03_users.md)
- [03_account_delete.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\03_account_delete.md)
- [07_password_change.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\07_password_change.md)
