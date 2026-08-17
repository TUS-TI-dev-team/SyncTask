# OTP_SESSIONテーブルのカウンターカラム（ATTEMPT_COUNT, SEND_COUNT）におけるNOT NULL制約の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 15:20:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`OTP_SESSION` テーブルの `ATTEMPT_COUNT`（試行失敗回数）および `SEND_COUNT`（再送回数）カラムの制約定義が `INT / DEFAULT 0` となっており、`NOT NULL` 制約が明示されていません。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の第4章「OTPセッション管理 (OTP_SESSION)」において、以下の通り定義されています：
  - `ATTEMPT_COUNT`: `INT` / `DEFAULT 0`
  - `SEND_COUNT`: `INT` / `DEFAULT 0`
- 同文書の他テーブル（例: `LOGIN_IP_RATE_LIMIT.FAILED_COUNT` では `INT / NOT NULL, DEFAULT 0`、`LOGIN_ACCOUNT.IS_DELETED` では `BOOLEAN / NOT NULL, DEFAULT FALSE`）では、デフォルト値を持つ数値/フラグカラムに明示的に `NOT NULL` が付与されています。
- `OTP_SESSION` の `ATTEMPT_COUNT` や `SEND_COUNT` に `NOT NULL` 制約がない場合、予期せぬ `NULL` 値の挿入・更新により、アプリケーション層でのインクリメント処理（`attempt_count + 1`）や比較処理で不具合が生じるリスクがあります。

## 3. 推奨される修正案
- `OTP_SESSION` テーブルの `ATTEMPT_COUNT` および `SEND_COUNT` のデータ型 / 制約定義を、明示的に `INT / NOT NULL, DEFAULT 0` に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:22:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の第4章「OTPセッション管理 (OTP_SESSION)」テーブル定義において、`ATTEMPT_COUNT` および `SEND_COUNT` カラムのデータ型/制約を `INT / NOT NULL, DEFAULT 0` に修正し、NULL値混入による不具合リスクを排除しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
