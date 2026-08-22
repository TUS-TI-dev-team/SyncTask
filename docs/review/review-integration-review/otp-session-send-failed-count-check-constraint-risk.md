# OTP_SESSION.SEND_FAILED_COUNT の CHECK制約（0〜4）による更新時エラーリスク

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 13:49:30
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)

## 1. 問題の概要

`docs/design/database_design.md` 7.2節において、`OTP_SESSION` テーブルの連続送信失敗回数カラム（`SEND_FAILED_COUNT`）に対する CHECK 制約が `CHECK (SEND_FAILED_COUNT BETWEEN 0 AND 4)` と定義されています。

一方、テーブル定義（4節）およびAPI設計書（`02_auth.md` 3.1節）では「初回送信・手動再送・自動再送を含め、連続送信失敗回数が5回に到達した時点で対象OTPセッションを物理削除する」と規定されています。

もしアプリケーションの実装において、メール送信失敗時に `SEND_FAILED_COUNT` をインクリメント（`SEND_FAILED_COUNT = SEND_FAILED_COUNT + 1`、4→5）して DB を UPDATE した後に「5回到達」を検知して物理削除（DELETE）するトランザクション順序をとった場合、5 への UPDATE 時にこの CHECK 制約に違反して DB 例外（Check constraint violation）が発生し、エラーハンドリングや補償削除処理が失敗するリスクがあります。

## 2. 詳細な指摘内容

### 該当箇所の定義

1. **`docs/design/database_design.md` (4節 テーブル定義)**:
   ```markdown
   | 連続送信失敗回数 | SEND_FAILED_COUNT | INT / NOT NULL, DEFAULT 0 | 初回送信・手動再送・自動再送を含む連続失敗回数。成功時に0へリセットし、5回到達時は対象セッションを物理削除 |
   ```

2. **`docs/design/database_design.md` (7.2節 CHECK制約)**:
   ```sql
   ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_send_failed_count
   CHECK (SEND_FAILED_COUNT BETWEEN 0 AND 4);
   ```

3. **`docs/design/api_design/02_auth.md` (3.1節 OTP API共通規則)**:
   ```markdown
   実メール送信に失敗した場合は DELIVERY_STATUS='sendable'、SEND_FAILED_COUNT+=1 とし、...
   初回送信、手動再送、5回照合失敗による自動再送を通じて5回連続で送信に失敗した場合、対象OTPセッションを物理削除し、410 Gone（code: OTP_SESSION_INVALIDATED）を返します。
   ```

4. **問題点**:
   - `SEND_FAILED_COUNT += 1` という一般的なインクリメント更新処理を行う場合、4回目の失敗状態から5回目の失敗が発生した際に `SEND_FAILED_COUNT = 5` への UPDATE が試行されます。
   - 制約が `BETWEEN 0 AND 4` だと、値 `5` の書き込みが拒絶され、DBエラー（ロールバック）となって正常な 410 エラー応答やセッション削除が行えなくなる恐れがあります。

## 3. 推奨される修正案

以下のいずれかの対応を実施してください：

1. **CHECK制約の緩和（推奨）**:
   `database_design.md` 7.2節の制約定義を以下のように修正し、カウント上限 5 まで（または非負値）を許容するようにします：
   ```sql
   ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_send_failed_count
   CHECK (SEND_FAILED_COUNT BETWEEN 0 AND 5);
   ```
   または
   ```sql
   ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_send_failed_count
   CHECK (SEND_FAILED_COUNT >= 0);
   ```

2. **実装・処理順序の明記**:
   5回目の失敗時は UPDATE を行わずに直ちに DELETE する設計を前提とする場合は、その旨を DB 設計書および API 設計書の処理順序に明示してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
`OTP_SESSION.SEND_FAILED_COUNT` の CHECK 制約を `CHECK (SEND_FAILED_COUNT BETWEEN 0 AND 5)` に緩和し、5回目の失敗記録および物理削除トランザクションで制約エラーが発生しないよう修正しました。

### 変更したファイル
- [database_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\database_design.md)
