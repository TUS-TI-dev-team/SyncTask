# EMAILカラムの小文字正規化保存方針がDB設計に未記載

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 14:08:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書にて「メールアドレスはシステム内部で一律小文字（`toLowerCase()`）に正規化して保持・比較し、Case-Insensitive な一意性を保証する」と定義されていますが、`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムの備考にこの正規化ルールが記載されていません。

## 2. 詳細な指摘内容
- `docs/req-def/requirements.md` の259行目:
  > メールアドレスはシステム内部（登録・ログイン・重複判定・認証処理等）で一律小文字（`toLowerCase()`）に正規化して保持・比較し、Case-Insensitive な一意性を保証する
- `docs/design/database_design.md` の17行目:
  > `EMAIL` | `VARCHAR(255)` / `UNIQUE, NOT NULL` | 認証用メール（論理削除時は衝突回避のため退避形式へ更新）
- DB設計書の備考には「認証用メール」「退避形式」の記載のみで、**小文字正規化して格納する**というルールが明記されていません。
- この記載漏れにより、実装者がDB設計書のみを参照して開発した場合、大文字混在のまま格納してしまい、`UNIQUE` 制約による重複検知が正しく機能しない（PostgreSQLのデフォルトでは大文字小文字を区別するため、`User@Example.com` と `user@example.com` が別レコードとして登録される）リスクがあります。

## 3. 推奨される修正案
`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムの備考を以下のように補足してください：

- 現在: `認証用メール（論理削除時は衝突回避のため退避形式へ更新）`
- 修正後: `認証用メール（登録・更新時に一律小文字へ正規化して保存。論理削除時は衝突回避のため退避形式へ更新）`

あるいは、DB側でも `CITEXT` 型の使用やトリガーによる正規化を検討する場合はその旨を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 14:26:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` の `LOGIN_ACCOUNT` テーブル定義における `EMAIL` カラムの備考に「登録・更新時に一律小文字へ正規化して保存」するルールを明記しました。
- これにより、要件定義書（`requirements.md`）の小文字正規化仕様と整合させ、Case-Insensitive な一意性保証を担保しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
