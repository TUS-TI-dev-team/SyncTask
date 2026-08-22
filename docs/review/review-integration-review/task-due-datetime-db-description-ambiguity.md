# TASKテーブルのDUE_DATETIMEカラム備考における未指定時挙動の記述曖昧さ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [02_task_management.md](docs/req-def/requirements/02_task_management.md)

## 1. 問題の概要
データベース設計書（`database_design.md`）の `TASK` テーブル定義において、`DUE_DATETIME` カラムの備考に「任意設定（未指定時は該当日 23:59 JST を適用）」と記載されている。
しかし、要件定義書（`02_task_management.md` 1節）およびAPI設計書（`04_tasks.md` 3.3.2 124行目）では、締切日時の指定自体を行わない場合は「締切未設定（`NULL`）」として保存され、「日付のみ指定されて時刻が省略された場合」にデフォルト締切時刻として `23:59:00+09:00` が適用される仕様となっている。
DB設計書の備考記述は「締切未設定（`NULL`）」と「日付のみ指定（時刻省略で23:59設定）」が混同されており、DB永続化仕様において誤認が生じるリスクがある。

## 2. 詳細な指摘内容
1. **DB設計書の記述**:
   - `docs/design/database_design.md` 第2節 `TASK` テーブル定義（49行目）：
     > | 締切日時 | `DUE_DATETIME` | `TIMESTAMPTZ` | 任意設定（未指定時は該当日 23:59 JST を適用） |
   - 「未指定時は該当日 23:59 JST を適用」と記載されているため、「締切を設定しないタスクを作成した場合でも何らかの該当日（作成日など？）の 23:59 が設定されるのか？」という誤読を招く表現になっている。
2. **要件定義書およびAPI設計書の仕様**:
   - `docs/req-def/requirements/02_task_management.md` 1節：
     > | **締切日時** | 任意 | 未設定 (null) | ・形式: `YYYY/MM/DD hh:mm`（日本標準時: JST）。<br>・時刻部分（`hh:mm`）が未指定（省略）の場合はデフォルト締切時刻として **`23:59 JST`** を適用。 |
   - `docs/design/api_design/04_tasks.md` 3.3.2（124行目）：
     > `due_datetime`: ISO 8601 日時文字列、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を設定）... 省略時または `null` 指定時は締切日時未設定（`null`）として作成。
   - 締切日時そのものを未指定（省略または `null`）とした場合は `DUE_DATETIME` は `NULL` として永続化されるのが正本仕様である。

## 3. 推奨される修正案
- `docs/design/database_design.md` 第2節 `TASK` テーブル定義の `DUE_DATETIME` カラムの備考欄を以下のように修正する：
  - 修正前: `任意設定（未指定時は該当日 23:59 JST を適用）`
  - 修正後: `任意設定（締切未設定時は NULL。日付のみ指定時は該当日 23:59:00+09:00 を設定して保存）`

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の `TASK.DUE_DATETIME` カラム備考を「任意設定（締切未設定時は NULL。日付のみ指定時は該当日 23:59:00+09:00 を設定して保存）」に修正しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
