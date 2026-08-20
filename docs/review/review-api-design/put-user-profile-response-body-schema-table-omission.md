# プロフィール更新 API (PUT users/{user_id}) におけるレスポンスボディ定義テーブルの欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`03_users.md` の 3.2.2 節 (`PUT users/{user_id}`) において、JSON レスポンス例は記述されているものの、3.2.1 節 (`GET users/{user_id}`) のようにレスポンスフィールドの型・必須有無・詳細説明を定めた「Response Body フィールド定義」テーブルが存在しません。

## 2. 詳細な指摘内容
- **`03_users.md` L74-L85 (3.2.2 Response 200 OK)**:
  JSON レスポンス例 (`{ "user": { "id": ..., "username": ..., "email": ..., "created_at": ..., "updated_at": ... } }`) のみが掲載されており、各フィールドのデータ型、Null許容性、日時フォーマット仕様のテーブル記述がありません。

### 問題点
1. **ドキュメント記述形式の非統一性**:
   `3.2.1 GET users/{user_id}` では Response Body フィールド定義テーブル（型・必須属性・説明）が追加されましたが、同様に `user` オブジェクトを返却する `3.2.2 PUT users/{user_id}` では JSON 例のみの掲載となっており、仕様書内のフォーマット統一性が損なわれています。
2. **API 実装・クライアント型生成時の不確実性**:
   更新後のレスポンスに含まれる各フィールド（`id`, `username`, `email`, `created_at`, `updated_at`）の型情報や非 Null 保証が明示されていないため、クライアント側の型定義作成や実装時に仕様の揺れが生じる原因となります。

## 3. 推奨される修正案
`03_users.md` の 3.2.2 節 (`PUT users/{user_id}`) の `##### Response (200 OK)` 配下に、3.2.1 節と同様の Response Body フィールド定義テーブルを追加してください。

**修正案の例 (`03_users.md` 3.2.2)**:
```markdown
##### Response (200 OK)

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `user` | object | ○ | ユーザー情報オブジェクト |
| `user.id` | string | ○ | ユーザーID（例: `usr_987654321`） |
| `user.username` | string | ○ | 更新後のユーザー名 |
| `user.email` | string | ○ | 登録メールアドレス |
| `user.created_at` | string | ○ | アカウント作成日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |
| `user.updated_at` | string | ○ | アカウント最終更新日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:18:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.2 節 (`PUT users/{user_id}`) の `Response (200 OK)` セクションに、Response Body フィールド定義テーブル（`user`, `user.id`, `user.username`, `user.email`, `user.created_at`, `user.updated_at`）を追加しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)

