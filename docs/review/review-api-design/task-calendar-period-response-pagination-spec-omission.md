# カレンダー期間取得（start_date / end_date 指定）時のレスポンス pagination オブジェクトの定義不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` においてカレンダー表示用に `start_date` および `end_date` を指定した場合に、レスポンス内の `pagination` オブジェクトの各フィールド（`page`, `limit`, `total_pages`）がどのように返却されるかの定義が明確ではありません。

## 2. 詳細な指摘内容
`04_tasks.md` (L20) では、`start_date` および `end_date` 指定時の動作について以下のように説明されています。
> 指定時はページネーション limit を解除し期間内の全タスクを返却。最大許容期間幅は 42日間（6週間）

しかし、共通レスポンスフォーマット（L41-L46）には通常一覧取得時の例として以下が示されています。
```json
"pagination": {
  "page": 1,
  "limit": 20,
  "total_count": 45,
  "total_pages": 3
}
```
ページネーション `limit` が解除されたカレンダー表示モードにおいて、`limit` 値に何が設定されるか（例: 取得された全件数 `total_count` と同値が設定されるのか、`limit: null` / `limit: 0` となるのか）、および `total_pages` が一律 `1` に固定されるのかが明確になっていません。フロントエンドの共通レスポンスパーサーや型定義で不整合や誤作動を引き起こすリスクがあります。

## 3. 推奨される修正案
`GET tasks` のレスポンス補足注記（`04_tasks.md` L49 付近）に以下の説明を追記してください。

> ※ `start_date` / `end_date` を指定したカレンダー期間取得時は、ページネーションが無効化されて期間内の全タスクが一括返却されるため、`pagination` オブジェクトは `page: 1`, `limit: total_count`（取得件数と同値）, `total_pages: 1` として返却されます（該当タスクが0件の場合は `limit: 0, total_count: 0, total_pages: 1`）。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 (`GET tasks`) のレスポンス注記に、`start_date` / `end_date` 指定時の `pagination` オブジェクトの値（`page: 1`, `limit: total_count`, `total_pages: 1`）に関する詳細な動作仕様を追記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
