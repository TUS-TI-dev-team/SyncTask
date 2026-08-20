# POST tasks における is_pinned パラメータの型・デフォルト値およびリクエストボディ定義の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:07:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` (3.3.2 節) のリクエストボディフィールド定義表において、ピン留めフラグ `is_pinned` パラメータの定義が欠落しており、新規タスク作成時にピン留め状態を指定可能かどうかの仕様が曖昧です。

## 2. 詳細な指摘内容
1. **リクエストボディ定義表における `is_pinned` の欠落**:
   - `04_tasks.md` L109-L120 の「Request Body フィールド定義」テーブルには `title`, `comment`, `priority`, `due_datetime`, `is_recurring`, `recurring_rule` のみが掲載されており、`is_pinned` フィールドが含まれていません。
2. **レスポンス仕様との整合性の欠如**:
   - 一方で `POST tasks` のレスポンス例（L138）には `"is_pinned": false` が返却値として記載されています。
   - タスク作成モーダル・UIにおいて、ユーザーが作成時に「ピン留めする」チェックボックスを有効化して `{"title": "...", "is_pinned": true}` を送信した場合に、サーバーがこれを許容して作成するか、無視して常に `false` とするか、あるいはバリデーションエラーとするかが仕様上未定義です。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 節の「Request Body フィールド定義」テーブルに `is_pinned` パラメータを追加し、以下の仕様を明記してください：

```markdown
| `is_pinned` | boolean | × | ピン留めフラグ（`true` / `false`）。省略時はデフォルト `false` として作成 |
```

また、`is_recurring: true` 時の `is_pinned` の扱い（一括生成されるすべてのタスクに適用されるか等）についても注記を補足してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.2 節 (`POST tasks`) の Request Body フィールド定義テーブルに `is_pinned` パラメータ（boolean, 任意, デフォルト `false`）を追加し、`is_recurring: true` 時に指定された場合は一括生成されるすべてのタスクに適用される旨の注記を補足しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
