# ローカルDB導入およびE2Eテスト環境構築 計画書

本ドキュメントでは、以下の2タスクに関する設計・実装計画を定義します。

1. **Task 1**: ローカル開発・テスト用 PostgreSQL（Docker）の導入
2. **Task 2**: フロントエンド ↔ バックエンド ↔ DB が連携するE2Eテスト環境の構築

---

## 1. 現状分析

### 1.1 バックエンド（Go / Gin）

- 最小スケルトン状態（`GET /` の Hello World + 開発モード限定の `/health-check` + Swagger UI のみ）。
- **DBドライバなし**: `go.mod` に `pgx` / `database/sql` / `gorm` 等のDB関連依存が一切存在しない。
- **環境変数読み込みなし**: `os.Getenv` / `godotenv` / `viper` 等の設定読み込みコードがない。参照されている環境変数は `GIN_MODE`（docker-compose.yml で設定）のみ。
- **マイグレーションなし**: `.sql` ファイル、`migrations/` ディレクトリ、スキーマ定義コードが一切存在しない。
- **リポジトリ層なし**: モデル定義、DBアクセス層が存在しない。

### 1.2 フロントエンド（Next.js 16 / React 19）

- UIページは実装済み（`login`, `signup`, `home`, `tasks`, `profile`, `reset-password`, `dev-sample`）。
- **状態はすべてクライアント側**（React Context `lib/store.tsx` のメモリ上、ハードコードされたシードデータ）。
- **バックエンドAPI呼び出しなし**: `fetch` / `axios` 等のAPIクライアントコードが存在しない。`NEXT_PUBLIC_API_URL` 等の環境変数も未定義。
- ルートページ (`/`) は `/login` へリダイレクトするのみ。

### 1.3 Docker Compose

- `frontend`（port 3000）と `backend`（port 8080）の2サービスのみ。
- **DBサービスなし**。`depends_on` / ヘルスチェック / ボリュームも未定義。
- 各サービスはバインドマウントでホットリロード対応。

### 1.4 E2Eテスト（Playwright）

- `frontend/e2e/dev-sample.spec.ts` のみ存在（3シナリオ）。すべて `frontend/app/dev-sample/` ページ（自己完結したテストフィクスチャページ）を対象とする。
- **フロントエンド ↔ バックエンド ↔ DB の結合を検証するテストではない**。
- `playwright.config.ts` の `webServer` は `npm run dev` のみ起動（バックエンド・DBは起動しない）。
- `globalSetup` なし（バックエンド起動待機の仕組みがない）。

### 1.5 ドキュメント

- `docs/design/database_design.md`: 全7テーブルのスキーマ定義が完成済み（設計のみ）。
- `docs/design/api_design/01_overview.md`: 全エンドポイントのAPI仕様が完成済み（実装なし）。
- `docs/design/tech_stack.md`: Supabase / PostgreSQL を「BaaS / DB」として記載。ローカル開発DBの記述なし。
- `docs/backend/test-guide/`: 空ディレクトリ（ファイルなし）。
- `frontend/TESTING_GUIDE.md`: テスト作成ガイド。E2Eセクションは `dev-sample` ページ前提の記述。

---

## 2. 設計決定事項

| 項目 | 決定内容 | 理由 |
|---|---|---|
| **Task 1 実装スコープ** | インフラ基盤のみ（DBドライバ導入、設定読み込み、マイグレーション仕組、スキーマSQL作成まで。APIハンドラ/リポジトリ層の実装は含まない） | バックエンドが最小スケルトンのため、まずDB接続基盤を整えることが優先。API実装は別タスクで対応。 |
| **Task 2 実装スコープ** | インフラ構築 + スモークテスト（docker-compose で全サービス起動 → Playwright が全サービスに対してスモークE2Eテスト実行。DB初期化/リセット仕組とglobalSetupによる起動待機を含む） | API実装前でも、インフラが正しく立ち上がることと3層が連携可能な環境が構築できていることを検証する。 |
| **Docker Compose 構成** | 単一 `docker-compose.yml` に `profiles` を使わず、DBサービスを追加してデフォルトで全サービス起動。テスト用DBは環境変数 `DB_NAME` で切り替え | シンプルさを優先。開発もテストも同じ `docker-compose.yml` を使い、`DB_NAME` 環境変数でDB名を切り替える。 |
| **マイグレーションツール** | `golang-migrate/migrate/v4` | Go エコシステムで標準的。CLI + Go ライブラリ両対応、Dockerイメージも提供。`iofs` ソースで `embed.FS` から直接実行可能。 |
| **DBドライバ** | `github.com/jackc/pgx/v5/stdlib`（`database/sql` インターフェース経由） | `database/sql` 標準インターフェースを維持しつつ、PostgreSQL専用の高性能ドライバを利用。将来的に `pgx` ネイティブAPIにも移行可能。`job_design.md` が前提とする `db.Conn(ctx)` による Advisory Lock（`pg_try_advisory_lock`）も `*sql.DB` 経由で対応可能。 |
| **クエリレイヤ（将来タスク）** | `sqlc`（SQL-first コード生成）+ `Masterminds/squirrel`（動的クエリビルダ） | 設計書（`database_design.md` / `job_design.md`）に PostgreSQL 固有機能（`pg_trgm` GIN インデックス、Advisory Lock、CTE チャンク削除、部分一意インデックス、条件付き CHECK 制約、`ILIKE` / `NULLS LAST` / `TIMESTAMPTZ`）が多用されており、ORM（GORM / ent）の抽象を通すメリットが薄い。生 SQL を書く `sqlc` が設計と最も整合する。命名規約が `UPPER_SNAKE_CASE`（`USER_ID`, `CREATED_AT` 等）のため ORM は全フィールド `column:` タグが必要で保守負荷が高い。動的 WHERE / ORDER BY を持つ `GET tasks` エンドポイントのみ `squirrel` で補完。**本Task 1では導入せず、将来のAPI実装タスクで導入する。** |
| **既存 `dev-sample.spec.ts`** | 削除して新規E2E（`smoke.spec.ts`）に置換 | ユーザー指示。`dev-sample` ページ本体とコンポーネントテスト（`__tests__/page.test.tsx`）はユニット/コンポーネントテストとして残置する。 |
| **本番環境** | Supabase（PostgreSQL マネージド）を継続利用 | ローカル開発・テストは Docker PostgreSQL、本番は Supabase。接続情報は環境変数で切り替え。 |

