# `GET tasks` における `due_date` パラメータと `include_completed` / `status` の相互作用および過去完了タスク除外ルールの論理不備

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-17 18:05:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`04_tasks.md` の `GET tasks` (L19) において、`due_date`（締切日絞り込み）パラメータの抽出仕様が「指定日当日の全タスク、および指定日より過去の未完了タスクを抽出対象とし、過去の完了済みタスクは常に除外」と定義されています。
この仕様は、`include_completed`（完了タスク含めるフラグ）や `status=completed`（ステータス明示指定）パラメータとの相互作用において論理的矛盾を生じさせており、過去の完了タスクを検索表示したいユースケースを阻害すると同時に、デフォルト検索時の動作に曖昧さをもたらしています。

## 2. 詳細な指摘内容
1. **`include_completed=true` / `status=completed` 指定時の過去完了タスク不当除外**:
   - `requirements.md` L131-L134 では、「指定日当日までの締め切り（〜指定日 23:59:59 JST）に該当する締切日時を持つタスクを検索。ステータス指定時はその条件に従う」と定義されています。
   - しかし `04_tasks.md` L19 では「過去の完了済みタスクは常に除外」と一律規定されているため、クライアントが `due_date=2026-08-20` と同時に `include_completed=true` や `status=completed` を指定して指定日までの完了タスク一覧を確認しようとしても、過去日（例: 2026-08-19以前）に完了したタスクが不当に除外されてしまいます。

2. **デフォルト `include_completed=false` 指定時の指定日当日完了タスクの過剰含める問題**:
   - `04_tasks.md` L15 では `include_completed` のデフォルト値は `false`（完了タスク除外）と定義されています。
   - しかし L19 では「指定日当日の 00:00:00+09:00 <= due_datetime <= 23:59:59+09:00 の全タスク」と規定されているため、`include_completed` を省略（デフォルト `false`）した状態であっても、指定日当日に締切を持つ完了タスクが抽出対象に含まれてしまうのか、それとも `include_completed=false` により除外されるのかの定義が競合・曖昧になっています。

## 3. 推奨される修正案
`04_tasks.md` L19 の `due_date` パラメータの説明および「リクエスト評価順序」におけるフィルタリングロジックを以下のように修正し、日付範囲による絞り込み（`due_datetime <= 指定日 23:59:59+09:00`）と、ステータスによる絞り込み（`include_completed` または `status`）を分離・統合した明確な定義へ変更してください。

```markdown
| `due_date` | string | × | - | 締切日絞り込み（`YYYY-MM-DD`。指定日当日の `23:59:59+09:00` 以前 `due_datetime <= YYYY-MM-DD 23:59:59+09:00` のタスクを抽出対象とする。締切日時未設定 `null` のタスクは除外。なお、完了タスクの含める/除外制御は `include_completed` または `status` パラメータの指定に従う。未指定時は絞り込みを行わない。※`start_date` / `end_date` との同時指定は不可） |
```
