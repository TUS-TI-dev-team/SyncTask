# ログイン処理設計書におけるアカウントロックアウト時429返却記述の残存

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:58:00
- **Target Files**:
  - [04_login.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/04_login.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)

## 1. 問題の概要

ユーザー列挙脆弱性対策として、API基本設計書（`01_overview.md`）および認証API設計書（`02_auth.md` 3.1.4 リクエスト評価順序）において「メールアドレス単位のアカウントロックアウト（`LOGIN_LOCK_UNTIL`）発生時はアカウントの登録有無を秘匿するため一貫して `401 Unauthorized`（code: `"UNAUTHORIZED"`, 遅延 1.0s ± 0.1s）を返却し、`429 Too Many Requests`（code: `"RATE_LIMIT_EXCEEDED"`）は IP アドレス単位レートリミット超過時のみに限定する」方針への修正が行われました。

しかし、ログイン処理設計書（`docs/design/process_design/04_login.md`）の本文・表・シーケンス図、および認証API設計書（`02_auth.md` 3.1.4）の Errors 一覧定義において、アカウントロックアウト時に `429 Too Many Requests` を返却する旧仕様の記述が残存しており、設計書間で重大な不整合が発生しています。

## 2. 詳細な指摘内容

1. **`docs/design/process_design/04_login.md` における残存箇所**:
   - **35行目 表**:
     `| 429 / RATE_LIMIT_EXCEEDED | メール単位ロック中、IP単位遮断中 | 1.0s ± 0.1s | 登録有無を示さない試行制限メッセージ |`
     → 「メール単位ロック中」が含まれたままとなっています。
   - **64-66行目 シーケンス図**:
     ```mermaid
     alt メール単位ロック中
         Backend->>DB: LOGIN_LOG、ACCESS_LOG記録
         Backend-->>Frontend: 429（1.0s ± 0.1s）
     ```
     → ロック中に 429 を返すフローとなっています。
   - **95行目 (4.4.1節 評価順序 ステップ5)**:
     `5. 有効アカウントの LOGIN_LOCK_UNTIL > NOW() なら、正しいパスワードでも 429。ロック中の要求で回数・期限を延長しない。`
     → 正しいパスワードでも 429 を返す旨が記述されています。
   - **107行目 (4.4.2節 認証失敗)**:
     `- 制限到達となった今回の照合失敗は 401。次回以降の制限中要求は 429 とする。`
     → メール単位ロックアウト中も次回以降 429 となる記述になっています。

2. **`docs/design/api_design/02_auth.md` 3.1.4節 Errors 定義**:
   - **234行目**:
     `- 429 Too Many Requests: メールアドレス単位ロックアウト（直近15分間に5回連続失敗で30分ロック）またはIPレートリミット超過（直近5分間に30回失敗で15分遮断）。遅延 1.0s ± 0.1s（code: "RATE_LIMIT_EXCEEDED"）`
     → 同ファイル 225行目のリクエスト評価順序ステップ2では「メールアドレスのアカウントロックアウト（`LOGIN_LOCK_UNTIL`）中である場合は...一貫して `401 Unauthorized`（code: `"UNAUTHORIZED"`）を返却」と修正されているのに対し、Errors の一覧表が旧仕様のまま残っています。

## 3. 推奨される修正案

1. `docs/design/process_design/04_login.md` を以下のとおり修正する：
   - 35行目の表: `429 / RATE_LIMIT_EXCEEDED` の条件を「IP単位遮断中（直近5分間に30回失敗で15分遮断）」のみとし、メール単位ロックアウト中は `401 / UNAUTHORIZED`（未登録・削除済み・パスワード不一致・アカウントロックアウト中）に含める。
   - 64-66行目のシーケンス図: メール単位ロック中であってもパスワード照合へ進むか、または内部的に認証拒否として `401（1.0s ± 0.1s）` を返すフローへ変更する。
   - 95行目（評価順序）: `LOGIN_LOCK_UNTIL > NOW()` の場合は、アカウント存在秘匿のため直ちにエラーとはせず、パスワード照合において正しいパスワードであっても認証拒否として `401 UNAUTHORIZED`（遅延 1.0s ± 0.1s）を返却するよう修正する。
   - 107行目: メール単位ロックアウト時は次回以降の要求も一貫して `401` を返却する旨を明記する。
2. `docs/design/api_design/02_auth.md` 234行目の Errors 定義から「メールアドレス単位ロックアウト」を削除し、IPレートリミット超過のみとする。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/process_design/04_login.md` において、メール単位ロックアウト時の返却ステータスを `401 UNAUTHORIZED`（1.0s ± 0.1s遅延）に統一し、表・シーケンス図・評価順序・失敗時処理を更新しました。
- `docs/design/api_design/02_auth.md` 3.1.4節 Errors 定義からメールアドレス単位ロックアウトの記述を削除し、IPレートリミット超過のみに限定しました。

### 変更したファイル
- [04_login.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/04_login.md)
- [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
