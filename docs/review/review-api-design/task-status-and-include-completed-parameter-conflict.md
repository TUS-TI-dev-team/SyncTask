# GET tasksにおける status パラメータと include_completed パラメータの矛盾・優先度未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` において、ステータス直接絞り込みパラメータ `status`（`not_started`, `in_progress`, `completed`）と、完了フラグ切り替えパラメータ `include_completed`（デフォルト: `false`）が同時に送信された際の優先順位および評価ルールが定義されていません。

## 2. 詳細な指摘内容
1. **`status=completed` と `include_completed=false` の評価矛盾**:
   - `include_completed` のデフォルト値は `false`（完了タスクを除外）です。
   - クライアントが `GET /api/tasks?status=completed` と明示的に完了タスクの取得を要求してリクエストした場合、`include_completed` のデフォルト値 `false` が同時に適用されると、`WHERE status = 'completed' AND status != 'completed'` と同義になり、検索結果が常に0件となる不具合が発生します。

2. **`status` と `include_completed` の優先関係の未記述**:
   - クライアントが `status=in_progress` かつ `include_completed=true` を送信した場合に、`in_progress` のみを取得するのか、完了タスクも含めて取得するのかの条件結合ロジックが曖昧です。

## 3. 推奨される修正案
1. `GET tasks` のパラメータ仕様において、以下の優先関係・評価ルールを明記してください:
   - `status` パラメータが明示的に指定された場合（`not_started`, `in_progress`, `completed`）は、`include_completed` パラメータの指定（またはデフォルト値 `false`）を無視し、`status` の指定条件を最優先で適用する。
   - `status` パラメータが未指定（省略）の場合にのみ、`include_completed` パラメータの指定に従って完了タスクの含め/除外を制御する（`false` 時は `status IN ('not_started', 'in_progress')`）。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`GET tasks` において `status` パラメータが明示指定された場合（`not_started`, `in_progress`, `completed`）は、`include_completed` の設定を無視して `status` 条件を最優先評価する評価規則を明記しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
