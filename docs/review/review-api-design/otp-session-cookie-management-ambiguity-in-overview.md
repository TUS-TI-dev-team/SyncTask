# OverviewにおけるOTPセッション管理方式の（または一時Cookie管理）記述による曖昧性

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:00:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 「1.1 セッション管理 & 認証方式」において、OTPセッションの送受信方式として `手続きごとの otp_session_id をリクエストボディで送受信（または一時Cookie管理）します。` と記載されています。しかし、`02_auth.md` に定義されている12個の認証API（新規登録、パスワードリセット、メールアドレス変更）では、すべてJSONレスポンスボディおよびリクエストボディで `otp_session_id` を明示的に送受信する設計となっており、CookieによるOTPセッション管理は一切使用されません。括弧書きの `（または一時Cookie管理）` という表現は実装方針に曖昧さを生じさせる可能性があります。

## 2. 詳細な指摘内容
1. **01_overview.md L23**:
   > `アカウント新規作成、パスワードリセット、メールアドレス変更の手続き中は、手続きごとの otp_session_id をリクエストボディで送受信（または一時Cookie管理）します。`

2. **02_auth.md の仕様との乖離**:
   - `auth/register/request-otp` (L26): レスポンスボディ `{"otp_session_id": "otp_sess_a1b2c3d4e5", ...}`
   - `auth/register/verify-otp` (L46): リクエストボディ `{"otp_session_id": "otp_sess_a1b2c3d4e5", "otp": "..."}`
   - `auth/password-reset/request-otp` / `verify-otp` / `resend-otp` / `reset`: すべてJSONボディにて `otp_session_id` を送受信。
   - `auth/change-email/request-otp` / `verify-otp` / `resend-otp`: すべてJSONボディにて `otp_session_id` を送受信。

Cookieによる一時管理は提供されず、JSONボディでの送受信に一貫して統一されているため、`（または一時Cookie管理）` の記述は不要な混乱を招く表現となっています。

## 3. 推奨される修正案
`01_overview.md` の L23 の記述から `（または一時Cookie管理）` を削除し、JSONリクエスト/レスポンスボディでの送受信であることを明確化してください。

```markdown
- **OTPセッション**:
  - アカウント新規作成、パスワードリセット、メールアドレス変更の手続き中は、手続きごとの `otp_session_id` をリクエストボディおよびレスポンスボディで送受信します。
  - OTP有効期限は発行から5分（手続き全体の最大有効期限は15分）です。
```

## 修正完了報告

- **Resolved At**: 2026-08-17 16:56:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.1 節の記述から `（または一時Cookie管理）` を削除し、`otp_session_id` がリクエストボディおよびレスポンスボディで送受信される仕様であることを明確化しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
