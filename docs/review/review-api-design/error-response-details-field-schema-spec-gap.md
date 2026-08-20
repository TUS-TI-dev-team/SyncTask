# 共通エラーレスポンスにおける details フィールドの型・省略・空状態の仕様未記述

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` 1.3 節「共通エラーレスポンス構造」で示されているエラー応答の JSON サンプルには `details` 配列が含まれているが、個別フィールドの入力バリデーションエラーを伴わないエラー（例: `401 UNAUTHORIZED`, `404 NOT_FOUND`, `410 GONE`, `429 RATE_LIMIT_EXCEEDED`, `500 INTERNAL_SERVER_ERROR` 等）において `details` フィールドがどのように返却されるか（空配列 `[]` として返却されるのか、`null` になるのか、キー自体が省略されるのか）のスキーマ定義・動作仕様が記述されていない。

## 2. 詳細な指摘内容
1. **共通エラーレスポンスのスキーマ定義不足**:
   `01_overview.md` L42-55 では、`details` にフィールドエラー情報（`field`, `message`）を含む JSON 例が掲載されているのみで、`details` フィールドの型・必須性・空時の挙動が未定義である。

2. **実装のブレとクライアント側の例外発生リスク**:
   - バックエンド実装者によって、フィールドエラーがない場合に `details: []` を返す場合、`details: null` を返す場合、`details` を省略する場合のブレが発生する。
   - フロントエンド側で `error.details.length` や `error.details.map(...)` などを実行した際、`details` が `undefined` や `null` の場合に Null Pointer レベルの JavaScript ランタイム例外（`TypeError`）が発生するリスクがある。

## 3. 推奨される修正案
`01_overview.md` 1.3 節に共通エラーレスポンスのフィールド定義テーブルを追記し、`details` フィールドの挙動を明確化してください。

```markdown
#### エラーレスポンス スキーマ定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `error.code` | string | ○ | エラーコード（例: `BAD_REQUEST`, `UNAUTHORIZED`, `SAME_AS_CURRENT_PASSWORD` 等） |
| `error.message` | string | ○ | ユーザー向け汎用エラーメッセージ |
| `error.details` | array | × | フィールド単位のバリデーション詳細情報リスト。対象フィールドが存在しないエラー応答の場合は空配列 `[]` を返却（`null` やキー省略は不可） |
| `error.details[].field` | string | ○ | エラー対象のフィールド名 |
| `error.details[].message` | string | ○ | フィールド固有のエラーメッセージ |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.3 節に「エラーレスポンス スキーマ定義」テーブルを追加し、フィールドエラーが存在しない場合は `details: []`（空配列）を返却する（`null` やキー省略不可）ルールを規定しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
