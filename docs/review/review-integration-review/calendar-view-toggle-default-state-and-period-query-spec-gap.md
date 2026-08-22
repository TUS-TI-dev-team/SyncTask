# カレンダー表示切り替え時における完了タスク表示トグルのデフォルト初期値および期間クエリ仕様の曖昧さ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L33)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L41)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L49-L55)
  - [04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md#L14)
  - [04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md#L19-L21)

## 1. 問題の概要
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) において、リスト表示とカレンダー表示の切り替え時における「完了タスク表示/非表示」トグルの初期値（デフォルトON/OFF）、およびカレンダー表示（月全体表示 vs 週表示）におけるAPIリクエスト（`start_date` / `end_date`）の送信仕様が明確に定義されていません。

## 2. 詳細な指摘内容
1. **完了タスク表示トグルの初期値の食い違い**:
   - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 41行目: 高優先度・ピン留め・カレンダー等で「完了タスク表示/非表示」切り替えトグルを提供する旨が書かれていますが、リスト表示時とカレンダー表示時でのトグルの初期状態について言及がありません。
   - [04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md) 14行目では、`include_completed` のデフォルト値が「通常一覧: `false`」「カレンダー（start_date/end_date指定時）: `true`」と定義されています。
   - 画面上でリスト表示からカレンダー表示に切り替えた際、トグルが自動的にON（完了表示）になるのか、あるいはトグルの状態が保持されるのかが曖昧です。
2. **月全体表示（グリッド）と週表示での期間クエリパラメータ**:
   - カレンダー月全体表示（5〜6週グリッド、最大42日）および週表示（7日間）において、フロントエンドがバックエンドに発行する `GET /api/tasks?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` の期間算出ルール（前月末日/翌月初日を含むグリッド全範囲を送信する等）が画面補足に明記されていません。

## 3. 推奨される修正案
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の「カレンダー表示および日付詳細ポップアップ」（49〜55行目）に以下の仕様を補足・明記してください。

```markdown
- **カレンダー表示時の完了タスク初期値と期間取得API仕様**:
  - カレンダー表示切り替え時は、API仕様（`include_completed` のカレンダー時デフォルト `true`）に合わせ、「完了タスク表示/非表示」トグルスイッチは初期状態で ON（完了タスクを含む）となる。
  - カレンダーグリッド描画時は、表示対象月/週のグリッドに含まれる全日付範囲（月全体表示時は前月末・翌月初を含む最大42日間、週表示時は日〜土の7日間）を算出し、`GET /api/tasks?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&include_completed={true|false}` を発行してタスクデータを一括取得する。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` の「カレンダー表示および日付詳細ポップアップ」補足欄に、カレンダー表示切り替え時の「完了タスク表示/非表示」トグル初期値（ON）および月/週表示時のグリッド全日付範囲期間取得API仕様を明記しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)

