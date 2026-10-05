# SyncTask Backend

SyncTask のバックエンド API サーバーです。Go および Gin Web Framework を使用して構築されています。

## ディレクトリ構成

```
backend/
├── cmd/seed/           # 開発用シードデータ投入コマンド (main.go)
├── db/
│   ├── migrations/     # DB マイグレーション SQL
│   └── seed/           # シードデータ定義・投入ロジック
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

テストコードを作成するときは、必ず **[テストコード作成ガイド](./TESTING_GUIDE.md)** を参照してください。

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

## シードデータ (開発用)

開発時の動作テスト用に、`LOGIN_ACCOUNT` / `TASK` / `LOGIN_SESSION` のテストデータを投入するコマンドを用意しています。

### 実行方法

`backend/` ディレクトリで以下を実行します (DB 接続先は環境変数 `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` / `DB_SSLMODE` で指定。未指定時はローカルの `synctask_dev` に接続)。

```bash
# 既存データを残し、不足分のみ投入 (冪等。何度実行しても OK)
go run ./cmd/seed
```

### 全削除 & 再投入

投入済みのシードデータを一旦すべて削除して作り直したい場合は `--force` を付けます。
`TASK` / `LOGIN_SESSION` / `LOGIN_ACCOUNT` の 3 テーブルが全件削除され、その後シードデータが再投入されます。
**開発用のコマンドなので、本番環境や共有 DB に対しては絶対に実行しないでください。**

```bash
go run ./cmd/seed --force
```

### 投入されるデータ

| メールアドレス | パスワード | 内容 |
|---|---|---|
| `alice@example.com` | `Password123!` | タスク 5 件 (高優先度・ピン留め・完了済みを含む) + 有効なログインセッション 1 件 |
| `bob@example.com` | `Password123!` | タスク 3 件 (日本語タイトルを含む) + 有効なログインセッション 1 件 |
| `carol@example.com` | `Password123!` | タスク 1 件 (期限切れ) |

- パスワードはすべて `Password123!` で、投入時に bcrypt (`bcrypt.DefaultCost`) でハッシュ化されます。
- ログインセッションは 30 日有効な固定 ID (`seed-session-alice-...` / `seed-session-bob-...`) を発行するため、ログイン処理を省略して API の動作テストが可能です。
- タスクの `SEARCH_TEXT` は `util.NormalizeSearchText` により自動生成されます。
- データ定義は `db/seed/seed.go` にあるので、追加・変更はそこを編集してください。

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
