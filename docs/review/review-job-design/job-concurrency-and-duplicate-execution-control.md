# ジョブ多重起動および重複実行防止（排他制御）設計の不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:10:30
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)

## 1. 問題の概要
[job_design.md](docs/design/job_design.md) において、ジョブの多重起動（重複実行）を防止するための排他制御機構（PostgreSQL Advisory Lock、分散ロック、実行フラグ等）が規定されていません。マルチインスタンス・コンテナ構成や処理遅延時に同一ジョブが多重起動された場合の競合対策が不足しています。

## 2. 詳細な指摘内容
- クリーンアップSQL自体は冪等性（Idempotency）を持つものの、バッチ処理に時間がかかり次のCron起動間隔と重複した場合や、Web/Workerサーバーが複数台稼働（マルチコンテナ・水平スケール）している環境では、同一のクリーンアップジョブが同時に発火する可能性があります。
- 排他ロック制御がない場合、同一レコードに対する複数ワーカーからの同時 `DELETE` が発生し、無駄なDBリソース消費やロック競合、最悪の場合はデッドロック（Deadlock）を引き起こす要因となります。

## 3. 推奨される修正案
1. **Advisory Lock または分散ロックによる重複実行スキップ機構の追加**:
   - PostgreSQLのAdvisory Lock（例: `pg_try_advisory_lock(job_id_hash)`）等を利用し、既に該当ジョブが他インスタンスまたは同一プロセスで実行中である場合は、待機せず即座にスキップして正常終了（またはスキップログ出力）とする設計を追加してください。
   ```sql
   -- Advisory Lock による多重起動防止の例
   SELECT pg_try_advisory_lock(:job_lock_id);
   -- 取得できた場合のみバッチ削除を実行し、終了時に pg_advisory_unlock(:job_lock_id) を実行
   ```
2. **シーケンス図およびエラー・ログへの反映**:
   - ジョブ開始時にロック取得を試行し、ロックが取得できなかった場合はスキップする処理フローをシーケンス図およびログ仕様（例: `[INFO] Job CLEANUP_OTP_SESSIONS is already running. Skipped.`）に追記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- `job_design.md` の第2章「2-1. 多重起動防止・排他制御（Advisory Lock）」に PostgreSQL の `pg_try_advisory_lock` を用いた共通排他制御機構を策定しました。
- ジョブごとに固有のロックキー識別子を定義し、ロック未取得時は待機せずスキップログを出力して即座に正常終了する仕様としました。
- 第4章のシーケンス図にロック取得・解放のフローを反映しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
