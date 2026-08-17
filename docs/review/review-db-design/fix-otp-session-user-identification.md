# OTP_SESSIONテーブルにおける認証種別および既存ユーザー識別カラムの不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`OTP_SESSION` テーブルに、新規登録・パスワードリセット・メールアドレス変更という各ユースケース（認証種別）を識別するカラムや、既存ユーザー（パスワードリセット・メールアドレス変更）を特定するための `USER_ID` カラムが不足しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の74〜86行目における `OTP_SESSION` テーブル定義では、新規アカウント登録用の `PENDING_USERNAME`, `PENDING_EMAIL`, `PENDING_PASSWORD_HASH` が配置されています。
- しかし、要件定義書（`docs/req-def/requirements.md`）では、メール認証（OTP）は以下の3つの用途で使用されます：
  1. アカウント新規作成（26〜33行目）
  2. メールアドレス変更（41〜49行目）
  3. パスワードリセット（68〜75行目）
- 現行のテーブル定義では、どの用途（`PURPOSE` または `TYPE`）で発行されたOTPセッションであるかが判別できず、パスワードリセット時に対象ユーザー（`USER_ID`）を特定してパスワード更新権限を付与する仕組みや、メールアドレス変更時に新旧どちらのメールアドレスおよび対象ユーザーを紐づけているのかがDBスキーマ上曖昧になっています。

## 3. 推奨される修正案
- 認証種別を示すカラム（例: `PURPOSE VARCHAR(20) NOT NULL`。値: `SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE` など）を追加してください。
- 既存ユーザーに対する操作（パスワードリセットやメールアドレス変更）を追跡・紐づけるため、`USER_ID VARCHAR(36)`（新規登録時はNULL許容）を追加することを検討・定義してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`OTP_SESSION` テーブルに、認証用途を示す `PURPOSE VARCHAR(20) NOT NULL`（`SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）と、既存ユーザー紐付け用の `USER_ID VARCHAR(36)`（新規登録時はNULL）を追加・定義しました。これにより、パスワードリセット時やメールアドレス変更時に対象ユーザーを確実に識別・追跡できるようになりました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
