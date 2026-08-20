# 概要書のエラーコード一覧におけるINVALID_PASSWORD_CONTENTの記載漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:33:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の「1.3 共通エラーレスポンス構造」における「代表的なエラーコード一覧」テーブルにおいて、認証・ユーザー管理機能の詳細設計書（`02_auth.md`）で定義されている `INVALID_PASSWORD_CONTENT`（HTTPステータス `422 Unprocessable Entity`：パスワード変更・リセット時にユーザー名やメールアドレスのローカル部が含まれる違反エラー）が記載漏れとなっている。

## 2. 詳細な指摘内容
`02_auth.md`（3.1.9 `POST auth/password-reset/reset`）および `03_users.md`（3.2.4 `PATCH users/{user_id}/password`）では、新パスワードにユーザー名やメールアドレスのローカル部（4文字以上）が含まれる場合のエラーコードとして `INVALID_PASSWORD_CONTENT`（HTTP `422`）を返却することが明記されています。

しかし、全体共通仕様書である `01_overview.md` の「代表的なエラーコード一覧」テーブル（74〜90行目）には、`SAME_AS_CURRENT_PASSWORD`, `SAME_AS_CURRENT_USERNAME`, `SAME_AS_CURRENT_EMAIL` などの `422` エラーコードは列挙されているものの、`INVALID_PASSWORD_CONTENT` が含まれていません。

これにより、API全体のエラー分類体系の共通理解が妨げられ、フロントエンド開発者が一覧表のみを参照した際に本エラーコードを見落とすリスクが存在します。

## 3. 推奨される修正案
`01_overview.md` の「1.3 共通エラーレスポンス構造」の「代表的なエラーコード一覧」テーブルに、以下の行を追加してください。

```markdown
| 422 | `INVALID_PASSWORD_CONTENT` | パスワード変更/リセット時のユーザー名・メールアドレスローカル部含有違反 |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:38:45
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/01_overview.md` の「1.3 共通エラーレスポンス構造」内の「代表的なエラーコード一覧」テーブルに、HTTPステータス `422` のエラーコード `INVALID_PASSWORD_CONTENT`（パスワード変更/リセット時のユーザー名・メールアドレスローカル部含有違反）を追加・追記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)

