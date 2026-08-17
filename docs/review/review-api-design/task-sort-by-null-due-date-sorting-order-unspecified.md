# sort_by パラメータにおける null 締切日時タスクのソート順序（NULLS LAST）の未規定

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`GET tasks` の `sort_by` パラメータで `due_date_asc`（締切昇順）または `due_date_desc`（締切降順）が指定された際、締切日時が未設定（`due_datetime: null`）のタスクがリストの先頭か末尾のどちらに配置されるか（`NULLS LAST`）が仕様書に規定されていません。

## 2. 詳細な指摘内容
`04_tasks.md` (L22) ではソート種別として `due_date_asc`（締切昇順）、`due_date_desc`（締切降順）が定義されています。
データベース設計書 `database_design.md` の推奨インデックス（L186, L189）では、`DUE_DATE ASC NULLS LAST` と `NULLS LAST` オプションが指定されています。

しかし、API仕様書 `04_tasks.md` には `null` 値の取り扱い（末尾配置）についての明記がありません。
RDBMS（PostgreSQL等）のデフォルト挙動では、`ORDER BY due_date ASC` の場合は `NULLS LAST` ですが、`ORDER BY due_date DESC` の場合は `NULLS FIRST` となり、締切降順ソート指定時に締切日時未設定のタスクが最上部に並んでしまう問題が発生する可能性があります。

## 3. 推奨される修正案
`04_tasks.md` (L22) の `sort_by` パラメータ説明に以下の補足を追記してください。

```markdown
| `sort_by` | string | × | `default` | ソート種別。指定可能値: `default`（ピン留め優先→締切昇順→作成日時降順）、`due_date_asc`（締切昇順）、`due_date_desc`（締切降順）、`created_at_desc`（作成日時降順）、`priority_desc`（優先度降順）。※`due_date_asc` および `due_date_desc` 指定時、締切日時未設定（`null`）のタスクは常に末尾に配置されます（`NULLS LAST`）。 |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 (`GET tasks`) の `sort_by` パラメータ説明欄に、「`due_date_asc` および `due_date_desc` 指定時、締切日時未設定（`null`）のタスクは常に末尾に配置されます（`NULLS LAST`）」というソート動作仕様を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
