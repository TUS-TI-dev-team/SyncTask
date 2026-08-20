# 各テーブルの日時カラムにおけるDEFAULT NOW()制約定義の欠落

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:15:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` において、全テーブルの作成日時・更新日時・アクセス日時などの `TIMESTAMPTZ` 型カラムに `NOT NULL` 制約が付与されているものの、`DEFAULT CURRENT_TIMESTAMP`（または `DEFAULT NOW()`）のデフォルト値が定義されておらず、DB層での自動タイムスタンプ補完が明記されていません。

## 2. 詳細な指摘内容
`database_design.md` の各テーブル定義において、以下のようなタイムスタンプカラムが存在します：

- `LOGIN_ACCOUNT`: `CREATED_AT` (`TIMESTAMPTZ / NOT NULL`), `UPDATED_AT` (`TIMESTAMPTZ / NOT NULL`)
- `TASK`: `CREATED_AT` (`TIMESTAMPTZ / NOT NULL`), `UPDATED_AT` (`TIMESTAMPTZ / NOT NULL`)
- `LOGIN_SESSION`: `CREATED_AT` (`TIMESTAMPTZ / NOT NULL`)
- `OTP_SESSION`: `CREATED_AT` (`TIMESTAMPTZ / NOT NULL`), `LAST_SENT_AT` (`TIMESTAMPTZ / NOT NULL`)
- `LOGIN_IP_RATE_LIMIT`: `LAST_FAILED_AT` (`TIMESTAMPTZ / NOT NULL`), `UPDATED_AT` (`TIMESTAMPTZ / NOT NULL`)
- `LOGIN_LOG`: `ACCESS_AT` (`TIMESTAMPTZ / NOT NULL`)
- `ACCESS_LOG`: `ACCESS_AT` (`TIMESTAMPTZ / NOT NULL`)
- `MAIL_AUTH_LOG`: `ACCESS_AT` (`TIMESTAMPTZ / NOT NULL`)

### 問題点：
- レコード挿入（INSERT）時にアプリケーション側でタイムスタンプの指定が漏れた場合、`null value in column "created_at" violates not-null constraint` のエラーが発生します。
- DBスキーマのベストプラクティスとして、作成日時やアクセス日時には `DEFAULT NOW()`（または `DEFAULT CURRENT_TIMESTAMP`）を明記しておくことで、データ登録時の堅牢性が向上します。

## 3. 推奨される修正案
テーブル定義の「データ型 / 制約」列において、タイムスタンプカラムに `DEFAULT NOW()` を追記してください：

例（`LOGIN_ACCOUNT` テーブル）:
```markdown
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL, DEFAULT NOW()` | |
| 更新日時 | `UPDATED_AT` | `TIMESTAMPTZ` / `NOT NULL, DEFAULT NOW()` | |
```
（`TASK`, `LOGIN_SESSION`, `OTP_SESSION`, `LOGIN_IP_RATE_LIMIT`, 各種ログテーブルも同様に `DEFAULT NOW()` を定義）
