# OTP検証・再送API群における OTP_SESSION ステータス前提条件および状態遷移違反エラー仕様の不透明さ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`02_auth.md` の OTP 検証API（3.1.2, 3.1.7, 3.1.11）および OTP 再送API（3.1.3, 3.1.8, 3.1.12）において、リクエスト対象となる `OTP_SESSION` の必須ステータス（前提条件 `STATUS == 'active'`）や、既に検証済み（`verified`）・失効済み（`expired`/`locked`）・完了済み（`completed`）のセッションに対して検証/再送リクエストが送信された場合のエラーコード仕様（`400 Bad Request` または `403 Forbidden` / `410 Gone`）が明示されていません。

## 2. 詳細な指摘内容
1. **ステータス前提条件の不透明さ**:
   - `database_design.md` L86 で `OTP_SESSION.STATUS` は `active`, `verified`, `expired`, `locked`, `completed` のいずれかをとると定義されています。
   - `02_auth.md` 3.1.7 (`password-reset/verify-otp`) では、検証成功時にステータスが `verified` に変更され、3.1.9 (`password-reset/reset`) でのみ使用可能になると規定されています。
   - しかし、以下のようなエッジケースにおける API の挙動およびエラー返却仕様が `02_auth.md` に規定されていません:
     - 既に `verified` に遷移した `otp_session_id` に対して再度 `verify-otp` や `resend-otp` を呼び出した場合。
     - 5回失敗により `locked` となった、あるいは期限切れで `expired` となった `otp_session_id` に対して `resend-otp` や `verify-otp` を呼び出した場合（なお 410 Gone は全体最大15分超過の記述が主となっており、ステータス非 `active` 時の扱いが不透明）。

2. **不正な重複検証・再送リスク**:
   - ステータスチェックの条件が曖昧な場合、一度検証に成功した OTP セッションに対して再送リクエストが通ってしまい、新しい OTP が再発行されてステータスが `active` に上書きされるなどの仕様矛盾が発生するおそれがあります。

## 3. 推奨される修正案
`02_auth.md` の 3.1.2, 3.1.3, 3.1.7, 3.1.8, 3.1.11, 3.1.12 の説明および Errors セクションに、以下の前提条件とエラー定義を明記してください。

1. **前提条件の明記**:
   - `POST .../verify-otp` および `POST .../resend-otp` のリクエスト処理時、対象 `OTP_SESSION` の `STATUS` が `active` であることを必須で検証する（`password-reset/reset` のみ `verified` であること）。
2. **Errors セクションへの追記**:
   - 既に `verified`, `locked`, `completed` 等の非 `active` ステータスである `otp_session_id` が指定された場合のエラーコードを明記する:
     - `400 Bad Request`: 無効なステータスのOTPセッション（`STATUS != 'active'`、code: `"BAD_REQUEST"`）または `410 Gone`（失効済みセッション、code: `"GONE"`）。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:56:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.2, 3.1.3, 3.1.7, 3.1.8, 3.1.11, 3.1.12 において、検証/再送対象の `OTP_SESSION` の前提条件として `STATUS == 'active'` であること（password-reset/reset のみ `verified` であること）を明記し、非 active ステータスの場合に `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"`）を即座に返却するエラー仕様を追加しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
