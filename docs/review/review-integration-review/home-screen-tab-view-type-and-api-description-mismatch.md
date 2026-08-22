# ホーム画面のview_typeパラメータ表記誤記およびAPI設計書説明文との仕様不整合

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:58:00
- **Target Files**:
  - [screen_design.md](docs/design/screen_design.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
画面設計書（`screen_design.md`）において、ホーム画面のタブ切り替えに伴うAPI呼び出しパラメータの例示で `view_type={priority|near_deadline|pinned}` と誤って記載されている（API設計書上の正しい値は `high_priority`）。また、API設計書（`04_tasks.md`）の `GET tasks` の説明文に、タブ切り替えUI化前の「3つのAPIリクエストを個別に並行発行して各セクションごとに表示する」という旧仕様の記述が残存している。

## 2. 詳細な指摘内容
1. **画面設計書のクエリパラメータ値の誤記**:
   - `docs/design/screen_design.md` 20行目:
     > タスク一覧取得API（`GET /api/tasks?view_type={priority|near_deadline|pinned}`）から選択中のタブのタスクデータを全件取得し...
   - しかし、`docs/design/api_design/04_tasks.md` 3.3.1 節で定義されている `view_type` の有効値は `high_priority`, `near_deadline`, `pinned` であり、`priority` ではない（`priority` は優先度絞り込み用パラメータ `priority=high` 等として予約されている）。このため画面設計書のパラメータ名がAPI仕様と不整合になっている。
2. **API設計書におけるホーム画面説明文の残存旧仕様**:
   - `docs/design/api_design/04_tasks.md` 5行目:
     > ホーム画面においては、フロントエンドが `view_type=high_priority`, `view_type=near_deadline`, `view_type=pinned` の3つのAPIリクエストを個別に並行発行し、各セクションごとに独立して20件単位のページ分割制御およびページネーションUI操作を行います。
   - 画面設計書ではホーム画面がタブ切り替えUI（選択中のタブのAPIのみ取得してページネーション制御）に改訂されたため、API設計書側の「3つのAPIリクエストを個別に並行発行し、各セクションごとに独立して...」という説明文と乖離が生じている。

## 3. 推奨される修正案
1. **`docs/design/screen_design.md` の修正**:
   - 20行目のAPIクエリパラメータ表記を `GET /api/tasks?view_type={high_priority|near_deadline|pinned}` に修正する。
2. **`docs/design/api_design/04_tasks.md` の修正**:
   - 5行目のホーム画面に関する説明文を「ホーム画面においては、フロントエンドが選択中のタブに応じて `view_type=high_priority`、`view_type=near_deadline`、`view_type=pinned` のAPIリクエストを発行し、20件単位のページ分割制御およびページネーションUI操作を行います」等に更新し、タブ切り替えUI仕様と整合させる。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` 20行目のパラメータ例を `high_priority` に修正し、`docs/design/api_design/04_tasks.md` 5行目のホーム画面説明文を選択中タブに応じたAPI発行仕様に更新しました。

### 変更したファイル
- [screen_design.md](docs/design/screen_design.md)
- [04_tasks.md](docs/design/api_design/04_tasks.md)

