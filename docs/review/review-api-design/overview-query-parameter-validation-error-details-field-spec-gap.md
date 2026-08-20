# GET リクエスト等のクエリパラメータバリデーションエラーにおける error.details[].field 記法の規定不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:43:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` 1.3 節の「エラーレスポンス スキーマ定義」において、`error.details[].field` の説明欄に「リクエストボディがネストされたオブジェクト構造を持つ場合は、`recurring_rule.due_time` や `recurring_rule.days_of_week` のようにドット記法で指定」と記載されているが、`GET tasks` などのクエリパラメータバリデーション違反（例: `priority`, `page`, `start_date` の不備）時に `field` フィールドへ設定すべき文字列に関する規定が明確に示されていない。

## 2. 詳細な指摘内容
1. **スキーマ定義テーブルの現在の記述（`01_overview.md` L70）**:
   `| error.details[].field | string | ○ | エラー対象のフィールド名。リクエストボディがネストされたオブジェクト構造を持つ場合は、recurring_rule.due_time や recurring_rule.days_of_week のようにドット記法（親オブジェクト.子フィールド）で指定 |`

2. **問題点**:
   `GET tasks` のような `GET` リクエストではリクエストボディが存在せず、クエリパラメータ（`priority`, `status`, `view_type`, `sort_by`, `page`, `limit`, `start_date`, `end_date`, `due_date`, `keyword`）の形式・範囲不正に伴う `400 BAD_REQUEST` が返却される。
   この際、`error.details[].field` にクエリパラメータ名が設定されるべきであることが全共通仕様に明記されていない場合、フロントエンド側でリクエストボディエラーとクエリパラメータエラーのハンドリングロジックを統一する際に曖昧さが生じる。

## 3. 推奨される修正案
`01_overview.md` 1.3 節の `error.details[].field` の説明文に `GET` リクエストのクエリパラメータに関する記述を追加し、仕様を補正してください。

```markdown
| `error.details[].field` | string | ○ | エラー対象のフィールド名またはクエリパラメータ名。リクエストボディがネストされたオブジェクト構造を持つ場合は `recurring_rule.due_time` のようにドット記法（`親オブジェクト.子フィールド`）で指定し、`GET` リクエスト等のクエリパラメータバリデーション違反時は対象のクエリパラメータ名（例: `priority`, `page`, `start_date`）を指定 |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.3 節の `error.details[].field` の説明文に、`GET` リクエスト等でクエリパラメータバリデーション違反が発生した場合は対象のクエリパラメータ名（例: `priority`, `page`, `start_date`）を指定する旨を明確に規定追加しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
