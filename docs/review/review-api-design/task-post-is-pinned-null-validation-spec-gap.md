# POST tasks における is_pinned パラメータへの null および非 boolean 値指定時のバリデーションエラー仕様欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`PATCH tasks/{task_id}` (3.3.4 節) では非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）に明示的に `null` が指定された場合、400 Bad Request を返却する仕様が明記されていますが、`POST tasks` (3.3.2 節) の `is_pinned` フィールドにおいて `null` や非 boolean 値（数値・文字列等）が送信された場合のエラー挙動が明記されていません。

## 2. 詳細な指摘内容
1. **`POST tasks` における `is_pinned` の型制約記述**:
   - `04_tasks.md` L116 の `is_pinned` パラメータ定義は `boolean`, `任意`（省略時デフォルト `false`）と記載されています。
2. **`null` および非 boolean 値の取り扱い**:
   - クライアントが `{"title": "タスク", "is_pinned": null}` や `{"title": "タスク", "is_pinned": "true"}` のように `null` または非 boolean 型の値を送信した場合に、サーバーが 400 Bad Request を返却するのか、あるいはデフォルト値 `false` にフォールバックするのかが明確ではありません。
3. **`PATCH` との記述不均衡**:
   - `PATCH tasks/{task_id}` (L224) では `is_pinned: null` が送信された場合に 400 Bad Request を返す旨が明記されており、`POST` と `PATCH` 間で入力バリデーションポリシーの記述に不均衡が存在します。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 節 (`POST tasks`) のフィールド定義およびリクエスト評価順序 2（入力バリデーション）に、以下の仕様を明記してください：

- `is_pinned` に明示的に `null` または非 boolean 型（数値・文字列等）が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却する。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:35
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.2 節 (`POST tasks`) のパラメータ定義、リクエスト評価順序 2、および Errors セクションに、`is_pinned` へ明示的に `null` または非 boolean 型（数値・文字列等）が指定された場合は 400 Bad Request (`code: "BAD_REQUEST"`) を返却する旨を明記し、`PATCH tasks/{task_id}` とバリデーションポリシーを統一しました。

### 変更したファイル
- [04_tasks.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
