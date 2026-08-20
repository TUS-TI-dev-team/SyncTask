# フロントエンド テストコード作成ガイド (Testing Guide)

このドキュメントでは、SyncTask フロントエンドにおける**テストコード作成の一連の流れ・書き方・ベストプラクティス**を解説します。
初心者の方でも迷わずテストを追加・運用できるようにステップ形式でまとめています。

---

## 📑 目次

1. [テストの全体像と使い分け](#1-テストの全体像と使い分け)
2. [ファイル配置と命名規約](#2-ファイル配置と命名規約)
3. [Code-as-Docs の基本方針](#3-code-as-docs-の基本方針)
4. [テスト作成の流れ（ステップバイステップ）](#4-テスト作成の流れステップバイステップ)
   - [Step 1: ロジック関数の単体テストを作成する](#step-1-ロジック関数の単体テストを作成する-vitest)
   - [Step 2: UIコンポーネントのテストを作成する](#step-2-uiコンポーネントのテストを作成する-react-testing-library)
   - [Step 3: ブラウザ結合テスト（E2E）を作成する](#step-3-ブラウザ結合テストe2eを作成する-playwright)
5. [テストの実行・デバッグ方法](#5-テストの実行デバッグ方法)
6. [テスト作成時のベストプラクティス](#6-テスト作成時のベストプラクティス)

---

## 1. テストの全体像と使い分け

SyncTask では、役割に応じて2種類のテストツールを使い分けています。

| テスト分類 | 対象 | 使用ツール | 目的・特徴 |
| :--- | :--- | :--- | :--- |
| **単体テスト (Unit)** | 純粋な関数・ビジネスロジック | Vitest | 計算や状態遷移などの境界値・異常系を網羅的かつ超高速に検証 |
| **コンポーネントテスト** | React コンポーネント | Vitest + React Testing Library | DOM構造、初期描画、イベント発火、ユーザー操作時の状態変化を検証 |
| **ブラウザ結合テスト (E2E)** | 実ブラウザでの画面遷移・シナリオ | Playwright | 実際に Next.js サーバーと Chromium ブラウザを起動し、複数画面やエンドツーエンドの操作を検証 |

---

## 2. ファイル配置と命名規約

テスト対象ファイルと同じ階層、または近接した場所にテストコードを配置します。

```
frontend/
├── lib/
│   ├── sample.ts                   # 実装コード
│   └── __tests__/                  # 単体テスト用ディレクトリ
│       └── sample.test.ts          # テストファイル: <対象名>.test.ts
├── app/
│   └── (機能名)/
│       ├── page.tsx                # ページコンポーネント
│       └── __tests__/              # コンポーネントテスト用ディレクトリ
│           └── page.test.tsx       # テストファイル: <対象名>.test.tsx
└── e2e/
    └── (機能名).spec.ts            # E2Eテストファイル: <対象名>.spec.ts
```

- **単体・コンポーネントテスト**: `__tests__/<name>.test.ts(x)`
- **E2Eブラウザテスト**: `e2e/<name>.spec.ts`

---

## 3. Code-as-Docs の基本方針

SyncTask では**「コードとテスト自体が仕様書（ドキュメント）になる」**という方針（Code-as-Docs）を採用しています。

1. **実装コード**: TSDoc / JSDoc の `@spec` を使い、事前条件・事後条件・例外仕様を記述する。
2. **テストコード**: `describe`（仕様のカテゴリ）と `it`（個別の振る舞い・期待値）に、日本語で明確な仕様文を記述する。

```typescript
// 良いテストコードの記述例
describe('createTask() 仕様', () => {
  it('正常系: 有効なタイトルと優先度で初期状態 todo のタスクが生成されること', () => { ... });
  it('異常系: タイトルが空文字列の場合は Error を送出すること', () => { ... });
});
```

---

## 4. テスト作成の流れ（ステップバイステップ）

新しい機能（例: フィルター機能や新規画面）を追加する際の一連のテスト作成手順です。

### Step 1: ロジック関数の単体テストを作成する (Vitest)

純粋関数（ユーティリティや状態操作ロジック）を作成したら、まず単体テストを書きます。

#### 1. 実装ファイル (`lib/task-filter.ts`)
```typescript
/**
 * @spec
 * 1. keyword が指定された場合、タイトルに部分一致するタスクのみを返す。
 * 2. 大文字・小文字は区別しない（case-insensitive）。
 */
export function searchTasks(tasks: Task[], keyword: string): Task[] {
  const trimmed = keyword.trim().toLowerCase();
  if (!trimmed) return tasks;
  return tasks.filter((t) => t.title.toLowerCase().includes(trimmed));
}
```

#### 2. テストファイル (`lib/__tests__/task-filter.test.ts`)
```typescript
import { describe, it, expect } from 'vitest';
import { searchTasks } from '../task-filter';

describe('task-filter.ts - タスク検索ロジック仕様', () => {
  const mockTasks = [
    { id: '1', title: 'React 学習', status: 'todo' },
    { id: '2', title: 'Next.js 開発', status: 'done' },
  ];

  it('正常系: 指定したキーワードに部分一致するタスクを抽出できること', () => {
    const result = searchTasks(mockTasks as any, 'React');
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('1');
  });

  it('正常系: 大文字・小文字を区別せず検索できること', () => {
    const result = searchTasks(mockTasks as any, 'react');
    expect(result).toHaveLength(1);
  });

  it('境界値: 空文字が渡された場合は全タスクを返すこと', () => {
    const result = searchTasks(mockTasks as any, '   ');
    expect(result).toHaveLength(2);
  });
});
```

---

### Step 2: UIコンポーネントのテストを作成する (React Testing Library)

ユーザー操作や画面の表示を検証するコンポーネントテストを作成します。

#### ポイント:
- 要素の特定には `data-testid` またはアクセシビリティロール（`getByRole` / `getByLabelText`）を使用します。
- クリックやキーボード入力は `@testing-library/user-event` の `userEvent.setup()` を使用します。

#### テストファイル例 (`app/dev-sample/__tests__/page.test.tsx`)
```tsx
import React from 'react';
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DevSamplePage from '../page';

describe('DevSamplePage UI コンポーネント仕様', () => {
  it('初期レンダリング: タイトルと初期要素が表示されること', () => {
    render(<DevSamplePage />);
    expect(screen.getByTestId('page-title')).toHaveTextContent('開発テスト用サンプルページ');
  });

  it('タスク追加: フォームに入力して送信するとリストに追加されること', async () => {
    const user = userEvent.setup();
    render(<DevSamplePage />);

    const input = screen.getByTestId('task-title-input');
    const submitBtn = screen.getByTestId('add-task-button');

    await user.type(input, '新しいタスク');
    await user.click(submitBtn);

    expect(screen.getByText('新しいタスク')).toBeInTheDocument();
  });
});
```

---

### Step 3: ブラウザ結合テスト（E2E）を作成する (Playwright)

ページ全体の結合動作や、ブラウザ上でのエンドツーエンドシナリオをテストします。

#### テストファイル例 (`e2e/task-flow.spec.ts`)
```typescript
import { test, expect } from '@playwright/test';

test.describe('タスク管理一連のユーザーフロー (E2E)', () => {
  test.beforeEach(async ({ page }) => {
    // 対象ページへ遷移
    await page.goto('/dev-sample');
  });

  test('シナリオ: タスクを追加して状態を完了に変更する一連の操作', async ({ page }) => {
    // 1. タスクを入力して追加
    await page.getByTestId('task-title-input').fill('E2E テストタスク');
    await page.getByTestId('add-task-button').click();

    // 2. 一覧に追加されたことを検証
    const taskItem = page.getByTestId('task-item').filter({ hasText: 'E2E テストタスク' });
    await expect(taskItem).toBeVisible();

    // 3. ステータスをトグルして完了にする
    const toggleBtn = taskItem.getByTestId('toggle-status-button');
    await toggleBtn.click(); // todo -> in_progress
    await toggleBtn.click(); // in_progress -> done

    // 4. 完了状態の反映を確認
    await expect(taskItem.getByTestId('task-status')).toHaveText('done');
    await expect(page.getByTestId('stat-done')).toHaveText('1');
  });
});
```

---

## 5. テストの実行・デバッグ方法

開発フェーズに合わせて使い分けます。

### 1. 開発中の高速フィードバック (Watch モード)
コード修正時に自動で単体テストを再実行します。

```bash
npm run test:unit:watch
```

### 2. ブラウザの動きを目視しながらデバッグ (Playwright UI モード)
ブラウザ画面のスクリーンショット、タイムライン、各ステップのDOMツリーを確認しながらデバッグできます。

```bash
npm run test:e2e:ui
```

### 3. コミット前・PR作成前の一括チェック
単体テストとE2Eテストの両方がすべてパスすることを確認します。

```bash
npm test
```

---

## 6. テスト作成時のベストプラクティス

1. **テストを独立させる (Isolation)**:
   - 各テストケース（`it`）は他のテストの実行順序や状態に依存しないように独立させてください。
2. **実装の詳細ではなく「振る舞い」をテストする**:
   - `useState` の内部変数などを直接チェックするのではなく、ユーザーが見る画面（DOM）の変化を検証します。
3. **`data-testid` を適切に配置する**:
   - デザイン変更（CSSクラス名の変更など）でテストが壊れないよう、テスト対象の主要要素には `data-testid="..."` を付与します。
4. **非同期処理の待機**:
   - Playwright の `expect(locator).toHaveText(...)` などの web-first assertions は自動的に要素の描画を待機（リトライ）するため、手動の `sleep` は避けましょう。
