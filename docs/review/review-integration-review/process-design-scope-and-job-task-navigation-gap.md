# process_design.md における定期ジョブ・タスク機能の全体導線および適用範囲の曖昧性

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [process_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design.md)
  - [process_design/README.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/README.md)
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)

## 1. 問題の概要
`docs/design/process_design.md` は「システムの主要ユースケースにおける処理フローおよびシーケンス図です（Notion最新版と完全同期）」と定義されていますが、リンクされている目次は `01_account_creation.md` 〜 `07_password_change.md` のアカウント・認証系機能のみに限定されており、バックグラウンド定期バッチ処理（`job_design.md`）やタスク管理機能（CRUD・繰り返し一括生成・検索等）の処理フローへの位置づけおよび導線が欠落しています。
これにより、システム全体の処理設計における本ドキュメントの適用範囲や他設計書との関係性が曖昧になっています。

## 2. 詳細な指摘内容
1. **定期ジョブ（パージバッチ）処理フローへの導線欠落**:
   - システム内の定期実行ジョブ（OTP・セッション・ログ・レートリミットのパージ）の処理フローおよびシーケンス図は `docs/design/job_design.md` に詳細化されていますが、`process_design.md` のトップ目次から参照リンクされておらず、処理設計全体の中でジョブ処理がどこに位置づけられているかが把握しづらくなっています。
2. **タスク管理機能の処理フローの位置づけの曖昧性**:
   - `process_design/README.md` では「本ディレクトリが対象とするのは、現在配置されているアカウント作成、アカウント編集、アカウント削除、ログイン、ログアウト、パスワードリセットおよびパスワード変更の処理設計です。タスク管理等、目次に文書が存在しない機能の処理設計は本ディレクトリでは定義しません」と注記されていますが、トップの `process_design.md` にはそのスコープ制約や、タスク管理の処理設計（`api_design/04_tasks.md` や要件定義書等）を参照すべき旨の案内が記載されていません。

## 3. 推奨される修正案
1. **`process_design.md` の更新**:
   - 本ドキュメントの適用範囲（アカウント・認証系ユースケースを詳細化していること）を明記する。
   - 定期バッチジョブ（パージ処理）のシーケンス・処理フローについては `[ジョブ詳細設計書 (job_design.md)](job_design.md)` を参照する旨の関連リンク・導線を追加する。
   - タスク管理の処理フロー（同期APIによる一括生成・検索等）については `[タスクAPI詳細仕様書 (api_design/04_tasks.md)](api_design/04_tasks.md)` および要件定義書を参照する旨の案内を追記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/process_design.md` にアカウント・認証系ユースケースを扱うスコープ注記を追記し、定期バッチジョブ（`job_design.md`）およびタスク管理機能（`api_design/04_tasks.md`）への参照リンク・導線を追加しました。

### 変更したファイル
- [process_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design.md)
