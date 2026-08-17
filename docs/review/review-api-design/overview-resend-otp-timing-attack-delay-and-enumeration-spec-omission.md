# `01_overview.md` における OTP 再送（`resend-otp`）時の Timing Attack 遅延制御およびアカウント列挙防止対象の明記漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1.2 節「セキュリティ & CSRF・アカウント列挙対策」において、アカウント列挙防止および Timing Attack 遅延制御（1.0s ± 0.1s）の対象として OTP 初回発行リクエスト（`request-otp`）は明記されているが、OTP 再送エンドポイント（`resend-otp`）における応答遅延およびアカウント列挙防止についての記述が欠落している。

## 2. 詳細な指摘内容
1. **アカウント列挙防止の対象エンドポイント記述漏れ（L32）**:
   L32 では「新規登録（`auth/register/request-otp`）、パスワードリセット（`auth/password-reset/request-otp`）、メールアドレス変更（`auth/change-email/request-otp`）において...」と `request-otp` の 3 エンドポイントのみが言及されている。
   しかし、個別仕様書 `02_auth.md`（L127, L325, L507等）では、`resend-otp` エンドポイントにおいても実セッション・ダミーセッションを問わず一律 1.0s ± 0.1s の応答遅延を適用してアカウント存在有無の推測を防止している。

2. **遅延制御対象一覧における表記の曖昧さ（L39）**:
   L39 の遅延制御適用条件にて「ログイン失敗、OTP発行処理（正常成功時およびアカウント存在有無のダミー処理時を含む一括）、OTP検証失敗、パスワード再認証失敗時」と記載されているが、「OTP発行処理」に `resend-otp`（OTP再送処理）が含まれるかが曖昧である。

## 3. 推奨される修正案
`01_overview.md` 1.2 節の記述を以下のように修正し、OTP再送処理（`resend-otp`）もアカウント列挙防止および Timing Attack 遅延制御の共通対象であることを明確化してください。

```markdown
- **アカウント列挙防止 (User Enumeration 対策)**:
  - 新規登録（`auth/register/request-otp`）、パスワードリセット（`auth/password-reset/request-otp`）、メールアドレス変更（`auth/change-email/request-otp`）および各種 OTP 再送（`resend-otp`）において、指定されたメールアドレスの登録有無...

- **遅延制御 (Timing Attack 対策)**:
  - ログイン失敗、OTP発行・再送処理（正常成功時およびアカウント存在有無のダミー処理時を含む一括）、OTP検証失敗、パスワード再認証失敗時は、一律 `1.0s ± 0.1s` のレスポンス遅延を適用します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:18:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.2 節の「アカウント列挙防止 (User Enumeration 対策)」および「遅延制御 (Timing Attack 対策)」の共通仕様記述を更新し、OTP初回発行（`request-otp`）に加えて各種 OTP 再送処理（`resend-otp`）もアカウント列挙防止および `1.0s ± 0.1s` レスポンス遅延制御の明確な対象として明記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)

