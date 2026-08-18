# 要件定義書の「優先タスク表示」に対する専用複合インデックス定義の不足

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:15:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書（`requirements.md`）の「優先タスク表示」機能では、優先度「高」のタスクをピン留めに関係なく締切日時昇順・作成日時降順でソートして取得しますが、`docs/design/database_design.md` の推奨インデックス定義（Section 7.1）にこのクエリを最適化する複合インデックスが含まれていません。

## 2. 詳細な指摘内容
1. **要件定義書の仕様 (`docs/req-def/requirements.md` L106-110)**:
   - 対象: 優先度「高」のタスク（`WHERE USER_ID = :user_id AND PRIORITY = 'HIGH'`）
   - 並び順: 「締切日時が早い順（昇順）。締切日時が同一または未設定の場合は作成日時の新しい順（降順）で表示（**ピン留めによる並び替えは行わない**）」
   - 抽出条件・ソート: `ORDER BY DUE_DATE ASC NULLS LAST, CREATED_AT DESC`

2. **DB設計書のインデックス定義 (`docs/design/database_design.md` L186, L189)**:
   ```sql
   CREATE INDEX idx_task_user_status_sort ON TASK (USER_ID, STATUS, IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC);
   CREATE INDEX idx_task_user_sort ON TASK (USER_ID, IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC);
   ```

### 問題点：
- 既存のインデックスはすべて `IS_PINNED DESC` がソート順の先頭に含まれており、かつ `PRIORITY` カラムが含まれていません。
- 「優先タスク表示」のクエリを実行する際、インデックスの順序（`IS_PINNED` 優先）と要求されるソート順（`IS_PINNED` を無視した `DUE_DATE` 昇順）が一致しないため、インデックススキャンによるソートスキップができず、DB側で追加のインメモリ／ファイルソート（Filesort）が発生し、性能要件（2秒以下）への影響が生じる可能性があります。

## 3. 推奨される修正案
「7.1 タスク管理 (`TASK`)」に「優先タスク表示」用の複合インデックスを追加してください：

```sql
-- 優先タスク一覧の高速化 (ユーザー別、優先度指定時: 締切日時昇順 NULLS LAST、作成日時降順)
CREATE INDEX idx_task_user_priority_sort ON TASK (USER_ID, PRIORITY, DUE_DATE ASC NULLS LAST, CREATED_AT DESC);
```
