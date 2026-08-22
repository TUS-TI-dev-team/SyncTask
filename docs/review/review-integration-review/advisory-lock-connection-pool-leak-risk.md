# Go コネクションプール環境下における PostgreSQL Advisory Lock 保持・解放仕様の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
  - [tech_stack.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/tech_stack.md)

## 1. 問題の概要
`job_design.md` では、定期クリーンアップジョブの多重起動防止・排他制御として PostgreSQL のセッションレベル・アドバイザリロック（`pg_try_advisory_lock` / `pg_advisory_unlock`）を採用し、ジョブ内でチャンク分割ループ（チャンクごとに個別のトランザクションを BEGIN/COMMIT）を実行する設計となっています。
しかし、バックエンドに採用されている Go 言語（`database/sql`）はコネクションプーリングを行うため、単に `sql.DB` から直接クエリを発行すると、ロック取得、各チャンク削除、ロック解放がそれぞれ別々の物理 DB コネクションで実行されてしまう致命的な不具合（アンロック失敗によるコネクションプール内でのロック残留・他インスタンスの永久ブロック）が発生します。

## 2. 詳細な指摘内容
1. **セッションレベル・アドバイザリロックのスコープ**:
   PostgreSQL の `pg_try_advisory_lock(int64)` は「同一の物理 DB 接続（セッション）」に紐づくロックです。
2. **Go の `database/sql` コネクションプールの挙動**:
   `db.Exec(...)` や `db.QueryRow(...)` を直接呼び出すと、プールから空いている任意のコネクションが一時的に払い出され、実行後に即座にプールへ返却されます。
3. **発生する不具合**:
   - `SELECT pg_try_advisory_lock(1002)` がコネクション A で実行され、ロックがコネクション A に紐づきます。
   - ループ内の `DELETE` 処理がコネクション B や C で実行されます。
   - 終了時の `SELECT pg_advisory_unlock(1002)` がコネクション D で実行された場合、コネクション D はロックを保持していないため `false` を返し、ロックは解放されません。
   - ロックを保持したままのコネクション A がプールに戻ると、次回以降の同一ジョブや他インスタンスからの実行がロック取得失敗で永久にスキップされ続ける（またはコネクション破棄まで多重実行が遮断される）障害を引き起こします。

## 3. 推奨される修正案
`job_design.md` の「2-1. 多重起動防止・排他制御（Advisory Lock）」および「4. 処理フロー & シーケンス図」に、Go 実装におけるコネクション固定要件を明記してください。

1. **`db.Conn(ctx)` による専用コネクションの排有取得**:
   ジョブ開始時に `conn, err := db.Conn(ctx)` でプールから専用の単一物理コネクションを確保し、ジョブ終了までこの `conn` インスタンスを排有する。
2. **同一コネクション上での一連処理の実行**:
   - `conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockKey)`
   - ロック獲得成功時: 同一 `conn` 上でチャンクごとのトランザクション開始（`conn.BeginTx`）・コミット、およびアンロック（`pg_advisory_unlock`）を実行する。
3. **`defer conn.Close()` による確実なリソース解放**:
   アンロック処理とコネクションのプール返却（`conn.Close()`）を `defer` / `finally` で確実に実行する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
Go 言語の `database/sql` コネクションプール環境において、ジョブ開始時に `db.Conn(ctx)` で専用コネクションを排有取得し、同一 `conn` インスタンス上でアドバイザリロックの取得・各チャンクトランザクション実行・アンロック・`defer conn.Close()` を行う要件を明記しました。

### 変更したファイル
- [job_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\job_design.md)
