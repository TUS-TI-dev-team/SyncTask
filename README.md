# SyncTask
練習用タスク管理Webアプリ

## 起動方法

### 開発環境（フロントエンド + バックエンド + PostgreSQL）

```bash
docker compose up --build
```

- フロントエンド: http://localhost:3000
- バックエンド: http://localhost:8080
- ヘルスチェック: http://localhost:8080/health-check
- Swagger UI: http://localhost:8080/swagger/index.html

### テスト用DBでの起動

```bash
DB_NAME=synctask_e2e docker compose up --build
```

## Skills インストール

```bash
npx skills experimental_install
```

## AI バックエンド自動実装オーケストレーション

Herdr および Antigravity CLI を用いた 3 階層マルチエージェントにより、未実装エンドポイントの Issue 取得から実装計画策定、TDD 実装、PR 作成、レビュー修正ループ、マージまでを完全自律実行できます。

### 起動コマンド
メイン Workspace (`w1`) で以下を実行します：

```bash
# Antigravity CLI のプロンプトで実行
/orchestrate-backend
```

### 人間向け運用ガイド
- **詳細設計書**: `docs/plans/be-dev-automation-by-gemini.md`
- **対話ポイント**:
  - 実装・レビュー中の質問：該当ワーカーの Herdr タブに直接回答
  - PR マージ承認：最上位司令塔がモーダル（`ask_question`）で提示
