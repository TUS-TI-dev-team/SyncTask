# TASKテーブルのUSER_IDカラムにおけるNOT NULL制約の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 15:40:00
- **Target Files**:
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
  - [requirements.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements.md)

## 1. 問題の概要
`TASK` テーブルの `USER_ID`（所有ユーザーID）カラムの制約定義が `VARCHAR(36) / FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` となっており、`NOT NULL` 制約が明示されていません。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の第2章「タスク管理 (TASK)」において、以下のように定義されています：
  - `USER_ID`: `VARCHAR(36)` / `FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)`
- `docs/req-def/requirements.md` のセキュリティ要件（187〜189行目）では、「本人以外は、そのアカウントの情報および所有タスクを見られない・操作できない」「タスク操作・閲覧の認可（BOLA/IDOR対策）: リクエスト元ログインユーザーのユーザーIDと操作対象タスクの所有ユーザーIDの一致を必須で検証する」と規定されています。
- リレーショナルデータベースの外部キー制約（Foreign Key）は、デフォルトで `NULL` 値の格納を許容します。そのため、`USER_ID` に `NOT NULL` 制約が付与されていない場合、所有者が存在しない孤立タスク（`USER_ID = NULL`）が作成・混入する余地が生じ、認可チェックやマルチテナントのデータ整合性にリスクをもたらします。

## 3. 推奨される修正案
`docs/design/database_design.md` の `TASK` テーブルにおける `USER_ID` カラムの制約定義を、明示的に `VARCHAR(36) / NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:41:40
- **Status**: Resolved

### 実施した修正内容
`docs/design/database_design.md` の `TASK` テーブル定義において、`USER_ID` カラムのデータ型 / 制約定義を `VARCHAR(36) / NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` に更新し、所有者を持たない孤立タスクの作成・混入を防ぐよう修正しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
