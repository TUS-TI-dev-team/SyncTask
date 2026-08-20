# SyncTask Frontend

SyncTask のフロントエンドアプリケーション（Next.js / React / TypeScript / Tailwind CSS）。

---

## 🚀 はじめに (Getting Started)

開発サーバーの起動:

```bash
npm run dev
```

ブラウザで [http://localhost:3000](http://localhost:3000) を開くとトップページが表示されます。
また、開発・テスト動作確認用ページとして [http://localhost:3000/dev-sample](http://localhost:3000/dev-sample) が利用可能です。

---

## 🧪 テスト環境・実行方法 (Testing)

本プロジェクトでは **単体・コンポーネントテスト** および **ブラウザ実動作による結合・E2Eテスト** を統合したテスト環境を提供しています。

### 1. テストフレームワーク構成
- **単体・コンポーネントテスト**: [Vitest](https://vitest.dev/) + [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/) + [jsdom](https://github.com/jsdom/jsdom)
  - 高速な TypeScript ネイティブ実行
  - ドメインロジック、ユーティリティ関数、UIコンポーネントの挙動検証
- **ブラウザ結合テスト (E2E)**: [Playwright](https://playwright.dev/)
  - Chromium などの実ブラウザを起動し、UI操作、画面遷移、状態更新をトレース検証

### 2. 初回セットアップ (Playwright ブラウザの導入)
Playwright の初回実行前にブラウザバイナリをインストールします:

```bash
npx playwright install chromium
```

### 3. テストコマンド一覧

| コマンド | 説明 |
| :--- | :--- |
| `npm test` | **すべてのテストを一括実行**（単体テスト → ブラウザ結合テスト） |
| `npm run test:unit` | 単体・コンポーネントテストを実行（CLI / CI向け） |
| `npm run test:unit:watch` | 単体テストをウォッチモードで実行（コード変更時に自動再実行） |
| `npm run test:e2e` | Playwright によるブラウザ結合テストを実行（自動で dev サーバーを立ち上げてテスト） |
| `npm run test:e2e:ui` | Playwright UI モードで実行（ブラウザ操作の可視化・タイムライン確認・デバッグ用） |

### 4. テスト作成ガイド
新規機能やコンポーネントに対するテストコードの具体的な作成手順・テンプレート・ベストプラクティスについては、[TESTING_GUIDE.md](./TESTING_GUIDE.md) を参照してください。

---

## 📐 Code-as-Docs 思想と実装例

本プロジェクトでは仕様書と実装の乖離を防ぐため、**コードおよびテストコード自体に仕様を明記する「Code-as-Docs」** を採用しています。

- **ロジック仕様 (`lib/sample.ts`)**:
  - TSDoc / JSDoc にて `@spec` タグを用い、引数の事前条件・事後条件・例外仕様・戻り値を記述。
- **仕様テスト (`lib/__tests__/sample.test.ts`, `app/dev-sample/__tests__/page.test.tsx`)**:
  - `describe` / `it` を用いて、関数の振る舞い・境界値・バリデーションエラーの仕様をテストケースとして網羅。
- **ブラウザ結合テスト (`e2e/dev-sample.spec.ts`)**:
  - 実際のユーザー操作シナリオ（タスク追加、優先度設定、ステータストグルによる統計反映）をシナリオ形式で定義・検証。

---

## 📂 ディレクトリ構成 (テスト関連)

```
frontend/
├── app/
│   └── dev-sample/             # テスト検証用のサンプルページ
│       ├── *.tsx
│       └── __tests__/          # コンポーネント単体テスト
│           └── *.test.tsx
├── lib/
│   ├── *.ts               # サンプルロジック (TSDoc仕様付き)
│   └── __tests__/              # ロジック単体テスト
│       └── *.test.ts
├── e2e/
│   └── *.spec.ts      # Playwright ブラウザ結合テスト
├── vitest.config.mts           # Vitest 設定
├── vitest.setup.ts             # Vitest セットアップ (@testing-library/jest-dom)
├── playwright.config.ts        # Playwright 設定 (WebServer自動起動等)
└── README.md
```