---

## 3. Task 1: ローカルDB（PostgreSQL）インフラ基盤の構築

### 3.1 目標

- Docker Compose でローカル PostgreSQL を立ち上げ、バックエンドから接続できるようにする。
- 本番は Supabase（PostgreSQL マネージド）を継続利用。ローカル開発・テスト時は Docker PostgreSQL。
- 環境変数で接続先を切り替え可能にする。
- `docs/design/database_design.md` のスキーマに基づくマイグレーションSQLを作成し、アプリ起動時に自動実行する。

### 3.2 Docker Compose 設計

#### 変更ファイル: `docker-compose.yml`

現状の `frontend` + `backend` に `db` サービスを追加する。

**追加内容:**

```yaml
services:
  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-synctask}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-synctask_pass}
      POSTGRES_DB: ${DB_NAME:-synctask_dev}
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-synctask} -d ${DB_NAME:-synctask_dev}"]
      interval: 5s
      timeout: 5s
      retries: 10
      start_period: 10s

  frontend:
    # ...（既存のまま）
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080  # 新規追加

  backend:
    # ...（既存のまま）
    depends_on:
      db:
        condition: service_healthy
    environment:
      - GIN_MODE=debug
      - DB_HOST=db                    # 新規追加（Dockerネットワーク内のホスト名）
      - DB_PORT=5432                  # 新規追加
      - DB_USER=${DB_USER:-synctask}  # 新規追加
      - DB_PASSWORD=${DB_PASSWORD:-synctask_pass}  # 新規追加
      - DB_NAME=${DB_NAME:-synctask_dev}  # 新規追加
      - DB_SSLMODE=disable            # 新規追加（ローカル開発時はSSL無効）
      - FRONTEND_URL=http://localhost:3000  # 新規追加（CORS用、将来的なCORSミドルウェア実装時に使用）

volumes:
  db_data:
```

**設計ポイント:**

- `postgres:17-alpine`: 軽量イメージ。PostgreSQL 17 は Supabase が対応する最新メジャーバージョンと整合する想定。Supabase プロジェクト作成時に選択した PG メジャーバージョンに合わせて `postgres:N-alpine` を調整可能（`pg_trgm` 拡張や Advisory Lock は PG 14 以降で動作するためバージョン非依存）。本計画では 17 をデフォルトとするが、本番 Supabase の PG バージョンが異なる場合はイメージタグを合わせること。
- `healthcheck`: `pg_isready` でDB起動完了を待機。`backend` の `depends_on.condition: service_healthy` により、DBがreadyになるまでバックエンドは起動しない。
- `db_data` ボリューム: DBデータを永続化。`docker compose down` では消えず、`docker compose down -v` で消去可能。
- **テスト用DBの切り替え**: `DB_NAME=synctask_e2e docker compose up` でテスト専用DBを使用可能。`docker-compose.test.yml` 等のファイル分割は行わない（環境変数で切り替え）。
- デフォルト値（`${DB_USER:-synctask}` 等）により、`.env` ファイルがなくても `docker compose up` で起動可能。

### 3.3 バックエンド変更

#### 3.3.1 新規: `backend/config/config.go`

環境変数から設定を読み込む `Config` 構造体と `Load()` 関数を定義する。

**構造:**

```go
package config

type Config struct {
    GinMode     string  // GIN_MODE (debug/release/test)
    FrontendURL string  // FRONTEND_URL (CORS許可オリジン)
    DB          DBConfig
}

type DBConfig struct {
    Host     string  // DB_HOST
    Port     string  // DB_PORT
    User     string  // DB_USER
    Password string  // DB_PASSWORD
    Name     string  // DB_NAME
    SSLMode  string  // DB_SSLMODE (disable/require/verify-full)
}

// Load は環境変数から Config を構築して返します。
func Load() *Config { ... }

// DSN は PostgreSQL 接続文字列を返します。
func (c DBConfig) DSN() string {
    return "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
}
```

**設計ポイント:**

- `os.Getenv` のみ使用（`godotenv` 等の外部ライブラリに依存しない）。`.env` ファイルの読み込みは Docker Compose の `env_file` または `docker compose --env-file` で対応。
- デフォルト値はDocker Composeのデフォルトと一致させる（`DB_HOST=localhost`, `DB_PORT=5432`, `DB_USER=synctask`, `DB_NAME=synctask_dev`, `DB_SSLMODE=disable`）。ただしDocker内ではdocker-compose.ymlで上書きされる。
- ローカル（Docker外）で `go run` する場合も考慮し、`DB_HOST` デフォルトは `localhost` とする。

