# `01_overview.md` におけるメールアドレス変更時の同一メールアドレス指定に対するダミーセッション返却記述の矛盾

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`01_overview.md` 1.2「アカウント列挙防止 (User Enumeration 対策)」の文章（L30）において、メールアドレス変更（`auth/change-email/request-otp`）の際に「現在と同一メールアドレスの指定」であっても一貫してダミーOTPセッションを発行して `200 OK` を返却すると記述されている。
しかし、`02_auth.md` 3.1.10 および `requirements.md` L45 では、ログイン中のユーザーが現在と同一のメールアドレスを変更先として指定した場合は `422 Unprocessable Entity` を返却すると規定されており、ドキュメント間で記述が矛盾している。

## 2. 詳細な指摘内容
- `01_overview.md` L30:
  `新規登録（auth/register/request-otp）、パスワードリセット（auth/password-reset/request-otp）、メールアドレス変更（auth/change-email/request-otp）において、指定されたメールアドレスの登録有無、他ユーザーとの重複、現在と同一メールアドレスの指定、または他ユーザーの有効なOTPセッション期間中（手続き中）の指定にかかわらず、一貫してダミーOTPセッションを発行して 200 OK（応答遅延 1.0s ± 0.1s）を返却します。`
- `02_auth.md` L310, L326:
  `現在のメールアドレスと同一の場合は 422 エラー`
  `- 422 Unprocessable Entity: 現在のメールアドレスと同一`
- `requirements.md` L45:
  `変更先メールアドレスが現在設定中の自分自身のメールアドレスと同一の場合は「現在のメールアドレスと同じです」とエラー表示し拒否する`

ログイン済みのユーザー自身の現在のメールアドレス指定については、既に画面等で本人に開示されている情報であるため、アカウント列挙防止（User Enumeration 対策）の対象外であり、`422` エラー（`SAME_AS_CURRENT_EMAIL`）を返却するのが正しい設計である。`01_overview.md` の記述は誤解を招く内容となっている。

## 3. 推奨される修正案
`01_overview.md` L30 から「現在と同一メールアドレスの指定、」というフレーズを削除し、以下のように修正してください：

`新規登録（auth/register/request-otp）、パスワードリセット（auth/password-reset/request-otp）、メールアドレス変更（auth/change-email/request-otp）において、指定されたメールアドレスの登録有無、他ユーザーとの重複、または他ユーザーの有効なOTPセッション期間中（手続き中）の指定にかかわらず、一貫してダミーOTPセッションを発行して 200 OK（応答遅延 1.0s ± 0.1s）を返却します。`

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.2 節のアカウント列挙防止対象から「現在と同一メールアドレスの指定」の記述を削除し、ログイン中ユーザーの同一メールアドレス指定時は 422 エラー（`SAME_AS_CURRENT_EMAIL`）を返却する要件仕様とドキュメント記述を統一しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
