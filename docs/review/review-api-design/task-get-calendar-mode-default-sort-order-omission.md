# GET tasks のカレンダー期間取得（start_date/end_date）時におけるデフォルトソート順序の仕様定義漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`04_tasks.md` の `GET tasks` (L84) において、`start_date` / `end_date` を指定したカレンダー期間取得時の動作仕様が記載されていますが、`sort_by` パラメータが省略された場合に返却される `items` 配列のデフォルトソート順序が明記されていません。

## 2. 詳細な指摘内容
1. **要件定義とAPI仕様書の記載漏れ**:
   - `requirements.md` L140 では、「同一日付セル内に複数のタスクが存在する場合の並び順は、1. ピン留めされているタスクを優先、2. 締切時刻（hh:mm:ss）が早い順（昇順）、3. 作成日時の新しい順（降順）で表示する」と規定されています。
   - `04_tasks.md` の L79-L82 では `view_type` 省略時や指定時のデフォルトソート順序が詳しく記載されていますが、`start_date` / `end_date` 指定（カレンダー表示用取得）時に `sort_by` を省略した場合のデフォルトソート順序について注記がありません。
   - バックエンドが何順で `items` を返却すべきか（`due_datetime ASC` なのか、ピン留め優先 `is_pinned DESC → due_datetime ASC → created_at DESC` なのか）が不明確であり、フロントエンドでカレンダー描画・グループ化を行う際に予期せぬソート順の不一致が発生する恐れがあります。

## 3. 推奨される修正案
`04_tasks.md` の L84（カレンダー期間取得に関する注記）またはソート順序説明ブロックに、`start_date` / `end_date` 指定時に `sort_by` が省略された場合のデフォルトソート順序を以下のように明記してください。

```markdown
※ `start_date` / `end_date` 指定（カレンダー期間取得）時に `sort_by` が省略された場合のデフォルトソート順序は、`is_pinned DESC`（ピン留め優先） → `due_datetime ASC`（締切日時昇順） → `created_at DESC`（作成日時降順） → `id DESC`（タイブレーク）となります。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 節 (`GET tasks`) のカレンダー期間取得注記に、`start_date` / `end_date` 指定時に `sort_by` が省略された場合のデフォルトソート順序が `is_pinned DESC`（ピン留め優先） → `due_datetime ASC`（締切日時昇順） → `created_at DESC`（作成日時降順） → `id DESC`（タイブレーク）となる旨を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
