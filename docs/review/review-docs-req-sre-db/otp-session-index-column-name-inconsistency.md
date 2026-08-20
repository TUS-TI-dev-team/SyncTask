# OTP_SESSIONインデックス定義における存在しないカラム名参照の不備

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-18 22:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「7.2 セッション管理 (`LOGIN_SESSION`, `OTP_SESSION`)」において、`OTP_SESSION` テーブルに対するインデックス作成SQLで定義されていないカラム名（`MAX_EXPIRES_AT` および `EXPIRES_AT`）が参照されており、DDL実行時にSQLエラーが発生します。

## 2. 詳細な指摘内容
`database_design.md` の L208 および L211 に以下のインデックス定義があります：

```sql
-- OTPセッションのメールアドレス照会および状態遷移用
CREATE INDEX idx_otp_session_pending_email ON OTP_SESSION (PENDING_EMAIL, STATUS, EXPIRES_AT);

-- OTPセッションの15分間隔Cronパージ用（全体最大有効期限または失効レコードのクリーンアップ）
CREATE INDEX idx_otp_session_purge ON OTP_SESSION (MAX_EXPIRES_AT, STATUS, EXPIRES_AT);
```

しかし、同ドキュメントの「4. OTPセッション管理 (OTP_SESSION)」テーブル定義（L74-91）では、有効期限に関するカラム名は以下のように定義されています：

- `OTP_EXPIRES_AT`: 単発OTPの有効期限（発行から5分）
- `SESSION_EXPIRES_AT`: 手続き・セッション全体の最大有効期限（初回発行から15分）

テーブル定義に `EXPIRES_AT` や `MAX_EXPIRES_AT` というカラム名は存在しないため、上記DDLをそのまま実行すると PostgreSQL で `column "expires_at" does not exist` や `column "max_expires_at" does not exist` のエラーが発生します。

## 3. 推奨される修正案
インデックス定義のカラム名を、テーブル定義の正式なカラム名（`OTP_EXPIRES_AT` および `SESSION_EXPIRES_AT`）に修正してください：

```sql
-- OTPセッションのメールアドレス照会および状態遷移用
CREATE INDEX idx_otp_session_pending_email ON OTP_SESSION (PENDING_EMAIL, STATUS, OTP_EXPIRES_AT);

-- OTPセッションの15分間隔Cronパージ用（全体最大有効期限または失効レコードのクリーンアップ）
CREATE INDEX idx_otp_session_purge ON OTP_SESSION (SESSION_EXPIRES_AT, STATUS, OTP_EXPIRES_AT);
```
