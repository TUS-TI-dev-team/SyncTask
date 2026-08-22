# E2Eテスト環境構築 + スモークテスト 計画書

本ドキュメントでは、フロントエンド ↔ バックエンド ↔ DB が連携するE2Eテスト環境の構築に関する設計・実装計画を定義します。ローカルDBインフラ基盤の構築については [01-setup-local-db.md](./01-setup-local-db.md) を参照してください。

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

## 3. E2Eテスト環境構築 + スモークテスト

### 3.1 目標

- `docker compose up` で frontend + backend + DB が連携して起動する環境を構築する。
- Playwright が全サービス（フロントエンド ↔ バックエンド ↔ DB）に対してスモークE2Eテストを実行できるようにする。
- `globalSetup` でバックエンド起動を待機し、未起動時は分かりやすいエラーメッセージを出力する。
- DB初期化/リセットの仕組みを提供する（テスト実行前にDBをクリーンな状態にする）。

### 3.2 削除ファイル

#### 3.2.1 削除: `frontend/e2e/dev-sample.spec.ts`

ユーザー指示により削除。新規E2E（`smoke.spec.ts`）に置換。

> **注**: `frontend/app/dev-sample/` ページ本体とそのコンポーネントテスト（`frontend/app/dev-sample/__tests__/page.test.tsx`）は**残置**する。これらはユニット/コンポーネントテストとして有効であり、E2Eテストとは独立して実行される。

### 3.3 新規作成ファイル

#### 3.3.1 新規: `frontend/e2e/global-setup.ts`

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

#### 3.3.2 新規: `frontend/e2e/smoke.spec.ts`

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

### 3.4 変更ファイル

#### 3.4.1 変更: `frontend/playwright.config.ts`

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

#### 3.4.2 変更: `frontend/package.json`

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

#### 3.4.3 変更: `frontend/TESTING_GUIDE.md`

E2Eセクションを「フルスタック結合テスト」向けに書き換える。

**変更内容:**

- セクション3「ブラウザ結合テスト（E2E）を作成する (Playwright)」を更新:
  - Docker Compose で全サービス（frontend + backend + DB）を起動してからテスト実行するワークフローを記述。
  - `globalSetup` の役割（バックエンド + DB の起動待機）を説明。
  - テストDBの概念（`DB_NAME=synctask_e2e` でテスト専用DBを使用可能）を記述。
- セクション5「テストの実行・デバッグ方法」のE2E部分を更新:
  - `docker compose up` → `npm run test:e2e` の2ステップワークフローを明記。
  - 未起動時のエラーメッセージと対処法を記述。

### 3.5 新規ドキュメント

#### 3.5.1 新規: `docs/backend/test-guide/README.md`

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
```

---

## 4. ファイル変更一覧

### 4.1 新規作成ファイル

| ファイルパス | 内容 |
|---|---|
| `frontend/e2e/global-setup.ts` | Playwright globalSetup: バックエンド + DB 起動待機 |
| `frontend/e2e/smoke.spec.ts` | フルスタックスモークE2Eテスト（4シナリオ） |
| `docs/backend/test-guide/README.md` | バックエンド + 結合テストガイド |

### 4.2 変更ファイル

| ファイルパス | 変更内容 |
|---|---|
| `frontend/playwright.config.ts` | `globalSetup` 追加、`webServer.reuseExistingServer: true` に変更 |
| `frontend/TESTING_GUIDE.md` | E2Eセクションをフルスタック結合テスト向けに更新 |

### 4.3 削除ファイル

| ファイルパス | 理由 |
|---|---|
| `frontend/e2e/dev-sample.spec.ts` | ユーザー指示により削除。新規E2E（`smoke.spec.ts`）に置換。 |

---

## 5. 実行順序

### Phase 2: E2Eテスト環境構築 + スモークテスト

1. `frontend/e2e/dev-sample.spec.ts` 削除
2. `frontend/e2e/global-setup.ts` 新規作成
3. `frontend/e2e/smoke.spec.ts` 新規作成
4. `frontend/playwright.config.ts` 変更（`globalSetup` 追加、`reuseExistingServer` 変更）
5. `frontend/TESTING_GUIDE.md` 更新
6. `docs/backend/test-guide/README.md` 新規作成
7. **検証**: `docker compose up` → `cd frontend && npm run test:e2e` でスモークテスト通過確認 → `npm run test:unit` でユニットテスト通過確認

---

## 6. 検証コマンド

実装完了後に以下のコマンドで動作確認を行う。

### 6.1 E2Eテスト

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

### 6.2 テスト用DBでのE2Eテスト

```bash
# テスト用DBで全サービス起動
DB_NAME=synctask_e2e docker compose up --build

# E2Eテスト実行
cd frontend && npm run test:e2e
```

---

## 7. 注意事項・制約

1. **フロントエンドのAPI呼び出しは実装しない**: Task 2のスコープは「E2Eテストインフラ構築 + スモークテスト」。フロントエンドからバックエンドAPIを呼び出す実装（APIクライアント、fetchラッパー等）は含まない。スモークテストは Playwright の `request` API でバックエンドを直接叩いて検証する。
2. **`dev-sample` ページは残置**: `frontend/app/dev-sample/` とコンポーネントテスト（`__tests__/page.test.tsx`）はユニット/コンポーネントテストとして残す。削除対象はE2Eの `dev-sample.spec.ts` のみ。

---

## 8. 将来の拡張ポイント

- **フロントエンドAPI統合**: フロントエンドからバックエンドAPIを呼び出すAPIクライアントを実装し、状態管理をクライアントメモリからAPI経由に移行する。
- **E2Eテスト拡充**: API実装後に、ログイン → タスク作成 → 一覧表示 → ステータス更新 → 削除などの実データフローE2Eテストを `smoke.spec.ts` に追加する。
- **CI/CD**: GitHub Actions 等で `docker compose up` → `npm run test:e2e` のワークフローを自動化する。
