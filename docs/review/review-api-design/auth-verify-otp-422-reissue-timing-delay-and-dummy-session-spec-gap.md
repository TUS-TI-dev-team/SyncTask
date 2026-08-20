# OTP検証API (`verify-otp`) における 5回連続失敗時422エラーの応答遅延 (1.0s ± 0.1s) およびダミーセッション挙動の記述不足

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の OTP検証API群（3.1.2, 3.1.7, 3.1.11）において、OTP誤入力5回連続到達時に返却される `422 Unprocessable Entity`（`code: "OTP_REISSUED_DUE_TO_FAILURES"`）について、レスポンス遅延（Timing Attack対策 `1.0s ± 0.1s`）の適用や、パスワードリセット・メールアドレス変更時のダミーセッションにおける同一応答挙動の明記が不足している。

## 2. 詳細な指摘内容
- **要件定義書および `01_overview.md` の規定**:
  - `requirements.md` L254: 「失敗またはダミー処理の場合は、一律1s±0.1sの範囲で応答速度を調整する」
  - `01_overview.md` L31: 「ダミーOTPセッションに対する後続の検証（verify-otp）や再送（resend-otp）に対しても、実セッションと全く同一のエラーコード（400, 410, 422, 429）および応答遅延（1.0s ± 0.1s）を適用します」
- **`02_auth.md` の現状記述**:
  - 3.1.2 (`register/verify-otp` L74): ダミーセッションについての言及はあるが、遅延 `1.0s ± 0.1s` の表記が抜けている。
  - 3.1.7 (`password-reset/verify-otp` L223): `- 422 Unprocessable Entity: 5回連続失敗に伴う自動再送実行通知（code: "OTP_REISSUED_DUE_TO_FAILURES"）` とのみあり、遅延 `1.0s ± 0.1s` および「未登録アドレス等のダミーセッション時も実際のメール再送は行わず同一の422応答を返還する」ことの明記がない。
  - 3.1.11 (`change-email/verify-otp` L365): 3.1.7 と同様に遅延表記およびダミーセッション時の同一応答（メール再送なし）についての記載が不足している。

## 3. 推奨される修正案
3.1.2, 3.1.7, 3.1.11 の各 `422 Unprocessable Entity` のエラー定義を以下のように統一して明確化してください：

```markdown
- `422 Unprocessable Entity`: 5回連続失敗に伴う自動再送実行通知（応答遅延 1.0s ± 0.1s、code: `"OTP_REISSUED_DUE_TO_FAILURES"`。ダミーセッション時も実際のメール再送を行わずに全く同一のレスポンスを返却）
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.2, 3.1.7, 3.1.11 における `422 Unprocessable Entity` のエラー定義に、応答遅延（1.0s ± 0.1s）の適用およびダミーセッション時も実際のメール再送を行わず同一レスポンスを返却する仕様を明記統合しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
