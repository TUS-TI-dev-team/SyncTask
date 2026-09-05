# 01. 共通基盤セットアップ手順 (Foundation Setup)

本手順書は、Phase 2 の並列タスクをスムーズに実行するための **Phase 1（共通基盤セットアップ）** の具体的な手順を定義します。

---

## 1. 依存ライブラリの追加 (`package.json`)

`SyncTask/frontend/package.json` に、デザインシステムおよびUI実装に必要なライブラリを追加してインストールします。

### 追加する依存関係

```json
{
  "dependencies": {
    "@base-ui/react": "^1.5.0",
    "class-variance-authority": "^0.7.1",
    "clsx": "^2.1.1",
    "lucide-react": "^1.16.0",
    "next": "16.2.6",
    "next-themes": "^0.4.6",
    "react": "19.2.4",
    "react-dom": "19.2.4",
    "shadcn": "^4.8.0",
    "sonner": "^2.0.7",
    "tailwind-merge": "^3.3.1",
    "tw-animate-css": "^1.4.0"
  },
  "devDependencies": {
    "@tailwindcss/postcss": "^4",
    "@types/node": "^20",
    "@types/react": "^19",
    "@types/react-dom": "^19",
    "eslint": "^9",
    "eslint-config-next": "16.2.6",
    "postcss": "^8.5",
    "tailwindcss": "^4",
    "typescript": "^5"
  }
}
```

### インストールコマンド
```bash
cd frontend
npm install
```

---

## 2. shadcn 設定およびスタイリングの構成

### 2.1 `components.json` の配置
プロジェクトルート（`frontend/components.json`）を作成します。

```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "base-nova",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "",
    "css": "app/globals.css",
    "baseColor": "neutral",
    "cssVariables": true,
    "prefix": ""
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "lib": "@/lib",
    "hooks": "@/hooks"
  },
  "iconLibrary": "lucide"
}
```

### 2.2 `app/globals.css` の更新
Blueprint ダークテーマ（OKLCHカラー変数、タイポグラフィ、アニメーション設定）を `SyncTask-Design-Idea/app/globals.css` から移植します。

---

## 3. ユーティリティとモックストアの移植

### 3.1 `lib/utils.ts`
クラス名結合ヘルパー `cn()` を配置します。

```typescript
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

### 3.2 `lib/store.tsx`
全画面のモック動作に必要なデータ構造・React Context（`StoreProvider`, `useStore`）を `SyncTask-Design-Idea/lib/store.tsx` から移植します。

- **型定義**: `Task`, `Priority`, `Status`, `Profile` 等
- **初期モックデータ**: `SEED_TASKS`, `SEED_PROFILE`
- **操作メソッド**:
  - `addTask(task)`
  - `updateTask(id, patch)`
  - `deleteTask(id)`
  - `togglePin(id)`
  - `updateProfile(patch)`

---

## 4. shadcn/UI プリミティブコンポーネントの配置 (`components/ui/`)

`SyncTask-Design-Idea/components/ui/` から以下のコンポーネントを `SyncTask/frontend/components/ui/` にコピー・配置します。

- `alert-dialog.tsx`
- `badge.tsx`
- `button.tsx`
- `card.tsx`
- `dialog.tsx`
- `dropdown-menu.tsx`
- `input.tsx`
- `label.tsx`
- `select.tsx`
- `separator.tsx`
- `sonner.tsx`
- `tabs.tsx`
- `textarea.tsx`

---

## 5. 共通レイアウト & ルートレイアウトの配置

### 5.1 `components/layouts/`
- `app-header.tsx`: ログイン後のヘッダー（ナビゲーションリンク、ユーザー表示、ログアウト）
- `app-shell.tsx`: ログイン後画面の共通ラッパー（Header + Container + Toast）
- `auth-shell.tsx`: 未ログイン画面の共通ラッパー（認証カードスタイル）

### 5.2 `components/common/`
- `pagination.tsx`: 要件定義に準拠したページネーション操作UI（N > 10 時の省略表示等）

### 5.3 `app/layout.tsx`
`RootLayout` 内で `StoreProvider` および `Toaster`（sonner）を配置し、全ページで状態管理・トースト通知が機能するようにします。

---

## 6. Phase 1 完了チェックリスト

- [ ] `npm install` がエラーなく完了している
- [ ] `lib/utils.ts`, `lib/store.tsx` が配置されている
- [ ] `components/ui/` 配下の全13コンポーネントが配置されている
- [ ] `components/layouts/` 配下のレイアウト部品が配置されている
- [ ] `app/layout.tsx` が `StoreProvider` と `Toaster` を含んで更新されている
- [ ] `npx tsc --noEmit` で共通基盤に型エラーがないことを確認
