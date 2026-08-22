# 定期タスク（繰り返しタスク）生成における同期API即時生成とバックグラウンドジョブの責務境界・設計方針の明確化

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md)
  - [02_task_management.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/02_task_management.md)

## 1. 問題の概要
要件定義書（`02_task_management.md` 2.2）およびタスクAPI設計書（`04_tasks.md` 3.3.2）では、繰り返しタスクは「Cron定期バッチによる未来タスク生成（定期ルールマスター保持）」ではなく、「同期API（`POST tasks`、`is_recurring: true`）による指定期間内（最大1年・最大100件）の即時一括生成」として設計されています。
しかし、ジョブ設計書（`job_design.md`）およびデータベース設計書（`database_design.md`）において、「繰り返しタスク生成ジョブや定期タスク定義テーブル（`RECURRING_TASKS` 等）が存在しない設計判断の根拠（同期APIによる即時一括生成方式を採用し、生成後は個別独立した `TASK` として管理するため）」が明記されていません。
このため、後続の実装者や運用者が「定期タスク自動生成バッチの実装漏れではないか？」と誤認したり、二重生成防止や責務分担についての疑問が生じるリスクがあります。

## 2. 詳細な指摘内容
1. **要件定義・API・DBの一貫性**:
   - `02_task_management.md` 2.2.2: 「生成後の個別独立性: 生成されたタスクはそれぞれ独立した通常タスクとしてDBに個別登録される。以降は個別に編集・状態変更・削除が可能であり、繰り返しタスク間の一括連動更新機能は提供しない。」
   - `04_tasks.md` 3.3.2: `POST tasks` で `is_recurring: true` 時に最大100件のタスクを即時一括生成して返却。
   - `database_design.md`: ルールテーブル（`RECURRING_TASKS`）は持たず、`TASK` テーブルに直接レコードが登録される。
2. **ジョブ設計書・DB設計書での境界説明の不足**:
   - `job_design.md` にはクリーンアップジョブのみが記載されており、タスク生成系のジョブについての言及が一切ありません。
   - バックグラウンドで定期タスクを生成・更新する責務が存在しない旨（同期API即時生成への責務集約）が設計書上に明示されていないため、ドメイン結合レビュー時に仕様の曖昧さが生じる原因となっています。

## 3. 推奨される修正案
1. **`job_design.md` への設計方針の追記**:
   - 第1章冒頭または補足節に、以下のような「定期タスク生成に関する設計方針」を明記してください：
     > **【定期タスク生成に関する設計方針】**  
     > 本システムでは、繰り返しタスクの生成についてバックグラウンド定期バッチによる未来タスク生成方式を採用せず、タスク作成API（`POST /api/v1/tasks`）における**同期即時一括生成方式**（期間・曜日に基づき最大100件を即時作成し、単一DBトランザクションで通常タスクとしてコミット）を採用しています。生成された各タスクは独立した通常レコードとして `TASK` テーブルに登録・管理されるため、Cronスケジューラにおける定期タスク生成ジョブおよび二重生成防止の定期排他バッチはスコープ外（不要）です。
2. **`database_design.md` への補足追記**:
   - `TASK` テーブルのセクションに、定期タスクはルールマスターテーブル（`RECURRING_TASKS` 等）を持たず、即時一括生成された独立レコードとして永続化される旨のノートを追記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
繰り返しタスクはバックグラウンドバッチによる未来生成方式ではなく、タスク作成API（`POST /api/tasks`）における「同期即時一括生成方式（最大100件）」を採用し、独立した通常 `TASK` レコードとして管理する旨の設計方針を `job_design.md` および `database_design.md` に追記しました。

### 変更したファイル
- [job_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\job_design.md)
- [database_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\database_design.md)
