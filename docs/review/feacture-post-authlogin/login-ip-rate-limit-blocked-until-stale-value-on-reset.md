# LOGIN_IP_RATE_LIMIT のブロック期限経過・カウンターリセット時における BLOCKED_UNTIL 残存

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-29 22:27:00
- **Target Files**:
  - [login.go](backend/repository/login.go)

## 1. 問題の概要
`backend/repository/login.go` の IP レートリミット失敗回数 UPSERT クエリ（`upsertIPFailureQuery`）において、前回の IP 遮断期間が満了し、かつ前回の失敗日時から 5 分以上経過した後の初回失敗時に、`FAILED_COUNT` は 1 にリセットされますが、`BLOCKED_UNTIL` カラムに過去のブロック満了日時がクリアされずにそのまま保持される挙動となっています。

## 2. 詳細な指摘内容
1. `backend/repository/login.go` の `upsertIPFailureQuery` では以下のように更新処理を行っています。
   ```sql
   BLOCKED_UNTIL = CASE
       WHEN (
           CASE
               WHEN EXCLUDED.LAST_FAILED_AT - LOGIN_IP_RATE_LIMIT.LAST_FAILED_AT > INTERVAL '5 minutes' THEN 1
               ELSE LOGIN_IP_RATE_LIMIT.FAILED_COUNT + 1
           END
       ) >= 30 THEN EXCLUDED.LAST_FAILED_AT + INTERVAL '15 minutes'
       ELSE LOGIN_IP_RATE_LIMIT.BLOCKED_UNTIL
   END
   ```
2. 5分経過によるカウンターリセット（1件目）が発生した場合、`ELSE LOGIN_IP_RATE_LIMIT.BLOCKED_UNTIL` にフォールバックするため、過去にブロックされていた IP の場合、過去の古い日時（期限切れの `BLOCKED_UNTIL`）が DB 上に残り続けます。
3. `selectIPLimit` においては `blockedUntil.Time.After(attempt.Now)` により現在時刻との比較が行われているため、実際の誤遮断（誤った 429 応答）は発生しません。しかし、DB 上のデータ状態として「リセット後にもかかわらず過去の遮断期限が保持されている」状態となり、監査やバッチ処理（クリーンアップジョブやログ集計等）において紛らわしいレコード状態となります。

## 3. 推奨される修正案
1. 失敗カウンターがリセットされる場合（`EXCLUDED.LAST_FAILED_AT - LOGIN_IP_RATE_LIMIT.LAST_FAILED_AT > INTERVAL '5 minutes'`）には、`BLOCKED_UNTIL` を明示的に `NULL` に初期化する条件分岐を追加することを検討してください。
   ```sql
   BLOCKED_UNTIL = CASE
       WHEN (
           CASE
               WHEN EXCLUDED.LAST_FAILED_AT - LOGIN_IP_RATE_LIMIT.LAST_FAILED_AT > INTERVAL '5 minutes' THEN 1
               ELSE LOGIN_IP_RATE_LIMIT.FAILED_COUNT + 1
           END
       ) >= 30 THEN EXCLUDED.LAST_FAILED_AT + INTERVAL '15 minutes'
       WHEN EXCLUDED.LAST_FAILED_AT - LOGIN_IP_RATE_LIMIT.LAST_FAILED_AT > INTERVAL '5 minutes' THEN NULL
       ELSE LOGIN_IP_RATE_LIMIT.BLOCKED_UNTIL

---

## 修正完了報告

- **Resolved At**: 2026-08-29 22:42:14
- **Status**: Resolved

### 実施した修正内容

- 失敗間隔が5分を超えてカウンターをリセットする場合、`BLOCKED_UNTIL` も `NULL` にクリアする条件を UPSERT クエリへ追加しました。
- 30回到達時の新規ブロック設定がリセット条件より優先されること、およびリセット条件で期限がクリアされることを単体テストで検証しました。

### 変更したファイル

- [login.go](backend/repository/login.go)
- [login_test.go](backend/repository/login_test.go)
   END
   ```
2. 単体テストまたは結合テストで、ブロック満了＋5分経過後の初回失敗時に `BLOCKED_UNTIL` がクリアされることを確認するケースを拡充します。
