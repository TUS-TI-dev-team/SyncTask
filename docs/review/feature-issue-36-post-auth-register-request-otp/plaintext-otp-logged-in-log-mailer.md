# LogMailer における平文OTPログ出力の禁止違反

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-09-05 21:47:00
- **Target Files**:
  - [mailer.go](backend/service/mailer.go)
  - [mailer_test.go](backend/service/mailer_test.go)

## 1. 問題の概要
`backend/service/mailer.go` の `LogMailer.SendOTP` において、`log.Printf("[LogMailer] Send OTP to %s: %s", toEmail, otp)` として平文OTPが標準ログに出力されています。

## 2. 詳細な指摘内容
`docs/design/process_design/01_account_creation.md` 1.2.1 にて次のように規定されています：
> 平文パスワードおよび平文OTPはDB・アプリケーションログへ記録しない。

開発用モック・ログ用Mailerであっても、平文OTPをログ出力することはセキュリティ要件に違反します。

## 3. 推奨される修正案
`LogMailer.SendOTP` において、平文OTPをマスクする（例: `***` や長さを出力するのみ）、あるいはメール送信イベントのみを出力するように修正してください。
単体テストでもログに平文OTPが出力されないことを確認してください。

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:50:30
- **Status**: Resolved

### 実施した修正内容
- `LogMailer.SendOTP` において、平文OTPを出力せず `[REDACTED]` とOTP文字列長のみログ出力するよう変更しました。
- 単体テスト `TestLogMailer_SendOTP` を作成し、ログ出力に平文OTPが含まれないことを検証しました。

### 変更したファイル
- [mailer.go](backend/service/mailer.go)
- [mailer_test.go](backend/service/mailer_test.go)