#### 3.3.2 新規: `backend/db/db.go`

`database/sql` + `pgx/v5/stdlib` によるDB接続プール初期化とマイグレーション実行。

**構造:**

```go
package db

import (
    "database/sql"
    "embed"
    "fmt"

    "github.com/golang-migrate/migrate/v4"
    "iofs" "github.com/golang-migrate/migrate/v4/source/iofs"
    _ "github.com/jackc/pgx/v5/stdlib"

    "synctask/backend/config"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Connect は PostgreSQL に接続し、*sql.DB を返します。
func Connect(cfg config.DBConfig) (*sql.DB, error) {
    db, err := sql.Open("pgx", cfg.DSN())
    if err != nil { return nil, err }
    if err := db.Ping(); err != nil { return nil, err }
    return db, nil
}

// Migrate は埋め込まれたマイグレーションファイルを順次実行します。
func Migrate(db *sql.DB) error {
    src, err := iofs.New(migrationFS, "migrations")
    if err != nil { return err }
    // ソースとターゲットを指定してマイグレーション実行
    // databaseターゲットには pgx で開いた *sql.DB を使用
    m, err := migrate.NewWithInstance("iofs", src, "pgx", databaseInstance)
    if err != nil { return err }
    defer m.Close()
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}
```

**設計ポイント:**

- `pgx/v5/stdlib` を使用し、`database/sql` インターフェース経由で接続。将来的に `pgx` ネイティブAPIへの移行も容易。
- `embed.FS` でマイグレーションSQLをバイナリに埋め込み、`golang-migrate` の `iofs` ソースから読み込み。別途マイグレーション用バイナリやDockerイメージ不要。
- `migrate.ErrNoChange` はエラー扱いしない（スキーマ変更がない場合は正常終了）。
- 接続プール設定（`db.SetMaxOpenConns`, `db.SetMaxIdleConns`, `db.SetConnMaxLifetime`）はデフォルトのまま（必要に応じて後日調整）。ただし `job_design.md` が前提とする **`db.Conn(ctx)` による専用コネクション取得 + Advisory Lock（`pg_try_advisory_lock` / `pg_advisory_unlock`）** は `*sql.DB` の標準APIで完結可能であり、`pgx/v5/stdlib` でも動作する。将来のAPI実装タスクで Cron ジョブを `*sql.DB` に依存させる設計と整合することを前提とする。

#### 3.3.3 新規: `backend/db/migrations/000001_init.up.sql`

`docs/design/database_design.md` の全テーブル（7テーブル）+ `pg_trgm` 拡張 + 推奨インデックス + CHECK制約を作成する。

**含める内容:**

1. `CREATE EXTENSION IF NOT EXISTS pg_trgm;` — 日本語検索用
2. `LOGIN_ACCOUNT` テーブル — ユーザーアカウント（論理削除、ログインロックアウト）
3. `TASK` テーブル — タスク（FK to LOGIN_ACCOUNT、優先度/ステータス/ピン留め/コメント/検索テキスト）
4. `LOGIN_SESSION` テーブル — Cookieベースセッション
5. `OTP_SESSION` テーブル — OTP（新規登録/パスワードリセット/メール変更）
6. `LOGIN_IP_RATE_LIMIT` テーブル — IPベースレートリミット
7. `LOGIN_LOG` テーブル — ログイン情報ログ
8. `ACCESS_LOG` テーブル — APIアクセスログ
9. `MAIL_AUTH_LOG` テーブル — メール認証ログ
10. 全インデックス（`database_design.md` 7.1〜7.4 に記載の推奨インデックス）
11. 全CHECK制約（`database_design.md` に記載のOTP_SESSION CHECK制約群、TASK priority/status CHECK制約）

> **参照**: `docs/design/database_design.md` の各テーブル定義および7. 推奨インデックス設計をそのままSQL化する。

#### 3.3.4 新規: `backend/db/migrations/000001_init.down.sql`

`000001_init.up.sql` の逆順で全テーブル・拡張をDROPする。

**構造:**

```sql
DROP TABLE IF EXISTS MAIL_AUTH_LOG;
DROP TABLE IF EXISTS ACCESS_LOG;
DROP TABLE IF EXISTS LOGIN_LOG;
DROP TABLE IF EXISTS LOGIN_IP_RATE_LIMIT;
DROP TABLE IF EXISTS OTP_SESSION;
DROP TABLE IF EXISTS LOGIN_SESSION;
DROP TABLE IF EXISTS TASK;
DROP TABLE IF EXISTS LOGIN_ACCOUNT;
DROP EXTENSION IF EXISTS pg_trgm;
```

#### 3.3.5 変更: `backend/main.go`

起動シーケンスを「設定読込 → DB接続 → マイグレーション → サーバー起動」に変更する。

**変更後の構造:**

```go
package main

import (
    "log"

    "synctask/backend/config"
    "synctask/backend/db"
    "synctask/backend/router"
)

func main() {
    cfg := config.Load()

    database, err := db.Connect(cfg.DB)
    if err != nil {
        log.Fatalf("DB接続に失敗しました: %v", err)
    }
    defer database.Close()

    if err := db.Migrate(database); err != nil {
        log.Fatalf("マイグレーションに失敗しました: %v", err)
    }

    r := router.SetupRouter(database)
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("サーバーの起動に失敗しました: %v", err)
    }
}
```

#### 3.3.6 変更: `backend/router/router.go`

