# POST tasks における is_recurring: true 指定時の recurring_rule 欠落・型不正エラー詳細の定義漏れ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`04_tasks.md` の `POST tasks` において、`is_recurring: true` 設定時に必須となる `recurring_rule` オブジェクト自体が欠落していたり、`null` や非オブジェクト型（数値・文字列等）で送信された場合のレスポンス `error.details` のフィールド名およびメッセージ定義が欠落しています。

## 2. 詳細な指摘内容
1. **エラー詳細スキーマにおける `recurring_rule` 自体の未定義**:
   - `04_tasks.md` L190-L194 の Errors セクションでは、生成件数0件時（`field: "recurring_rule"`）、生成件数100件超過時（`field: "recurring_rule"`）、`due_time` 不正時（`field: "recurring_rule.due_time"`）のエラー詳細形式やドット記法ルールが記載されています。
   - しかし、`is_recurring: true` でありながら `recurring_rule` キーが省略されていたり、`"recurring_rule": null` や `"recurring_rule": "invalid"` が送信された場合の `error.details` 項目（例: `field: "recurring_rule"`, `message: "is_recurring: true の場合、recurring_rule は必須です"` 等）についての具体的な指定がありません。

2. **フロントエンドのエラー表示・ハイライトへの影響**:
   - `01_overview.md` L70 では `error.details[].field` によるエラー箇所の特定が定められているため、トップレベルオブジェクト `recurring_rule` 自体の検証エラー表現を明確化しておく必要があります。

## 3. 推奨される修正案
`04_tasks.md` L190 以降の `POST tasks` の `Errors (400 Bad Request)` セクションに以下の定義を追加してください。

```markdown
- `is_recurring: true` 時の `recurring_rule` 欠落・null・非オブジェクト指定時:
  `error.details: [{ "field": "recurring_rule", "message": "is_recurringがtrueの場合、recurring_ruleオブジェクトの指定は必須です" }]`
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.2 節 (`POST tasks`) の Errors (400 Bad Request) セクションに、`is_recurring: true` 時の `recurring_rule` オブジェクト欠落・null・非オブジェクト型指定時の `error.details` 定義（`field: "recurring_rule"`, `message: "is_recurringがtrueの場合、recurring_ruleオブジェクトの指定は必須です"`）を追記・明確化しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
