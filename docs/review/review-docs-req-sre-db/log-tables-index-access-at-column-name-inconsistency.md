# ログテーブル群のインデックス定義におけるカラム名（CREATED_AT vs ACCESS_AT）不整合

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-18 22:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「7.4 ログテーブル (`LOGIN_LOG`, `ACCESS_LOG`, `MAIL_AUTH_LOG`)」において、各ログテーブルに対するインデックス定義で存在しないカラム名 `CREATED_AT` が参照されており、テーブル定義の `ACCESS_AT` と乖離しているため DDL 実行時エラーになります。

## 2. 詳細な指摘内容
`database_design.md` の L223-232 に以下のインデックス定義が記述されています：

```sql
-- ログインログのIP/メール別照会および日次パージ
CREATE INDEX idx_login_log_ip ON LOGIN_LOG (IP_ADDRESS, CREATED_AT DESC);
CREATE INDEX idx_login_log_email ON LOGIN_LOG (EMAIL, CREATED_AT DESC);
CREATE INDEX idx_login_log_purge ON LOGIN_LOG (CREATED_AT);

-- アクセスログのパージ用
CREATE INDEX idx_access_log_purge ON ACCESS_LOG (CREATED_AT);

-- メール認証ログのメール別照会およびパージ用
CREATE INDEX idx_mail_auth_log_email ON MAIL_AUTH_LOG (EMAIL, CREATED_AT DESC);
CREATE INDEX idx_mail_auth_log_purge ON MAIL_AUTH_LOG (CREATED_AT);
```

しかし、同ドキュメントの「6. ログ管理 (LOGS)」におけるテーブル定義では、各テーブルの日時カラムは以下のように定義されています：

- `LOGIN_LOG` (L133): `ACCESS_AT` (`TIMESTAMPTZ` / `NOT NULL`)
- `ACCESS_LOG` (L151): `ACCESS_AT` (`TIMESTAMPTZ` / `NOT NULL`)
- `MAIL_AUTH_LOG` (L172-173): `ACCESS_AT` (`TIMESTAMPTZ` / `NOT NULL`)

いずれのログテーブルにも `CREATED_AT` カラムは存在しないため、上記DDL文を実行すると `column "created_at" does not exist` エラーとなりインデックスが作成できません。

## 3. 推奨される修正案
インデックス定義の `CREATED_AT` を、テーブル定義のカラム名である `ACCESS_AT` に統一・修正してください：

```sql
-- ログインログのIP/メール別照会および日次パージ
CREATE INDEX idx_login_log_ip ON LOGIN_LOG (IP_ADDRESS, ACCESS_AT DESC);
CREATE INDEX idx_login_log_email ON LOGIN_LOG (EMAIL, ACCESS_AT DESC);
CREATE INDEX idx_login_log_purge ON LOGIN_LOG (ACCESS_AT);

-- アクセスログのパージ用
CREATE INDEX idx_access_log_purge ON ACCESS_LOG (ACCESS_AT);

-- メール認証ログのメール別照会およびパージ用
CREATE INDEX idx_mail_auth_log_email ON MAIL_AUTH_LOG (EMAIL, ACCESS_AT DESC);
CREATE INDEX idx_mail_auth_log_purge ON MAIL_AUTH_LOG (ACCESS_AT);
```