`SetupRouter` のシグネチャを `SetupRouter(db *sql.DB) *gin.Engine` に変更し、`HealthCheckHandler` に `db` を渡す。

**変更内容:**

- `import "database/sql"` 追加
- `func SetupRouter(db *sql.DB) *gin.Engine` に変更
- `r.GET("/health-check", handler.HealthCheckHandler(db))` に変更

#### 3.3.7 変更: `backend/handler/health.go`

`HealthCheckHandler` を `func HealthCheckHandler(db *sql.DB) gin.HandlerFunc` に変更し、DB接続状態をレスポンスに含める。

**変更内容:**

- `HealthResponse` に `Database string` フィールド追加（例: `"connected"` / `"disconnected"`）
- `db.Ping()` を実行し、成功時に `"connected"`、失敗時に `"disconnected"` を設定

```go
type HealthResponse struct {
    Status   string `json:"status" example:"ok"`
    Message  string `json:"message" example:"healthy"`
    Database string `json:"database" example:"connected"`
}

func HealthCheckHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        dbStatus := "connected"
        if err := db.Ping(); err != nil {
            dbStatus = "disconnected"
        }
        c.JSON(http.StatusOK, HealthResponse{
            Status:   "ok",
            Message:  "healthy",
            Database: dbStatus,
        })
    }
}
```

#### 3.3.8 変更: `backend/handler/health_test.go`

新しいハンドラシグネチャ（`HealthCheckHandler(db)` が `gin.HandlerFunc` を返す）に対応する。

- `sql.DB` のモックまたは実際のテスト用DBを使用
- `httptest` で `gin.HandlerFunc` を呼び出すよう調整

**テスト方針:**

- `sql.Open("pgx", ...)` でテスト用DBに接続（`DB_HOST` 環境変数または `localhost` を使用）
- テスト用DBが未起動の場合は `t.Skip` でスキップ（CI環境以外ではDB起動を必須としない）
- または `database/sql` のモック（`sqlmock`）を使用してDB不要でテスト

> **推奨**: `sqlmock` (`github.com/DATA-DOG/go-sqlmock`) を使用してDB不要でハンドラのユニットテストを維持。DB依存のテストはE2Eテスト側でカバー。

#### 3.3.9 変更: `backend/router/router_test.go`

`SetupRouter(db)` の新しいシグネチャに対応。

- `sql.DB` のモック（`sqlmock`）を渡す
- またはテスト用DB接続を渡す

#### 3.3.10 変更: `backend/go.mod`

以下の依存を `go mod tidy` で追加:

| パッケージ | 用途 |
|---|---|
| `github.com/jackc/pgx/v5` | PostgreSQL ドライバ（`stdlib` 経由で `database/sql` として使用） |
| `github.com/golang-migrate/migrate/v4` | マイグレーション実行（`iofs` ソース使用） |
| `github.com/DATA-DOG/go-sqlmock` | ハンドラテスト用モック（dev依存） |

#### 3.3.11 変更: `backend/Dockerfile`

レイヤーキャッシュ効率化のため、`go.mod` / `go.sum` を先にCOPYする構成に変更。

```dockerfile
FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]
```

### 3.4 環境変数ファイル

#### 3.4.1 新規: `backend/.env.example`

バックエンド単体（Docker外）で `go run` する場合の環境変数テンプレート。

```env
# Gin mode (debug/release/test)
GIN_MODE=debug

# Frontend URL (CORS allowed origin)
FRONTEND_URL=http://localhost:3000

# Database connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=synctask
DB_PASSWORD=synctask_pass
DB_NAME=synctask_dev
DB_SSLMODE=disable
```

#### 3.4.2 新規: `.env.example`（ルート）

Docker Compose 用環境変数テンプレート。

```env
# Docker Compose 環境変数（オプション: .env ファイルで上書き可能）
DB_USER=synctask
DB_PASSWORD=synctask_pass
DB_NAME=synctask_dev
```

#### 3.4.3 変更: `.gitignore`（ルート）

`.env` ファイルをGitignoreに追加（既存の `tmp/` と `node_modules/` に加えて）。

**追加内容:**

```
# Environment variables
.env
.env.local
.env.*.local
```

### 3.5 ドキュメント変更

#### 3.5.1 変更: `docs/design/tech_stack.md`

以下の内容を追加・更新:

1. **セクション4「データベース & インフラストラクチャ」の更新**:
   - ローカル開発・テスト環境: Docker PostgreSQL（`postgres:17-alpine`）
   - 本番環境: Supabase（PostgreSQL マネージド）
   - 環境変数による接続先切り替え方式の記述

2. **新規セクション「バックエンド DB アクセス層（将来計画）」の追加**:
   - クエリレイヤ: `sqlc`（SQL-first コード生成）+ `Masterminds/squirrel`（動的クエリビルダ）
   - 選定理由: 設計書における PostgreSQL 固有機能（`pg_trgm` GIN, Advisory Lock, CTE, 部分インデックス, 条件付き CHECK 制約, `ILIKE` / `NULLS LAST`）の多用、`UPPER_SNAKE_CASE` 命名規約に対する ORM の不適合
   - 現状（Task 1）では DB 接続基盤のみ。`sqlc` / `squirrel` は将来のAPI実装タスクで導入する計画と明記

