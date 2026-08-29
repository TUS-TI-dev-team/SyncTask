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
