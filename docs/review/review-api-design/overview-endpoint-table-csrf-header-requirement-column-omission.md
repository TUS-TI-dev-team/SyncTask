# `01_overview.md` 2. エンドポイント一覧テーブルにおける CSRF ヘッダー（`X-CSRF-Token`）要否カラムの欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 2. エンドポイント一覧テーブルにおいて、各エンドポイントの「認証要否」のみが記載されており、認証必須の状態変更エンドポイント（`POST`, `PUT`, `PATCH`, `DELETE`）で必須となる CSRF ヘッダー（`X-CSRF-Token`）の要否カラムおよび明記が欠落している。

## 2. 詳細な指摘内容
1. **CSRF ヘッダー要否情報の欠落（L96-118）**:
   1.2 節において「認証を必要とするすべての状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）において CSRFトークン（`X-CSRF-Token` ヘッダー）の検証を必須とする」と規定されているが、2. エンドポイント一覧テーブルには「認証要否」カラムしか存在しない。
2. **未認証 POST と認証必須 POST の誤認リスク**:
   `POST auth/login` や `POST auth/register/request-otp` などの未認証 POST（CSRF ヘッダー不要）と、`POST auth/logout` や `POST auth/change-email/request-otp` などの認証必須 POST（CSRF ヘッダー必須）の相違がエンドポイント一覧表から直接判別できず、一覧性・可読性を損ねている。

## 3. 推奨される修正案
`01_overview.md` 2. エンドポイント一覧テーブルに「CSRFヘッダー」カラムを追加するか、各エントリに CSRF トークン検証の要否を明記してください。

```markdown
| カテゴリ | メソッド | URI | 役割・機能 | 認証要否 | CSRFヘッダー |
| :--- | :--- | :--- | :--- | :---: | :---: |
| **認証 (Auth)** | `POST` | `auth/register/request-otp` | 新規登録情報のバリデーション・OTP発行・メール送信 | 不要 | 不要 |
| | `POST` | `auth/login` | メールアドレス・パスワードによるログイン認証 | 不要 | 不要 |
| | `POST` | `auth/logout` | ログインセッションの破棄・ログアウト | 必須 | 必須 |
| | `POST` | `auth/change-email/request-otp` | メールアドレス変更用OTP作成・送信 | 必須 | 必須 |
...
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:18:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 2. エンドポイント一覧テーブルに `CSRFヘッダー` カラムを追加し、各エンドポイントの CSRF ヘッダー要否（状態変更かつ認証必須エンドポイントは「必須」、それ以外は「不要」）を明確化しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)

