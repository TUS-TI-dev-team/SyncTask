# メールアドレス変更OTP発行APIにおけるクールダウン検証とビジネスルール検証の評価順序不備

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/change-email/request-otp`（メールアドレス変更OTP発行API）において、現在のメールアドレスと同一かどうかのビジネスルール検証（`422 SAME_AS_CURRENT_EMAIL`）が、60秒クールダウン検証（`429 OTP_RESEND_COOLDOWN`）よりも前のステップで評価される記述となっています。

## 2. 詳細な指摘内容
`docs/design/api_design/02_auth.md` の `3.1.10 POST auth/change-email/request-otp` のリクエスト評価順序は以下のように記述されています：

> 2. **リクエスト構文・入力バリデーション (`400 Bad Request` / `422 SAME_AS_CURRENT_EMAIL`)**:
> リクエストボディの `new_email` 形式を検証（不備時は 400 `BAD_REQUEST`）。現在のメールアドレスと同一かを検証（同一時は 422 `SAME_AS_CURRENT_EMAIL`）。
> 3. **アクティブセッション・クールダウン検証 (`429 Too Many Requests`)**:
> 該当ユーザーに対して既に有効な `active` のメールアドレス変更用OTPセッションが存在し、前回の発行から60秒未満である場合は `429 Too Many Requests`（code: `"OTP_RESEND_COOLDOWN"`）を返却します。

この評価順序の場合、前回のOTP発行から60秒以内のクールダウン期間中にユーザーが現在と同じメールアドレスを送信した場合、429（クールダウン制限）ではなく 422（`SAME_AS_CURRENT_EMAIL`）が先に返却されてしまいます。
一般的なセキュリティ・レートリミット設計では、インフラ/連打防止制御（429）をビジネスルール検証（422）よりも先に評価すべきであり、`3.2.2` や `3.2.4` とも評価順序の層（入力検証 400 → 認可 404 → レートリミット 429 → ビジネスルール 422）が不一致となっています。

## 3. 推奨される修正案
`3.1.10` の「リクエスト評価順序」を整理し、Step 2 を純粋な単一入力バリデーション（`400 BAD_REQUEST`）とし、Step 3 にクールダウン検証（`429 OTP_RESEND_COOLDOWN`）、Step 4 にビジネスルール検証（`422 SAME_AS_CURRENT_EMAIL`）を配置する順序に変更してください。
