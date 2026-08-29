# Technology Stack (技術スタック)

本ドキュメントでは、SyncTask プロジェクトにおいて採用されている技術スタック、選定理由、アーキテクチャ構成、および開発・運用環境を定義します。

---

## 1. 全体アーキテクチャ概要

SyncTask は、フロントエンドとバックエンドが分離された RESTful Web アプリケーション構成を採用しています。

```mermaid
graph TD
    Client["Web Browser / User Agent"]
    
    subgraph Frontend_Layer ["Frontend (Next.js App Router)"]
        FE["Next.js 16 + React 19 + TypeScript"]
        Tailwind["Tailwind CSS v4"]
    end
    
    subgraph Backend_Layer ["Backend (Go Gin Framework)"]
        BE["Gin Web Framework / Go 1.26"]
        Air["Air Hot Reload"]
        Config["Config (os.Getenv)"]
        DBConn["DB Pool (pgx/v5 stdlib)"]
        Migrate["Migrate (golang-migrate + embed.FS)"]
    end
    
    subgraph Data_Layer ["Database & Infrastructure"]
        LocalDB[("Local: Docker PostgreSQL 17")]
        ProdDB[("Prod: Supabase / PostgreSQL")]
        Docker["Docker / Docker Compose"]
    end

    Client -->|HTTP/HTTPS JSON API| FE
    FE -->|HTTP/HTTPS REST API| BE
    BE -->|SQL / pgx stdlib| LocalDB
    BE -.->|SQL / pgx stdlib (本番)| ProdDB
```

---

## 2. フロントエンド (Frontend)

