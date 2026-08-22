# 設計ファイル全体 結合レビュー総合結果サマリ (Integration Review Summary)

- **Status**: Passed (全指摘解消・整合性検証完了)
- **Reviewed At**: 2026-08-22
- **Branch**: `review/integration-review`
- **Target Files**:
  - [tech_stack.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/tech_stack.md)
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
  - [api_design/01_overview.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/01_overview.md)
  - [api_design/02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [api_design/03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md)
  - [api_design/04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md)
  - [process_design/README.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/README.md)
  - [process_design/01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [process_design/02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md)
  - [process_design/03_account_delete.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/03_account_delete.md)
  - [process_design/04_login.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/04_login.md)
  - [process_design/05_logout.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/05_logout.md)
  - [process_design/06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [process_design/07_password_change.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/07_password_change.md)
  - [req-def/requirements/01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md)
  - [req-def/requirements/02_task_management.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/02_task_management.md)
  - [req-def/requirements/03_security_and_system.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/03_security_and_system.md)
  - [req-def/requirements/04_non_functional.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/04_non_functional.md)

---

## 1. 結合査読の実施概要

`docs/design/` 配下の全設計ドキュメントおよび上位仕様である `docs/req-def/` の要件定義書を対象に、4つの機能ドメイン別の子エージェント（Worker）による並列査読（`herdr-review-loop`）および対話型ヒアリング（`/grill-me`）を実施しました。

### 査読対象ドメインと担当
1. **認証・ユーザー管理ドメイン**: 会員登録、ログイン、ログアウト、メール変更、パスワードリセット、パスワード変更、アカウント削除における画面・API・処理・DBの整合性
2. **タスク管理・検索・カレンダー・期限ドメイン**: タスクCRUD、ソート、フィルタ、ページネーション、カレンダー期間取得、バリデーション、DBカラム型の整合性
3. **定期タスク・自動生成・バッチ・ジョブドメイン**: バッチ削除クエリ（CTE）、定期タスク自動生成（同期API即時生成）とCronジョブの責務境界
4. **アーキテクチャ・共通基盤・セキュリティ・非機能要件ドメイン**: CORS仕様、共通セキュリティレスポンスヘッダー、Retry-Afterヘッダー、Cookieセキュリティ属性、DB CHECK制約およびインデックス設計

---

## 2. 対話ヒアリング (/grill-me) で合意・確定した設計方針

| No | 設計論点 | 合意・確定した仕様内容 |
| :--- | :--- | :--- |
| 1 | **タスク検索における日本語同一視** | アプリケーション層で正規化（小文字化＋NFKC正規化＋ひらがなカタカナ統一）を行い、`TASKS.SEARCH_TEXT` カラムと `pg_trgm` GINインデックスによる部分一致検索を採用。 |
| 2 | **ホーム画面のタスク一覧レイアウト** | 「優先タスク」「締切間近タスク」「ピン留めタスク」を切り替える**タブ切り替えUI**を採用し、選択中のタブのタスクを単一のページネーションUI（上下配置、20件単位）で表示。 |

---

## 3. 指摘事項と対応結果一覧（全18件 Resolved）

| No | 指摘ファイル | 概要と対応結果 |
| :--- | :--- | :--- |
| 1 | [account-lockout-429-user-enumeration-vulnerability.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/account-lockout-429-user-enumeration-vulnerability.md) | ロックアウト・レート制限時のユーザー列挙防止のため、未登録ユーザーに対しても一貫して同一の429レスポンスを返却する仕様を明記。 |
| 2 | [advisory-lock-connection-pool-leak-risk.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/advisory-lock-connection-pool-leak-risk.md) | PostgreSQL アドバイザリロック取得時にコネクションプール枯渇を防ぐため、専用コネクション排有取得と確実なアンロック/クローズ手順を明記。 |
| 3 | [batch-retry-scope-and-lock-release-ambiguity.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/batch-retry-scope-and-lock-release-ambiguity.md) | バッチ処理の一時エラーリトライスコープをチャンク単位（最大3回）と定義し、恒久エラー時の即時中断とdeferアンロックを明記。 |
| 4 | [calendar-date-detail-to-task-create-modal-transition-gap.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/calendar-date-detail-to-task-create-modal-transition-gap.md) | カレンダー日付詳細ポップアップからタスク作成を開く際、親ポップアップを閉じ締切日時に該当日付+23:59をプリフィルする仕様を明記。 |
| 5 | [cors-missing-access-control-expose-headers-retry-after.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/cors-missing-access-control-expose-headers-retry-after.md) | CORS仕様に `Access-Control-Expose-Headers: Retry-After` を追加し、フロントエンドから429待機秒数を読み取り可能に改善。 |
| 6 | [home-screen-multi-view-pagination-and-api-ambiguity.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/home-screen-multi-view-pagination-and-api-ambiguity.md) | ホーム画面をタブ切り替えUIとし、選択中タブのAPI呼び出しと単一ページネーションUI制御仕様を明記。 |
| 7 | [job-design-document-reference-path-inconsistency.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/job-design-document-reference-path-inconsistency.md) | job_design.md 内の参照先パスを分割ディレクトリ構成に合わせて更新。 |
| 8 | [mail-auth-log-missing-cancel-event-type-in-db-and-process-design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/mail-auth-log-missing-cancel-event-type-in-db-and-process-design.md) | `MAIL_AUTH_LOG.EVENT_TYPE` に `CANCELLED` を追加し、キャンセル時のログ記録仕様を整合。 |
| 9 | [otp-api-cooldown-seconds-response-schema-inconsistency.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/otp-api-cooldown-seconds-response-schema-inconsistency.md) | 全OTP発行・再送APIレスポンスに `cooldown_seconds: 60` を共通定義。 |
| 10 | [otp-session-cleanup-query-index-mismatch.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/otp-session-cleanup-query-index-mismatch.md) | `OTP_SESSION` のパージ用インデックスをクエリ条件（STATUS, OTP_EXPIRES_AT）に最適化。 |
| 11 | [otp-session-send-failed-count-check-constraint-risk.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/otp-session-send-failed-count-check-constraint-risk.md) | `OTP_SESSION.SEND_FAILED_COUNT` のCHECK制約を `BETWEEN 0 AND 5` に緩和。 |
| 12 | [password-allowed-symbols-inconsistency-between-req-and-design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/password-allowed-symbols-inconsistency-between-req-and-design.md) | 要件定義書のパスワード記号定義をASCII印字可能半角記号全32種類に同期。 |
| 13 | [password-reset-completion-missing-active-otp-sessions-purge.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/password-reset-completion-missing-active-otp-sessions-purge.md) | パスワードリセット完了時に対象ユーザーの全OTPセッションおよび全ログインセッションを一括破棄する仕様を明記。 |
| 14 | [process-design-missing-otp-session-cancel-and-screen-design-inconsistency.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/process-design-missing-otp-session-cancel-and-screen-design-inconsistency.md) | 処理設計書に `POST /api/auth/otp-session/cancel` フローを追加し、画面の戻る・キャンセル記述を統一。 |
| 15 | [reauth-failed-count-time-interval-reset-specification-omission.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/reauth-failed-count-time-interval-reset-specification-omission.md) | 再認証失敗カウンターについて直前失敗から15分経過で 0 にリセットする仕様を明記。 |
| 16 | [recurring-task-generation-boundary-and-responsibility.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/recurring-task-generation-boundary-and-responsibility.md) | 繰り返しタスクはAPIでの同期即時一括生成方式を採用し、定期バッチは生成を行わない境界設計を明記。 |
| 17 | [task-detail-screen-definition-and-api-mapping-gap.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/task-detail-screen-definition-and-api-mapping-gap.md) | 単体詳細専用画面は設けず「タスク編集ポップアップ」が閲覧兼編集の役割を担う仕様を明記。 |
| 18 | [task-search-japanese-normalization-and-db-design-gap.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/review/review-integration-review/task-search-japanese-normalization-and-db-design-gap.md) | 日本語同一視検索のため、アプリ側正規化文字列を `TASKS.SEARCH_TEXT` に格納し `pg_trgm` GINインデックスで検索する設計を反映。 |

---

## 4. テスト検証結果

- **Backend テスト (Go)**: `go test ./...` -> **All PASS**
- **Frontend 単体テスト (Vitest)**: `npm run test:unit` -> **14 passed**
- **Frontend E2Eテスト (Playwright)**: `npm run test:e2e` -> **3 passed**
