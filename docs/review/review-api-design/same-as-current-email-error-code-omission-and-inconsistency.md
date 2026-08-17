# メールアドレス変更時の同一アドレスエラーコード SAME_AS_CURRENT_EMAIL の未定義と不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`01_overview.md` 1.3 の「代表的なエラーコード一覧」において、`SAME_AS_CURRENT_PASSWORD`（パスワード変更時の同一パスワード指定）および `SAME_AS_CURRENT_USERNAME`（ユーザー名変更時の同一ユーザー名指定）が `422` ステータスの固有エラーコードとして定義されている。しかし、メールアドレス変更（`POST auth/change-email/request-otp`）において現在のメールアドレスと同一のアドレスを指定した場合の固有エラーコード `SAME_AS_CURRENT_EMAIL` が `01_overview.md` のエラーコード一覧に記載されていない。また、`02_auth.md` の 3.1.10 節のエラー定義でも汎用コード `UNPROCESSABLE_ENTITY` が用いられており、同一値指定エラーに対するエラーコード体系の整合性が崩れている。

## 2. 詳細な指摘内容
1. **`01_overview.md` 1.3 代表的なエラーコード一覧での欠落**:
   `01_overview.md` L72-73 にて `SAME_AS_CURRENT_PASSWORD` と `SAME_AS_CURRENT_USERNAME` が明記されているが、`SAME_AS_CURRENT_EMAIL` が一覧に含まれていない。

2. **`02_auth.md` 3.1.10 節との不整合**:
   `02_auth.md` 3.1.10 節（L326）において、`422 Unprocessable Entity: 現在のメールアドレスと同一（code: "UNPROCESSABLE_ENTITY"）` と汎用エラーコード `UNPROCESSABLE_ENTITY` が指定されている。
   ユーザー名・パスワード変更時には固有コード（`SAME_AS_CURRENT_*`）が割り振られているにも関わらず、メールアドレス変更時のみ汎用コードになると、フロントエンド側で特定のエラー状態（現在のメールアドレスと同じ入力であること）を他のビジネスルール違反と区別してフォームやメッセージに反映できなくなる。

## 3. 推奨される修正案
1. `01_overview.md` 1.3 の代表的エラーコード一覧に以下を追加してください:
   ```markdown
   | 422 | `SAME_AS_CURRENT_EMAIL` | メールアドレス変更時の同一メールアドレス指定 |
   ```

2. `02_auth.md` 3.1.10 節のエラー定義を以下のように修正してください:
   ```markdown
   - `422 Unprocessable Entity`: 現在のメールアドレスと同一（code: `"SAME_AS_CURRENT_EMAIL"`）
   ```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.3 の一覧テーブルに `SAME_AS_CURRENT_EMAIL` を追加し、`02_auth.md` 3.1.10 のエラー定義コードを `SAME_AS_CURRENT_EMAIL` に更新しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
- [02_auth.md](docs/design/api_design/02_auth.md)
