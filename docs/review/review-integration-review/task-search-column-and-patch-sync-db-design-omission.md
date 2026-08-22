# TASKテーブルにおけるSEARCH_TEXT定義欠落およびPATCH時の検索テキスト再生成仕様の不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:58:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
要件定義書（`02_task_management.md` 3.2）およびタスクAPI設計書（`04_tasks.md` 3.3.1, 3.3.2）において、日本語同一視（全角/半角、ひらがな/カタカナ、英大文字/小文字）検索を実現するために `TASK.SEARCH_TEXT` カラムおよび `pg_trgm` GINインデックスを利用する方針が定められている。
しかし、データベース設計書（`database_design.md`）の `TASK` テーブル定義およびインデックス設計に `SEARCH_TEXT` カラムと `pg_trgm` GINインデックスの定義が欠落している。また、タスク更新API（`PATCH tasks/{task_id}`）においてタイトルやコメントが更新された際の `SEARCH_TEXT` 自動再生成・更新仕様が明記されていない。

## 2. 詳細な指摘内容
1. **DB設計書における `SEARCH_TEXT` カラム定義の欠落**:
   - `docs/design/database_design.md` 第2節の `TASK` テーブル定義（カラム一覧）に、正規化文字列を保持する `SEARCH_TEXT` カラム（`TEXT` / `NOT NULL`）が存在しない。
2. **DB設計書における `pg_trgm` 拡張および GIN インデックス定義の欠落**:
   - `docs/design/database_design.md` 第7.1節（インデックス設計）において、`SEARCH_TEXT` に対する部分一致検索インデックス（`CREATE EXTENSION IF NOT EXISTS pg_trgm;` および `CREATE INDEX idx_task_search_text ON TASK USING gin (SEARCH_TEXT gin_trgm_ops);`）が定義されていない。
3. **タスク更新API（`PATCH tasks/{task_id}`）における再生成責務の記述漏れ**:
   - `docs/design/api_design/04_tasks.md` 3.3.2（`POST tasks`）では作成時にタイトル・コメントから `SEARCH_TEXT` を自動生成して保存する旨が記載されているが、3.3.4（`PATCH tasks/{task_id}`）には `title` または `comment` が更新された際に `SEARCH_TEXT` を小文字化＋NFKC正規化＋ひらがな→カタカナ変換ルールで再生成・更新する旨の記載が欠落している。更新時に再生成されない場合、検索インデックスとデータの実値が乖離するバグの原因となる。
4. **定期タスク永続化方針ノートのDB設計書への反映漏れ**:
   - `docs/design/database_design.md` の `TASK` テーブルセクションにおいて、繰り返しタスクは別テーブルを持たず即時一括生成された独立レコードとして永続化される旨の補足ノートが記載されていない。

## 3. 推奨される修正案
1. **`docs/design/database_design.md` の修正**:
   - 第2節 `TASK` テーブル定義に `SEARCH_TEXT`（`TEXT` / `NOT NULL` / 検索用正規化文字列）カラムを追加する。
   - `TASK` テーブル直下に「繰り返しタスクはルールマスターテーブルを持たず、即時一括生成された独立通常タスクレコードとして永続化される」旨のノートを追記する。
   - 第7.1節に `CREATE EXTENSION IF NOT EXISTS pg_trgm;` および `CREATE INDEX idx_task_search_text ON TASK USING gin (SEARCH_TEXT gin_trgm_ops);` を追記する。
2. **`docs/design/api_design/04_tasks.md` の修正**:
   - 3.3.4 `PATCH tasks/{task_id}` の概要および説明に「`title` または `comment` が更新された場合、バックエンドで検索用正規化文字列（小文字化＋NFKC正規化＋ひらがな→カタカナ変換）を再生成し、`TASK.SEARCH_TEXT` カラムを自動更新する」旨を明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` に `TASK.SEARCH_TEXT` カラム（`TEXT NOT NULL`）、定期タスク永続化方針ノート、および `pg_trgm` GINインデックス定義を追加しました。
- `docs/design/api_design/04_tasks.md` 3.3.4 `PATCH tasks/{task_id}` において、`title` または `comment` 更新時に `SEARCH_TEXT` を自動再生成・更新する仕様を明記しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
- [04_tasks.md](docs/design/api_design/04_tasks.md)
