# GET tasksにおけるdue_date絞り込みとinclude_completedパラメータ併用時の抽出仕様の曖昧さ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [02_task_management.md](docs/req-def/requirements/02_task_management.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
API設計書（`04_tasks.md` 3.3.1）において、`due_date` クエリパラメータは「指定日当日の全タスク、および指定日より過去の未完了タスク（`status != 'completed'`）を抽出対象とし、過去の完了済みタスクは常に除外」と説明されている。
しかし、`include_completed` パラメータ（通常一覧時デフォルト `false`）と併用された場合に、「指定日当日の完了タスク」が `include_completed=false` のときに含まれるのか除外されるのかが、パラメータ説明の組み合わせとして曖昧である。

## 2. 詳細な指摘内容
1. **API設計書の記述**:
   - `docs/design/api_design/04_tasks.md` 3.3.1 `GET tasks` Query Parameters：
     - `include_completed`: 「通常一覧時はデフォルト `false`、`start_date` / `end_date` 指定（カレンダー表示）時はデフォルト `true`。なお `status` が明示指定された場合は本パラメータは無視されます」
     - `due_date`: 「締切日絞り込み（`YYYY-MM-DD`。指定日当日の `00:00:00+09:00 <= due_datetime <= 23:59:59+09:00` の全タスク、および指定日より過去 `due_datetime < 00:00:00+09:00` の未完了タスク（`status != 'completed'`）を抽出対象とし、過去の完了済みタスクは常に除外。締切日時未設定 `null` のタスクは除外。未指定時は絞り込みを行わない...）」
   - `due_date` 指定時に `include_completed` を省略（デフォルト `false`）した場合、`due_date` の説明にある「指定日当日の全タスク」のうち完了タスクが含まれるのか、それとも `include_completed=false` が適用されて指定日当日の完了タスクも除外されるのかが明確に定義されていない。
2. **要件定義書および画面設計書の仕様との対照**:
   - `docs/req-def/requirements/02_task_management.md` 3.2 節（94行目）：
     > なお、ステータス絞り込み（完了）や完了タスク表示トグルとの併用時も、過去日付の完了タスクは抽出されず、指定日当日の完了タスクのみが表示対象となる。
   - `docs/design/screen_design.md` タスク操作補足（44行目）：
     > 「完了タスク表示/非表示」切り替えトグルがONの場合も、表示される完了タスクは指定日当日に締切を迎えた完了タスクのみとなり、過去日付の完了タスクは含まれない。
   - 要件上は、「完了表示トグルがOFF（`include_completed=false`）のときは指定日当日も含め未完了タスクのみを表示」「完了表示トグルがON（`include_completed=true`）のときは指定日当日の完了タスクを含めて表示（過去の完了タスクは常に除外）」というフィルタリングが意図されている。

## 3. 推奨される修正案
- `docs/design/api_design/04_tasks.md` の 3.3.1 `GET tasks` において、`due_date` および `include_completed` の説明欄を以下のように明確化する：
  1. `due_date` 指定時の抽出ロジック：
     - `include_completed=true` の場合: 指定日当日の全タスク（未完了＋完了）＋ 指定日より過去の未完了タスク（過去の完了タスクは除外）
     - `include_completed=false`（デフォルト）の場合: 指定日当日の未完了タスク ＋ 指定日より過去の未完了タスク（当日・過去を問わずすべての完了タスクを除外）
  2. この関係性を `due_date` パラメータの説明および `include_completed` の説明欄に明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/04_tasks.md` 3.3.1 `GET tasks` において、`due_date` と `include_completed` 併用時の抽出ロジック（`include_completed=true` 時は当日完了を含め過去完了は除外、`false` 時は当日・過去とも完了除外）を明確化しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
