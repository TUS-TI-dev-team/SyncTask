# LOGIN_SESSIONテーブルのUSER_IDカラムにおけるNOT NULL制約の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 15:43:00
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [requirements.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_SESSION` テーブルの `USER_ID`（ログインユーザーID）カラムの制約定義が `VARCHAR(36) / FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` となっており、`NOT NULL` 制約が明示されていません。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の第3章「ログインセッション管理 (LOGIN_SESSION)」において、以下のように定義されています：
  - `USER_ID`: `VARCHAR(36)` / `FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)`
- 一方で、同文書の `TASK` テーブル（第2章）では、`USER_ID` カラムに対して `VARCHAR(36) / NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` と明示的な `NOT NULL` 制約が付与されています。
- `docs/req-def/requirements.md` のセッション管理要件において、ログインセッションはログイン認証が成功した正規ユーザーに対してのみ発行され、常に特定のユーザー（`USER_ID`）に紐づくものです（未認証状態の手続きは `OTP_SESSION` テーブルで管理されます）。
- リレーショナルデータベースの外部キー制約（Foreign Key）は、デフォルトで `NULL` 値の格納を許容します。そのため、`USER_ID` に明示的な `NOT NULL` 制約が付与されていない場合、所有者が存在しない孤立セッション（`USER_ID = NULL`）が作成・混入する余地が生じます。
- また、パスワード変更時、アカウント削除時、メールアドレス変更時、パスワードリセット完了時、パスワード再認証5回失敗時などに実行される「該当ユーザーのすべての既存ログインセッションを一括無効化（物理削除）」する処理（`DELETE FROM LOGIN_SESSION WHERE USER_ID = :user_id`）や認可制御のデータ整合性を担保するためにも、`USER_ID` は必須（`NOT NULL`）である必要があります。

## 3. 推奨される修正案
`docs/design/database_design.md` の `LOGIN_SESSION` テーブルにおける `USER_ID` カラムのデータ型 / 制約定義を、明示的に `VARCHAR(36) / NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:44:30
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の第3章「ログインセッション管理 (`LOGIN_SESSION`)」において、`USER_ID` カラムのデータ型/制約定義を `VARCHAR(36) / NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` に更新し、明示的な `NOT NULL` 制約を付与しました。これにより、NULL値を持つ孤立セッションレコードの混入を防ぎ、ユーザー別セッション一括削除や認可制御のデータ整合性を担保しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