| カテゴリ | 技術 / ライブラリ | バージョン | 選定理由・目的 |
| --- | --- | --- | --- |
| **フレームワーク** | [Next.js](https://nextjs.org/) (App Router) | `16.2.6` | サーバーコンポーネント (RSC) による高速描画、ファイルベースルーティング、および高いパフォーマンスの実現 |
| **UIライブラリ** | [React](https://react.dev/) | `19.2.4` | コンポーネント指向 UI 構築、宣言的 UI 開発 |
| **言語** | [TypeScript](https://www.typescriptlang.org/) | `^5` | 型安全性、エディタ補完の強化、開発効率向上とランタイムエラー防止 |
| **スタイリング** | [Tailwind CSS](https://tailwindcss.com/) | `^4` | ユーティリティファーストなスタイリング、高速なレスポンシブ UI 実装 |
| **リンター / フォーマッタ** | [ESLint](https://eslint.org/) | `^9` | コード品質の保持、チーム間でのコーディングスタイルの統一 |

---

## 3. バックエンド (Backend)

| カテゴリ | 技術 / ライブラリ | バージョン | 選定理由・目的 |
| --- | --- | --- | --- |
| **言語** | [Go](https://go.dev/) | `1.26.1` | 高い実行速度、軽量なメモリ使用量、強力な並行処理能力 (Goroutine) |
| **Web フレームワーク** | [Gin](https://gin-gonic.com/) | `v1.12.0` | 高速で軽量な HTTP ルーティング、ミドルウェアサポート、レスポンス処理の容易さ |
| **DBドライバ** | `github.com/jackc/pgx/v5` (`stdlib`) | `v5.10.0` | `database/sql` 標準インターフェースを維持しつつ、PostgreSQL専用の高性能ドライバを利用。Advisory Lock連携にも対応 |
| **マイグレーション** | `github.com/golang-migrate/migrate/v4` | `v4.19.1` | Go標準的なマイグレーションツール。`embed.FS` + `iofs` ソースにより起動時自動マイグレーションに対応 |
| **バリデーション** | `go-playground/validator` | `v10.30.3` | リクエストボディやパラメータの厳格かつ宣言的なバリデーション |
| **テストモック** | `github.com/DATA-DOG/go-sqlmock` | `v1.5.2` | DB非依存のハンドラ・ルーター単体テスト用モック |
| **CORS ミドルウェア** | `gin-contrib/cors` | - | 環境変数 `FRONTEND_URL` に基づくクロスオリジン制御、認証情報（Credentials）送受信、プリフライト（OPTIONS）制御、`Access-Control-Expose-Headers: Retry-After` 出力 |
| **ジョブスケジューラ** | `robfig/cron/v3` | `v3.0.1` | 定期パージバッチ（OTP・セッション・ログ・レートリミット削除）の定期実行制御（JST基準・バックエンドプロセス内常駐スケジューラ。PostgreSQL Advisory Lock［`pg_try_advisory_lock` / `db.Conn`］と連携した水平スケール時の安全な多重起動防止制御を含む） |
| **開発環境** | [Air](https://github.com/air-verse/air) | - | バックエンド Go コードの変更検知と自動再ビルド (Hot Reload) |

---

## 4. データベース & インフラストラクチャ (Database & Infrastructure)

| カテゴリ | 技術 / サービス | 選定理由・目的 |
| --- | --- | --- |
| **ローカル開発・テスト DB** | Docker PostgreSQL (`postgres:17-alpine`) | ローカル環境およびE2Eテストでの完全再現性と軽量・高速起動の確保 |
| **本番 DB** | [Supabase](https://supabase.com/) (PostgreSQL マネージド) | PostgreSQL データベースのマネージド環境の提供 |
| **RDBMS** | PostgreSQL | 高い信頼性、ACID 補償、リレーショナルデータ構造 (アカウント、タスク、セッション管理) の堅牢な保持、`pg_trgm` 拡張による検索 |
| **コンテナ化** | [Docker](https://www.docker.com/) / [Docker Compose](https://docs.docker.com/compose/) | 開発メンバー間における同一実行環境の容易な構築・再現性の確保 |

### Docker Compose 構成
- **3サービス構成**: `db` (PostgreSQL 17), `backend` (Go Gin / Air), `frontend` (Next.js 16)
- **起動順序制御**: `backend` は `depends_on.db` の `condition: service_healthy`（`pg_isready`）により、DBの健全起動を待機してから起動
- **データ永続化**: `db_data` 名前付きボリュームでDBデータを永続化（`docker compose down -v` でリセット可能）
- **テスト用DB切替**: `DB_NAME=synctask_e2e docker compose up` でテスト専用DBを即座に利用可能

### マイグレーション方式
- **ローカル開発・テスト環境**: アプリ起動時に `embed.FS` に埋め込まれた SQL（`backend/db/migrations/`）を `golang-migrate` 経由で自動実行
- **本番環境 (Supabase)**: 将来的に CI/CD パイプラインから `migrate` CLI を直接実行する方式へ移行（スキーマ競合・起動ブロック回避）

---

## 5. バックエンド DB アクセス層（将来計画）

- **クエリレイヤ**: `sqlc`（SQL-first コード生成）+ `Masterminds/squirrel`（動的クエリビルダ）
- **選定理由**: 設計書における PostgreSQL 固有機能（`pg_trgm` GIN インデックス、Advisory Lock、CTE チャンク削除、部分一意インデックス、条件付き CHECK 制約、`ILIKE` / `NULLS LAST` / `TIMESTAMPTZ`）が多用されており、ORM の抽象を通すより生 SQL を型安全に扱う `sqlc` が最も整合する。動的 WHERE / ORDER BY を持つ `GET tasks` エンドポイントのみ `squirrel` で補完する。
- **現状**: DB 接続基盤（`*sql.DB` + `pgx/v5/stdlib`）のみ構築済み。`sqlc` / `squirrel` は API 実装タスクにて導入予定。

---

## 6. 環境変数一覧

| 環境変数 | 用途 | デフォルト値 | 対象 |
|---|---|---|---|
| `GIN_MODE` | Gin実行モード (`debug`/`release`/`test`) | `debug` | backend |
| `FRONTEND_URL` | CORS許可オリジン | `http://localhost:3000` | backend |
| `DB_HOST` | DBホスト名 | `localhost`（Docker内: `db`） | backend |
| `DB_PORT` | DBポート番号 | `5432` | backend |
| `DB_USER` | DBユーザー名 | `synctask` | db, backend |
| `DB_PASSWORD` | DBパスワード | `synctask_pass` | db, backend |
| `DB_NAME` | DB名（開発: `synctask_dev` / テスト: `synctask_e2e`） | `synctask_dev` | db, backend |
| `DB_SSLMODE` | SSL接続モード (`disable`/`require`/`verify-full`) | `disable` | backend |
| `NEXT_PUBLIC_API_URL` | バックエンドAPIベースURL | `http://localhost:8080` | frontend |

---

## 7. 開発環境・ツールチェーン

- **パッケージマネージャ**: Node.js (`npm`), Go Modules (`go.mod`)
- **コード管理**: Git / GitHub
- **API通信フォーマット**: `JSON` (application/json)
- **認証・セッション方式**: Cookie ベースセッション (`session_id`) + ワンタイムパスワード (OTP)
