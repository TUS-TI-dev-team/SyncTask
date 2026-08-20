# `POST auth/password-reset/verify-otp` における検証成功時のステータス遷移（`verified`）および有効期限延長記述の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の 3.1.7 `POST auth/password-reset/verify-otp` のレスポンス説明（L213-L218）において、検証成功時に `OTP_SESSION.STATUS` が `verified` へ更新され、仮セッションの有効期限がその時点から15分間に延長されるという重要仕様の記述が欠落している。

## 2. 詳細な指摘内容
- **要件定義書およびDB設計書における定義**:
  - `requirements.md` L74: 「OTP検証成功時に当該OTPセッションのステータスを「検証済み」に変更し、そこから仮セッションの有効期限を15分間とする（この仮セッションは新パスワード設定APIへのアクセスのみを許可する制限付きセッションとし...）」
  - `database_design.md` L90: `EXPIRES_AT`: 「発行から5分（パスワードリセットの検証成功時はその時点から15分間に延長）」
- **`02_auth.md` 3.1.7 の記述**:
  ```json
  {
    "message": "OTP verified successfully."
  }
  ```
  ※レスポンスボディとメッセージのみが定義されており、検証成功後のDBステータス遷移（`verified`）および有効期限の15分延長についての注記が存在しない（3.1.9 のエラー欄 L291 に「検証成功後15分経過」と間接的に触れられているのみである）。

## 3. 推奨される修正案
`3.1.7` の `Response (200 OK)` の下に、以下の補足注記を追加してください：

`※検証成功時、当該OTPセッション（OTP_SESSION）のステータスを verified に変更し、有効期限（EXPIRES_AT）を検証成功時点から15分間に延長します（この検証済みOTPセッションは後続の POST auth/password-reset/reset エンドポイントでのみ使用可能となります）。`

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.7 (`POST auth/password-reset/verify-otp`) の `Response (200 OK)` の下に、検証成功時のステータス遷移（`verified`）および有効期限15分間延長の補足注記を追加しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
