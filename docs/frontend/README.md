# SyncTask フロントエンド移行計画書 (Migration Plan)

本ディレクトリは、`SyncTask-Design-Idea`（デザインプロトタイプ）から `SyncTask/frontend`（本番リポジトリ）へ Next.js フロントエンドを移植するための詳細計画およびタスク定義書を格納しています。

---

## 📌 移行の基本方針

1. **スコープ**:
   - 今回の移行では **UI・画面遷移・ポップアップ・モック状態（見た目と動作確認）** を最優先とします。
   - バックエンドAPIとの実通信や詳細なビジネスロジックは後回しとし、まずは `npm run dev` で全画面・ダイアログのUIが正常に表示・操作できる状態を構築します。
2. **デザインと構造**:
   - `SyncTask-Design-Idea` のデザイン・スタイル（Tailwind v4 + shadcn/UI + Base UI）を忠実に再現・移行します。
   - 要件定義（`docs/req-def/`）および画面設計書（`docs/design/screen_design.md`）に合致するよう、適宜ディレクトリ構造やコンポーネント構成を整理します。
3. **並列作業体制**:
   - 複数のAIエージェント（`herdr` / `antigravity` 等）による並列作業を前提とし、依存関係のないドメイン機能グループ単位でタスクを分割しています。

---

## 📂 ドキュメント構成

| ファイル / ディレクトリ | 内容説明 |
| --- | --- |
| [00_overview.md](./00_overview.md) | 移行計画全体の概要、設計思想、フェーズ構成、および完了条件 |
| [01_foundation_setup.md](./01_foundation_setup.md) | **Phase 1**: 共通基盤（依存パッケージ、Tailwind/CSS変数、UIコンポーネント、モックStore、共通レイアウト）の構築手順 |
| [tasks/](./tasks/) | **Phase 2**: 各AIエージェントに割り当てる独立した並列タスク指示書 |
| ├ [task_01_auth.md](./tasks/task_01_auth.md) | 並列タスク1: 認証系画面・コンポーネント群（ログイン、登録、OTP、パスワードリセット等） |
| ├ [task_02_home.md](./tasks/task_02_home.md) | 並列タスク2: ホーム画面（優先度・締切・ピン止めタブ、リスト表示、ページネーション） |
| ├ [task_03_tasks_view.md](./tasks/task_03_tasks_view.md) | 並列タスク3: タスク一覧画面（検索・絞り込み・リスト/カレンダー表示切り替え） |
| ├ [task_04_task_dialogs.md](./tasks/task_04_task_dialogs.md) | 並列タスク4: タスク操作ダイアログ群（作成、編集、削除確認、日付詳細ポップアップ） |
| └ [task_05_profile.md](./tasks/task_05_profile.md) | 並列タスク5: プロフィール・アカウント管理画面群（表示、編集、メール変更OTP、パスワード変更、アカウント削除） |
| [02_future_implementation_and_tests.md](./02_future_implementation_and_tests.md) | **Phase 3以降**: 今後の機能追加（API連携・状態管理）およびテスト実装計画とToDoリスト |

---

## 🤖 herdr を用いた並列エージェント実行ガイド

本移行作業の Phase 2（各タスクの並列実行）は、`herdr` CLI を使用して複数エージェントを立ち上げて実施することを推奨します。

### 実行手順の例

```bash
# 1. 共通基盤 (Phase 1) の完了を確認
# Phase 1 は単一エージェントで実施し、npm run dev / npm run build が通ることを確認します。

# 2. 並列ワーカーの起動
# 各タスク指示書をプロンプトとして herdr agent を起動
herdr tab new --title "worker-auth"
herdr agent run --tab "worker-auth" --prompt "docs/frontend/migration-from-design-repo/tasks/task_01_auth.md に従って認証系画面を frontend/ に移植してください。"

herdr tab new --title "worker-home"
herdr agent run --tab "worker-home" --prompt "docs/frontend/migration-from-design-repo/tasks/task_02_home.md に従ってホーム画面を frontend/ に移植してください。"

herdr tab new --title "worker-tasks-view"
herdr agent run --tab "worker-tasks-view" --prompt "docs/frontend/migration-from-design-repo/tasks/task_03_tasks_view.md に従ってタスク一覧画面を frontend/ に移植してください。"

herdr tab new --title "worker-dialogs"
herdr agent run --tab "worker-dialogs" --prompt "docs/frontend/migration-from-design-repo/tasks/task_04_task_dialogs.md に従ってタスクダイアログ群を frontend/ に移植してください。"

herdr tab new --title "worker-profile"
herdr agent run --tab "worker-profile" --prompt "docs/frontend/migration-from-design-repo/tasks/task_05_profile.md に従ってプロフィール画面群を frontend/ に移植してください。"

# 3. 全タスクの完了待機と動作検証
# npm run dev を起動し、ブラウザ上で各ページ・モーダルの表示確認を行います。
```
