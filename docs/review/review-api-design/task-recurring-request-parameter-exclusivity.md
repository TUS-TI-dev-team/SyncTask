# タスク作成APIにおける単一作成パラメータと繰り返し作成ルールの適用条件の曖昧さ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 13:17:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`POST tasks`（新規タスク作成API）において、単一タスク作成時に使用する `due_datetime` と、繰り返し一括作成時に使用する `recurring_rule`（およびその内部の `due_time`）の双方が同一リクエスト例に混在しており、`is_recurring: true` / `false` 時の各パラメータの適用優先度やバリデーション規則が曖昧です。

## 2. 詳細な指摘内容
1. **リクエスト例におけるパラメータ混在**:
   - `docs/design/api_design.md` L562-575 の `POST tasks` リクエスト例において、`due_datetime: "2026-08-22T18:00:00+09:00"` と `is_recurring: true`、`recurring_rule`（内部に `due_time: "18:00"`）が同時に指定されています。
   - `is_recurring: true` の場合、各タスクの締切日時は `recurring_rule.start_date` 〜 `end_date`、`days_of_week`、および `recurring_rule.due_time` から算出されます。トップレベルの `due_datetime` が同時に渡された場合に無視されるのか、バリデーションエラーとなるのかが不明確です。
2. **単一作成時と繰り返し一括作成時のスキーマ注記不足**:
   - `is_recurring: false`（または未指定）時は `due_datetime`（任意、省略時は当日 23:59:00）を使用し `recurring_rule` は無視または禁止。
   - `is_recurring: true` 時は `recurring_rule` が必須となり、トップレベルの `due_datetime` は指定不要（または無視）であることを明確にする必要があります。

## 3. 推奨される修正案
1. `POST tasks` のスキーマ定義において、単一作成（`is_recurring: false`）と繰り返し一括作成（`is_recurring: true`）のパラメータ適用条件・排他制御ルールを明記してください。
2. リクエスト例を「単一作成時」と「繰り返し一括作成時」の2パターンに分けて記載するか、注記を追加して実装者が誤解しないように整理してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` の `POST tasks` において、リクエストボディ例を「単一タスク作成時」と「毎週繰り返し一括作成時」の2つに明確に分離しました。
- フィールド定義表において、`is_recurring: true` の際はトップレベルの `due_datetime` が指定されていても無視され、`recurring_rule`（および内部の `due_time`）に従って各タスクの締切日時が生成される仕様を明記しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
