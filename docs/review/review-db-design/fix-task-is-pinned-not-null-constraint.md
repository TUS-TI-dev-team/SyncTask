# TASKテーブルのIS_PINNEDカラムにおけるNOT NULL制約の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 15:07:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`TASK` テーブルの `IS_PINNED`（ピン留めフラグ）カラムの制約が `BOOLEAN / DEFAULT FALSE` となっており、`NOT NULL` 制約が明記されていません。これにより NULL 値が混入する余地があり、ソート処理やピン留めタスク絞り込み時の予期せぬ挙動につながるリスクがあります。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の51行目において、`TASK` テーブルの `IS_PINNED` カラムの制約が以下のように定義されています：
  ```markdown
  | ピン止めフラグ | `IS_PINNED` | `BOOLEAN` / `DEFAULT FALSE` | | ✅ |
  ```
- 他のテーブルの真偽値カラム（例: `LOGIN_ACCOUNT.IS_DELETED` 19行目 `BOOLEAN / NOT NULL, DEFAULT FALSE`）では `NOT NULL` 制約が付与されていますが、`IS_PINNED` には `NOT NULL` が欠落しています。
- `docs/req-def/requirements.md` ではピン留めタスクの優先表示（121行目、140行目）やピン留めタスク表示（113〜117行目 `is_pinned = true`）が規定されています。データベース上で `IS_PINNED` に NULL が許容されている場合、3値論理（TRUE / FALSE / UNKNOWN）による検索漏れや、`ORDER BY IS_PINNED DESC` 実行時に DBMS のデフォルトソート順（PostgreSQL 等では `NULLS FIRST` となり NULL レコードが最上位に来る等）によって一覧の表示順序が崩れる原因となります。

## 3. 推奨される修正案
`docs/design/database_design.md` の `TASK` テーブルにおける `IS_PINNED` カラムの制約を `BOOLEAN / NOT NULL, DEFAULT FALSE` に更新してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:11:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の `TASK` テーブル定義において、`IS_PINNED` カラムの制約を `BOOLEAN / NOT NULL, DEFAULT FALSE` に更新し、NULL値の混入を防止しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
