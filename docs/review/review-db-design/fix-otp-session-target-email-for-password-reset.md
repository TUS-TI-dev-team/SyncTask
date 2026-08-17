# OTP_SESSIONにおけるパスワードリセット時等の対象メールアドレス保持と排他制御の明確化

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 14:00:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`OTP_SESSION` テーブルの `PENDING_EMAIL` カラムの備考が「メール変更時 / 新規作成時」となっており、パスワードリセット時（`PURPOSE = 'PASSWORD_RESET'`）は `USER_ID` のみ保持して `PENDING_EMAIL` が `NULL` になる想定となっています。
しかし要件定義書では「OTPセッション有効期間中は、新規登録・メール変更・パスワードリセット手続き中においてメールアドレスに対する重複登録・変更リクエストを排他維持する」と規定されており、`PENDING_EMAIL` が `NULL` の場合、対象メールアドレスに対する排他確認のために `LOGIN_ACCOUNT` テーブルとの JOIN または追加クエリが常に必要となり、整合性チェックや排他制御が複雑化・漏洩するリスクがあります。

## 2. 詳細な指摘内容
- `requirements.md` の33行目、49行目、201行目:
  - 「OTPセッション有効期間中（初回発行から5分間、再送信を含め全体最大15分間）は、作成予定のメールアドレスに対する他ユーザーからの重複登録・変更リクエストを拒否（ダミーOTP処理対象）とし、認証確定または全体有効期限切れまで排他維持する」
- `database_design.md` の80〜82行目:
  - `USER_ID`: `VARCHAR(36)` / 既存ユーザー識別用（パスワードリセット・メール変更時。新規登録時はNULL）
  - `PENDING_EMAIL`: `VARCHAR(255)` / メール変更時 / 新規作成時
- パスワードリセット実行時に `PENDING_EMAIL` が `NULL` になると、あるメールアドレスに対して有効な `OTP_SESSION` が存在するかをチェックする際に `PENDING_EMAIL = ?` だけでなく `USER_ID IN (SELECT USER_ID FROM LOGIN_ACCOUNT WHERE EMAIL = ?)` の検索が必要となり、設計の一貫性が損なわれます。

## 3. 推奨される修正案
1. `PENDING_EMAIL` カラムの用途を拡張し、新規登録・メール変更だけでなく、**パスワードリセット時も含めた全認証種別における対象メールアドレス（排他制御用）** を格納する仕様に統一する（またはカラム名を `TARGET_EMAIL` / `PENDING_EMAIL` として備考を「認証対象 / 変更予定メールアドレス（全種別の排他制御に利用）」と明記する）。
2. これにより、あらゆるOTP認証フローにおいて `OTP_SESSION` テーブル単体で「指定メールアドレスに対する有効な認証セッションの有無（排他状態）」を高速かつ確実にクエリ・検証可能とします。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 14:26:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` において、`OTP_SESSION` テーブルの `PENDING_EMAIL` カラム（項目名: `認証対象/変更予定メールアドレス`）の備考を更新しました。
- 新規登録・メール変更だけでなく、パスワードリセット時も含めた全認証種別において対象メールアドレスを格納する仕様とし、`OTP_SESSION` テーブル単体での高速・確実な排他維持およびインデックス照会（`idx_otp_session_pending_email`）を可能としました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
