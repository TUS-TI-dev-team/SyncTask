# 大量データ削除時のチャンク分割とトランザクション境界・ロック競合防止設計の不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:10:30
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
[job_design.md](docs/design/job_design.md) において、各クリーンアップジョブが単一の `DELETE FROM ...` クエリで全件一括削除を行う設計となっています。大量レコードが存在する場合に長時間ロックの保持、WAL（Write-Ahead Logging）の急増、APIリクエスト（ユーザーの認証やログインセッション更新等）とのロック競合・デッドロック、およびDBタイムアウトを招く危険性があります。

## 2. 詳細な指摘内容
1. **一括 `DELETE` によるロック競合・デッドロックリスク**:
   - `job_design.md` のクリーンアップSQL（例: lines 27-30, lines 40-42）では、対象件数の上限（LIMIT）やチャンク分割の考慮が一切なく、1つのSQLで条件に一致する全レコードを一括削除する方針になっています。
   - `LOGIN_SESSION` やログテーブル、`OTP_SESSION` で削除対象件数が数千〜数万件規模に達した場合、大量の行ロックが長時間保持され、通常のオンラインAPI（ログイン認証、セッション検証、OTP発行等）で発生する `UPDATE` / `DELETE` / `INSERT` と競合し、ロック待ちタイムアウトやデッドロックを誘発します。
2. **トランザクション境界とコミット単位の未定義**:
   - 削除処理が1トランザクションで実行されるのか、チャンク単位でコミットされるのかが明記されていません。1つの巨大トランザクションで実行すると、万が一途中でタイムアウトや通信エラーが発生した場合にすべての削除処理がロールバックされ、リトライ時にも同じタイムアウトを繰り返す恐れがあります。

## 3. 推奨される修正案
1. **チャンク単位でのバッチ削除方式の採用**:
   - ジョブ設計書に「バッチ削除方式（チャンク処理）」の共通設計を追加してください。
   - 例: 1回あたりの削除件数（`BATCH_SIZE`、例: 500〜1,000件）を定め、削除対象が0件になるまでループ処理し、各チャンクごとにコミット（トランザクション分離）を行う方式を規定する。
   ```sql
   -- チャンク削除の例（PostgreSQL CTE または サブクエリ LIMIT による削除）
   DELETE FROM LOGIN_SESSION
   WHERE SESSION_ID IN (
       SELECT SESSION_ID FROM LOGIN_SESSION
       WHERE EXPIRES_AT < NOW()
       LIMIT 1000
   );
   ```
2. **負荷軽減インターバル（スリープ）の考慮**:
   - チャンク間で微小なウェイト（例: 50〜100ms）を設けることで、レプリケーション遅延やCPU/I/O負荷集中を抑制する方針を明記する。
3. **設定パラメータの追加**:
   - 第5章 環境変数一覧に `JOB_CLEANUP_BATCH_SIZE`（デフォルト: 1000）等を追加する。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- `job_design.md` の「2-2. チャンク分割バッチ削除方式」に全ジョブ共通のバッチ削除仕様を策定しました。
- 1回あたり上限1,000件（`JOB_CLEANUP_BATCH_SIZE`）でのCTEによるLIMIT付き削除、各チャンクごとの独立したトランザクションコミット、削除件数0件までのループ処理、およびチャンク間待機（`JOB_CLEANUP_INTERVAL_MS=50ms`）による負荷軽減設計を反映しました。
- 各個別ジョブのクリーンアップSQLおよびシーケンス図にもチャンク削除パターンを適用しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