3. **新規セクション「環境変数一覧」の追加**:

   | 環境変数 | 用途 | デフォルト値 | 対象 |
   |---|---|---|---|
   | `GIN_MODE` | Gin実行モード | `debug` | backend |
   | `FRONTEND_URL` | CORS許可オリジン | `http://localhost:3000` | backend |
   | `DB_HOST` | DBホスト | `localhost`（Docker内: `db`） | backend |
   | `DB_PORT` | DBポート | `5432` | backend |
   | `DB_USER` | DBユーザー名 | `synctask` | db, backend |
   | `DB_PASSWORD` | DBパスワード | `synctask_pass` | db, backend |
   | `DB_NAME` | DB名（開発: `synctask_dev` / テスト: `synctask_e2e`） | `synctask_dev` | db, backend |
   | `DB_SSLMODE` | SSL接続モード | `disable`（本番: `require`） | backend |
   | `NEXT_PUBLIC_API_URL` | バックエンドAPIベースURL | `http://localhost:8080` | frontend |

4. **Docker Compose 構成の記述**:
   - 3サービス構成（`db`, `backend`, `frontend`）
   - `depends_on` とヘルスチェックによる起動順序制御
   - 開発用とテスト用のDB切り替え方法（`DB_NAME` 環境変数）

5. **マイグレーション方式の記述**:
   - `golang-migrate/migrate/v4` + `embed.FS` によるアプリ起動時自動実行（ローカル開発・テスト環境向け）
   - マイグレーションファイル配置場所: `backend/db/migrations/`
   - 本番（Supabase）では将来的に CI/CD パイプラインから `migrate` CLI 実行に移行する旨を注記

#### 3.5.2 変更: `README.md`（ルート）

ローカル開発の起動手順を更新。

```markdown
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
```

---

## 4. Task 2: E2Eテスト環境構築 + スモークテスト

### 4.1 目標

- `docker compose up` で frontend + backend + DB が連携して起動する環境を構築する。
- Playwright が全サービス（フロントエンド ↔ バックエンド ↔ DB）に対してスモークE2Eテストを実行できるようにする。
- `globalSetup` でバックエンド起動を待機し、未起動時は分かりやすいエラーメッセージを出力する。
- DB初期化/リセットの仕組みを提供する（テスト実行前にDBをクリーンな状態にする）。

### 4.2 削除ファイル

#### 4.2.1 削除: `frontend/e2e/dev-sample.spec.ts`

ユーザー指示により削除。新規E2E（`smoke.spec.ts`）に置換。

> **注**: `frontend/app/dev-sample/` ページ本体とそのコンポーネントテスト（`frontend/app/dev-sample/__tests__/page.test.tsx`）は**残置**する。これらはユニット/コンポーネントテストとして有効であり、E2Eテストとは独立して実行される。

### 4.3 新規作成ファイル

#### 4.3.1 新規: `frontend/e2e/global-setup.ts`

Playwright の `globalSetup` スクリプト。バックエンドの起動を待機する。

**構造:**

```typescript
import { request } from '@playwright/test';

async function globalSetup() {
  const apiURL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';
  const maxRetries = 30;
  const retryInterval = 1000; // 1秒

  for (let i = 0; i < maxRetries; i++) {
    try {
      const context = await request.newContext();
      const response = await context.get(`${apiURL}/health-check`);
      if (response.ok()) {
        const body = await response.json();
        if (body.database === 'connected') {
          console.log('✅ バックエンド + DB が起動済み（E2Eテスト開始可能）');
          return;
        }
        console.log(`⚠️ バックエンド起動済みですがDB接続状態: ${body.database}`);
      }
      await context.dispose();
    } catch {
      // バックエンド未起動
    }
    await new Promise((resolve) => setTimeout(resolve, retryInterval));
  }

  throw new Error(
    `バックエンドが起動していません（${apiURL}/health-check に接続できません）。\n` +
    'E2Eテストを実行する前に `docker compose up` で全サービスを起動してください。'
  );
}

export default globalSetup;
```

**設計ポイント:**

- バックエンドの `/health-check` エンドポイント（開発モード限定）をポーリング。
- レスポンスの `database` フィールドが `"connected"` であることも確認（DB接続済みであることの検証）。
- 最大30秒待機（1秒 × 30リトライ）。
- 未起動時は `docker compose up` を促すエラーメッセージを出力。
- `NEXT_PUBLIC_API_URL` 環境変数でバックエンドURLを切り替え可能（デフォルト: `http://localhost:8080`）。

#### 4.3.2 新規: `frontend/e2e/smoke.spec.ts`

フルスタックスモークE2Eテスト。3層（フロントエンド ↔ バックエンド ↔ DB）が連携していることを検証する。

**テストシナリオ:**

```typescript
import { test, expect, request } from '@playwright/test';

test.describe('フルスタックスモークテスト (Frontend ↔ Backend ↔ DB)', () => {
  test('シナリオ1: フロントエンドのログインページが正常にレンダリングされること', async ({ page }) => {
    await page.goto('/login');
    // ページタイトル要素または見出しの確認
    await expect(page).toHaveURL(/\/login/);
    // メールアドレス入力フィールドの存在確認
    await expect(page.getByLabel('メールアドレス')).toBeVisible();
    // パスワード入力フィールドの存在確認
    await expect(page.getByLabel('パスワード')).toBeVisible();
    // ログインボタンの存在確認
    await expect(page.getByRole('button', { name: 'ログイン' })).toBeVisible();
  });

  test('シナリオ2: バックエンドAPI（/health-check）が200 OKで応答すること', async () => {
    const apiURL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';
    const context = await request.newContext();
    const response = await context.get(`${apiURL}/health-check`);
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.status).toBe('ok');
    expect(body.message).toBe('healthy');
    await context.dispose();
  });

  test('シナリオ3: バックエンドがDBに接続済みであること（ヘルスチェック経由）', async () => {
    const apiURL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';
    const context = await request.newContext();
    const response = await context.get(`${apiURL}/health-check`);
    const body = await response.json();
    expect(body.database).toBe('connected');
    await context.dispose();
  });

  test('シナリオ4: ルートページ（/）が/loginへリダイレクトすること', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveURL(/\/login/);
  });
});
```

