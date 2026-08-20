# アカウント削除処理における未登録/仮登録OTPセッションの物理削除範囲の曖昧さ

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-17 18:00:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`03_users.md` の `DELETE users/{user_id}` (3.2.3) において、アカウント論理削除に伴うデータ消去仕様として「該当ユーザーのアクティブなOTPセッション（`OTP_SESSION`）は即座に物理削除されます」と記述されている。
しかし、新規登録（`SIGNUP`）用 OTP セッションは `database_design.md` 上 `USER_ID` が `NULL` であり、`PENDING_EMAIL` カラムでメールアドレスが管理されているため、`USER_ID` 条件のみで削除処理を行った場合に該当メールアドレスの新規登録中OTPセッションが残存する懸念がある。

## 2. 詳細な指摘内容
1. **`OTP_SESSION` テーブルの構造**:
   - `database_design.md` Section 4 の `OTP_SESSION` テーブル定義（L81）によると、`USER_ID` カラムはパスワードリセット・メールアドレス変更時のみ設定され、新規登録（`SIGNUP`）時は `NULL` となる。
   - 新規登録OTPセッションでは対象メールアドレスが `PENDING_EMAIL` カラムに保存される。
2. **アカウント削除時の削除漏れリスク**:
   - アカウント削除時に `OTP_SESSION.USER_ID = <user_id>` のみ条件として物理削除を行うと、同じメールアドレスに対して別途発行されていたアクティブな `SIGNUP` OTPセッション（`USER_ID` が `NULL` のレコード）が削除されずに DB 上に残存する可能性がある。
   - `OTP_SESSION` テーブルには `uq_otp_session_active_pending_email`（`PENDING_EMAIL` に対する部分一意インデックス）が定義されているため、古い `SIGNUP` OTPセッションが残存していると、削除直後に同メールアドレスで再登録を試みた際に一意性制約違反や重複エラー等の不整合が発生するリスクがある。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 (`DELETE users/{user_id}`) の仕様説明において、物理削除対象となる `OTP_SESSION` レコードの抽出条件として `USER_ID = <user_id>` に加え、削除対象ユーザーの退避前メールアドレスと一致する `PENDING_EMAIL` を持つ `OTP_SESSION` レコードも含まれることを明記してください。

**修正案の例**:
> なお、所有タスクデータおよび該当ユーザーの全ログインセッション（`LOGIN_SESSION`）と、該当ユーザーIDまたは該当メールアドレス（`PENDING_EMAIL`）に紐づく全アクティブなOTPセッション（`OTP_SESSION`）は即座に物理削除されます。
