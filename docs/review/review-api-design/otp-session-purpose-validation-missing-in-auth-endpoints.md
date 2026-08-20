# OTP検証・再送・リセットAPI群における `PURPOSE`（認証種別）検証の規定漏れ

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の OTP 検証・再送・リセット系の各エンドポイント（3.1.2, 3.1.3, 3.1.7, 3.1.8, 3.1.9, 3.1.11, 3.1.12）における「リクエスト評価順序」のセッション状態・期限検証ステップにおいて、`otp_session_id` の存在・ステータス（`active` / `verified`）・有効期限のみが検証条件として規定されており、`OTP_SESSION.PURPOSE`（`SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）が当該エンドポイントの要求種別と一致していることの検証規定が欠落している。

## 2. 詳細な指摘内容
`database_design.md` 4章（`OTP_SESSION`）では、OTPセッションの用途を識別するカラム `PURPOSE`（`SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）が定義されている。

しかし、`02_auth.md` の以下の各エンドポイントの「リクエスト評価順序」および `Errors` セクションでは、`PURPOSE` の照合が規定されていない：
- **3.1.2 `POST auth/register/verify-otp`**: L102 では `otp_session_id` の存在・`active` ステータス・有効期限のみ評価。`PURPOSE = 'SIGNUP'` の検証が未規定。
- **3.1.3 `POST auth/register/resend-otp`**: L151 では `otp_session_id` の存在・`active` ステータス・最大有効期限のみ評価。`PURPOSE = 'SIGNUP'` の検証が未規定。
- **3.1.7 `POST auth/password-reset/verify-otp`**: L332 では `otp_session_id` の存在・`active` ステータス・有効期限のみ評価。`PURPOSE = 'PASSWORD_RESET'` の検証が未規定。
- **3.1.8 `POST auth/password-reset/resend-otp`**: L382 では `otp_session_id` の存在・`active` ステータス・最大有効期限のみ評価。`PURPOSE = 'PASSWORD_RESET'` の検証が未規定。
- **3.1.9 `POST auth/password-reset/reset`**: L434 では `otp_session_id` の存在・`verified` ステータス・仮セッション有効期限のみ評価。`PURPOSE = 'PASSWORD_RESET'` の検証が未規定。
- **3.1.11 `POST auth/change-email/verify-otp`**: L542 では `otp_session_id` の存在・ユーザー紐づき・`active` ステータス・有効期限のみ評価。`PURPOSE = 'EMAIL_CHANGE'` の検証が未規定。
- **3.1.12 `POST auth/change-email/resend-otp`**: L597 では `otp_session_id` の存在・ユーザー紐づき・`active` ステータス・最大有効期限のみ評価。`PURPOSE = 'EMAIL_CHANGE'` の検証が未規定。

### セキュリティおよび実装上のリスク:
1. **異種用途セッションの不正流用リスク**: パスワードリセット用に発行された `otp_session_id` (`PURPOSE='PASSWORD_RESET'`) を `auth/register/verify-otp` に送信した場合、`PENDING_USERNAME` が `NULL` のまま登録処理が進み、DBレベルでの非Null制約違反（500エラー）や不正なアカウント作成が発生する可能性がある。
2. **アカウント乗っ取り・権限昇越リスク**: 新規登録用に発行された `otp_session_id` (`PURPOSE='SIGNUP'`) を `auth/password-reset/verify-otp` や `auth/password-reset/reset` に送信して検証状態へ遷移させ、他ユーザーのパスワードを不正上書き・変更される危険性がある。

## 3. 推奨される修正案
`02_auth.md` の対象7エンドポイント（3.1.2, 3.1.3, 3.1.7, 3.1.8, 3.1.9, 3.1.11, 3.1.12）の「リクエスト評価順序」におけるセッション検証ステップに、対象 `OTP_SESSION` の `PURPOSE` が該当エンドポイントの目的に合致していること（`SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）を検証する条件を明確に追記してください。不一致の場合は Timing Attack 対策として遅延 1.0s ± 0.1s を適用の上、`400 Bad Request`（code: `"BAD_REQUEST"`）を返却する仕様として明記することを推奨します。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の対象7エンドポイント（3.1.2, 3.1.3, 3.1.7, 3.1.8, 3.1.9, 3.1.11, 3.1.12）の「リクエスト評価順序」および注記・Errorsセクションにおいて、`OTP_SESSION.PURPOSE` がエンドポイントの目的に合致していること（`SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）を評価条件として追記しました。不一致・セッション不在時は Timing Attack 対策として遅延 `1.0s ± 0.1s` を適用した上で 400 Bad Request (code: `"BAD_REQUEST"`) または 403/410 を返却する仕様に統一しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
