# パスワードリセット完了時における対象ユーザーのアクティブOTPセッション一括物理削除の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md)
  - [07_password_change.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/07_password_change.md)
  - [03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md)

## 1. 問題の概要
ログイン状態でのパスワード変更（`PATCH users/{user_id}/password`）やアカウント削除（`DELETE users/{user_id}`）では、確定トランザクション時に全ログインセッションだけでなく「対象ユーザーに紐づくアクティブな `OTP_SESSION`」を即座に一括物理削除することが明記されています。
しかし、パスワードリセット完了処理（`POST auth/password-reset/reset`）においては、「当該パスワードリセット用OTPセッション」と全ログインセッションの削除のみが記載されており、当該ユーザーに紐づく他のアクティブな `OTP_SESSION`（例: 並行して要求されていたメールアドレス変更用OTPセッション等）の物理削除が漏れています。

## 2. 詳細な指摘内容

1. **ドキュメント間のセッションパージ方針の不一致**:
   - `docs/design/process_design/07_password_change.md` 7.4.4節では、「`5. 対象ユーザーに紐づくアクティブな OTP_SESSION を物理削除し、並行中の認証手続きを失効させる。`」と明確に規定されています。
   - `docs/design/api_design/03_users.md` 3.2.4節でも「`全ログインセッション（LOGIN_SESSION）およびアクティブなOTPセッション（OTP_SESSION）を一括物理削除します`」とされています。
   - 一方、`docs/design/process_design/06_password_reset.md` 6.5.2節では「`3. 使用したパスワードリセット用の検証済み OTP_SESSION を物理削除する。`」とだけ記載されており、当該ユーザーの他のアクティブな `OTP_SESSION` のパージが明記されていません。
   - 同様に、`docs/design/api_design/02_auth.md` 3.1.9節および `docs/req-def/requirements/01_account_and_auth.md` 3.3節でも「当該パスワードリセット用OTPセッションをDBから物理削除」とのみ記載されています。

2. **セキュリティ上の影響**:
   - パスワードリセットはアカウントリカバリのための重要操作であり、リセット完了時には当該ユーザーに関連するすべての進行中手続き（他の端末等で発行されていたメールアドレス変更OTPセッション等）を無効化しなければ、古い手続きが悪用されるリスクが残ります。

## 3. 推奨される修正案

- `docs/design/api_design/02_auth.md`（3.1.9節）、`docs/design/process_design/06_password_reset.md`（6.5.2節）、および `docs/req-def/requirements/01_account_and_auth.md`（3.3節）において、パスワードリセット確定トランザクションで物理削除する対象として「当該パスワードリセット用OTPセッションおよび対象ユーザーに紐づくすべてのアクティブな `OTP_SESSION`」を明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
パスワードリセット完了確定トランザクションにおいて、使用したパスワードリセット用 `OTP_SESSION` だけでなく、当該ユーザーに紐づくすべてのアクティブな `OTP_SESSION` および全ログインセッション（`LOGIN_SESSION`）を一括物理削除する仕様を明記しました。

### 変更したファイル
- [02_auth.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\02_auth.md)
- [06_password_reset.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\process_design\06_password_reset.md)
- [01_account_and_auth.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\req-def\requirements\01_account_and_auth.md)
