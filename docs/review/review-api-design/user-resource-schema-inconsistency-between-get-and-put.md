# ユーザー情報取得 (`GET`) および更新 (`PUT`) レスポンス間における `user` オブジェクトのスキーマ不整合と日時フィールド欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`03_users.md` において、`GET users/{user_id}` (3.2.1) のレスポンスに `created_at` が含まれるのに対し、`PUT users/{user_id}` (3.2.2) のレスポンス `user` オブジェクトでは `created_at` が省略されており、同一のリソース表現（`user` オブジェクト）でスキーマの不整合が生じている。

また、データベース設計 (`database_design.md` L27) で定義されている `LOGIN_ACCOUNT` テーブルの `UPDATED_AT`（`updated_at`）が `GET` および `PUT` のいずれのレスポンススキーマからも漏れている。

## 2. 詳細な指摘内容
- `GET users/{user_id}` (3.2.1 L11-L20) レスポンス:
  ```json
  {
    "user": {
      "id": "usr_987654321",
      "username": "exampleUser",
      "email": "user@example.com",
      "created_at": "2026-08-01T10:00:00+09:00"
    }
  }
  ```
- `PUT users/{user_id}` (3.2.2 L47-L56) レスポンス:
  ```json
  {
    "user": {
      "id": "usr_987654321",
      "username": "newUsername",
      "email": "user@example.com"
    }
  }
  ```

`PUT` レスポンスで `created_at` が削除されているため、フロントエンドでユーザー更新処理（`PUT`）の戻り値を用いて状態（State/Store）をそのまま上書き更新した際に `created_at` が消失するリスクがある。
さらに、プロフィール情報の最終更新日時を示す `updated_at` が返却されないため、クライアント側で更新日時の表示や状態同期の確認を行えない。

## 3. 推奨される修正案
`GET users/{user_id}` および `PUT users/{user_id}` のレスポンスに含まれる `user` オブジェクト構造を統一し、`created_at` および `updated_at` を常に返却する仕様に変更してください。

```json
{
  "user": {
    "id": "usr_987654321",
    "username": "newUsername",
    "email": "user@example.com",
    "created_at": "2026-08-01T10:00:00+09:00",
    "updated_at": "2026-08-17T16:35:00+09:00"
  }
}
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`GET users/{user_id}` (3.2.1) および `PUT users/{user_id}` (3.2.2) のレスポンス `user` オブジェクトの構造を統一し、`created_at` と `updated_at` の両方を常に返却する仕様に更新しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
