# 新規登録OTP検証およびログイン API レスポンスの user オブジェクトにおける created_at / updated_at フィールドの欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`03_users.md`（および `01_overview.md`）で定義されている `user` リソースオブジェクトの標準スキーマには `id`, `username`, `email`, `created_at`, `updated_at` が含まれます。しかし、`02_auth.md` の 3.1.2 節 (`POST auth/register/verify-otp`) および 3.1.4 節 (`POST auth/login`) のレスポンス JSON 例では、`user` オブジェクトに `created_at` および `updated_at` が含まれておらず、`id`, `username`, `email` のみとなっています。

## 2. 詳細な指摘内容
1. `03_users.md` 3.2.1 (`GET users/{user_id}`) および 3.2.2 (`PUT users/{user_id}`) のレスポンスでは、`user` オブジェクトに登録日時 `created_at` と更新日時 `updated_at`（ISO 8601 JST 形式）が含まれる設計となっています。
2. 一方、`02_auth.md` の 3.1.2 節（L71-76）および 3.1.4 節（L163-169）のレスポンスでは：
   ```json
   {
     "user": {
       "id": "usr_987654321",
       "username": "exampleUser",
       "email": "user@example.com"
     }
   }
   ```
   となり、`created_at` と `updated_at` が省略されています。
3. 同一の `user` ドメインリソースモデルにおいて、認証 API 経由でログイン/登録した場合とユーザープロフィール取得 API 経由で取得した場合でプロパティ構造が異なると、フロントエンドでの状態管理クラスや TypeScript 型定義を共通化できず、不必要な分岐やオプショナル扱いの乱用が発生します。

## 3. 推奨される修正案
`02_auth.md` の 3.1.2 (`POST auth/register/verify-otp`) および 3.1.4 (`POST auth/login`) のレスポンス JSON 例に `created_at` および `updated_at` フィールドを追加し、`03_users.md` の `user` リソーススキーマと完全同期させてください。

```json
{
  "user": {
    "id": "usr_987654321",
    "username": "exampleUser",
    "email": "user@example.com",
    "created_at": "2026-08-17T12:00:00+09:00",
    "updated_at": "2026-08-17T12:00:00+09:00"
  }
}
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:30
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.2 (`POST auth/register/verify-otp`) および 3.1.4 (`POST auth/login`) のレスポンス JSON 例およびレスポンス定義テーブルに `created_at` と `updated_at` フィールドを追加し、`03_users.md` の `user` 標準スキーマと完全同期させました。

### 変更したファイル
- [02_auth.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