**設計ポイント:**

- シナリオ1: フロントエンド（Next.js）が正常に配信されていることを確認。
- シナリオ2: バックエンド（Go/Gin）が正常に起動していることを確認。
- シナリオ3: バックエンドがDB（PostgreSQL）に接続済みであることを確認。3層が連携していることを1つのテストファイルで検証。
- シナリオ4: フロントエンドのルーティングが正常に動作していることを確認。
- `request` API（PlaywrightのAPI testing機能）を使用してバックエンドAPIを直接呼び出す。
- API URLは `NEXT_PUBLIC_API_URL` 環境変数で切り替え可能（デフォルト: `http://localhost:8080`）。

### 4.4 変更ファイル

#### 4.4.1 変更: `frontend/playwright.config.ts`

`globalSetup` を追加し、`webServer` の挙動を調整する。

**変更内容:**

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'list',
  globalSetup: './e2e/global-setup.ts',  // 新規追加
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: true,  // 変更: !process.env.CI → true（Docker Compose起動済みサーバーを再利用）
    timeout: 120 * 1000,
  },
});
```

**設計ポイント:**

- `globalSetup`: バックエンド + DB の起動を待機してからテスト開始。
- `webServer.reuseExistingServer: true`: Docker Compose で既に起動しているフロントエンドを再利用。未起動の場合は `npm run dev` で起動を試みる。
- `testDir: './e2e'`: `smoke.spec.ts` を含む `e2e/` ディレクトリ配下の `.spec.ts` ファイルを実行。

#### 4.4.2 変更: `frontend/package.json`

テストスクリプトを整理する。

**変更内容:**

```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "eslint",
    "test": "npm run test:unit && npm run test:e2e",
    "test:unit": "vitest run",
    "test:unit:watch": "vitest",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui"
  }
}
```

> 既存のスクリプトから変更なし（`test:e2e` はそのまま `playwright test` を実行）。`globalSetup` と `playwright.config.ts` の変更により、`npm run test:e2e` が自動的にフルスタックE2Eテストとして動作する。

#### 4.4.3 変更: `frontend/TESTING_GUIDE.md`

E2Eセクションを「フルスタック結合テスト」向けに書き換える。

**変更内容:**

- セクション3「ブラウザ結合テスト（E2E）を作成する (Playwright)」を更新:
  - Docker Compose で全サービス（frontend + backend + DB）を起動してからテスト実行するワークフローを記述。
  - `globalSetup` の役割（バックエンド + DB の起動待機）を説明。
  - テストDBの概念（`DB_NAME=synctask_e2e` でテスト専用DBを使用可能）を記述。
- セクション5「テストの実行・デバッグ方法」のE2E部分を更新:
  - `docker compose up` → `npm run test:e2e` の2ステップワークフローを明記。
  - 未起動時のエラーメッセージと対処法を記述。

### 4.5 新規ドキュメント

#### 4.5.1 新規: `docs/backend/test-guide/README.md`

バックエンド + 結合テストのガイド。

**構造:**

```markdown
# バックエンド テストガイド

## 1. テスト構成

| テスト分類 | 対象 | 使用ツール | 実行コマンド |
|---|---|---|---|
| ユニットテスト | ハンドラ、ルーター、ビジネスロジック | Go testing + testify | `go test ./...` |
| E2E結合テスト | フロントエンド ↔ バックエンド ↔ DB | Playwright | `npm run test:e2e`（frontend ディレクトリ内） |

## 2. ユニットテスト

### 実行方法

```bash
cd backend
go test ./...
go test -v ./...
go test -cover ./...
```

### DB依存テスト

- ハンドラのDB依存部分は `go-sqlmock` でモック化。
- DB接続が必要なテストは `DB_HOST` 環境変数が設定されている場合のみ実行（未設定時は `t.Skip`）。

## 3. E2E結合テスト

### 事前準備

```bash
# テスト用DBで全サービス起動
DB_NAME=synctask_e2e docker compose up --build
```

### テスト実行

```bash
cd frontend
npm run test:e2e
```

### ワークフロー

1. `docker compose up` で frontend + backend + DB が起動。
2. `globalSetup` がバックエンドの `/health-check` をポーリングし、DB接続済みを確認してからテスト開始。
3. Playwright がフロントエンド（`localhost:3000`）とバックエンド（`localhost:8080`）に対してテスト実行。

### トラブルシューティング

- **「バックエンドが起動していません」エラー**: `docker compose up` を実行して全サービスを起動してください。
- **DB接続エラー**: `DB_NAME` 環境変数でテスト用DBを指定してください。


---

## 5. ファイル変更一覧

### 5.1 新規作成ファイル

