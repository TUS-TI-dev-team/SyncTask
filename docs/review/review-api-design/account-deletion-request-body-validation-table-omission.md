# アカウント削除API (`DELETE users/{user_id}`) におけるリクエストボディのパラメータ仕様テーブルの欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`03_users.md` の 3.2.3 節 (`DELETE users/{user_id}`) において、リクエストボディ例 (`{"password": "Password123!"}`) は記述されているが、パラメータの型・必須チェック・入力バリデーション条件を定義したテーブルが存在しない。

`PUT users/{user_id}` (3.2.2) にはパラメータ定義テーブルが存在するのに対し、`DELETE` および `PATCH password` では記述形式が統一されておらず仕様定義が漏れている。

## 2. 詳細な指摘内容
- `03_users.md` 3.2.3 節 L73-L78:
  ```json
  ##### Request Body
  {
    "password": "Password123!"
  }
  ```
  この直後にパラメータの型・制約を示すテーブルが存在しない。
  これにより、以下のルールが明示されていない:
  1. `password`: 必須項目（○）、文字列（string）。
  2. 未入力（空文字やプロパティ欠落）送信時のバリデーションエラー（`400 Bad Request` / `code: "BAD_REQUEST"`）に関する制約記述。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 節 (`DELETE users/{user_id}`) に以下の Request Body パラメータテーブルを追加してください。

```markdown
##### Request Body
```json
{
  "password": "Password123!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `password` | string | ○ | 再認証用パスワード。未入力時は 400 エラー、パスワード不一致時は 401 エラー |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`DELETE users/{user_id}` 節に `password` フィールドの型・必須チェック（○）・制約条件（再認証用パスワード、未入力時 400 エラー、不一致時 401 エラー）を定義した Request Body パラメータテーブルを追加しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
