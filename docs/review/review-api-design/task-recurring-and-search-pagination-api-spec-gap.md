# タスク繰り返し作成および検索・ページネーション・カレンダー取得仕様の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
要件定義書および画面設計書で中核機能として定められている「毎週タスク自動作成（即時一括生成方式）」、「タスク一覧のサーバーサイド絞り込み・ページネーション（20件/ページ）」、および「カレンダー表示（グリッド期間指定取得）」に対するリクエストパラメータおよびレスポンス仕様が `api_design.md` で欠落しています。

## 2. 詳細な指摘内容
1. **タスク作成（`POST tasks/`）における繰り返し作成パラメータの欠落**:
   - `docs/req-def/requirements.md` L83-92 および `docs/design/screen_design.md` L25, L33-34 では、期間（開始日・終了日）および曜日を指定して最大100件までタスクを一括即時生成する機能が定義されています。
   - しかし `docs/design/api_design.md` L28 では「入力: セッションID、タスク情報 (名前, 優先度, 締切, コメント)」のみとなっており、繰り返し作成用フラグや期間・曜日指定パラメータが一切定義されていません。また、複数件生成時のレスポンス形式（作成されたタスク配列 or 生成件数）も不明です。
2. **タスク一覧取得（`GET tasks/`）におけるクエリパラメータ・ページネーション仕様の欠落**:
   - `docs/design/api_design.md` L27 では「フロント側でソート・検索」「出力: タスク一覧JSON」と記載されています。
   - しかし `docs/req-def/requirements.md` L107, L111, L117, L124, L137 では「1ページあたり20件表示」のページネーションや、キーワード（部分一致・Case-Insensitive）・優先度・締切日・ステータス・完了タスク表示切替などの詳細な絞り込み要件が定義されています。
   - 全件をフロントエンドに返却してクライアント側で処理する設計の場合、タスク件数が増加した際にターンアラウンドタイム2秒以下（L158）を満たせなくなる恐れがあります。
3. **カレンダー表示用取得クエリの欠落**:
   - `docs/req-def/requirements.md` L147 には「バックエンドAPIはグリッド全体の開始日から終了日までのタスクを取得する」と明記されていますが、期間範囲（`from_date`, `to_date` 等）を指定するクエリパラメータが設計されていません。

## 3. 推奨される修正案
1. `POST tasks/` の入力仕様に、単一作成だけでなく繰り返し一括作成用のパラメータ（`is_recurring`, `start_date`, `end_date`, `days_of_week` 等）を追加し、レスポンス構造を明確化してください。
2. `GET tasks/` に、サーバーサイドでの絞り込みおよびページネーション用クエリパラメータ（`page`, `limit`, `keyword`, `priority`, `status`, `include_completed`, `due_date`, `from_date`, `to_date`, `sort_by` 等）を定義してください。
3. レスポンス仕様に、タスク配列（`items`）に加えてページネーションメタデータ（`total_count`, `page`, `total_pages` 等）を含めてください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` の `POST tasks` において、単一作成に加え `is_recurring`, `recurring_rule`（開始日・終了日・曜日配列・締切時刻）による即時一括生成（最大100件）パラメータおよびレスポンススキーマを定義しました。
- `GET tasks` において、サーバーサイドページネーション（`page`, `limit`）、ビュー指定（`view_type`）、検索・絞り込み（`keyword`, `priority`, `status`, `include_completed`, `due_date`）、カレンダー期間取得（`start_date`, `end_date`）をカバーする包括的なクエリパラメータとレスポンスメタデータ（`pagination` オブジェクト）を網羅的に定義しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
