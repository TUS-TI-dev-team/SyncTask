# TASKテーブルにおけるCOMMENTカラムのNULL許容性と空文字永続化方針の未確定性

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
API設計書（`04_tasks.md`）では、タスクのコメントについて「未入力時は空文字 `""` として登録」「レスポンス時は未入力の場合空文字 `""` として返却」「PATCH時に `""` または `null` 指定でコメントをクリア」と定義されている。
一方で、データベース設計書（`database_design.md`）の `TASK.COMMENT` カラム定義は `TEXT` とのみ記載されており、`NOT NULL` 制約や `DEFAULT ''`（空文字デフォルト）の有無、または `NULL` 許容でDB上はNULL保持しAPI層で空文字変換するのかといった永続化方針が明記されていない。

## 2. 詳細な指摘内容
1. **DB設計書の定義**:
   - `docs/design/database_design.md` 第2節 `TASK` テーブル定義（52行目）：
     > | コメント | `COMMENT` | `TEXT` | 補足メモ（0〜1000文字） |
   - `NOT NULL` 制約や `DEFAULT` 句の記載がなく、NULL許容（Nullable）カラムとなっている。
2. **API設計書の仕様**:
   - `docs/design/api_design/04_tasks.md` 3.3.1 / 3.3.2 / 3.3.3 / 3.3.4：
     - `POST tasks`: 「未入力時は空文字 `""` として登録」
     - `GET tasks` / `GET tasks/{task_id}`: 「コメントが未入力の場合、`comment` は空文字 `""` として返却されます」
     - `PATCH tasks/{task_id}`: 「空文字 `""` または `null` 指定でコメントをクリア（削除）」
   - API設計書では「未入力は空文字 `""`」と明記されているが、DB側で `TEXT NOT NULL DEFAULT ''` として空文字統一で永続化する設計なのか、あるいは DB上は `NULL` を許容して未入力時は `NULL` を格納し、APIシリアライザ側で `null` を `""` に置換して返却する設計なのかがDB定義から判別できない。

## 3. 推奨される修正案
- `docs/design/database_design.md` において、`TASK.COMMENT` カラムのデータ型・制約または備考欄に永続化方針を明記する：
  - **案A（DB上で空文字統一）**: `TEXT` / `NOT NULL, DEFAULT ''` とし、備考に「未入力時は空文字 `''` で保存」と記載する。
  - **案B（DB上でNULL許容）**: `TEXT`（NULL許容）を維持し、備考に「未入力・クリア時は NULL または空文字を保存（API返却時に空文字 `""` へ正規化）」と明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の `TASK.COMMENT` カラム定義を `TEXT` / `NOT NULL, DEFAULT ''`（備考: 補足メモ（0〜1000文字、未入力時は空文字で保存））に更新し、NULLを排除して空文字統一で永続化する方針を明記しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
