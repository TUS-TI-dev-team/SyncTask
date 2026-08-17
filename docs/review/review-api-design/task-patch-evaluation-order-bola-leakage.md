# PATCH tasks/{task_id} におけるリクエスト評価順序の不備による IDOR/BOLA 存在検証サイドチャネルの脆弱性

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`04_tasks.md` の `PATCH tasks/{task_id}` (L263-L270) において、リクエスト評価順序のステップ2に「リクエスト構文・入力バリデーション (`400 Bad Request`)」、ステップ3に「認可チェック・IDOR/BOLA検証 (`404 Not Found`)」が定義されています。バリデーション（400）を認可・存在検証（404）より先に実行すると、攻撃者が他ユーザーの `task_id` に対して不正なペイロード（例: `{"title": ""}`）を送信することで、リソースの存在有無（400か404か）を判別できるサイドチャネル攻撃が可能となります。

## 2. 詳細な指摘内容
1. **認可・存在検証（404）前の入力バリデーション（400）実行による情報漏洩**:
   - `requirements.md` L189 では、「他ユーザー所有のタスクIDへのアクセス・操作要求に対しては、タスクの存在有無を秘匿するため一律 `404 Not Found` を返却し処理を拒否する」と規定されています。
   - しかし `04_tasks.md` L263-L270 の評価順序では、ステップ2で入力バリデーション（400）が評価され、ステップ3で認可チェック（404）が評価されます。
   - これにより、攻撃者が存在する他ユーザーの `task_id` に対しバリデーションエラーとなるリクエスト（例: 空タイトル `{"title": ""}` や不正な列挙値 `{"priority": "invalid"}`）を送信した場合、404 ではなく `400 Bad Request` が返却されます。一方、存在しない `task_id` に対して同様のリクエストを送信した場合は、ステップ3まで到達しないか、あるいは存在確認が後回しになることでレスポンス結果が異なり、他ユーザーの `task_id` が存在するか否かを特定されてしまいます。

2. **他のパスパラメータ型エンドポイント（`GET /DELETE tasks/{task_id}`）との不整合**:
   - `GET tasks/{task_id}` (L208-L211) および `DELETE tasks/{task_id}` (L346-L349) では、認証・CSRF検証に続き、まず認可チェック・IDOR/BOLA検証（`404 Not Found`）を行う順序となっており、`PATCH tasks/{task_id}` の評価順序のみが不整合を生んでいます。

## 3. 推奨される修正案
`04_tasks.md` の `PATCH tasks/{task_id}` における「リクエスト評価順序」を以下のように修正し、認可チェック・IDOR/BOLA検証（404）をリクエストボディのバリデーション（400）より先に実行するよう定義変更してください。

```markdown
##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `task_id` の存在およびセッションユーザーの所有タスクかを検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
3. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）への `null` 指定有無、文字数・列挙値制約を検証。不備がある場合は 400 `BAD_REQUEST` を返却。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.4 節 (`PATCH tasks/{task_id}`) の「リクエスト評価順序」において、認可チェック・IDOR/BOLA検証 (`404 Not Found`) をステップ2とし、リクエスト構文・入力バリデーション (`400 Bad Request`) をステップ3に配置変更しました。これにより他ユーザー所有リソースへのバリデーションエラー応答によるリソース存在漏洩（BOLA/IDORサイドチャネル攻撃）を防止し、他のパスパラメータエンドポイント（GET / DELETE）と評価順序の整合性を統一しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
