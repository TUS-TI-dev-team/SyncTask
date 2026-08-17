# ログテーブルおよびレートリミットテーブルの定期パージジョブ定義の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:10:30
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
[database_design.md](docs/design/database_design.md) にて定期パージ（物理削除）が規定されているログテーブル（`ACCESS_LOG`, `LOGIN_LOG`, `MAIL_AUTH_LOG`）およびレートリミット管理テーブル（`LOGIN_IP_RATE_LIMIT`）のクリーンアップジョブが、[job_design.md](docs/design/job_design.md) のジョブ一覧、詳細仕様、処理フロー、環境変数に一切定義されていません。

## 2. 詳細な指摘内容
[database_design.md](docs/design/database_design.md) では以下の通り、日次Cronジョブによる定期パージ方針およびそのためのインデックスが明記されています：

1. **`LOGIN_IP_RATE_LIMIT` (lines 115-116, 215-216)**:
   - 保持期間・パージ方針: 遮断解除日時（`BLOCKED_UNTIL`）を経過し、かつ `LAST_FAILED_AT` から1日（24時間）以上経過した不要レコードを日次Cronジョブ（毎日 03:00 JST / Cron: `0 3 * * *`）にて物理削除。
   - インデックス: `idx_login_ip_rate_limit_purge (BLOCKED_UNTIL, LAST_FAILED_AT)`
2. **`ACCESS_LOG` (lines 153-154, 227-228)**:
   - 保持期間・パージ方針: 90日間保持し、日次Cronジョブ（毎日 01:00 JST / Cron: `0 1 * * *`）にて経過レコードを物理削除。
   - インデックス: `idx_access_log_purge (CREATED_AT)`
3. **`LOGIN_LOG` (lines 135-136, 224)**:
   - 保持期間・パージ方針: 1年間（365日間）保持し、日次Cronジョブ（毎日 02:00 JST / Cron: `0 2 * * *`）にて経過レコードを物理削除。
   - インデックス: `idx_login_log_purge (CREATED_AT)`
4. **`MAIL_AUTH_LOG` (lines 174-175, 231)**:
   - 保持期間・パージ方針: 1年間（365日間）保持し、日次Cronジョブ（毎日 02:00 JST / Cron: `0 2 * * *`）にて経過レコードを物理削除。
   - インデックス: `idx_mail_auth_log_purge (CREATED_AT)`

しかし、[job_design.md](docs/design/job_design.md) では `CLEANUP_OTP_SESSIONS` と `CLEANUP_EXPIRED_SESSIONS` の2つしか定義されておらず、これら4つのテーブルに対するパージジョブが完全に欠落しています。このままではログデータや一時的なレートリミットレコードが無制限に肥大化し、ストレージ圧迫およびDB性能劣化を引き起こす重大なリスクがあります。

## 3. 推奨される修正案
[job_design.md](docs/design/job_design.md) に以下の内容を追加・反映してください：

1. **第1章 ジョブ種別一覧**:
   - `CLEANUP_ACCESS_LOGS`（毎日 01:00 JST / 対象: `ACCESS_LOG`）
   - `CLEANUP_AUTH_LOGS`（または `CLEANUP_LOGIN_LOGS`, `CLEANUP_MAIL_AUTH_LOGS`）（毎日 02:00 JST / 対象: `LOGIN_LOG`, `MAIL_AUTH_LOG`）
   - `CLEANUP_RATE_LIMITS`（毎日 03:00 JST / 対象: `LOGIN_IP_RATE_LIMIT`）
2. **第2章 ジョブ詳細仕様**:
   - 各ジョブの目的、実行頻度、クリーンアップSQL（例: `DELETE FROM ACCESS_LOG WHERE CREATED_AT < NOW() - INTERVAL '90 days';` 等）、ログ記録仕様を追加。
3. **第3章 処理フロー & シーケンス図**:
   - 各パージジョブの実行フロー・シーケンスを追加。
4. **第5章 環境変数・設定値一覧**:
   - 各Cron式（`CRON_ACCESS_LOG_CLEANUP_SCHEDULE`, `CRON_AUTH_LOG_CLEANUP_SCHEDULE`, `CRON_RATE_LIMIT_CLEANUP_SCHEDULE` 等）や保持期間設定値を追加。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- `job_design.md` に欠落していた以下の3つのクリーンアップジョブを追加定義しました：
  - `CLEANUP_ACCESS_LOGS`: 毎日 01:00 JST (`0 1 * * *`) / `ACCESS_LOG` (90日経過レコード)
  - `CLEANUP_AUTH_LOGS`: 毎日 02:00 JST (`0 2 * * *`) / `LOGIN_LOG`, `MAIL_AUTH_LOG` (365日経過レコード)
  - `CLEANUP_RATE_LIMITS`: 毎日 03:00 JST (`0 3 * * *`) / `LOGIN_IP_RATE_LIMIT` (解除日時経過かつ24時間以上経過レコード)
- 各ジョブの目的、チャンク削除SQL、実行ログ記録仕様、および環境変数設定（`CRON_ACCESS_LOG_CLEANUP_SCHEDULE`, `CRON_AUTH_LOG_CLEANUP_SCHEDULE`, `CRON_RATE_LIMIT_CLEANUP_SCHEDULE`）を反映しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
