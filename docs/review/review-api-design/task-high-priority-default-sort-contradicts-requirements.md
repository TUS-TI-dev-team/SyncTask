# 優先タスク (high_priority) ビューのデフォルトソート順における要件定義書との相違

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`04_tasks.md` L58 では `high_priority` ビュー指定時のデフォルトソート順が `ピン留め優先（is_pinned DESC） → 締切日時昇順（due_datetime ASC NULLS LAST） → 作成日時降順（created_at DESC）` と定義されていますが、`requirements.md` L106 の業務要件では「優先タスク表示: 並び順: 締切日時が早い順（昇順）。締切日時が同一または未設定の場合は作成日時の新しい順（降順）で表示（ピン留めによる並び替えは行わない）」と規定されており、要件定義書とAPI設計書の間で `high_priority` ビューでのピン留め優先ソートの有無について明確な不一致が存在します。

## 2. 詳細な指摘内容
1. **要件定義書 (`requirements.md` L106) の記述**:
   - 優先タスク表示 (`high_priority`) の並び順について `並び順: 締切日時が早い順（昇順）。締切日時が同一または未設定の場合は作成日時の新しい順（降順）で表示（ピン留めによる並び替えは行わない）` と明記されています。
2. **API設計書 (`04_tasks.md` L58) の記述**:
   - `GET tasks` のレスポンス補足注記 L58 において、`high_priority: ピン留め優先（is_pinned DESC） → 締切日時昇順（due_datetime ASC NULLS LAST） → 作成日時降順（created_at DESC）` と記述されています。
3. **影響**:
   - 実装者が API 設計書 (`04_tasks.md`) の通りにバックエンドのデフォルト SQL ソート条件を組み込むと、優先タスク一覧画面においてピン留めされたタスクが最上部に優先配置され、業務要件で規定されている「締切日時が早い順（ピン留めによる並び替えを行わない）」という仕様から逸脱した挙動となります。

## 3. 推奨される修正案
`04_tasks.md` 3.3.1 節 (`GET tasks`) の L58 にある `high_priority` のデフォルトソート定義から `ピン留め優先（is_pinned DESC）` を取り除き、要件定義書に合わせて以下のように修正してください：

```markdown
- `high_priority`: 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）※ピン留めによる並び替えは行いません
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:35
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 節 (`GET tasks`) の `high_priority` ビューにおけるデフォルトソート定義から `ピン留め優先（is_pinned DESC）` を取り除き、要件定義書 (`requirements.md`) の業務規定通り `締切日時昇順（due_datetime ASC NULLS LAST） → 作成日時降順（created_at DESC）※ピン留めによる並び替えは行いません` へ修正しました。

### 変更したファイル
- [04_tasks.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
