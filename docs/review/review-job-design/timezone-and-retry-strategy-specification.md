# Cronタイムゾーン仕様の曖昧さおよびリトライ対象例外・バックオフ方針の具体化不足

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:10:30
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
[job_design.md](docs/design/job_design.md) において、Cronスケジュールの基準タイムゾーン（JSTかUTCか）が明記されておらず、[database_design.md](docs/design/database_design.md) の「JST」表記との整合性が不明瞭です。また、エラーハンドリングにおける「DB接続一時失敗」の判定基準やリトライ方式（固定間隔か指数バックオフか）の定義が不十分です。

## 2. 詳細な指摘内容
1. **タイムゾーンの曖昧さ**:
   - [database_design.md](docs/design/database_design.md) では一貫して `00:00 JST`, `01:00 JST`, `02:00 JST`, `03:00 JST` と JST（日本標準時）が明記されています。
   - 一方で [job_design.md](docs/design/job_design.md) では `0 0 * * *` とのみ記載され、第4章のエラーログ例では `2026-08-08T15:00:00Z`（UTC表記）となっています。サーバーやコンテナのデフォルトタイムゾーンが UTC の場合、意図した実行時刻（深夜帯など）から9時間ずれて稼働ピーク時に重なる危険性があります。
2. **リトライ対象エラーおよびバックオフ方針の具体化不足**:
   - 第4章 4-1 で「DB接続一時失敗時: 5秒間隔で最大3回リトライを行う」と記載されていますが、対象となる例外種別（コネクションプール枯渇、一時的なネットワーク切断、デッドロック検知等）とリトライ不可な恒久的エラー（SQL構文エラー、テーブル未存在等）の切り分けがありません。
   - また、5秒の固定間隔よりも指数バックオフ（Exponential Backoff with Jitter）を採用する方が、一時的なDB高負荷時の回復において安全です。

## 3. 推奨される修正案
1. **タイムゾーンの明記**:
   - Cronスケジューラーのタイムゾーンが `Asia/Tokyo (JST)` であることを第1章および第5章の環境変数仕様に明記してください（例: `CRON_TIMEZONE=Asia/Tokyo`）。
2. **リトライ方針の具体化**:
   - リトライ対象となる一時エラー（接続タイムアウト、一時的通信エラー、デッドロック `40P01`）を明示し、恒久エラー（スキーマ不整合など）は即座に失敗としてアラート通知する方針を記載してください。
   - リトライ間隔について、固定5秒間隔に加えて指数バックオフの推奨または適用を検討・明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- 第1章および第6章（環境変数一覧）にスケジューラのタイムゾーンとして `Asia/Tokyo (JST)` および `CRON_TIMEZONE=Asia/Tokyo` を明記しました。
- 第5章「5-1. エラーハンドリング & リトライ戦略」に、リトライ対象の一時エラー（DBタイムアウト、コネクション枯渇、デッドロック `40P01` 等）とリトライ不可の恒久エラー（構文エラー、スキーマ不一致等）の分類を定義しました。
- リトライ方式として指数バックオフ＋ジッター（Exponential Backoff with Full Jitter）を採用し、最大3回のリトライ計算式を規定しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
