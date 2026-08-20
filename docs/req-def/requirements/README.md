# SyncTask 要件定義書 (Requirements Definition)

本ディレクトリは、タスク管理Webアプリケーション「**SyncTask**」の要件定義書を格納しています。
システム全体の概要、各機能要件、セキュリティおよびシステム基盤要件、非機能要件を分野・観点別に構造化して管理しています。

---

## 構成一覧

| ファイル | タイトル | 概要 |
| :--- | :--- | :--- |
| [00_overview.md](docs/req-def/requirements/00_overview.md) | システム概要・業務要件 | プロジェクト背景、目的、システム概要、利用者の範囲、期待する成果 |
| [01_account_and_auth.md](docs/req-def/requirements/01_account_and_auth.md) | アカウント・認証機能要件 | ユーザー名・パスワード・メール要件、アカウント登録/編集/削除、ログイン/ログアウト、PWリセット |
| [02_task_management.md](docs/req-def/requirements/02_task_management.md) | タスク管理機能要件 | タスクデータ仕様、タスクCRUD、繰り返しタスク一括生成、各種表示・検索・カレンダー表示 |
| [03_security_and_system.md](docs/req-def/requirements/03_security_and_system.md) | セキュリティ・システム共通要件 | OTP仕様、ロックアウト/レートリミット、セッション管理、ログ管理・保持ポリシー、脆弱性対策 |
| [04_non_functional.md](docs/req-def/requirements/04_non_functional.md) | 非機能要件 | 性能（応答速度）、可用性、対応プラットフォーム・ブラウザ環境、運用・保守性方針 |

---

## ドキュメント運用方針

1. **整合性の維持**: 各要件間の依存関係や共通仕様（バリデーションやセキュリティ方針等）に変更が生じた場合は、関連する各ドキュメントを同期して更新してください。
