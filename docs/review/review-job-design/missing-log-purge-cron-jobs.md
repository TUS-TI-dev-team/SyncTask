# ログパージ用定期バッチジョブの定義漏れ

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書（`requirements.md`）の「ログ保持・パージ方針」において規定されている「DBへのアクセスログ削除ジョブ（毎日 01:00 JST）」および「ログイン情報・メール認証ログ削除ジョブ（毎日 02:00 JST）」の定義が、`job_design.md`（ジョブ種別一覧、詳細仕様、シーケンス図、環境変数一覧等）から完全に欠落しています。

## 2. 詳細な指摘内容
要件定義書（[requirements.md:L239-241](docs/req-def/requirements.md#L239-L241)）では以下のように定期パージ用 Cron ジョブの仕様が明確に定義されています：
- **DBへのアクセスログ**: 90日間保持し、日次 Cron ジョブ（毎日 01:00 JST / Cron: `0 1 * * *`）にて経過レコードを物理削除する。
- **ログイン情報・メール認証ログ**: 1年間（365日間）保持し、日次 Cron ジョブ（毎日 02:00 JST / Cron: `0 2 * * *`）にて経過レコードを物理削除する。

しかし、[job_design.md:L8-16](docs/design/job_design.md#L8-L16) の「1. ジョブ種別一覧」には `CLEANUP_OTP_SESSIONS` と `CLEANUP_EXPIRED_SESSIONS` の2種類しか定義されておらず、ログ削除に関するバッチジョブ（例: `CLEANUP_ACCESS_LOGS`, `CLEANUP_AUTH_LOGS` 等）の記述が一切ありません。第2章の詳細仕様、第3章の処理フロー、第5章の環境変数一覧にも反映されていません。

## 3. 推奨される修正案
1. [job_design.md](docs/design/job_design.md) の「1. ジョブ種別一覧」に以下の2つのジョブを追加してください。
   - `CLEANUP_ACCESS_LOGS`: DBアクセスログ（90日経過レコード）の物理削除（日次 `0 1 * * *` JST）
   - `CLEANUP_AUTH_LOGS`: ログイン情報・メール認証ログ（365日経過レコード）の物理削除（日次 `0 2 * * *` JST）
2. 「2. ジョブ詳細仕様」に、それぞれの削除クエリ（`DELETE FROM ... WHERE CREATED_AT < NOW() - INTERVAL '90 days'` 等）、実行ログ記録仕様を追加してください。
3. 「3. ジョブステータス & 処理フロー」のシーケンス図および「5. 環境変数・設定値一覧」に、対応するスケジュール設定値（`CRON_ACCESS_LOG_CLEANUP_SCHEDULE`, `CRON_AUTH_LOG_CLEANUP_SCHEDULE` 等）を追加してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- 要件定義書に合わせて、`job_design.md` に `CLEANUP_ACCESS_LOGS`（毎日 01:00 JST）および `CLEANUP_AUTH_LOGS`（毎日 02:00 JST）を定義しました。
- 各パージ処理のチャンク削除SQL、ログ記録仕様、シーケンス図、および環境変数定義（`CRON_ACCESS_LOG_CLEANUP_SCHEDULE`, `CRON_AUTH_LOG_CLEANUP_SCHEDULE`）を追加・反映しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
