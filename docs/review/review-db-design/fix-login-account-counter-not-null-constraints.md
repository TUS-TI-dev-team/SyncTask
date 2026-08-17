# LOGIN_ACCOUNTテーブルのカウンターカラム（LOGIN_FAILED_COUNT, REAUTH_FAILED_COUNT）におけるNOT NULL制約の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 15:40:00
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `LOGIN_FAILED_COUNT`（ログイン失敗回数）および `REAUTH_FAILED_COUNT`（再認証失敗回数）カラムの制約定義が `INT / DEFAULT 0` となっており、`NOT NULL` 制約が明示されていません。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の第1章「アカウント管理 (LOGIN_ACCOUNT)」において、以下の通り定義されています：
  - `LOGIN_FAILED_COUNT`: `INT` / `DEFAULT 0`
  - `REAUTH_FAILED_COUNT`: `INT` / `DEFAULT 0`
- 他テーブル（例: `OTP_SESSION.ATTEMPT_COUNT` および `SEND_COUNT` では `INT / NOT NULL, DEFAULT 0`、`LOGIN_IP_RATE_LIMIT.FAILED_COUNT` では `INT / NOT NULL, DEFAULT 0`、`LOGIN_ACCOUNT.IS_DELETED` では `BOOLEAN / NOT NULL, DEFAULT FALSE`）では明示的に `NOT NULL` が付与されています。
- これらのカウンターカラムに `NOT NULL` 制約がない場合、予期せぬ `NULL` 値の挿入・更新により、アプリケーション層またはDB層でのインクリメント処理（`LOGIN_FAILED_COUNT + 1`）や比較処理（`WHERE LOGIN_FAILED_COUNT >= 5`）で `NULL` 伝播が発生し、アカウントロックアウトやセッション強制破棄の判定に不具合が生じるリスクがあります。

## 3. 推奨される修正案
`docs/design/database_design.md` の `LOGIN_ACCOUNT` テーブルにおける `LOGIN_FAILED_COUNT` および `REAUTH_FAILED_COUNT` のデータ型 / 制約定義を、明示的に `INT / NOT NULL, DEFAULT 0` に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:41:40
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の `LOGIN_ACCOUNT` テーブル定義において、`LOGIN_FAILED_COUNT` および `REAUTH_FAILED_COUNT` カラムのデータ型 / 制約定義を `INT / NOT NULL, DEFAULT 0` に更新し、予期せぬ NULL 値の混入を防ぐよう修正しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
