# `DELETE users/{user_id}` API における論理削除時の `EMAIL` 退避処理および `DELETED_AT` 更新仕様の記載漏れ

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.3 節 (`DELETE users/{user_id}`) の概要説明文に、アカウント削除時の `DELETED_AT=NOW()` 設定および同一メールアドレスでの再登録を可能とする `EMAIL` カラムの退避処理（`deleted_<USER_ID>_<EMAIL>`）を明記しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
- **Severity**: Major
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`03_users.md` の `DELETE users/{user_id}` (3.2.3 節) の説明において、アカウント削除時の動作として「アカウントを論理削除（`IS_DELETED=true`）します。所有タスクデータおよび全セッションは物理削除されます。」と記載されているが、`DELETED_AT` のタイムスタンプ更新および同メールアドレスでの再登録を可能にするための `EMAIL` カラム退避フォーマット（`deleted_<USER_ID>_<EMAIL>`）への更新仕様の記載が欠落している。

## 2. 詳細な指摘内容
1. **DB UNIQUE 制約衝突リスク**:
   - `LOGIN_ACCOUNT` テーブルの `EMAIL` カラムには UNIQUE 制約 (`VARCHAR(320) / UNIQUE, NOT NULL`) が設定されている。
   - 要件定義書（`requirements.md` L53）およびデータベース設計書（`database_design.md` L31-L32）では、「論理削除後の同メールアドレスでの再登録を可能とするため、論理削除実行時に `EMAIL` カラムの値を退避フォーマット（例: `deleted_<USER_ID>_<EMAIL>`）に更新し、有効なアカウント間でのみ一意性を維持する」と明確に規定されている。
   - `03_users.md` 3.2.3 節にこの `EMAIL` 退避更新処理が明記されていないため、実装者が `IS_DELETED = TRUE` のみを行い `EMAIL` カラムを元のまま残す実装を行う懸念がある。その場合、退会ユーザーと同一のメールアドレスで新規登録を試みた際に DB の UNIQUE 制約違反エラーが発生し、再登録が不可能となる。

2. **`DELETED_AT` 更新の明記漏れ**:
   - 論理削除時に `DELETED_AT = NOW()` をセットする規定も 3.2.3 節の本文内に含まれていない。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 節 (`DELETE users/{user_id}`) の概要説明文に、論理削除時の `DELETED_AT` 設定および `EMAIL` 退避処理（`deleted_<USER_ID>_<EMAIL>`）を明記してください。

**修正案の例 (`03_users.md` 3.2.3 概要)**:
```markdown
#### 3.2.3 `DELETE users/{user_id}`
パスワード再認証を行い、アカウントを論理削除（`IS_DELETED=true`, `DELETED_AT=NOW()`）します。同メールアドレスでの将来の再登録を可能とするため、`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムを退避フォーマット（`deleted_<USER_ID>_<EMAIL>`）へ更新します。なお、所有タスクデータおよび該当ユーザーの全ログインセッションは物理削除されます。
```
