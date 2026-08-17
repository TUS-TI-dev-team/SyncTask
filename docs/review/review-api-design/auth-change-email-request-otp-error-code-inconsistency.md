# `POST auth/change-email/request-otp` における同一メールアドレスエラーコードの命名不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`02_auth.md` の 3.1.10 `POST auth/change-email/request-otp` において、変更先メールアドレスに現在と同一のメールアドレスを指定した場合のエラーコードが `code: "UNPROCESSABLE_ENTITY"`（L326）と定義されている。
これは、ユーザー名変更 (`PUT users/{user_id}`) の `SAME_AS_CURRENT_USERNAME` やパスワード変更 (`PATCH users/{user_id}/password`) の `SAME_AS_CURRENT_PASSWORD` と比較して命名規則の不整合がある。

## 2. 詳細な指摘内容
- `03_users.md` L66: `422 Unprocessable Entity: 現在のユーザー名と同一 (code: "SAME_AS_CURRENT_USERNAME")`
- `03_users.md` L151: `422 Unprocessable Entity: 新パスワードが現在のパスワードと同一 (code: "SAME_AS_CURRENT_PASSWORD")`
- `01_overview.md` L72-L73: 体系的に `SAME_AS_CURRENT_PASSWORD`, `SAME_AS_CURRENT_USERNAME` など具体的なビジネスルールエラーコードが挙げられている。
- `02_auth.md` L326: `- 422 Unprocessable Entity: 現在のメールアドレスと同一 (code: "UNPROCESSABLE_ENTITY")`

現在のメールアドレスと同一の入力を拒否するエラーについて、汎用的な HTTP ステータス名 `UNPROCESSABLE_ENTITY` をそのままエラーコードとして使用すると、フロントエンド側で他ケースの 422 エラー（例: 5回連続失敗自動再送など）と明確に区別して「現在のメールアドレスと同じです」という特定のフィールドエラーを表示する際の判定が難しくなる。

## 3. 推奨される修正案
`02_auth.md` 3.1.10 L326 のエラーコードを以下のように修正してください：

```markdown
- `422 Unprocessable Entity`: 現在のメールアドレスと同一（code: `"SAME_AS_CURRENT_EMAIL"`）
```

併せて、`01_overview.md` 1.3 の代表的なエラーコード一覧テーブルに `SAME_AS_CURRENT_EMAIL` を追記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.10 における同一メールアドレス指定エラーのエラーコードを `SAME_AS_CURRENT_EMAIL` に修正し、`01_overview.md` 1.3 代表的エラーコード一覧にも `SAME_AS_CURRENT_EMAIL` を追記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
- [02_auth.md](docs/design/api_design/02_auth.md)
