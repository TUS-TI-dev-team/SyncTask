# OTP再送API群（3.1.3, 3.1.8, 3.1.12）における試行失敗カウンターリセットおよび有効期限更新仕様の記載漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:07:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`02_auth.md` の OTP 再送 API 群（3.1.3 `register/resend-otp`, 3.1.8 `password-reset/resend-otp`, 3.1.12 `change-email/resend-otp`）において、ユーザー操作による手動 OTP 再送成功時に `OTP_SESSION` テーブルで実施されるべき具体的な状態更新処理（`ATTEMPT_COUNT` の 0 リセット、`EXPIRES_AT` の 5分延長、`SEND_COUNT` の加算、`LAST_SENT_AT` の更新）についての明示的な記載が漏れています。

## 2. 詳細な指摘内容
1. **`ATTEMPT_COUNT` （試行失敗カウンター）リセット処理の非明示**:
   - `database_design.md` L87 にて `ATTEMPT_COUNT` は「1つのOTPに対して最大5回（5回失敗時は自動再送・失効制御）」と定義されています。
   - ユーザーが旧OTPで検証に失敗（例: 4回失敗）した後に手動で OTP 再送 API（`resend-otp`）を実行して新しい 8桁 OTP を受信した場合、手動再送に伴い `ATTEMPT_COUNT` を `0` にリセットしなければ、新OTP入力時に初回誤入力（通算5回目）で直ちに 422 自動再送エラー（`OTP_REISSUED_DUE_TO_FAILURES`）が発動してしまいます。
   - `02_auth.md` の 3.1.3, 3.1.8, 3.1.12 のいずれの説明・注記にも「手動再送成功時に `ATTEMPT_COUNT` を 0 にリセットする」旨が記載されていません。
2. **`EXPIRES_AT` （5分間有効期限）の更新と `MAX_EXPIRES_AT` 制約**:
   - OTP を再送した際、単一OTPの有効期限 `EXPIRES_AT` が再送時点から 5分後（`NOW() + 300s`）に更新されること（ただし初回発行からの全体最大有効期限 `MAX_EXPIRES_AT` 15分を超過しない範囲）、および `SEND_COUNT` が +1 加算される仕様が設計書本文で具体的に説明されていません。

## 3. 推奨される修正案
`02_auth.md` の 3.1.3, 3.1.8, 3.1.12 の本文注記（`※`）に以下の仕様を明記してください：

```markdown
※再送処理成功時、対象の `OTP_SESSION` レコードにおいて新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、試行失敗回数（`ATTEMPT_COUNT`）を 0 にリセット、送信回数（`SEND_COUNT`）を +1 加算、直前送信日時（`LAST_SENT_AT`）を更新するとともに、有効期限（`EXPIRES_AT`）を再送信時点から5分間（全体最大有効期限 `MAX_EXPIRES_AT` の範囲内）へ更新延長します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.3, 3.1.8, 3.1.12 の注記およびリクエスト評価順序において、手動 OTP 再送成功時に `OTP_SESSION` テーブルの `ATTEMPT_COUNT` を 0 にリセットし、`SEND_COUNT` 加算、`LAST_SENT_AT` 更新、`EXPIRES_AT` を 5分延長（`MAX_EXPIRES_AT` 15分の範囲内）する DB 状態更新仕様を明記しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
