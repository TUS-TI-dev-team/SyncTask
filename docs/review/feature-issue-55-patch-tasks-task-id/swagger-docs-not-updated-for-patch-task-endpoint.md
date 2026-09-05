# Swagger ドキュメント（backend/docs）にタスク部分更新エンドポイントが未反映

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-09-05 21:47:00
- **Target Files**:
  - [task_patch.go](backend/handler/task_patch.go)
  - [docs.go](backend/docs/docs.go)
  - [swagger.json](backend/docs/swagger.json)
  - [swagger.yaml](backend/docs/swagger.yaml)

## 1. 問題の概要
`backend/handler/task_patch.go` に Swagger アノテーション（`@Summary タスク部分更新`, `@Router /api/tasks/{task_id} [patch]` 等）が記述されていますが、自動生成先である `backend/docs/` 配下のファイル（`docs.go`, `swagger.json`, `swagger.yaml`）が再生成されておらず、Swagger UI や OpenAPI 定義に `PATCH /api/tasks/{task_id}` が反映されていません。

## 2. 詳細な指摘内容
1. `backend/handler/task_patch.go` の `PatchTaskHandler` には Swag コメント（L17-L29）が付与されています。
2. しかし、`backend/docs/swagger.json` および `backend/docs/docs.go` は更新されておらず、`/api/tasks/{task_id}` の PATCH メソッドや `PatchTaskRequest`, `PatchTaskResponse` のモデル定義が含まれていません。
3. これにより、Swagger UI（`/swagger/index.html`）上で本エンドポイントの仕様確認やテスト実行ができず、API ドキュメントと実装に不整合が生じています。

## 3. 推奨される修正案
1. `swag init -g main.go -d ./,./handler,./model --output ./docs` を実行し、`backend/docs/` 配下のファイルを再生成・最新化してください。
2. 生成された `docs.go`, `swagger.json`, `swagger.yaml` に `PATCH /api/tasks/{task_id}` および関連モデルが含まれていることを確認してください。

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:51:00
- **Status**: Resolved

### 実施した修正内容
- `backend` ディレクトリにおいて `go run github.com/swaggo/swag/cmd/swag init -g main.go -d ./,./handler,./model --output ./docs` を実行し、Swagger ドキュメントファイルを再生成しました。
- 再生成された `backend/docs/docs.go`、`backend/docs/swagger.json`、`backend/docs/swagger.yaml` に `PATCH /api/tasks/{task_id}` のパス定義、および `model.PatchTaskRequest`, `model.PatchTaskResponse` のスキーマ定義が正しく反映されていることを確認しました。

### 変更したファイル
- [docs.go](backend/docs/docs.go)
- [swagger.json](backend/docs/swagger.json)
- [swagger.yaml](backend/docs/swagger.yaml)
