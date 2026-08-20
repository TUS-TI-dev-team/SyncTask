# GET tasks の sort_by における優先度ソート順・タイブレーク条件およびページネーション決定性の不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` API のクエリパラメータ `sort_by` において、`priority_desc` 指定時のビジネス優先度並び順（`high` > `medium` > `low`）および、単一のソートキーを指定した場合に同一値をもつタスクの確定的なタイブレーク順序（二次・三次ソートキー）の定義が不足しており、サーバーサイドページネーションの決定性が損なわれるリスクがあります。

## 2. 詳細な指摘内容
`04_tasks.md` L22 では `sort_by` パラメータとして以下が定義されています：
> `sort_by`: `default`（ピン留め優先→締切昇順→作成日時降順）、`due_date_asc`（締切昇順）、`due_date_desc`（締切降順）、`created_at_desc`（作成日時降順）、`priority_desc`（優先度降順）。※ `due_date_asc` および `due_date_desc` 指定時、締切日時未設定（`null`）のタスクは常に末尾に配置されます（`NULLS LAST`）。

この定義には以下の不足点が存在します：
1. **`priority_desc` のソート基準の未定義**: アルファベット順（`medium` > `low` > `high` 等）ではなく、ドメイン上の優先度（`high` → `medium` → `low`）で降順にソートする旨が明記されていません。
2. **タイブレーク（同順位判定）の未定義**: `priority_desc`, `due_date_asc`, `due_date_desc`, `created_at_desc` 等の指定時、複数のタスクが同じ優先度・締切日時・作成日時を持つ場合の二次・三次ソートキー（例: `created_at DESC`, `id DESC`）が規定されていません。

SQL データベースにおいてユニークな順序保証キー（タイブレークキー）が存在しない場合、`LIMIT / OFFSET` 方式によるページネーション（`page`, `limit`）においてページの境界付近でタスクの重複表示や要素の読み飛ばしが発生します。

## 3. 推奨される修正案
`04_tasks.md` L22 の `sort_by` パラメータ定義の説明文を以下のように修正し、優先度のランク順および決定論的タイブレーク順序を明記してください。

```markdown
| `sort_by` | string | × | `default` | ソート種別。指定可能値: `default`（ピン留め優先→締切昇順→作成日時降順）、`due_date_asc`（締切昇順）、`due_date_desc`（締切降順）、`created_at_desc`（作成日時降順）、`priority_desc`（優先度降順: `high` → `medium` → `low`）。※`due_date_asc` および `due_date_desc` 指定時、締切日時未設定（`null`）のタスクは常に末尾に配置されます（`NULLS LAST`）。なお、すべてのソート指定において第一ソート条件で同一値となるタスクについては、確定的なページネーション順序を保証するためタイブレーク条件として `created_at DESC` → `id DESC` が適用されます。 |
```

## 修正完了報告

- **Resolved At**: 2026-08-17 16:56:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 の `sort_by` クエリパラメータ説明文を修正し、`priority_desc` の降順指定時のドメイン優先度並び順（`high` → `medium` → `low`）および、第一ソート条件が同一値となるタスクに対する確定的なタイブレーク条件（`created_at DESC` → `id DESC`）を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
