# tech_stack.md における常駐ジョブスケジューラと PostgreSQL Advisory Lock のアーキテクチャ定義不足

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [tech_stack.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/tech_stack.md)
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)

## 1. 問題の概要
`docs/design/tech_stack.md` のバックエンド技術スタックにおいて、ジョブスケジューラとして `robfig/cron/v3`（`v3.0.1`）が記載されていますが、複数コンテナ・水平スケール構成下における多重起動防止策としての PostgreSQL Advisory Lock（`pg_try_advisory_lock`）や専用接続（`db.Conn`）の活用に関する技術選定の目的・アーキテクチャ連携の記載が不足しています。

## 2. 詳細な指摘内容
1. **プロセス内常駐スケジューラと分散排他制御の連携方針**:
   - `tech_stack.md` では `robfig/cron/v3` の選定理由として「定期パージバッチ（OTP・セッション・ログ・レートリミット削除）の定期実行制御（JST基準・バックエンドプロセス内常駐スケジューラ）」と記載されています。
   - バックエンドプロセス内で常駐実行する場合、単一ノード内では Goroutine スケジューリングで動作しますが、将来的な複数インスタンス展開時には各インスタンスで同一スケジュールでジョブがトリガーされます。
   - その際の排他制御（多重実行スキップ）として PostgreSQL Advisory Lock を活用し、Go の `database/sql` コネクションプール環境下で専用コネクション（`db.Conn`）を排有して制御する旨のアーキテクチャ方針について、技術スタック側にも簡潔な目的・選定理由を記載しておくことで、技術選定の整合性がより明確になります。

## 3. 推奨される修正案
1. **`tech_stack.md` 第3節（バックエンド）の更新**:
   - `robfig/cron/v3` の選定理由・目的欄に、「PostgreSQL Advisory Lock（`pg_try_advisory_lock` / `db.Conn`）と連携した安全な多重起動防止制御」を含む旨を追記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/tech_stack.md` 第3節（バックエンド）の `robfig/cron/v3` 選定理由に、PostgreSQL Advisory Lock（`pg_try_advisory_lock` / `db.Conn`）と連携した安全な多重起動防止制御の目的を追記しました。

### 変更したファイル
- [tech_stack.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/tech_stack.md)
