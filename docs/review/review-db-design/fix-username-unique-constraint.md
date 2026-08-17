# USERNAMEのUNIQUE制約による同名ユーザー登録要件との不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `USERNAME` カラムに `UNIQUE` 制約が指定されていますが、要件定義書に記載されている「同名のユーザー名登録は可」という要件と矛盾しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の16行目において、`USERNAME` カラムの制約が `VARCHAR(20) / UNIQUE, NOT NULL` と定義されています。
- 一方で、`docs/req-def/requirements.md` の249行目（ユーザー名要件）では以下のように明記されています：
  > 同名のユーザー名登録は可（ログイン識別子にはメールアドレスを使用するため）
- そのため、DBテーブル定義に `UNIQUE` 制約が存在すると、既存ユーザーと同名のユーザー名で新規登録またはユーザー名変更を行った際に一意性制約違反（エラー）が発生し、要件を満たすことができなくなります。

## 3. 推奨される修正案
`LOGIN_ACCOUNT` テーブルの `USERNAME` カラムのデータ型/制約から `UNIQUE` を削除し、`VARCHAR(20) / NOT NULL` に変更してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`LOGIN_ACCOUNT` テーブルの `USERNAME` カラムから `UNIQUE` 制約を削除し、`VARCHAR(20) / NOT NULL`（備考: 2〜20文字、英大小数字［同名登録可］）に修正しました。これにより、ログイン識別子をメールアドレスとする仕様に合致し、同名のユーザー名での登録が可能となりました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
