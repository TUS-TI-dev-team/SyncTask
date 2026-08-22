# タスク検索における日本語同一視要件とAPI・DB設計の乖離

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [02_task_management.md](docs/req-def/requirements/02_task_management.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
要件定義書（`02_task_management.md`）ではキーワード検索において「英大文字/小文字、日本語全角/半角、ひらがな/カタカナを区別せず部分一致で検索する」と定義されているが、API設計書（`04_tasks.md`）およびデータベース設計書（`database_design.md`）では「Case-Insensitive（大文字小文字無視）」としか記載されておらず、日本語の全角/半角やひらがな/カタカナを同一視して検索するためのアーキテクチャ・DB設計・インデックス方針が完全に欠落している。

## 2. 詳細な指摘内容
1. **要件定義書の定義**:
   - `docs/req-def/requirements/02_task_management.md` 3.2 節において以下のように規定されている：
     > **Case-Insensitive 部分一致検索**: 英大文字/小文字、日本語全角/半角、ひらがな/カタカナを区別せず部分一致で検索する。
2. **API設計書の定義不足**:
   - `docs/design/api_design/04_tasks.md` 3.3.1 節のクエリパラメータ `keyword` の説明では以下のようにしか書かれていない：
     > タスク名およびコメントの部分一致検索（Case-Insensitive、前後の空白文字（半角・全角スペース、タブ、改行）をトリム処理。... SQLワイルドカード特殊文字 % や _ や \ はリテラルエスケープして部分一致を検索...）
   - 日本語全角/半角（例: `１２３` vs `123`、`ｶﾀｶﾅ` vs `カタカナ`）やひらがな/カタカナ（例: `れぽーと` vs `レポート`）の同一視検索を行うための正規化処理（NFKC正規化、カタカナ変換等）の責務（API側でクエリと登録値を正規化するのか、DB拡張・照合順序で行うのか）が明記されていない。
3. **DB設計書の設計欠落**:
   - `docs/design/database_design.md` の第2節（`TASK` テーブル定義）および第7.1節（インデックス設計）において、`TITLE`（`VARCHAR(255)`）および `COMMENT`（`TEXT`）に対する日本語同一視・部分一致検索用の拡張（例: pg_trgm / pg_bigm、検索用正規化カラム `TITLE_SEARCH_NORM`、ICU Collation 等）が一切定義されていない。
   - 一般的なRDBMS（PostgreSQL等）の標準機能（`ILIKE` や `LOWER()`）のみではアルファベットの大文字小文字しか同一視できず、日本語の全角/半角やひらがな/カタカナを同一視した部分一致検索は実現不可能であるため、実装時に検索機能の不具合や性能劣化を引き起こす重大なリスクがある。

## 3. 推奨される修正案
以下のいずれかの方針を決定し、API設計書およびDB設計書に仕様を明記すること：
- **方針A（アプリケーション層での正規化方式）**:
  - `TASK` テーブルに検索専用の正規化カラム（例: `SEARCH_TEXT TEXT`）を追加、または登録・更新時に全角英数・半角カナをUnicode正規化（NFKC）し、ひらがなをカタカナへ統一した文字列を生成・格納する。
  - API検索時（`GET tasks`）も入力された `keyword` を同一ルールで正規化して部分一致クエリを実行する。
  - `SEARCH_TEXT` カラムに trigram インデックス（`gin_trgm_ops` 等）を設定する。
- **方針B（DB拡張・Collation方式）**:
  - RDBMS固有のICU Collation（全角半角・かなカナ無視設定）や全文検索エンジンの適用方針をDB設計書およびインデックス設計に具体的に追記する。
- **方針C（要件スコープの調整）**:
  - もし日本語の全角/半角・かな/カナ同一視を行わず英大文字小文字のみの無視にとどめる場合は、`docs/req-def/requirements/02_task_management.md` の記述を修正してドキュメント間の整合性を確保する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
タスク検索における日本語同一視（英大文字/小文字、日本語全角/半角、ひらがな/カタカナ同一視）を実現するため、アプリケーション層で小文字化＋NFKC正規化＋ひらがな→カタカナ変換した文字列を `TASK.SEARCH_TEXT` カラムに格納し、`pg_trgm` GIN インデックスによる部分一致検索を行う設計を DB設計書・API設計書・要件定義書に反映しました。

### 変更したファイル
- [database_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\database_design.md)
- [04_tasks.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\04_tasks.md)
- [02_task_management.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\req-def\requirements\02_task_management.md)
