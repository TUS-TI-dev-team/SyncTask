# パスワードリセット完了APIにおける Timing Attack 対策遅延仕様の記述漏れ

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/password-reset/reset`（パスワードリセット完了API）において、無効・未検証のOTPセッション指定（`403 Forbidden` / `410 Gone`）やビジネスルール違反（`422 Unprocessable Entity`）発生時のレスポンス遅延（`1.0s ± 0.1s`）に関する明記が漏れています。

## 2. 詳細な指摘内容
`docs/design/api_design/01_overview.md` 1.2 節（遅延制御 / Timing Attack 対策）では、認証失敗およびOTP検証失敗時には一律 `1.0s ± 0.1s` のレスポンス遅延を適用する方針が定められています。
また、`02_auth.md` の他のOTP検証エンドポイント（`3.1.2`, `3.1.7`, `3.1.11` 等）では、無効なセッション指定や検証失敗時に `一律 1.0s ± 0.1s の遅延を適用` する旨が明確に規定されています。

しかし、`3.1.9 POST auth/password-reset/reset` の「リクエスト評価順序」および Errors 節には遅延処理に関する記述が一切存在しません。
もし遅延なしで即座に `403` や `410` や `422` が返却された場合、攻撃者がレスポンス時間の差を計測することにより、提示した `otp_session_id` が実在・検証済みであるか否かを探索・推測（Timing Attack）される脆弱性のリスクが生じます。

## 3. 推奨される修正案
`3.1.9` の「リクエスト評価順序」ステップ2（`403` / `410`）、ステップ3（`422`）および Errors 節に、Timing Attack 対策として一律 `1.0s ± 0.1s` のレスポンス遅延を適用する旨を明記してください。
