# アカウント削除およびパスワード変更処理における OTP セッション (OTP_SESSION) 物理削除仕様の記述欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:05:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`03_users.md` の 3.2.3 節 (`DELETE users/{user_id}`) および 3.2.4 節 (`PATCH users/{user_id}/password`) の概要文および処理仕様において、アカウント削除およびパスワード変更成功時に「全ログインセッション (`LOGIN_SESSION`) が物理削除される」旨が明記されていますが、該当ユーザーに紐づくアクティブな OTP セッション (`OTP_SESSION`) の物理削除仕様の記載が欠落しています。

## 2. 詳細な指摘内容
- **`03_users.md` L87 (3.2.3 概要説明文)**:
  > なお、所有タスクデータおよび該当ユーザーの全ログインセッションは物理削除されます。
- **`03_users.md` L140 (3.2.4 概要説明文)**:
  > 現在のパスワードを検証した上で、新しいパスワードへ変更し、全セッションを一括物理削除します。

### 問題点
1. **`database_design.md` との仕様差分**:
   `database_design.md` (1. アカウント管理 NOTE L34-L35) では「セッション (`LOGIN_SESSION`, `OTP_SESSION`): ログアウト・アカウント削除時および期限切れ時は物理削除 (`DELETE`) されます。」と明確に定義されていますが、`03_users.md` では `LOGIN_SESSION` のみの削除と解釈できる表現になっており、ドキュメント間で整合性が取れていません。
2. **残留 OTP セッションによるセキュリティリスク**:
   アカウント削除またはパスワード変更前に発行されたメールアドレス変更用 OTP やパスワードリセット用 OTP などの `OTP_SESSION` が DB に物理削除されず残存した場合、論理削除後のアカウントやパスワード変更後のアカウントに対して旧 OTP セッションを利用した不正な検証試行や状態不整合が発生するリスクが生じます。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 節および 3.2.4 節の概要説明文および評価順序補足において、ログインセッション (`LOGIN_SESSION`) のみならず、該当ユーザー ID (`USER_ID`) に紐づくアクティブな OTP セッション (`OTP_SESSION`) もすべて即座に物理削除する旨を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.3 節 (`DELETE users/{user_id}`) および 3.2.4 節 (`PATCH users/{user_id}/password`) の概要説明文において、処理成功時に全ログインセッション（`LOGIN_SESSION`）に加えてアクティブな OTP セッション（`OTP_SESSION`）も即座に一括物理削除される旨を明記し、`database_design.md` との定義不整合を解消しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)

### 修正案の例

**3.2.3 概要説明文 (`DELETE users/{user_id}`)**:
```markdown
パスワード再認証を行い、アカウントを論理削除（`IS_DELETED=true`, `DELETED_AT=NOW()`）します。同メールアドレスでの将来の再登録を可能とするため、`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムを退避フォーマット（`deleted_<USER_ID>_<EMAIL>`）へ更新します。なお、所有タスクデータおよび該当ユーザーの全ログインセッション（`LOGIN_SESSION`）とアクティブなOTPセッション（`OTP_SESSION`）は即座に物理削除されます。
```

**3.2.4 概要説明文 (`PATCH users/{user_id}/password`)**:
```markdown
現在のパスワードを検証した上で、新しいパスワードへ変更し、該当ユーザーの全ログインセッション（`LOGIN_SESSION`）およびアクティブなOTPセッション（`OTP_SESSION`）を一括物理削除します。
```
