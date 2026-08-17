# OTP再送信APIの欠落およびパスワードリセット発行時OTPセッション識別子の不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
新規登録およびメールアドレス変更のOTP入力画面において画面仕様・要件定義上存在する「再送信」機能に対応するAPIエンドポイントが定義されていません。また、パスワードリセット用OTP発行APIのレスポンスにOTPセッションIDが含まれておらず、後続のOTP検証APIを呼び出すことができません。

## 2. 詳細な指摘内容
1. **新規登録・メール変更のOTP再送エンドポイントの欠落**:
   - `docs/req-def/requirements.md` L228 および `docs/design/screen_design.md` L9, L49 では新規登録・メール変更の各画面に「再送信」ボタンが存在し、手動再送処理（60秒クールダウン）が規定されています。
   - しかし `docs/design/api_design.md` には `auth/password-reset/resend-otp` のみが存在し、`auth/register/resend-otp` や `auth/change-email/resend-otp` が欠落しています。
2. **パスワードリセット発行時のセッション識別子欠落**:
   - `docs/design/api_design.md` L18 の `auth/password-reset/request-otp` は出力が「汎用200 OKメッセージ」のみとなっています。
   - 一方で、後続の `auth/password-reset/verify-otp`（L19）および `auth/password-reset/resend-otp`（L20）の入力には `OTPセッションID` が必須となっています。
   - 発行時に `OTPセッションID`（または一時セッションCookie）が返却されない場合、クライアントは後続APIにどのセッションIDを渡せばよいか判別できません。

## 3. 推奨される修正案
1. 新規登録用 `auth/register/resend-otp` およびメール変更用 `auth/change-email/resend-otp` エンドポイントを定義に追加してください。
2. `auth/password-reset/request-otp` のレスポンスに `otp_session_id`（またはCookie等のセッション管理識別子）を含めるか、セッションの受け渡し方法を仕様として統一・明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` に `POST auth/register/resend-otp` および `POST auth/change-email/resend-otp` エンドポイントを正式に追加定義しました（60秒クールダウン等のエラー仕様含む）。
- `POST auth/password-reset/request-otp` のレスポンス仕様に `otp_session_id`（未登録メール時はダミーID）を明示的に返却する定義を追加し、後続の検証・再送APIへ確実に連携できるように修正しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
