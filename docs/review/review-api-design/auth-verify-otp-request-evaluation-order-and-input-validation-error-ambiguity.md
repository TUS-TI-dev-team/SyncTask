# OTP検証API群におけるリクエスト評価順序および入力バリデーションエラー区別の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の OTP 検証 API 群（3.1.2 `register/verify-otp`, 3.1.7 `password-reset/verify-otp`, 3.1.11 `change-email/verify-otp`）において、エラー仕様 (`##### Errors`) セクションにおける `400 Bad Request` が「`OTP不一致`」としてのみ定義されており、リクエストパラメータ欠落や桁数・文字種不備などの入力値バリデーション違反との評価順序および挙動の区別が明記されていません。

## 2. 詳細な指摘内容
1. **リクエスト評価順序の未定義**:
   - `03_users.md` の 3.2.3 (`DELETE user`) や 3.2.4 (`PATCH password`) では、「リクエスト評価順序」セクションが用意され、1. リクエスト構文・入力バリデーション (`400 Bad Request`, 即時返却、遅延なし、試行カウンター加算なし) と 2. 認証・ハッシュ照合 (遅延 `1.0s ± 0.1s`, 失敗カウンター加算) の2段階が明確に切り分けられています。
   - 一方 `02_auth.md` の 3.1.2, 3.1.7, 3.1.11 では `400 Bad Request` に「`OTP不一致（入力試行5回未満...遅延 1.0s ± 0.1s、code: "BAD_REQUEST"）`」とのみ書かれています。
   - これにより、以下のような疑問や実装上のブレが発生します:
     - リクエストボディの JSON 形式不正や必須項目 `otp_session_id` / `otp` の欠落、または `otp` が 8 桁未満/記号を含む場合の入力バリデーション違反（`400 Bad Request`）で、人工遅延（`1.0s ± 0.1s`）を適用すべきか否か。
     - 入力形式違反のリクエストに対して `OTP_SESSION.ATTEMPT_COUNT`（試行失敗回数）をカウントアップすべきか否か。

2. **誤ったロックアウト発生リスク**:
   - 入力バリデーションエラー時にも `ATTEMPT_COUNT` を加算する実装となった場合、クライアント側のフォーマット入力ミス（桁不足やタイポ）によって正当な OTP 検証試行権限が不当に消化され、5回失敗による自動再送・カウンターリセットが意図せず誘発されるリスクがあります。

## 3. 推奨される修正案
`02_auth.md` の 3.1.2, 3.1.7, 3.1.11 の各節に「リクエスト評価順序」を追記し、入力構文・バリデーションエラー（即時 400 返却、遅延なし、試行回数加算なし）と OTP 値照合不一致（400 返却、遅延 1.0s ± 0.1s 適用、試行失敗回数 `ATTEMPT_COUNT` カウントアップ）を明確に区分してください。

```markdown
##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`, `otp`）の有無、および `otp` の形式（英数字8桁）を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし、試行回数 `ATTEMPT_COUNT` 加算なし）。
2. **OTPセッション状態・期限検証 (`410 Gone` / `403 Forbidden`)**:
   指定された `otp_session_id` の存在、ステータス（`active`）、および有効期限（`EXPIRES_AT` / `MAX_EXPIRES_AT`）を検証します。
3. **OTP照合検証 (`400 BAD_REQUEST` / `422 OTP_REISSUED_DUE_TO_FAILURES`)**:
   入力された `otp` のハッシュ照合を実施します。
   - 不一致（試行1〜4回目）: 失敗回数（`ATTEMPT_COUNT`）を+1加算し、`400 Bad Request`（code: `"BAD_REQUEST"`、遅延 1.0s ± 0.1s）を返却します。
   - 不一致（試行5回達成）: 失敗回数をリセットし、OTP自動再発行通知 `422 Unprocessable Entity`（code: `"OTP_REISSUED_DUE_TO_FAILURES"`、遅延 1.0s ± 0.1s）を返却します。
```
