# ログテーブル（LOGIN_LOG, MAIL_AUTH_LOG）におけるEMAILカラム桁数の不整合と桁あふれリスク

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 15:07:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムは論理削除時の退避形式（`deleted_<USER_ID>_<EMAIL>`）に対応するため `VARCHAR(320)` に拡張されていますが、`LOGIN_LOG` および `MAIL_AUTH_LOG` テーブルの `EMAIL` カラムは `VARCHAR(255)` のままとなっています。これにより論理削除済みアカウントへのログイン試行ログや長大なメールアドレスでの認証ログ記録時に桁あふれ例外が発生するリスクがあります。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の17行目および32行目において、論理削除後の同メール再登録対応のため `LOGIN_ACCOUNT.EMAIL` の型が `VARCHAR(320)` に設定されています。
- しかし、同ファイル内のログテーブル定義において：
  - `LOGIN_LOG.EMAIL`（130行目）: `VARCHAR(255) / NOT NULL`
  - `MAIL_AUTH_LOG.EMAIL`（166行目）: `VARCHAR(255) / NOT NULL`
  - （参考: `OTP_SESSION.PENDING_EMAIL` 83行目も `VARCHAR(255)`）
  と定義されています。
- 論理削除済みアカウントの退避フォーマット（例: `deleted_123e4567-e89b-12d3-a456-426614174000_user@example.com`）や、長大なメールアドレスを用いてログイン試行・認証要求が行われた際、ログ記録 INSERT 処理で `value too long for type character varying(255)` などの DB 実行時例外が発生し、監査ログの記録失敗や内部エラー（HTTP 500）を引き起こす可能性があります。

## 3. 推奨される修正案
`docs/design/database_design.md` において、`LOGIN_LOG` および `MAIL_AUTH_LOG`（必要に応じて `OTP_SESSION.PENDING_EMAIL` も含む）の `EMAIL` カラムのデータ型を `VARCHAR(320)` に拡張・統一してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:11:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` において、`LOGIN_LOG` および `MAIL_AUTH_LOG` テーブルの `EMAIL` カラムのデータ型を `VARCHAR(320)` に拡張し、論理削除アカウントや長大メールアドレスに対するログ記録時の桁あふれ例外を防止しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
