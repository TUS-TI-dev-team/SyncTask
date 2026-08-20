# SyncTask Backend

SyncTask のバックエンド API サーバーです。Go および Gin Web Framework を使用して構築されています。

## ディレクトリ構成

```
backend/
├── docs/             # swaggo により自動生成される API 仕様書 (docs.go, swagger.json, swagger.yaml)
├── handler/          # API ハンドラー (ビジネスロジック・エンドポイント処理)
├── router/           # ルーティング定義
├── main.go           # サーバーエントリーポイント・API メタ情報
├── go.mod / go.sum   # 依存パッケージ定義
├── Dockerfile        # コンテナ実行用 Dockerfile
├── .air.toml         # ホットリロード設定
└── README.md         # 本ドキュメント
```

---

## テスト

### テスト実行方法

バックエンドディレクトリ (`backend/`) で以下のコマンドを実行します。

```bash
# 全テストを実行
go test ./...

# 詳細な出力付きで実行
go test -v ./...

# カバレッジを測定して実行
go test -cover ./...
```

> **Note**: Go のインストールパスが環境変数 `PATH` に通っていない場合は、フルパス（例: `"C:\Program Files\Go\bin\go.exe" test -v ./...`）で実行するか、Go のインストールディレクトリを `PATH` に追加してください。

---

## API 仕様書 (Swagger)

`swaggo/gin-swagger` を用いた OpenAPI (Swagger) 準拠の仕様書生成と Web 表示に対応しています。

### API 仕様書の自動生成・更新
コード内の GoDoc アノテーションを更新した場合、以下のコマンドで `docs/` 配下の仕様書ファイルを再生成できます。

```bash
go run github.com/swaggo/swag/cmd/swag@latest init
```

### Swagger UI の確認
サーバーを開発モードで起動後、ブラウザで以下の URL にアクセスします。

- **URL**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## サーバー起動方法

### ローカルでの直接起動
```bash
go run main.go
```

### Docker Compose での起動 (プロジェクトルート)
```bash
docker compose up backend
```
