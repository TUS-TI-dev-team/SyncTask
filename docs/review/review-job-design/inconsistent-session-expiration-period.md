# ログインセッション有効期限の定義不一致

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:10:30
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`LOGIN_SESSION` の有効期限について、[database_design.md](docs/design/database_design.md) では「1ヶ月 (43200分)」と定義されているのに対し、[job_design.md](docs/design/job_design.md) では「1週間（10080分）」と記載されており、設計書間で矛盾・齟齬が生じています。

## 2. 詳細な指摘内容
- **[database_design.md](docs/design/database_design.md) line 66 (`LOGIN_SESSION.EXPIRES_AT`)**:
  > `1ヶ月 (43200分) / APIアクセス時にSliding Expirationで自動延長 / 期限切れは日次Cron（00:00 JST）で物理削除`
- **[job_design.md](docs/design/job_design.md) line 36 (`CLEANUP_EXPIRED_SESSIONS`)**:
  > `- **目的**: 1週間（10080分）の有効期限が経過した LOGIN_SESSION レコードを削除し、インデックスサイズおよびストレージ領域をクリーンに保つ。`

SQL自体は `WHERE EXPIRES_AT < NOW()` で実行されるためロジック上はカラムの値に従いますが、設計書の説明文で有効期限（1ヶ月 vs 1週間）に不一致があると、実装者やテスト設計者、運用者が混乱する原因となります。

## 3. 推奨される修正案
[job_design.md](docs/design/job_design.md) 第2章 2-2 の目的の記述を、[database_design.md](docs/design/database_design.md) の仕様に合わせて修正してください：

```markdown
- **目的**: 1ヶ月（43200分）の有効期限が経過した（`EXPIRES_AT < NOW()`）`LOGIN_SESSION` レコードを削除し、インデックスサイズおよびストレージ領域をクリーンに保つ。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- `job_design.md` の第3章「3-2. ログインセッションクリーンアップ (`CLEANUP_EXPIRED_SESSIONS`)」における目的の説明文を修正しました。
- `database_design.md` の記述（1ヶ月/43200分）と整合するように修正を行いました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
