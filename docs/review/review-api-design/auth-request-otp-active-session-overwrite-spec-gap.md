# OTP発行要求API群における既存アクティブセッションの上書き・更新挙動の明記漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の OTP 発行要求エンドポイント（3.1.1 `auth/register/request-otp`, 3.1.6 `auth/password-reset/request-otp`, 3.1.10 `auth/change-email/request-otp`）において、クールダウン期間（60秒）経過後に同一メールアドレス（または同一ユーザー）から再度 `request-otp` が呼び出された場合、DB上の既存 `active` OTPセッションをどのように上書き・更新・置換するかについての明確な記述が欠落している。

## 2. 詳細な指摘内容
`database_design.md` L204 では、部分ユニークインデックス `CREATE UNIQUE INDEX uq_otp_session_active_pending_email ON OTP_SESSION (PENDING_EMAIL) WHERE STATUS IN ('active', 'verified');` が定義されている。

これにより、同一 `PENDING_EMAIL` に対して `active` または `verified` ステータスのレコードが複数存在することはDB制約上許容されない。

しかし、`02_auth.md` の 3.1.1, 3.1.6, 3.1.10 のリクエスト評価順序ステップ 3/4 では「未登録時は新規OTP発行」「登録済み時は新規OTP発行」と記述されているのみであり、クールダウン（60秒）経過後の再要求時に以下のようなDB制約違反リスクが存在する：
1. **DBユニーク制約違反（HTTP 500）の懸念**: 実装者が古い `active` レコードを破棄または更新せずに単純な `INSERT` を実行した場合、PostgreSQL の `UniqueViolation` エラーが発生してサーバー内部エラー（500）となる。
2. **OTPセッションレコードの蓄積・整合性不備**: 旧 `active` セッションのステータスを `expired` に変更するか、レコード内容を上書き更新するか、あるいは物理削除して新規作成するかの制御方針が仕様書に規定されていないため、実装による揺れや不整合が生じる。

## 3. 推奨される修正案
`02_auth.md` の 3.1.1, 3.1.6, 3.1.10 各節の仕様記述（「リクエスト評価順序」および注記）において、以下の挙動を明記してください：
- クールダウン期間（60秒）経過後に同一 `PENDING_EMAIL`（または同一ユーザー）に対して `request-otp` が呼び出された場合、DBの一意制約（`uq_otp_session_active_pending_email`）の競合を回避するため、既存の `active` レコードの属性（`OTP_HASH`, `ATTEMPT_COUNT=0`, `SEND_COUNT+=1`, `LAST_SENT_AT`, `EXPIRES_AT` 等）を上書き更新するか、旧 `active` レコードを物理削除/無効化した上で新レコードを発行することを明確に規定してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.1 (`auth/register/request-otp`), 3.1.6 (`auth/password-reset/request-otp`), 3.1.10 (`auth/change-email/request-otp`) の各リクエスト評価順序ステップ3/4の仕様記述において、60秒のクールダウン経過後に同一メールアドレス（または同一ユーザー）に対して再度 `request-otp` が呼び出された場合、DBの一意制約（`uq_otp_session_active_pending_email`）競合を回避するため、既存の `active` レコードの属性（`OTP_HASH`, `ATTEMPT_COUNT=0`, `SEND_COUNT+=1`, `LAST_SENT_AT`, `EXPIRES_AT` 等）を上書き更新して新OTPコードを発行・送信する挙動を明記しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
