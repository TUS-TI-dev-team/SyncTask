# レビュー結果サマリ

- **Status**: Resolved (指摘なし・全項目クリア)
- **Reviewed At**: 2026-08-17
- **Target Sections**: 
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
    - 3. ログインセッション管理 (`LOGIN_SESSION`)
    - 4. OTPセッション管理 (`OTP_SESSION`)
    - 7.2 セッション管理推奨インデックス (`INDEXES`)

## 査読結果詳細

[requirements.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements.md) の各要件項目と照合し、必要なカラム項目・データ型・制約・デフォルト値・備考・推奨インデックスの精査を行いました。

1. **ログインセッション管理 (`LOGIN_SESSION`)**:
   - **セッション管理・Cookie保存**: `SESSION_ID` (`VARCHAR(64)` / `PRIMARY KEY`, Cookie `sync_task_sid` 保存)
   - **ユーザー紐付け**: `USER_ID` (`VARCHAR(36)` / `NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)`)
   - **有効期限・Sliding Expiration**: `EXPIRES_AT` (`TIMESTAMPTZ` / `NOT NULL`, 1ヶ月/43200分、APIアクセス時自動延長、日次00:00 JST Cronで物理削除)
   - **端末ログアウト / 複数セッション一括無効化**: PK検索による単一セッション削除、および `USER_ID` インデックス（`idx_login_session_user`）による全セッション一括削除に対応。
   - **アクセス元情報・作成日時**: `IP_ADDRESS` (`VARCHAR(45)`), `USER_AGENT` (`TEXT`), `CREATED_AT` (`TIMESTAMPTZ` / `NOT NULL`) を網羅。

2. **OTPセッション管理 (`OTP_SESSION`)**:
   - **認証種別 (PURPOSE)**: `SIGNUP`, `PASSWORD_RESET`, `EMAIL_CHANGE` の各目的に対応（`VARCHAR(20)` / `NOT NULL`）。
   - **ユーザー紐付け**: `USER_ID` (`VARCHAR(36)` / `FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)`, 新規登録時はNULL、パスワードリセット・メール変更時は対象UID）。
   - **登録・変更予定データ**: `PENDING_USERNAME` (`VARCHAR(20)`), `PENDING_EMAIL` (`VARCHAR(320)` / `NOT NULL`, 一律小文字正規化), `PENDING_PASSWORD_HASH` (`VARCHAR(255)`)。
   - **OTPハッシュ**: `OTP_HASH` (`VARCHAR(255)` / `NOT NULL`, 8桁英数字・Case-Insensitive対応ハッシュ)。
   - **ステータス管理**: `STATUS` (`VARCHAR(20)` / `NOT NULL`: `active`, `verified`, `expired`, `locked`, `completed`)。
   - **有効期限・ライフサイクル**: 単発OTP用 `EXPIRES_AT` (発行から5分、リセット検証成功時は15分に延長) および 手続き全体上限 `MAX_EXPIRES_AT` (初回発行から15分) で完全管理。
   - **試行失敗回数 & 手動/自動再送**: `ATTEMPT_COUNT` (`INT` / `NOT NULL, DEFAULT 0`, 最大5回), `SEND_COUNT` (`INT` / `NOT NULL, DEFAULT 0`), `LAST_SENT_AT` (`TIMESTAMPTZ` / `NOT NULL`, 60秒クールダウン判定)。
   - **排他制御**: `PENDING_EMAIL` および 部分一意インデックス `uq_otp_session_active_pending_email` (`WHERE STATUS IN ('active', 'verified')`) によりDBレベルで排他担保。
   - **ダミーOTP処理**: ダミー処理時はDBセッションを発行せずログ（`MAIL_AUTH_LOG`）のみで追跡し秘匿。
   - **パージ方針**: 確定時即時物理削除、および15分周期Cron（`*/15 * * * *` JST）による一括物理削除が要件と完全整合。

