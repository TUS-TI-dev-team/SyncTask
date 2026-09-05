# OTPセッション登録時の一意制約競合におけるダミー処理フォールバックの未実装

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-09-05 21:47:00
- **Target Files**:
  - [register_request_otp.go](backend/service/register_request_otp.go)
  - [register_request_otp.go](backend/repository/register_request_otp.go)
  - [register_request_otp_test.go](backend/service/register_request_otp_test.go)
  - [register_request_otp_test.go](backend/repository/register_request_otp_test.go)

## 1. 問題の概要
`backend/service/register_request_otp.go` および `backend/repository/register_request_otp.go` において、並行リクエスト等により同一メールアドレスに対する有効なOTPセッションの一意制約（`uq_otp_session_active_pending_email`）違反などの競合が発生した際、トランザクションをロールバックしてダミーセッション作成にフォールバックする制御が実装されておらず、クライアントへ 500 エラーが返却されてしまう可能性があります。

## 2. 詳細な指摘内容
`docs/design/process_design/01_account_creation.md` 1.2.2「競合処理」において、次のように明記されています：
> 判定後に同一メールアドレスのアカウント登録またはOTP作成が競合し、一意制約に抵触した場合はトランザクションをロールバックしてダミー処理へ切り替える。内部競合や登録状況をクライアントへ露出しない。

事前確認（`FindActiveUserByEmail`, `FindActiveOtpSessionByEmail`）で未登録と判定されても、同時に同じメールでリクエストが走った場合は DB レベルの一意制約に抵触します。このとき内部エラーを返さず、トランザクションをロールバックしてダミーセッション（`IS_DUMMY=true`）として再作成し、遅延後に同一構造の `200 OK` を返却する必要があります。

## 3. 推奨される修正案
1. `repository.SaveSessionWithLogs` または専用メソッドで、一意制約違反エラーを検知できるようにする（例: PostgreSQL の unique constraint violation / エラー判定、または専用のエラー型 `ErrConflict` を定義して返却）。
2. `service.RequestOtp` において、通常セッション作成で一意制約違反等の競合が発生した場合、ロールバック後にダミーセッションとして保存を再試行し、メール送信は行わずに遅延後に `200 OK` を返却するようフォールバック処理を実装する。
3. `service` および `repository` の単体テストにおいて、一意制約競合時にダミーセッションへフォールバックして 200 OK と遅延が返るテストケースを追加する。

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:50:30
- **Status**: Resolved

### 実施した修正内容
- `backend/repository/register_request_otp.go` に `ErrConflict` および一意制約違反判定（`isUniqueViolation`）を追加し、一意制約違反時にトランザクションをロールバックして `ErrConflict` を返すよう実装しました。
- `backend/service/register_request_otp.go` において、通常セッション保存時に `ErrConflict` が返却された場合に自動でダミーセッション作成・保存へ切り替え、実メール送信を行わずに遅延を適用して 200 OK を返却するようフォールバック処理を実装しました。
- `repository_test.go` および `service_test.go` に一意制約競合時ロールバック・ダミーフォールバックの単体テストケースを追加し、全パスを確認しました。

### 変更したファイル
- [register_request_otp.go](backend/repository/register_request_otp.go)
- [register_request_otp_test.go](backend/repository/register_request_otp_test.go)
- [register_request_otp.go](backend/service/register_request_otp.go)
- [register_request_otp_test.go](backend/service/register_request_otp_test.go)
