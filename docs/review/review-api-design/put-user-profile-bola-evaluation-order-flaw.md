# `PUT users/{user_id}` API における BOLA/IDOR 検証の評価順序の欠陥

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-17 18:00:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`03_users.md` の `PUT users/{user_id}` (3.2.2) における「リクエスト評価順序」にて、認可チェック・IDOR/BOLA検証（`404 Not Found`）がリクエスト構文・入力バリデーション（`400 Bad Request`）よりも後ろのステップ（ステップ3）に配置されている。

## 2. 詳細な指摘内容
1. **認可検証の遅延に伴うレスポンスの不整合と情報漏洩リスク**:
   - 現在の `3.2.2` の評価順序は以下の通り定義されている：
     1. **認証・CSRF検証 (`401` / `403`)**
     2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**
     3. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**
     4. **ビジネスルール検証 (`422 Unprocessable Entity`)**
   - この評価順序に従うと、攻撃者が他ユーザーの `user_id` （`path user_id != session user_id`）を指定して `PUT users/{other_user_id}` を呼び出した際、送信したリクエストボディが形式不正（例: `username` が1文字、または必須フィールド欠落）である場合、ステップ2の入力バリデーションにより `400 Bad Request` が返却される。
   - 一方、形式上正当なリクエストボディを送信した場合は、ステップ3に到達して `404 Not Found` が返却される。
   - これにより、認可権限のないリソースへのアクセス試行において、リクエストボディの妥当性によって応答ステータスコード（400 vs 404）が変化し、内部の入力検証挙動や認可チェック前の処理状態が外部に推測可能となる。
2. **`04_tasks.md` やセキュリティ基本方針との不整合**:
   - `01_overview.md` 1.2 節の IDOR/BOLA 方針、および `04_tasks.md` 3.3.4 （`PATCH tasks/{task_id}`）等では、認証・CSRF検証直後のステップ2で IDOR/BOLA 検証（`404 Not Found`）を実行し、他者所有リソースへのアクセスを最優先で遮断した後にリクエストボディのバリデーション（ステップ3）を行う構造になっている。
   - `PUT users/{user_id}` においても同様に、IDOR/BOLA 検証をリクエストボディのバリデーションより前に実行し、一貫性を保つ必要がある。

## 3. 推奨される修正案
`03_users.md` の 3.2.2 (`PUT users/{user_id}`) における「リクエスト評価順序」を更新し、認証・CSRF検証（ステップ1）の直後に認可チェック・IDOR/BOLA検証（`404 Not Found`）を行うよう修正してください。

```markdown
##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
3. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須フィールド `username` の有無および型（`null` 指定不可）、トリム後の文字数（2〜20文字）・使用可能文字（英数字）を検証（不備時は 400 `BAD_REQUEST` を即座に返却）。なお、読み取り専用・変更不可フィールド（`id`, `email`, `created_at`, `updated_at` 等）がボディに含まれている場合は、エラーとせず更新対象外として単に無視します。
4. **ビジネスルール検証 (`422 Unprocessable Entity`)**:
   トリム後の `username` が現在のユーザー名と同一か検証（同一の場合は 422 `SAME_AS_CURRENT_USERNAME`）。
```
