# ログインセッション有効期限の要件定義書との不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`job_design.md` のログインセッションクリーンアップ処理（`CLEANUP_EXPIRED_SESSIONS`）の目的において「1週間（10080分）の有効期限が経過した」と記載されていますが、要件定義書（`requirements.md`）ではログインセッションの有効期限は「1ヶ月（43200分）」と規定されており、矛盾しています。

## 2. 詳細な指摘内容
- 要件定義書（[requirements.md:L191](docs/req-def/requirements.md#L191), [requirements.md:L267-274](docs/req-def/requirements.md#L267-L274)）では、「ログイン認証のセッション期限は1ヶ月とする（APIアクセスごとに自動延長される Sliding Expiration 方式）」「有効期限 43200分(1ヶ月)」と明確に定められています。
- しかし、[job_design.md:L36](docs/design/job_design.md#L36) では以下のように記載されています：
  > 目的: 1週間（10080分）の有効期限が経過した LOGIN_SESSION レコードを削除し、インデックスサイズおよびストレージ領域をクリーンに保つ。

クリーンアップ SQL 自体は `WHERE EXPIRES_AT < NOW()` となっているものの、設計書内の仕様説明で「1週間（10080分）」と誤った期間が記載されているため、実装担当者や運用担当者に誤解を与え、セッション発行側の実装値との不整合や混乱を招く恐れがあります。

## 3. 推奨される修正案
[job_design.md:L36](docs/design/job_design.md#L36) の説明文を要件定義書に合わせて修正してください。

```markdown
- **目的**: 1ヶ月（43200分、Sliding Expirationにより更新された最終有効期限）が経過した `LOGIN_SESSION` レコードを削除し、インデックスサイズおよびストレージ領域をクリーンに保つ。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- `job_design.md` の第3章「3-2. ログインセッションクリーンアップ (`CLEANUP_EXPIRED_SESSIONS`)」における目的の説明文を修正しました。
- 要件定義書およびDB設計書に合わせ、「1ヶ月（43200分、Sliding Expiration により延長された最終有効期限 `EXPIRES_AT < NOW()`）が経過した `LOGIN_SESSION` レコードを削除する」仕様に統一しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
