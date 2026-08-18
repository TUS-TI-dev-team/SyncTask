# LOGIN_IP_RATE_LIMITの日次パージ実行スケジュールにおけるDB設計書とジョブ設計書の不整合

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [job_design.md](docs/design/job_design.md)

## 1. 問題の概要
`LOGIN_IP_RATE_LIMIT` テーブルの保持期間・パージ方針において、`database_design.md` では「日次Cronジョブ（毎日 00:00 JST / `0 0 * * *`）」と記載されているのに対し、`job_design.md` では「毎日 03:00 JST / `0 3 * * *`」と定義されており、実行タイミングの記述に不整合があります。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` L115:
  ```markdown
  > **保持期間・パージ方針**: 遮断解除日時（`BLOCKED_UNTIL`）を経過し、かつ `LAST_FAILED_AT` から1日（24時間）以上経過した不要レコードは、日次Cronジョブ（毎日 00:00 JST / Cron: `0 0 * * *`）にて物理削除します。
  ```

- `docs/design/job_design.md` L19 & L260:
  ```markdown
  | `CLEANUP_RATE_LIMITS` | 解除日時経過かつ24時間以上経過したレートリミットレコードの削除 | 毎日 03:00 (`0 3 * * *`) | `LOGIN_IP_RATE_LIMIT` | チャンク分割 SQL DELETE | 冪等 / Advisory Lock |
  ```
  ```markdown
  | `CRON_RATE_LIMIT_CLEANUP_SCHEDULE` | レートリミットクリーンアップ Cron 式 | `0 3 * * *` | 毎日 03:00 JST |
  ```

### 問題点：
- 毎日 00:00 には `CLEANUP_EXPIRED_SESSIONS`（期限切れログインセッション削除）が割り当てられており、バッチ負荷分散のため `CLEANUP_RATE_LIMITS` は毎日 03:00 にスケジュールされています。
- `database_design.md` 側の記述が `毎日 00:00 JST / Cron: 0 0 * * *` のままとなっており、運用・実装時の混乱を招く原因となります。

## 3. 推奨される修正案
`docs/design/database_design.md` L115 の記述を `job_design.md` に合わせて「毎日 03:00 JST / Cron: `0 3 * * *`」に統一してください：

```markdown
> **保持期間・パージ方針**: 遮断解除日時（`BLOCKED_UNTIL`）を経過し、かつ `LAST_FAILED_AT` から1日（24時間）以上経過した不要レコードは、日次Cronジョブ（毎日 03:00 JST / Cron: `0 3 * * *`）にて物理削除します。
```
