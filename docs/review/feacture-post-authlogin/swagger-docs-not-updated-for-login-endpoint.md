# Swagger ドキュメント（backend/docs）にログインエンドポイントが未反映

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-29 22:27:00
- **Target Files**:
  - [login.go](backend/handler/login.go)
  - [docs.go](backend/docs/docs.go)
  - [swagger.json](backend/docs/swagger.json)
  - [swagger.yaml](backend/docs/swagger.yaml)

## 1. 問題の概要
`backend/handler/login.go` に Swagger アノテーション（`@Summary ログイン`, `@Router /api/auth/login [post]` 等）が追記されていますが、自動生成先である `backend/docs/` 配下のファイル（`docs.go`, `swagger.json`, `swagger.yaml`）が再生成されておらず、Swagger UI や OpenAPI 定義に `POST /api/auth/login` が反映されていません。

## 2. 詳細な指摘内容
1. `backend/handler/login.go` の `LoginHandler` には Swag コメントが適切に付与されています。
2. しかし、`backend/docs/swagger.json` および `backend/docs/docs.go` は前回コミット時のままとなっており、`/api/auth/login` エンドポイントや `LoginRequest`, `LoginResponse` のスキーマ定義が含まれていません。
3. 開発モードでサーバー起動時に `/swagger/index.html` にアクセスしてもログイン API が表示されず、クライアント開発や API ドキュメント検証で齟齬が生じます。

## 3. 推奨される修正案
1. `swag init -g main.go -d ./,./handler,./model --output ./docs` 等を実行して `backend/docs/` 配下のドキュメントファイルを最新化します。
2. 生成された `docs.go`, `swagger.json`, `swagger.yaml` に `/api/auth/login` が含まれていることを確認します。

---

## 修正完了報告

- **Resolved At**: 2026-08-29 22:42:14
- **Status**: Resolved

### 実施した修正内容

- `swag init` を実行し、Swagger ドキュメントを再生成しました。
- `POST /api/auth/login` と `LoginRequest`、`LoginResponse` の各スキーマが生成物に含まれることを確認しました。
- Swagger CLI の実行に必要な依存関係チェックサムを `go.sum` に反映しました。

### 変更したファイル

- [docs.go](backend/docs/docs.go)
- [swagger.json](backend/docs/swagger.json)
- [swagger.yaml](backend/docs/swagger.yaml)
- [go.sum](backend/go.sum)