| ファイルパス | 内容 |
|---|---|
| `backend/config/config.go` | 環境変数からConfigを読み込む構造体とLoad()関数 |
| `backend/db/db.go` | DB接続プール初期化（pgx/stdlib）+ マイグレーション実行（golang-migrate + embed.FS） |
| `backend/db/migrations/000001_init.up.sql` | 全テーブル（7テーブル）+ pg_trgm拡張 + インデックス + CHECK制約のCREATE |
| `backend/db/migrations/000001_init.down.sql` | 全テーブル・拡張のDROP（逆順） |
| `backend/.env.example` | バックエンド環境変数テンプレート |
| `.env.example` | Docker Compose用環境変数テンプレート |
| `frontend/e2e/global-setup.ts` | Playwright globalSetup: バックエンド + DB 起動待機 |
| `frontend/e2e/smoke.spec.ts` | フルスタックスモークE2Eテスト（4シナリオ） |
| `docs/backend/test-guide/README.md` | バックエンド + 結合テストガイド |

### 5.2 変更ファイル

| ファイルパス | 変更内容 |
|---|---|
| `docker-compose.yml` | `db` サービス追加、`depends_on` + ヘルスチェック、各サービスにDB/API env var追加、`db_data` ボリューム追加 |
| `backend/main.go` | 起動シーケンス変更: 設定読込 → DB接続 → マイグレーション → サーバー起動 |
| `backend/router/router.go` | `SetupRouter(db *sql.DB)` シグネチャ変更、`HealthCheckHandler(db)` に変更 |
| `backend/handler/health.go` | `HealthCheckHandler(db) gin.HandlerFunc` に変更、`Database` フィールド追加、`db.Ping()` 実行 |
| `backend/handler/health_test.go` | 新しいハンドラシグネチャに対応（sqlmock使用） |
| `backend/router/router_test.go` | 新しい `SetupRouter(db)` シグネチャに対応 |
| `backend/go.mod` | `pgx/v5`, `golang-migrate/v4`, `go-sqlmock` 追加（`go mod tidy` 実行） |
| `backend/Dockerfile` | レイヤーキャッシュ改善（`go.mod`/`go.sum` を先にCOPY） |
| `frontend/playwright.config.ts` | `globalSetup` 追加、`webServer.reuseExistingServer: true` に変更 |
| `frontend/TESTING_GUIDE.md` | E2Eセクションをフルスタック結合テスト向けに更新 |
| `docs/design/tech_stack.md` | ローカルDB、環境変数一覧、Docker Compose構成、マイグレーション方式の記述追加 |
| `README.md` | ローカル開発の起動手順を更新 |
| `.gitignore` | `.env` ファイルを追加 |

### 5.3 削除ファイル

| ファイルパス | 理由 |
|---|---|
| `frontend/e2e/dev-sample.spec.ts` | ユーザー指示により削除。新規E2E（`smoke.spec.ts`）に置換。 |

---

## 6. 実行順序

### Phase 1: Task 1（ローカルDBインフラ基盤）

1. `docker-compose.yml` に `db` サービスを追加
2. `backend/config/config.go` 新規作成
3. `backend/db/migrations/000001_init.up.sql` および `000001_init.down.sql` 新規作成
4. `backend/db/db.go` 新規作成（DB接続 + マイグレーション）
5. `backend/go.mod` に依存追加（`go mod tidy`）
6. `backend/main.go` 変更（起動シーケンス）
7. `backend/router/router.go` 変更（シグネチャ変更）
8. `backend/handler/health.go` 変更（DB状態追加）
9. `backend/handler/health_test.go` および `backend/router/router_test.go` 変更（新シグネチャ対応）
10. `backend/Dockerfile` 変更（レイヤーキャッシュ改善）
11. `backend/.env.example` および `.env.example` 新規作成
12. `.gitignore` 変更
13. `docs/design/tech_stack.md` および `README.md` 更新
14. **検証**: `docker compose up --build` で全サービス起動 → `localhost:8080/health-check` でDB接続確認 → `go test ./...` でバックエンドテスト通過確認

### Phase 2: Task 2（E2Eテスト環境構築 + スモークテスト）

1. `frontend/e2e/dev-sample.spec.ts` 削除
2. `frontend/e2e/global-setup.ts` 新規作成
3. `frontend/e2e/smoke.spec.ts` 新規作成
4. `frontend/playwright.config.ts` 変更（`globalSetup` 追加、`reuseExistingServer` 変更）
5. `frontend/TESTING_GUIDE.md` 更新
6. `docs/backend/test-guide/README.md` 新規作成
7. **検証**: `docker compose up` → `cd frontend && npm run test:e2e` でスモークテスト通過確認 → `npm run test:unit` でユニットテスト通過確認

---

## 7. 検証コマンド

実装完了後に以下のコマンドで動作確認を行う。

### 7.1 バックエンド

```bash
# バックエンドのビルド確認
cd backend && go build ./...

# バックエンドのユニットテスト
cd backend && go test ./...

# バックエンドのテストカバレッジ
cd backend && go test -cover ./...
```

### 7.2 Docker Compose 起動確認

```bash
# 全サービス起動
docker compose up --build

# 別ターミナルでヘルスチェック確認
curl http://localhost:8080/health-check
# 期待レスポンス: {"status":"ok","message":"healthy","database":"connected"}

# PostgreSQL直接接続確認
docker compose exec db psql -U synctask -d synctask_dev -c "\dt"
# 期待: 7テーブルが表示される
```

### 7.3 E2Eテスト

```bash
# 全サービス起動（別ターミナル）
docker compose up --build

# E2Eテスト実行
cd frontend && npm run test:e2e
# 期待: 4シナリオすべてパス

# ユニットテスト（フロントエンド）
cd frontend && npm run test:unit
# 期待: 既存のユニットテストがパス
```

### 7.4 テスト用DBでのE2Eテスト

