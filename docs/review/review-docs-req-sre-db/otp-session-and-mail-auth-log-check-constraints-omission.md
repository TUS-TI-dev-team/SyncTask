# OTP_SESSIONおよびMAIL_AUTH_LOGにおける列挙型カラムのCHECK制約定義の欠落

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:15:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「4. OTPセッション管理 (OTP_SESSION)」および「6.3 メール認証ログ (MAIL_AUTH_LOG)」において、列挙値（ENUMドメイン）を保持するカラム（`PURPOSE`, `STATUS`, `AUTH_TYPE`, `EVENT_TYPE`）に CHECK 制約が定義されておらず、DBレベルでのデータ整合性保証が不足しています。

## 2. 詳細な指摘内容
`database_design.md` では以下のカラムが `VARCHAR` 型で定義されていますが、許容される固定値に対する CHECK 制約の記載がありません：

1. **`OTP_SESSION` テーブル (L79, L85)**:
   - `PURPOSE`: `VARCHAR(20) / NOT NULL`（許容値: `SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）
   - `STATUS`: `VARCHAR(20) / NOT NULL, DEFAULT 'active'`（許容値: `active`, `verified`, `expired`, `completed`）

2. **`MAIL_AUTH_LOG` テーブル (L167, L169)**:
   - `AUTH_TYPE`: `VARCHAR(20) / NOT NULL`（許容値: `SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE`）
   - `EVENT_TYPE`: `VARCHAR(30) / NOT NULL`（許容値: `ISSUED`, `VERIFY_SUCCESS`, `VERIFY_FAILED`, `RESEND_REQUESTED`, `AUTO_RESEND`, `EXPIRED`）

### 問題点：
- アプリケーション側のバリデーション漏れやSQL直接実行、バッチ処理等の不具合により、未定義の文字列（例: `LOGIN`, `pending`, `INVALID` など）が混入するリスクがあります。
- `TASK` テーブルにおける `PRIORITY` / `STATUS` の CHECK 制約と同様に、DB層で `CHECK` 制約を明記しておくことで、不正データの混入を確実に防止できます。

## 3. 推奨される修正案
テーブル定義の「データ型 / 制約」列に各カラムの CHECK 制約を追記してください：

### `OTP_SESSION` テーブル
```markdown
| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| 認証種別 | `PURPOSE` | `VARCHAR(20)` / `NOT NULL, CHECK (PURPOSE IN ('SIGNUP', 'PASSWORD_RESET', 'EMAIL_CHANGE'))` | `SIGNUP` (新規登録), `PASSWORD_RESET` (パスワードリセット), `EMAIL_CHANGE` (メールアドレス変更) |
| ステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'active', CHECK (STATUS IN ('active', 'verified', 'expired', 'completed'))` | `active`, `verified`, `expired`, `completed` |
```

### `MAIL_AUTH_LOG` テーブル
```markdown
| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| 認証種別 | `AUTH_TYPE` | `VARCHAR(20)` / `NOT NULL, CHECK (AUTH_TYPE IN ('SIGNUP', 'PASSWORD_RESET', 'EMAIL_CHANGE'))` | `SIGNUP` (新規登録), `PASSWORD_RESET` (パスワードリセット), `EMAIL_CHANGE` (メールアドレス変更) |
| 処理イベント種別 | `EVENT_TYPE` | `VARCHAR(30)` / `NOT NULL, CHECK (EVENT_TYPE IN ('ISSUED', 'VERIFY_SUCCESS', 'VERIFY_FAILED', 'RESEND_REQUESTED', 'AUTO_RESEND', 'EXPIRED'))` | `ISSUED` (発行), `VERIFY_SUCCESS` (検証成功), `VERIFY_FAILED` (検証失敗), `RESEND_REQUESTED` (手動再送), `AUTO_RESEND` (5回失敗時自動処理), `EXPIRED` (有効期限切れ) |
```
