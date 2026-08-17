# TASKテーブルのPRIORITYにおける値定義およびデフォルト値の不備

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 13:16:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
`TASK` テーブルの `PRIORITY` カラムの備考に「LOW, MEDIUM, HIGH など」と曖昧に記載されており、要件定義書で指定されている「高、中、低（初期値：中）」という取り得る値の限定およびデフォルト値定義（`DEFAULT 'MEDIUM'`）が欠落しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の45行目において、`PRIORITY` カラムの制約が `VARCHAR(20) / NOT NULL`、備考が `LOW, MEDIUM, HIGH など` と定義されています。
- しかし、`docs/req-def/requirements.md` の80行目、95行目、130行目および `docs/design/screen_design.md` の25行目・26行目では以下のように定義されています：
  > 80: 優先度（高、中、低、初期値：中）
  > 95: タスクの優先度変更（高、中、低）
  > 130: 完全一致（高、中、低。未指定・すべて選択時は全優先度を対象とする）
- 備考の「など」という表記により取り得る値の定義が曖昧になっており、またタスク作成時の初期値である「中（`MEDIUM`）」に対応する `DEFAULT 'MEDIUM'` 制約がスキーマ上に定義されていません。

## 3. 推奨される修正案
- `PRIORITY` カラムのデータ型 / 制約を `VARCHAR(20) / NOT NULL, DEFAULT 'MEDIUM'` に変更してください。
- 備考欄から「など」を削除し、取り得る値（例: `LOW` [低], `MEDIUM` [中・初期値], `HIGH` [高]）を明確に限定して記載してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`TASK` テーブルの `PRIORITY` カラムの制約を `VARCHAR(20) / NOT NULL, DEFAULT 'MEDIUM'` に更新し、備考欄の記述を `LOW (低), MEDIUM (中・初期値), HIGH (高)` に限定・明確化しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