```bash
# テスト用DBで全サービス起動
DB_NAME=synctask_e2e docker compose up --build

# E2Eテスト実行
cd frontend && npm run test:e2e
```

---

## 8. 注意事項・制約

1. **本番環境（Supabase）への影響なし**: ローカル開発・テスト用のDocker PostgreSQL追加のみ。本番のSupabase設定は環境変数で切り替え（`DB_HOST`, `DB_SSLMODE` 等）。
2. **APIハンドラ/リポジトリ層は実装しない**: Task 1のスコープは「インフラ基盤のみ」。`database_design.md` のスキーマをマイグレーションで作成し、DB接続とヘルスチェックでのDB状態確認まで。CRUD APIの実装は別タスク。
3. **フロントエンドのAPI呼び出しは実装しない**: Task 2のスコープは「E2Eテストインフラ構築 + スモークテスト」。フロントエンドからバックエンドAPIを呼び出す実装（APIクライアント、fetchラッパー等）は含まない。スモークテストは Playwright の `request` API でバックエンドを直接叩いて検証する。
4. **`dev-sample` ページは残置**: `frontend/app/dev-sample/` とコンポーネントテスト（`__tests__/page.test.tsx`）はユニット/コンポーネントテストとして残す。削除対象はE2Eの `dev-sample.spec.ts` のみ。
5. **`gin-contrib/cors` は実装しない**: `tech_stack.md` にCORSミドルウェアの記載があるが、現状のコードには未実装。本計画ではCORS実装は含まない（別タスク）。`FRONTEND_URL` 環境変数の設定のみ追加し、将来のCORS実装に備える。
6. **本番（Supabase）でのマイグレーション実行方針**: 本Task 1ではアプリ起動時自動マイグレーション（`embed.FS` + `golang-migrate`）を採用するが、これはローカル開発・テスト環境向け。本番 Supabase ではアプリプロセス起動時の自動実行はスキーマ競合や起動ブロックのリスクがあるため、**将来的には CI/CD パイプラインから `migrate` CLI で実行する方式に移行する**（将来拡張ポイント参照）。Task 1 時点では起動時自動実行で一元化し、本番デプロイ前に見直す。
7. **クエリレイヤ技術（sqlc / squirrel）は本Task 1では導入しない**: 設計決定事項として `sqlc` + `squirrel` 方針を記録したが、Task 1 のスコープ（インフラ基盤のみ）には含めない。`*sql.DB` インターフェースで DB 接続基盤を構築しておくことで、将来タスクで `sqlc` が生成するコードがそのまま `*sql.DB` に依存する形で自然に統合できる前提を維持する。
8. **`job_design.md` の Advisory Lock 要件との整合**: 将来の Cron ジョブ実装（`robfig/cron/v3`）が前提とする `db.Conn(ctx)` による専用コネクション取得 + `pg_try_advisory_lock` は `*sql.DB` 標準APIで完結する。Task 1 で `pgx/v5/stdlib` + `database/sql` インターフェースを選定したことで、将来のジョブ実装も `*sql.DB` に依存する形で一貫して構築可能。Task 1 時点ではジョブ実装は含まない。

---

## 9. 将来の拡張ポイント

- **API実装**: Task 1で構築したDB基盤の上に、`api_design/` に基づくCRUDハンドラとリポジトリ層を実装する。
- **リポジトリ層技術（sqlc + squirrel）導入**: 設計決定事項として記録した `sqlc`（SQL-first コード生成）+ `Masterminds/squirrel`（動的クエリビルダ）を導入する。`backend/db/queries/*.sql` 配下にクエリファイルを配置し `sqlc generate` で型安全な Go コードを生成。動的 WHERE / ORDER BY を持つ `GET tasks` エンドポイントのみ `squirrel` で組み立てる。`UPPER_SNAKE_CASE` 列名は `sqlc.yaml` の `renames` でキャメルケースにマッピング。`backend/Makefile` に `sqlc-generate` / `migrate-up` / `migrate-down` ターゲットを追加し開発体験を整備する。`sqlc` は `*sql.DB` を前提として生成されるため、Task 1 で構築した `pgx/v5/stdlib` 接続基盤と自然に統合可能。
- **フロントエンドAPI統合**: フロントエンドからバックエンドAPIを呼び出すAPIクライアントを実装し、状態管理をクライアントメモリからAPI経由に移行する。
- **E2Eテスト拡充**: API実装後に、ログイン → タスク作成 → 一覧表示 → ステータス更新 → 削除などの実データフローE2Eテストを `smoke.spec.ts` に追加する。
- **CI/CD**: GitHub Actions 等で `docker compose up` → `npm run test:e2e` のワークフローを自動化する。
- **本番マイグレーションの CI/CD 移行**: 本番 Supabase ではアプリ起動時自動マイグレーションではなく、CI/CD パイプラインから `migrate` CLI（`golang-migrate/migrate`）で実行する方式に移行する。スキーマ競合・起動ブロック回避とデプロイの再現性確保が目的。
- **DB初期化/リセット自動化**: テスト実行前にDBをクリーンな状態にリセットする仕組み（`TRUNCATE` 全テーブル、または `migrate down` → `migrate up`）を `globalSetup` に追加する。現状はマイグレーション実行のみで、テスト間のデータリセットは各テストファイルで対応。`TRUNCATE` 方式は高速だが外部キー制約の順序に注意、`migrate down → up` 方式は確実だが遅い。テスト対象テーブル数が増えた段階で方式を選定する。
