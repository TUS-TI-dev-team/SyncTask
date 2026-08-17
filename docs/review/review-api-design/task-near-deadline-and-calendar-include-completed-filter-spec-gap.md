# near_deadline ビューおよびカレンダー期間取得における include_completed フィルタ挙動の不備

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:30:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`GET tasks` において、`view_type=near_deadline` 指定時の完了タスク除外ルール、および `start_date` / `end_date` 指定（カレンダーグリッド取得）時の `include_completed` フィルタのデフォルト挙動についての仕様が曖昧です。

## 2. 詳細な指摘内容
1. **`near_deadline` ビューにおける `include_completed` 優先度の未定義**:
   - `requirements.md` (L108-110) では、「締切間近タスク表示」の対象ステータスは「未着手」「進行中」のみであり、完了タスクは厳格に除外する（`include_completed` オプションの対象外）と定義されています。
   - `04_tasks.md` L14 では `view_type=near_deadline` と `include_completed`（デフォルト: `false`）が個別のパラメータとして並記されており、クライアントが `view_type=near_deadline&include_completed=true` を送信した場合に、完了タスクを含めるのか除外するのかの優先順位が仕様書に明記されていません。

2. **カレンダー期間取得における `include_completed` のデフォルト値の曖昧さ**:
   - `requirements.md` (L139-142) では、カレンダー表示において「締切日時が設定されている全ステータス（未着手・進行中・完了）のタスクを表示する」ことを基本ルールとしつつ、表示切り替えフィルタを提供するとしています。
   - `04_tasks.md` L15 では `include_completed` のデフォルト値が `false` と定義されています。このため、カレンダー画面が `start_date=2026-08-01&end_date=2026-08-31` のように期間クエリのみを指定してリクエストした場合、デフォルトで完了タスクがすべて除外されてしまい、カレンダーの初期表示要件を満たさなくなる恐れがあります。

## 3. 推奨される修正案
1. **`view_type=near_deadline` の制約追加**:
   - `04_tasks.md` の `view_type` パラメータ説明欄、およびレスポンス仕様注記に「`view_type=near_deadline` 指定時は `include_completed` パラメータの値に関わらず、常に完了タスク (`status=completed`) を除外します」と明確に追記してください。
2. **カレンダー期間取得時の `include_completed` 挙動の明確化**:
   - `start_date` および `end_date` を用いたカレンダーグリッド取得時の `include_completed` のデフォルト挙動（カレンダー表示時はデフォルト `true` とする、あるいはフロントエンドが明示的に `include_completed=true` を指定してリクエストする規約とする旨）をドキュメントに明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`view_type=near_deadline` 指定時は `include_completed` の設定に関わらず常に完了タスクを除外する優先ルールを追加し、カレンダー期間取得（`start_date`, `end_date` 指定）時は `include_completed` のデフォルトが `true` となる旨を追記しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
