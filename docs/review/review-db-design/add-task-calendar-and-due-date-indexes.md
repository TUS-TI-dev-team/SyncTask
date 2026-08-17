# カレンダー表示および締切日検索・締切間近タスク表示に対応する推奨インデックスの不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 15:07:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`database_design.md` の「7.1 タスク管理 (TASK)」に定義されている推奨インデックスは `(USER_ID, STATUS, IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC)` のみとなっており、カレンダー表示（期間指定による全ステータスタスク取得）や締切間近タスク表示・締切日検索のクエリにおいて `DUE_DATE` を効率的に絞り込むためのインデックスが不足しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の183〜187行目において、以下のインデックスのみが定義されています：
  ```sql
  -- タスク一覧の複合ソート高速化 (ユーザー別、ステータス別、ピン留め降順、締切日時昇順 NULLS LAST、作成日時降順)
  CREATE INDEX idx_task_user_status_sort ON TASK (USER_ID, STATUS, IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC);
  ```
- 一方、`docs/req-def/requirements.md` では以下の処理が規定されています：
  1. **カレンダー表示（138〜153行目）**:
     - 「締切日時が設定されている全ステータス（「未着手」「進行中」「完了」）のタスクについて、日本標準時（JST）における日付のセル上に配置して表示する」「バックエンドAPIはグリッド全体の開始日から終了日までのタスクを取得する」
     - 実行されるクエリ例: `SELECT * FROM TASK WHERE USER_ID = ? AND DUE_DATE >= ? AND DUE_DATE <= ?`
     - このクエリでは `STATUS` 条件が指定されないため、B-Tree インデックスのプレフィックスルールにより `idx_task_user_status_sort` の2番目のキー `STATUS` がスキップされ、`DUE_DATE` の範囲スキャンインデックスとして機能しません。
  2. **締切間近タスク表示（108〜112行目）およびタスク検索（131〜132行目）**:
     - `DUE_DATE <= ?` による範囲検索が行われます。
- また、一覧表示において「完了タスクを含める」オプションが選択された場合（全ステータス取得時）も、`STATUS` 条件なしで `(IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC)` のソートを行うため、`idx_task_user_status_sort` ではインデックスソートが適用されずファイルソートが発生します。

## 3. 推奨される修正案
`docs/design/database_design.md` の「7.1 タスク管理 (TASK)」セクションに、カレンダー表示および締切日範囲検索用のインデックス（および必要に応じて全ステータス一覧ソート用インデックス）を追加定義してください。

例:
```sql
-- タスク一覧の複合ソート高速化 (ユーザー別、ステータス指定時)
CREATE INDEX idx_task_user_status_sort ON TASK (USER_ID, STATUS, IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC);

-- カレンダー表示および締切日範囲検索用 (ユーザー別、締切日時)
CREATE INDEX idx_task_user_due_date ON TASK (USER_ID, DUE_DATE);
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:11:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の「7.1 タスク管理 (TASK)」において、以下のインデックスを追加定義しました：
1. 全ステータス取得時（完了を含む一覧表示等）の複合ソート高速化用インデックス: `idx_task_user_sort (USER_ID, IS_PINNED DESC, DUE_DATE ASC NULLS LAST, CREATED_AT DESC)`
2. カレンダー表示（期間指定全ステータス取得）および締切日範囲検索用インデックス: `idx_task_user_due_date (USER_ID, DUE_DATE)`

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
