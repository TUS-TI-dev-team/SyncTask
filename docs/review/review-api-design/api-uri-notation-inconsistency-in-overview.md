# 概要書内におけるAPI URIパス表記の表記揺れ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` 内で参照されている API の URI 表記において、1.2 セキュリティ仕様では `/api/users/{user_id}` や `/api/tasks/{task_id}` のように `/api/` プレフィックスを付与した絶対パス表記が用いられているのに対し、2. エンドポイント一覧テーブルでは `auth/register/request-otp` や `users/{user_id}` のように先頭のスラッシュ `/` および `/api/` を省略した相対パス表記が用いられており、表記揺れが発生している。

## 2. 詳細な指摘内容
1. **1.2 セキュリティ仕様（L31）**:
   `ユーザー情報（/api/users/{user_id}）およびタスク情報（/api/tasks/{task_id}）へのアクセス・変更・削除時は...` と記述。

2. **2. エンドポイント一覧テーブル（L74-94）**:
   - `auth/register/request-otp`
   - `users/{user_id}`
   - `tasks`
   - `tasks/{task_id}`

ベースURLとして L5 にて `https://<domain>/api/` が定義されているため、エンドポイント参照時のパス表記ルール（`/api/` 抜き相対パスで統一するか、先頭 `/` 付き相対パス `/users/{user_id}` とする等）を統一すべきである。

## 3. 推奨される修正案
`01_overview.md` 内の URI 表記をどちらか一方に統一してください。
ベースURL `https://<domain>/api/` からの相対パスとして表記する場合は、1.2 の記述を以下のように修整します:

```markdown
- **認可制御 (IDOR / BOLA 対策)**:
  - ユーザー情報（`users/{user_id}`）およびタスク情報（`tasks/{task_id}`）へのアクセス・変更・削除時は、セッション内のログインユーザーIDとリソースの所有ユーザーIDの一致を厳格に検証します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 内の URI 表記をベースURLからの相対パス（`/api/` プレフィックスなし）に統一し、表記揺れを解消しました。

### 変更したファイル
- [01_overview.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/01_overview.md)
