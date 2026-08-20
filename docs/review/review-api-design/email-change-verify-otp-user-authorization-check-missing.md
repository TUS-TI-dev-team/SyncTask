# メールアドレス変更確定API (`POST auth/change-email/verify-otp`) における異ユーザー所有OTPセッション検証保護の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`POST auth/change-email/verify-otp`（3.1.11）において、リクエストで指定された `otp_session_id` に紐づく `OTP_SESSION.USER_ID` と、現在ログイン中のセッションユーザーID（`session.user_id`）の一致を検証する認可制御の仕様が明記されていません。他ユーザーの `otp_session_id` を指定して不正な検証やセッション破壊を試みる攻撃（IDOR）を防ぐための明確なエラー定義（`403 Forbidden`）が必要です。

## 2. 詳細な指摘内容
1. **異ユーザー間でのOTPセッション操作リスク（IDOR）**:
   - `database_design.md` L83 において `OTP_SESSION` テーブルには `USER_ID` カラムが定義されています。
   - `3.1.11 POST auth/change-email/verify-otp` は認証必須（Cookie）APIですが、バックエンドが「送られてきた `otp_session_id` の `USER_ID` が現在ログイン中のユーザーIDと一致しているか」を検証する旨が明記されていません。
   - 万一、ログイン中のユーザーAが他ユーザーBの `otp_session_id` を送信した場合に、バックエンドで所有者チェックを行わずにOTP検証やアカウントのメールアドレス更新を処理してしまうと、他ユーザーのアカウント所有権やセッション状態を不当に更新・破棄できてしまう脆弱性（IDOR）につながる恐れがあります。

## 3. 推奨される修正案
1. `3.1.11 POST auth/change-email/verify-otp` の処理説明に以下を明記してください:
   ```markdown
   ※指定された `otp_session_id` に紐づくユーザーID（`OTP_SESSION.USER_ID`）が、現在認証中のログインユーザーIDと一致しない場合は、他者リソースへのアクセス拒否として `403 Forbidden`（"code": "FORBIDDEN"）を返却します。
   ```
2. `3.1.11` の `##### Errors` に以下を追加してください:
   - `403 Forbidden`: `otp_session_id` の所有ユーザー不一致または CSRF トークン不正

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST auth/change-email/verify-otp` (3.1.11) において、`OTP_SESSION.USER_ID` と現在ログイン中のユーザーIDの一致検証（認可制御）を明記し、不一致時に `403 Forbidden` (`FORBIDDEN`) を返却する仕様およびエラー定義を追加しました。

### 変更したファイル
- [02_auth.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
