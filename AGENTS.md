# SyncTask 開発エージェント・オーケストレーション規約

このリポジトリは、タスク管理 Web アプリケーション **SyncTask** のコードベースです。
バックエンドは Go / Gin / PostgreSQL / Supabase、フロントエンドは Next.js / TypeScript で構築されています。

---

## 🛠 コマンドと開発ルール

### バックエンド (Go)
- **テスト実行**:
  ```bash
  cd backend
  go test ./...         # 全単体テスト実行
  go test -v ./...      # 詳細表示
  go test -cover ./...   # カバレッジ測定
  ```
- **テスト規約**: `backend/TESTING_GUIDE.md` を厳守すること。
  - テスト関数・サブテスト名は日本語で「`<分類>: <期待される結果>`」とする。
  - Code-as-Docs 原則に従い、`@spec` を GoDoc とテストケース名に明記する。
  - 即時中断は `require`、値検証は `assert` を使い分ける。

### 仕様書・設計書
- API 設計: `docs/design/api_design/`
- DB 設計: `docs/design/database_design.md`
- 業務設計: `docs/design/process_design/`
- 実装計画書: `docs/plans/backend/`
- 自動化設計書: `docs/plans/be-dev-automation-by-gemini.md`
