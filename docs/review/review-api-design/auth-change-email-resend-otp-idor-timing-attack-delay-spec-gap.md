# メールアドレス変更OTP再送APIにおける他者セッション指定時の Timing Attack 対策遅延記述漏れ

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/change-email/resend-otp`（メールアドレス変更OTP再送API）において、他ユーザー所有の `otp_session_id` を指定して `403 Forbidden` を返却する際、Timing Attack 対策の応答遅延（`1.0s ± 0.1s`）を適用するか否かの記述が漏れています。

## 2. 詳細な指摘内容
`docs/design/api_design/02_auth.md` の `3.1.12 POST auth/change-email/resend-otp` のリクエスト評価順序ステップ3では以下のように記述されています：

> 3. **認可・OTPセッション状態・目的・期限検証 (`400 Bad Request` / `403 Forbidden` / `410 Gone`)**:
> 指定された `otp_session_id` の存在および現在ログイン中のユーザーに紐づくセッションであることを検証（他者所有の場合は 403 `FORBIDDEN`）。用途 `PURPOSE` がメールアドレス変更（`EMAIL_CHANGE`）であること、ステータス（`active` であること）、および全体最大有効期限（`MAX_EXPIRES_AT` 15分）を検証します。セッション不在・`PURPOSE`不一致・失効時の場合は Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request` または `410 Gone` を返却します（ダミーセッション時含む）。

上記テキストでは、「セッション不在・PURPOSE不一致・失効時の場合」には `1.0s ± 0.1s` の遅延を適用すると明記されていますが、他者所有セッション判定時の `403 Forbidden` に対する遅延適用の有無が明記されていません。
もし他者所有セッションに対する `403 Forbidden` が遅延なしで即座に返却される場合、攻撃者がレスポンス時間の差を計測することで、指定した `otp_session_id` が「実在し他ユーザーに所有されているセッション」か「存在しないセッション」かを識別できてしまうタイミング攻撃のリスクが存在します。

## 3. 推奨される修正案
`3.1.12` の「リクエスト評価順序」ステップ3および Errors 節において、他ユーザー所有の `otp_session_id` 指定による `403 Forbidden` 返却時にも一律 `1.0s ± 0.1s` のレスポンス遅延を適用する旨を明確に記述してください。
