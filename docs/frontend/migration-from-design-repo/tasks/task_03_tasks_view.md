# Task 03: タスク一覧・カレンダー画面の移植 (Tasks & Calendar View)

## 1. 担当概要

本タスクは、タスク管理のメイン画面となるタスク一覧・カレンダー画面（`app/tasks`）およびそのビューコンポーネント（`components/tasks/tasks-view.tsx`, `task-calendar.tsx`, `task-card.tsx`）を `SyncTask-Design-Idea` から `SyncTask/frontend` に移植する作業です。

---

## 2. 移行対象ファイル一覧

| 移行元 (Design-Idea) | 移行先 (SyncTask/frontend) | 概要 |
| --- | --- | --- |
| `components/tasks/task-card.tsx` | `components/tasks/task-card.tsx` | リスト表示用タスクカード（優先度バッジ、ステータス変更、ピン留め、編集/削除メニュー） |
| `components/tasks/task-calendar.tsx` | `components/tasks/task-calendar.tsx` | カレンダーコンポーネント（月表示/週表示、ナビゲーション、セル内タスク操作） |
| `components/tasks/tasks-view.tsx` | `components/tasks/tasks-view.tsx` | タスク一覧メインビュー（検索、フィルター、リスト/カレンダー表示切替、ページネーション） |
| `app/tasks/page.tsx` | `app/tasks/page.tsx` | タスク画面ページエントリーポイント |

---

## 3. 仕様書・設計との整合ポイント (`screen_design.md`)

- **検索 & フィルター機能**:
  - キーワード検索入力欄。
  - 優先度絞り込みプルダウン（高 / 中 / 低 / 全て）。
  - ステータス絞り込みプルダウン（未着手 / 進行中 / 完了 / 全て）。
  - 締切日絞り込み入力欄。
  - 「完了タスク表示/非表示」切り替えトグル。
- **表示形式切り替え**:
  - 「リスト表示」と「カレンダー表示」を切り替え可能。
- **カレンダー表示の仕様**:
  - 「全体表示（月グリッド）」と「週表示（7日間）」の切り替え。
  - ナビゲーションボタン（全体: 前月 / 次月 / 今日、週: 前週 / 翌週 / 今日）。
  - カレンダーセル上のタスク直接ステータス変更・ピン止め操作（クリック時に日付詳細モーダルが誤発火しないよう `e.stopPropagation()` 等で制御）。
  - 日付セルクリックで「日付詳細ポップアップ」を開くトリガー連携。
- **リスト表示時のページネーション**:
  - 上下に配置し、$N > 10$ のページネーション仕様を適用。

---

## 4. 作業手順

1. `components/tasks/task-card.tsx` を作成・配置する。
2. `components/tasks/task-calendar.tsx` を作成・配置する。
3. `components/tasks/tasks-view.tsx` を作成・配置する。
4. `app/tasks/page.tsx` を作成する。
5. リスト表示とカレンダー表示の切り替え、各種フィルタリングの連動を確認する。

---

## 5. 完了確認チェックリスト

- [ ] `npx tsc --noEmit` で型エラーが発生しない
- [ ] ブラウザで `/tasks` にアクセスし、タスク一覧が表示される
- [ ] 検索キーワード入力やプルダウン選択でタスク一覧がフィルタリングされる
- [ ] 「カレンダー表示」に切り替え、月グリッドおよび週表示が正常に描画される
- [ ] カレンダーの「前月/前週」「次月/次週」「今日」ボタンで日付範囲が遷移する
- [ ] カレンダー内のタスクステータスやピン止めを直接切り替えられる
