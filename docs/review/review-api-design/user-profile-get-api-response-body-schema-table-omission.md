# ユーザープロフィール取得 API (GET users/{user_id}) におけるレスポンスボディ定義テーブルの欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:05:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`03_users.md` の 3.2.1 節 (`GET users/{user_id}`) において、JSON レスポンス例は記載されているものの、他 API エンドポイントのようにレスポンスフィールドの型・必須有無・詳細説明を定めた「Response Body フィールド定義」テーブルが存在しません。

## 2. 詳細な指摘内容
- **`03_users.md` L16-L27 (3.2.1 Response 200 OK)**:
  JSON レスポンス例 (`{ "user": { "id": ..., "username": ..., "email": ..., "created_at": ..., "updated_at": ... } }`) のみが掲載されており、各フィールドのデータ型、Null許容性、日時フォーマット仕様のテーブル記述がありません。

### 問題点
1. ドキュメント記述形式の非統一性:
   `PUT users/{user_id}` や `GET tasks` などの他 API ドキュメントではリクエスト/レスポンスパラメータの型や制約が明示的なテーブル形式で定義されていますが、`GET users/{user_id}` では JSON 例のみとなっているため仕様書の完全性と統一性が欠けています。
2. API 実装・クライアント型生成時の誤認懸念:
   `created_at` / `updated_at` の日時フォーマット（ISO 8601 拡張形式 JST）や各フィールドの非 Null 保証について明文記述がないため、フロントエンド開発や自動コード生成ツール適用時に型解釈の曖昧さが生じる恐れがあります。

## 3. 推奨される修正案
`03_users.md` の 3.2.1 節 (`GET users/{user_id}`) の `##### Response (200 OK)` 配下に以下のレスポンスフィールド定義テーブルを追加してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.1 節 (`GET users/{user_id}`) の `##### Response (200 OK)` 配下に Response Body フィールド定義テーブル（型・必須属性・ISO 8601拡張日時フォーマット説明等）を追加しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)

### 修正案の例 (`03_users.md` 3.2.1)
```markdown
##### Response (200 OK)

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `user` | object | ○ | ユーザー情報オブジェクト |
| `user.id` | string | ○ | ユーザーID（例: `usr_987654321`） |
| `user.username` | string | ○ | ユーザー名 |
| `user.email` | string | ○ | 登録メールアドレス |
| `user.created_at` | string | ○ | アカウント作成日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |
| `user.updated_at` | string | ○ | アカウント最終更新日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |
```
