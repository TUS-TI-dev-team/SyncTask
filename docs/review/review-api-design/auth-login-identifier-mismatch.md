# ログインAPIおよびパスワードリセットAPIの識別子仕様の不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書（`requirements.md`）ではログイン識別子として「メールアドレス」を使用しユーザー名は同名登録可と定義されていますが、`api_design.md` ではログインAPI（`auth/login`）の入力が「ユーザー名、パスワード」となっており、仕様間で重大な不整合が発生しています。

## 2. 詳細な指摘内容
1. **ログインAPIの入力不整合**:
   - `docs/req-def/requirements.md` L62:
     > - メールアドレス（ユーザー名ではない）とパスワードでログイン
   - `docs/req-def/requirements.md` L249:
     > - 同名のユーザー名登録は可（ログイン識別子にはメールアドレスを使用するため）
   - 一方で、`docs/design/api_design.md` L16:
     > `auth/login` | `POST` | ... | 入力: ユーザー名、パスワード
   - 同名のユーザー名登録が許可されているシステムにおいてユーザー名でログインさせようとすると、ユーザーを一意に特定できず認証不能または意図しないアカウントへのログインが発生します。

2. **パスワードリセットAPIの入力不整合**:
   - `docs/design/api_design.md` L18 の `auth/password-reset/request-otp` では入力が「ユーザー名 または メールアドレス」となっていますが、`screen_design.md` L10 および `requirements.md` L69-70 ではメールアドレスのみによる特定・認証となっています。

## 3. 推奨される修正案
1. `auth/login` の入力仕様を「メールアドレス（`email`）、パスワード（`password`）」に修正してください。
2. `auth/password-reset/request-otp` の入力仕様を「メールアドレス（`email`）」に統一してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` の `POST auth/login` の入力フィールドを `email` と `password` に修正し、要件定義書（L62）の「メールアドレスでログイン」仕様と完全に整合させました。
- `POST auth/password-reset/request-otp` の入力フィールドも `email` に統一しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
