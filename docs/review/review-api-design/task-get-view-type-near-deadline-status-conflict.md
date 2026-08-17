# GET tasks における view_type=near_deadline と status=completed のパラメータ競合・優先順位未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`04_tasks.md` の `GET tasks` において、`view_type=near_deadline`（締切間近ビュー）は「常に完了タスクを除外する」と定義されている一方で、`status` パラメータは「明示指定時は include_completed を無視して指定ステータスを最優先適用」と定義されています。`view_type=near_deadline` と `status=completed` が同時に指定された場合の優先関係および挙動（400エラー、0件返却、またはstatus優先）が仕様書上で未定義となっています。

## 2. 詳細な指摘内容
1. **仕様定義の衝突**:
   - `view_type=near_deadline` の説明（L14）: 「72時間以内/期限超過。`include_completed` の指定に関わらず常に完了タスクは除外」
   - `status` の説明（L18）: 「明示指定時は `include_completed` の指定を無視して指定ステータスを最優先適用」
   - `requirements.md` L109: 「締切間近タスク表示の対象ステータス: 『未着手』『進行中』のタスクのみ（『完了』タスクは除外）」
   - クライアントが `GET /tasks?view_type=near_deadline&status=completed` をリクエストした場合、`view_type` の除外ルールと `status` の最優先適用の指定が直接バッティングします。

2. **フロントエンド・バックエンドの実装乖離リスク**:
   - 挙動が明記されていないため、バックエンド実装において「0件返却（AND検索扱い）」、「400 Bad Request 返却（パラメータ競合エラー）」、「`status=completed` の結果返却」など実装者によって処理が割れる原因となります。

## 3. 推奨される修正案
`04_tasks.md` の `view_type` および `status` または「リクエスト評価順序 / Errors」セクションに、パラメータ併用時の明確な仕様を追記してください。
業務要件（締切間近表示は未完了タスクのみが対象）に合わせ、`view_type=near_deadline` と `status=completed` の同時指定は 400 Bad Request（code: `"BAD_REQUEST"`）とするか、あるいは `view_type=near_deadline` 指定時は `status=completed` が指定されても完了タスク除外ルールが最優先され 0件を返却する旨を明記してください。
（推奨: パラメータ不整合として 400 Bad Request を返却するか、`view_type=near_deadline` において `status=completed` は不可である旨をパラメータ説明に明確化）

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 節 (`GET tasks`) の `view_type` パラメータ説明、`status` パラメータ説明、リクエスト評価順序、および Errors セクションに、`view_type=near_deadline`（完了タスク除外）と `status=completed` の同時指定はパラメータ競合・不整合エラーとして 400 Bad Request (code: `"BAD_REQUEST"`) を返却する旨を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
