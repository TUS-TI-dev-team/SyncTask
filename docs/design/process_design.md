# Process Design (処理設計)

## 概要

本ディレクトリおよび配下のドキュメントでは、SyncTask システムにおける**アカウントおよび認証系ユースケース**の詳細処理フローおよびシーケンス図を定義しています（Notion最新版と完全同期）。

各処理の詳細は以下のドキュメントを参照してください。

### アカウント・認証系ユースケース一覧
- [1. アカウント作成](process_design/01_account_creation.md)
- [2. アカウント編集](process_design/02_account_edit.md)
- [3. アカウント削除](process_design/03_account_delete.md)
- [4. ログイン](process_design/04_login.md)
- [5. ログアウト](process_design/05_logout.md)
- [6. パスワードリセット](process_design/06_password_reset.md)
- [7. パスワード変更](process_design/07_password_change.md)

### その他の処理フロー・関連設計書への導線
- **定期実行バッチ・パージ処理**:
  - セッション、OTP、ログ、レートリミット等の定期パージ処理フローおよび排他制御シーケンスについては、[ジョブ詳細設計書 (job_design.md)](job_design.md) を参照してください。
- **タスク管理処理フロー**:
  - タスクのCRUD操作、毎週繰り返しタスクの同期即時一括生成（最大100件）、日本語同一視検索等の処理仕様・シーケンスについては、[タスクAPI詳細仕様書 (api_design/04_tasks.md)](api_design/04_tasks.md) および [タスク管理要件定義書 (req-def/requirements/02_task_management.md)](../req-def/requirements/02_task_management.md) を参照してください。


