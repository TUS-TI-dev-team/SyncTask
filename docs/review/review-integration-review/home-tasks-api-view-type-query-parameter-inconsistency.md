# ホーム画面におけるタスク取得APIのクエリパラメータ名不整合（view_type=priority 誤記）

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L20)
  - [04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md#L5-L13)

## 1. 問題の概要
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の「ページネーションの補足」に記載されているタスク一覧取得APIのパラメータ例において、優先タスクの指定値が `view_type=priority` と誤記されています。バックエンドAPI仕様（[04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md)）では `view_type=high_priority` と定義されており、定義外の値を送信すると `400 Bad Request` になるため、フロントエンドとバックエンド間でパラメータ不整合が発生します。

## 2. 詳細な指摘内容
1. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 20行目**:
   - `タスク一覧取得API（GET /api/tasks?view_type={priority|near_deadline|pinned}）から選択中のタブのタスクデータを全件取得し...` と記載されています。
2. **[04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md) 5行目および13行目**:
   - `view_type=high_priority, view_type=near_deadline, view_type=pinned`
   - `view_type` パラメータの許容列挙値は `high_priority`（優先高）、`near_deadline`（72時間以内）、`pinned`（ピン留めのみ）であり、定義外の値（`priority` 等）が指定された場合は `400 Bad Request`（code: `"BAD_REQUEST"`）を返却すると定義されています。
3. **影響**:
   - 画面設計書の記述通りにフロントエンドを実装した場合、優先タスクタブの選択時に `view_type=priority` をリクエストしてしまい、APIが `400 Bad Request` エラーとなってタスク一覧が表示できなくなります。

## 3. 推奨される修正案
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 20行目の記載を、API仕様書に合わせて `high_priority` に修正してください。

```markdown
- タスク一覧取得API（`GET /api/tasks?view_type={high_priority|near_deadline|pinned}`）から選択中のタブのタスクデータを全件取得し、画面側（クライアントサイド）で20件単位のページ分割表示およびページネーション操作UIの制御を行う。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` 20行目のAPIクエリパラメータ表記を `view_type={high_priority|near_deadline|pinned}` に修正しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)

