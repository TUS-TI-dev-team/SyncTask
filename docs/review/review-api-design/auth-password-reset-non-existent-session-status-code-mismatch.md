# パスワードリセット完了APIにおける存在しないOTPセッションID指定時のエラーコード不一致

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/password-reset/reset`（パスワードリセット完了API）において、存在しない `otp_session_id` を指定した場合に `403 Forbidden`（code: `"FORBIDDEN"`）を返却する仕様となっており、他のすべてのOTP関連エンドポイント（`400 Bad Request`）とステータスコードが不一致となっています。

## 2. 詳細な指摘内容
`docs/design/api_design/02_auth.md` の `3.1.9 POST auth/password-reset/reset` リクエスト評価順序ステップ2において、以下のように記述されています：

> 2. **OTPセッション状態・目的・期限検証 (`403 Forbidden` / `410 Gone`)**:
> 指定された `otp_session_id` の存在、用途 `PURPOSE` がパスワードリセット（`PASSWORD_RESET`）であること、ステータスが `verified` であること、および仮セッションの有効期限（検証成功後15分）内であることを検証します。未検証・`PURPOSE`不一致・無効時は `403 Forbidden`（code: `"FORBIDDEN"`）、期限切れ時は `410 Gone`（code: `"GONE"`）を返却します。

一方、`3.1.2` (`register/verify-otp`)、`3.1.3` (`register/resend-otp`)、`3.1.7` (`password-reset/verify-otp`)、`3.1.8` (`password-reset/resend-otp`)、`3.1.11` (`change-email/verify-otp`)、`3.1.12` (`change-email/resend-otp`) では、存在しない `otp_session_id` や無効なセッション指定時は一律 `400 Bad Request`（code: `"BAD_REQUEST"`）を返却する定義となっています。

`403 Forbidden` は、DB上に存在するOTPセッションのステータスが `verified` でない（`active` のまま等）場合や、他用途のセッション・他ユーザーのセッションを指定した場合の認可エラー用として使用されるべきであり、存在しないIDを指定した場合に `403` を返却するとAPI全体の設計整合性が損なわれます。

## 3. 推奨される修正案
`3.1.9` のリクエスト評価順序ステップ2および Errors セクションを修正し、`otp_session_id` が存在しない場合の応答を `400 Bad Request`（code: `"BAD_REQUEST"`）に変更し、`403 Forbidden` は「実在するセッションのステータスが `verified` でない場合または PURPOSE 不一致の場合」に限定してください。
