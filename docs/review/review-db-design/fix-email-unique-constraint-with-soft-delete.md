# 論理削除アカウントとEMAIL一意性制約の再登録要件における整合性不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムに単純な `UNIQUE` 制約が指定されていますが、要件定義書に定められている「論理削除後の登録メールアドレスでの再登録は可能」「有効な他ユーザー（未削除アカウント）との重複なし」という仕様を満たすためのインデックス/制約設計の配慮が不足しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の17行目において、`EMAIL` カラムの制約が `VARCHAR(255) / UNIQUE, NOT NULL` と定義されています。また、19行目および31行目でアカウント削除時は論理削除（`IS_DELETED = TRUE`）を行うとされています。
- しかし、`docs/req-def/requirements.md` の53行目および260行目では以下のように定義されています：
  > 53: 論理削除後の登録メールアドレスでの再登録は可能
  > 260: 有効な他ユーザ（未削除アカウント）との重複なし
- テーブル全体に対する単純な `UNIQUE(EMAIL)` 制約が存在すると、一度論理削除されたアカウントのメールアドレスレコードがDBに残存しているため、同一メールアドレスで再登録しようとした際にDBの一意性制約エラーが発生してしまいます。

## 3. 推奨される修正案
- 単純なカラム制約 `UNIQUE` ではなく、有効な（未削除の）レコードのみを対象とする部分一意インデックス（例: PostgreSQL の `CREATE UNIQUE INDEX idx_unique_active_email ON LOGIN_ACCOUNT (EMAIL) WHERE IS_DELETED = FALSE`）を採用するか、
- あるいは論理削除時にメールアドレスにタイムスタンプやサフィックスを付与して退避する等の物理的な衝突回避方式について、DB設計書内に明確に設計方針・制約定義を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`LOGIN_ACCOUNT` テーブルの `EMAIL` カラム備考および削除方針の注記に、退会・論理削除時に `EMAIL` の値を退避形式（`deleted_<USER_ID>_<EMAIL>`）へ更新・退避する設計方針を明記しました。これにより、DB上の一意性制約 `UNIQUE(EMAIL)` を維持したまま、論理削除されたメールアドレスを用いた再登録が可能となりました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
