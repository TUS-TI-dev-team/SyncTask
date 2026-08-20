# プロフィール更新APIにおける評価順序の構文チェック記述漏れおよび読み取り専用フィールド扱いの未規定

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`03_users.md` の 3.2.2 節 (`PUT users/{user_id}`) において、「リクエスト評価順序」のステップ2（`400 Bad Request`）の記述が `username` カラムのトリム後文字数・文字種検証のみとなっており、JSON パラメータ構文チェック、必須フィールド `username` の存在・Null非許容チェック、およびリクエストボディに読み取り専用/システムフィールド（`id`, `email`, `created_at`, `updated_at` 等）が含まれた場合の扱いについての規定が欠落しています。

## 2. 詳細な指摘内容
1. **リクエスト構文および必須/Null検証の明記不足**:
   - `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) のステップ2では、「リクエストボディの JSON 形式、必須フィールドの有無を検証（不備時は 400 を返却）」と明記されています。
   - しかし 3.2.2 節のステップ2では `username` の文字列トリム後のバリデーションのみが記述されており、JSON形式不正、`username` フィールドの欠落、あるいは `username: null` や非文字列（数値・オブジェクト等）が送信された場合の挙動が明示されていません。

2. **読み取り専用/変更不可フィールド（`email` 等）送信時の仕様未規定**:
   - `PUT users/{user_id}` はユーザープロフィールを更新する API ですが、メールアドレスの変更は専用の OTP 認証フロー（`POST auth/change-email/request-otp`）を通じて行う仕様となっています。
   - `04_tasks.md` 3.3.4 節 (`PATCH tasks/{task_id}`) では、「システムの読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）が含まれている場合は、エラーとせず更新対象外として単に無視します」と明確に規定されています。
   - 3.2.2 節においても、クライアントがリクエストボディに `id`, `email`, `created_at`, `updated_at` などの変更不可/読み取り専用フィールドを含めて送信した場合に、エラーとせず単に無視して `username` のみを更新する仕様であるかを明文化する必要があります。

## 3. 推奨される修正案
`03_users.md` 3.2.2 節 (`PUT users/{user_id}`) の「リクエスト評価順序」ステップ2、および `##### Request Body` の補足注記を以下のように改修してください。

1. **ステップ2（リクエスト構文・入力バリデーション `400 Bad Request`）**:
   - リクエストボディの JSON 形式、必須フィールド `username` の有無および型（`null` 指定不可）、トリム後の文字数（2〜20文字）・使用可能文字（英数字）を検証（不備時は 400 `BAD_REQUEST` を即座に返却）。なお、読み取り専用・変更不可フィールド（`id`, `email`, `created_at`, `updated_at` 等）がボディに含まれている場合は、エラーとせず更新対象外として単に無視します。

2. **`##### Request Body` のパラメータテーブル `username` の制約説明**:
   - `必須。null 不可。前後の空白を自動トリム後、2〜20文字かつ英数字（大文字小文字可）であること。要件違反・未入力時は 400 エラー。トリム後のユーザー名が現在のユーザー名と同一の場合は 422 エラー。`

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.2 節 (`PUT users/{user_id}`) の「リクエスト評価順序」ステップ2および Request Body パラメータテーブル、Errors 一覧において、JSON 形式チェック、必須フィールド `username` の有無および null 不可指定の検証、文字数・文字種検証を明記しました。また、`id`, `email`, `created_at`, `updated_at` 等の読み取り専用・変更不可フィールドが送信された場合はエラーとせず単に無視して更新対象外とする仕様を追加しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
