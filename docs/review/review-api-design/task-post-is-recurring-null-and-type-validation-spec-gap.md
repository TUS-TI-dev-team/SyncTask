# POST tasks における is_recurring パラメータへの null および非 boolean 型指定時のバリデーション仕様欠落

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` (3.3.2 節) のリクエストボディにおいて、`is_recurring` パラメータに対し `null` や非 boolean 型（文字列 `"true"`, 数値 `1` 等）が指定された場合のバリデーション挙動およびエラー詳細（`error.details`）の定義が欠落しています。

## 2. 詳細な指摘内容
1. **他の boolean / Enum フィールドとの表記の不整合**:
   - `04_tasks.md` L138 の `is_pinned` や L136 の `priority` では、「明示的に `null` または非 boolean 型（数値・文字列等）が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却」とバリデーションエラー条件が厳格に明記されています。
   - しかし、L139 の `is_recurring` では「繰り返し一括作成フラグ（デフォルト: `false`）」と初期値のみが記載されており、`"is_recurring": null` や `"is_recurring": "true"`, `"is_recurring": 1` といった不正な型・`null` が明示送信された場合の挙動が未定義です。
2. **実装における曖昧性**:
   - 不正な型が指定された際、バックエンド実装で `BAD_REQUEST` エラーとして即時拒否すべきか、真偽値への型変換（coercion）を行うべきかが曖昧になっています。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 節 (`POST tasks`) の Request Body フィールド定義テーブルにおける `is_recurring` の項目を以下のように更新し、非 boolean 型および `null` 指定時の拒否挙動を明確化してください：

```markdown
| `is_recurring` | boolean | × | 繰り返し一括作成フラグ（`true` / `false`）。省略時はデフォルト `false` として作成。明示的に `null` または非 boolean 型（数値・文字列等）が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却 |
```

また、Errors セクション (400 Bad Request) の箇条書きにも `is_recurring` への `null` / 非 boolean 指定を追記してください。
