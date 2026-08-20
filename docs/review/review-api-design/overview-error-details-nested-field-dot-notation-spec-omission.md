# 概要書のエラーレスポンス定義におけるネストされたフィールドのドット記法仕様の記載漏れ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:33:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の「1.3 共通エラーレスポンス構造」におけるスキーマ定義で、`error.details[].field` は単なる「エラー対象のフィールド名」と説明されており、リクエストボディ内にネストされたオブジェクトが含まれる場合（例: `04_tasks.md` における `recurring_rule` オブジェクト）のフィールド名表記ルール（ドット記法）に関する共通仕様が明記されていない。

## 2. 詳細な指摘内容
`04_tasks.md`（3.3.2 `POST tasks`）等のタスク作成 API では、`recurring_rule` などのネストされたオブジェクト形式のリクエストボディが定義されています。これらのネストされたフィールドに対してバリデーションエラーが発生した場合、`error.details[].field` には `recurring_rule.due_time` や `recurring_rule.days_of_week` のようにドット記法（`parent.child`）を用いたパス文字列が格納されます。

しかし、`01_overview.md` の 1.3 共通エラーレスポンス構造の例示およびスキーマ定義テーブル（61〜67行目）には `email` のような単一階層のフィールド例しか掲載されておらず、ネストされたフィールドのエラー表現方法がドット記法に従う旨が共通仕様として明記されていません。

## 3. 推奨される修正案
`01_overview.md` の「1.3 共通エラーレスポンス構造」の `error.details[].field` の説明に以下の補足を追記してください。

```markdown
| `error.details[].field` | string | ○ | エラー対象のフィールド名。リクエストボディがネストされたオブジェクト構造を持つ場合は、`recurring_rule.due_time` や `recurring_rule.days_of_week` のようにドット記法（`親オブジェクト.子フィールド`）で指定 |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:38:45
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/01_overview.md` の「1.3 共通エラーレスポンス構造」におけるスキーマ定義テーブルの `error.details[].field` の説明に、リクエストボディがネストされたオブジェクト構造を持つ場合のドット記法表記ルール（例: `recurring_rule.due_time` や `recurring_rule.days_of_week`）についての補足を明記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)

