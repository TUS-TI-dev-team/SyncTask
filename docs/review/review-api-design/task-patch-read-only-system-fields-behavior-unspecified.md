# PATCH tasks/{task_id} におけるシステム読み取り専用フィールド指定時の挙動未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`PATCH tasks/{task_id}` (3.3.4 節) の部分更新処理において、クライアントがリクエストボディにシステムの読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）を含めて送信した場合のバックエンドでの扱い（無視するのか拒否するのか）が明記されていません。

## 2. 詳細な指摘内容
1. **読み取り専用フィールドの処理方針の曖昧さ**:
   - `04_tasks.md` L208 のリクエスト評価順序では「非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）への `null` 指定有無、文字数・列挙値制約を検証」と定義されていますが、`id`, `user_id`, `created_at`, `updated_at` といった変更不可能なシステムプロパティがペーロードに含まれていた場合の挙動が書かれていません。
   - クライアントが既存オブジェクト構造全体をそのままパッチリクエストとして送信した場合に、サーバーがこれらのフィールドを単に「無視して更新対象外とする」のか、あるいは「`400 Bad Request`（`code: "BAD_REQUEST"`）を返却する」のかが不透明です。

## 3. 推奨される修正案
`PATCH tasks/{task_id}` の注記またはリクエスト評価順序に、「リクエストボディに読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）が含まれている場合は、更新対象外として無視し、許可されたフィールドのみを正常に更新処理する（または拒否する場合は 400 Bad Request とする）」旨の挙動方針を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:18:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.4 節 (`PATCH tasks/{task_id}`) のリクエスト評価順序ステップ 2 および Request Body 注記に、リクエストボディにシステムの読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）が含まれている場合はエラーとせず更新対象外として単に無視し、許可されたフィールドのみを正常に更新処理する方針を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)

