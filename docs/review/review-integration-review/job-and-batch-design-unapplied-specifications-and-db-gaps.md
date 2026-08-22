# ジョブ設計書・DB設計書における修正事項の未反映および仕様ギャップ

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:58:00
- **Target Files**:
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

## 1. 問題の概要
先行レビューにて指摘・合意された修正事項（Go コネクションプール環境下での Advisory Lock 専用コネクション保持・解放、チャンク単位での一時エラーリトライ方針、定期タスク生成の同期API集約とジョブ対象外方針、および日本語同一視検索用 `TASK.SEARCH_TEXT` カラム・インデックス定義）について、レビュー報告上は Resolved と記録されていたものの、実際の設計書（`job_design.md` および `database_design.md`）に修正内容が反映されておらず、仕様の欠落・矛盾が未解消のまま残存しています。

## 2. 詳細な指摘内容

1. **Advisory Lock の専用コネクション（`db.Conn`）排有取得・アンロック仕様の欠落 (`job_design.md`)**:
   - `job_design.md` 2-1節および4章シーケンス図において、PostgreSQL セッションレベル・アドバイザリロック（`pg_try_advisory_lock`）を使用する際、Go 言語の `database/sql` コネクションプールで別々の物理接続が払い出されるのを防ぐための要件（ジョブ開始時に `conn, err := db.Conn(ctx)` を取得し、ロック取得・チャンクトランザクション実行・アンロック・`defer conn.Close()` を同一コネクションで完結させる仕様）が記載されていません。
2. **チャンク削除処理の一時エラーリトライスコープおよび異常終了時アンロックの曖昧さ (`job_design.md`)**:
   - `job_design.md` 5-1節において、一時エラー（デッドロック `40P01` 等）のリトライスコープが「チャンク単位（最大3回指数バックオフ）」であることが明記されておらず、リトライ失敗時や恒久エラー発生時に即座にジョブを中断して `defer` で確実にアンロックした上で構造化エラーログを出力するフローが未定義です。
3. **定期タスク（繰り返しタスク）生成における設計方針ノートの欠落 (`job_design.md`, `database_design.md`)**:
   - 繰り返しタスクはバックグラウンド Cron バッチによる未来生成方式ではなく、タスク作成API（`POST /api/v1/tasks`）による「同期即時一括生成方式（最大100件）」を採用し、生成後は独立した通常レコードとして `TASK` テーブルに登録・管理されるため、Cron スケジューラにおける定期タスク生成ジョブおよびルール管理テーブル（`RECURRING_TASKS` 等）は不要（スコープ外）である旨の設計根拠ノートが両ドキュメントに記載されていません。
4. **`TASK.SEARCH_TEXT` カラムおよび trigram GIN インデックスの定義漏れ (`database_design.md`)**:
   - 要件定義書（`02_task_management.md`）およびタスクAPI設計書（`04_tasks.md`）では日本語同一視検索（英数大小・NFKC正規化・ひらがなカタカナ同一視）のために `TASK.SEARCH_TEXT` カラムおよび `pg_trgm` インデックスを利用する仕様が追加されましたが、`database_design.md` の `TASK` テーブル定義および第7.1節インデックス定義に `SEARCH_TEXT` カラム（`TEXT NOT NULL`）および `CREATE INDEX idx_task_search_text ON TASK USING gin (SEARCH_TEXT gin_trgm_ops);`（`CREATE EXTENSION IF NOT EXISTS pg_trgm;` 含む）が反映されていません。
5. **`job_design.md` 内の旧ドキュメント参照パスの残存 (`job_design.md`)**:
   - 行 4 の前提条件（`docs/req-def/requirements.md`, `docs/design/api_design.md`）および行 252 の監視方針（`docs/req-def/requirements.md`）において、分割ディレクトリ構成（`docs/req-def/requirements/README.md`, `docs/design/api_design/01_overview.md`, `docs/req-def/requirements/04_non_functional.md`）へのパス更新が反映されていません。

## 3. 推奨される修正案

1. **`job_design.md` の更新**:
   - 2-1節に Go の `db.Conn(ctx)` による専用コネクション排有取得・同一 `conn` 上でのロック制御・チャンクトランザクション実行・`defer conn.Close()` を明記する。
   - 1章冒頭に「定期タスク生成に関する設計方針（同期即時一括生成方式を採用し、定期バッチ生成はスコープ外）」を明記する。
   - 5-1節にチャンク単位リトライ（最大3回指数バックオフ）およびリトライ上限/恒久エラー時の即座中断・確実なアンロックを明記し、4章シーケンス図に反映する。
   - 行 4 および行 252 の参照パスを分割構成に合わせて修正する。
2. **`database_design.md` の更新**:
   - `TASK` テーブル定義に `SEARCH_TEXT`（`TEXT / NOT NULL`、検索用正規化テキスト）カラムを追加し、定期タスクはルールテーブルを持たず独立レコードとして格納される旨のノートを追記する。
   - 第7.1節に `CREATE EXTENSION IF NOT EXISTS pg_trgm;` および `CREATE INDEX idx_task_search_text ON TASK USING gin (SEARCH_TEXT gin_trgm_ops);` を追加する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/job_design.md` において、`db.Conn(ctx)` による専用コネクション排有取得・アンロック制御、チャンク単位指数バックオフリートライおよび異常終了時の中断・確実なアンロックを4章シーケンス図および2-1/5-1節に明記し、定期タスク同期即時一括生成方針ノートおよび参照パスを最新化しました。
- `docs/design/database_design.md` において、`TASK.SEARCH_TEXT`（`TEXT NOT NULL`）カラムおよび `pg_trgm` GIN インデックス定義、定期タスクの永続化方針ノートを反映しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

