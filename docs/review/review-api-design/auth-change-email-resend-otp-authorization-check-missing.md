# `POST auth/change-email/resend-otp` における他ユーザーOTPセッション指定時の `403 Forbidden` 認可エラー定義欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の 3.1.12 `POST auth/change-email/resend-otp` のエラー仕様（L396-L400）において、`403 Forbidden` の理由として `CSRFトークン不正` のみが記載されており、3.1.11 (`change-email/verify-otp`) に記載されているような「他ユーザー所有の `otp_session_id` を指定した場合の認可エラー (403 Forbidden)」の定義が欠落している。

## 2. 詳細な指摘内容
- `3.1.11 POST auth/change-email/verify-otp` (L363):
  `- 403 Forbidden: 他ユーザー所有の otp_session_id 指定（認可不一致）または CSRFトークン不正（code: "FORBIDDEN"）`
- `3.1.12 POST auth/change-email/resend-otp` (L398):
  `- 403 Forbidden: CSRFトークン不正（code: "FORBIDDEN"）`

メールアドレス変更機能はログイン必須の認証付きAPIであり、`OTP_SESSION` レコードには所有者の `USER_ID` が紐づいている。
ログイン中のユーザー A がユーザー B の `otp_session_id` を指定して `resend-otp` を呼び出した場合、認可エラーとして処理すべきであるが、3.1.12 のエラー仕様にその旨が明記されていない。

## 3. 推奨される修正案
`3.1.12` の `Errors` セクションの `403 Forbidden`（L398）を以下のように修正してください：

```markdown
- `403 Forbidden`: 他ユーザー所有の `otp_session_id` 指定（認可不一致）または CSRFトークン不正（code: `"FORBIDDEN"`）
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.12 (`POST auth/change-email/resend-otp`) の `Errors` セクションにおける `403 Forbidden` の理由として、「他ユーザー所有の `otp_session_id` 指定（認可不一致）」を明記追記しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
