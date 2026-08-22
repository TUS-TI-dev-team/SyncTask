# Database Design (データベース設計)

## 概要

システムのデータベーステーブル構造およびカラム定義です。

---

## 1. アカウント管理 (LOGIN_ACCOUNT)

**Table Name**: `LOGIN_ACCOUNT`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID |
| ユーザー名 | `USER_NAME` | `VARCHAR(20)` / `NOT NULL` | 2〜20文字、英大小数字（同名登録可） |
| メールアドレス | `EMAIL` | `VARCHAR(320)` / `UNIQUE, NOT NULL` | 認証用メール（登録・更新時に一律小文字へ正規化して保存。論理削除時は衝突回避のため退避形式へ更新） |
| パスワードハッシュ | `PASSWORD_HASH` | `VARCHAR(255)` / `NOT NULL` | ソルトを付与 + ハッシュ化保存 |
| 削除フラグ (論理削除) | `IS_DELETED` | `BOOLEAN` / `NOT NULL, DEFAULT FALSE` | アカウント削除時は論理削除 |
| アカウント削除日時 | `DELETED_AT` | `TIMESTAMPTZ` | 削除処理タイムスタンプ |
| ログイン失敗回数 | `LOGIN_FAILED_COUNT` | `INT` / `NOT NULL, DEFAULT 0` | 15分間のインターバルを挟まずに5回連続失敗で30分間ロック / 最後の失敗から15分経過またはログイン成功時に0にリセット |
| ログイン最終失敗日時 | `LOGIN_LAST_FAILED_AT` | `TIMESTAMPTZ` | 最終失敗タイムスタンプ |
| ロック解除日時 | `LOGIN_LOCK_UNTIL` | `TIMESTAMPTZ` | 5回連続失敗時にロックアウト発生時刻から30分後（NOW() + INTERVAL '30 minutes'）を設定。ロック中の追加試行では延長しない（固定30分） |
| 再認証失敗回数 | `REAUTH_FAILED_COUNT` | `INT` / `NOT NULL, DEFAULT 0` | パスワード変更・アカウント削除時の再認証失敗回数。5回連続失敗でログインセッション物理削除（成功時、パスワード変更時、ログアウト時、5回失敗セッション破棄時に0リセット） |
| 再認証最終失敗日時 | `REAUTH_LAST_FAILED_AT` | `TIMESTAMPTZ` | 再認証最終失敗タイムスタンプ |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | |
| 更新日時 | `UPDATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | |

> [!NOTE]
> **削除方針およびメールアドレス重複回避・ログイン試行制御**
> - **アカウント (`LOGIN_ACCOUNT`)**: 退会・アカウント削除時は論理削除 (`IS_DELETED = TRUE`, `DELETED_AT = NOW()`) を行います。
>   - 論理削除後の同メールアドレスでの再登録を可能とするため、論理削除実行時に `EMAIL` カラムの値を退避フォーマット（例: `deleted_<USER_ID>_<EMAIL>`）に更新し、有効なアカウント間でのみ一意性を維持します（最大300文字超に対応するため `VARCHAR(320)` を確保）。
>   - 論理削除されたアカウントの元メールアドレスでログイン試行された場合、`EMAIL` 検索では「該当アカウントなし（未登録）」と同一パスを通るため、論理削除の有無は外部に秘匿されます。未登録アドレスや退避済みアドレスへの連続試行に対しては、IPアドレス単位レートリミット（`LOGIN_IP_RATE_LIMIT`）によりブルートフォース攻撃から保護します。
> - **タスク (`TASK`)**: アカウント論理削除に伴い、所有するタスクデータおよび関連データは即座にDBから物理削除されます。
> - **セッション (`LOGIN_SESSION`, `OTP_SESSION`)**: ログアウト・アカウント削除時および期限切れ時は物理削除 (`DELETE`) されます。

---

## 2. タスク管理 (TASK)

**Table Name**: `TASK`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| タスクID | `TASK_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` | 所有ユーザー |
| タスク名 | `TITLE` | `VARCHAR(255)` / `NOT NULL` | 1〜100文字（制御文字不可） |
| 優先度 | `PRIORITY` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'medium'` | `'low'` (低), `'medium'` (中・初期値), `'high'` (高) |
| 締切日時 | `DUE_DATETIME` | `TIMESTAMPTZ` | 任意設定（締切未設定時は NULL。日付のみ指定時は該当日 23:59:00+09:00 を設定して保存） |
| タスクステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'not_started'` | `'not_started'` (未着手・初期値), `'in_progress'` (進行中), `'completed'` (完了) |
| ピン止めフラグ | `IS_PINNED` | `BOOLEAN` / `NOT NULL, DEFAULT FALSE` | |
| コメント | `COMMENT` | `TEXT` / `NOT NULL, DEFAULT ''` | 補足メモ（0〜1000文字。未入力時は空文字 `''` で保存） |
| 検索用テキスト | `SEARCH_TEXT` | `TEXT` / `NOT NULL` | 日本語同一視・部分一致検索用正規化文字列（小文字化＋NFKC正規化＋ひらがな→カタカナ変換） |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | |
| 更新日時 | `UPDATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | |

> [!NOTE]
> **定期タスク（繰り返しタスク）の永続化方針**
> 繰り返しタスクはルールマスターテーブル（`RECURRING_TASKS` 等）を持たず、作成API（`POST /api/v1/tasks`）実行時に最大100件の独立した通常タスクレコードとして `TASK` テーブルへ同期即時一括生成・永続化されます。そのため、定期バッチ生成ジョブは不要（スコープ外）となります。

---

## 3. ログインセッション管理 (LOGIN_SESSION)

**Table Name**: `LOGIN_SESSION`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| セッションID | `SESSION_ID` | `VARCHAR(64)` / `PRIMARY KEY` | ランダムトークン / Cookie保存 (`sync_task_sid`) |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `NOT NULL, FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` | ログインユーザーID |
| 有効期限 | `EXPIRES_AT` | `TIMESTAMPTZ` / `NOT NULL` | 30日 (43200分) / APIアクセス時にSliding Expirationで自動延長 / 期限切れは日次Cron（00:00 JST）で物理削除 |
| User-Agent | `USER_AGENT` | `TEXT` | クライアント情報 |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | |

---

## 4. OTPセッション管理 (OTP_SESSION)

**Table Name**: `OTP_SESSION`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| OTPセッションID | `OTP_SESSION_ID` | `VARCHAR(64)` / `PRIMARY KEY` | Cookie/パラメータ保存 |
| 認証種別 | `PURPOSE` | `VARCHAR(20)` / `NOT NULL` | `SIGNUP` (新規登録), `PASSWORD_RESET` (パスワードリセット), `EMAIL_CHANGE` (メールアドレス変更) |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` | 既存ユーザー識別用（パスワードリセット・メール変更時。新規登録時はNULL） |
| 登録予定ユーザー名 | `PENDING_USERNAME` | `VARCHAR(20)` | アカウント作成時は登録予定値 |
| 認証対象/変更予定メールアドレス | `PENDING_EMAIL` | `VARCHAR(255)` | 実セッションでは認証対象 / 変更予定メールアドレスを小文字へ正規化して保存し、重複排除・排他制御に利用。ダミーセッションでは一意制約競合防止および情報秘匿のためNULLとする（表示には `MASKED_EMAIL` を使用） |
| マスク済みメールアドレス | `MASKED_EMAIL` | `VARCHAR(320)` / `NOT NULL` | API応答表示用。先頭4文字（ローカル部が4文字未満の場合は先頭1文字）とドメイン以外を固定10文字（`*`）でマスクした値のみを保存し、ダミーセッションの再送応答にも利用 |
| 登録予定パスワードハッシュ | `PENDING_PASSWORD_HASH` | `VARCHAR(255)` | メール変更時・パスワードリセット時はNULL<br>ソルト + ハッシュ化 |
| OTPハッシュ | `OTP_HASH` | `VARCHAR(255)` | 実セッションでは8桁英数字のソルト付きハッシュ。ダミーセッションではNULLとし、`IS_DUMMY` を先に判定して照合成功させない |
| ステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'active'` | `active`, `verified`（完了時・失効時は即時物理削除されるため、DB上に `completed` / `expired` として残存しない） |
| ダミー処理区分 | `IS_DUMMY` | `BOOLEAN` / `NOT NULL, DEFAULT FALSE` | ダミーセッション判定の正本。`TRUE` のセッションではOTP生成・送信・照合成功・確定処理を禁止 |
| 試行失敗回数 | `ATTEMPT_COUNT` | `INT` / `NOT NULL, DEFAULT 0` | 1つのOTPに対して最大5回（5回失敗時は自動再送・失効制御） |
| 再送回数 | `SEND_COUNT` | `INT` / `NOT NULL, DEFAULT 0` | 手動/自動再送回数 |
| 連続送信失敗回数 | `SEND_FAILED_COUNT` | `INT` / `NOT NULL, DEFAULT 0` | 初回送信・手動再送・自動再送を含む連続失敗回数。成功時に0へリセットし、5回到達時は対象セッションを物理削除 |
| 配信状態 | `DELIVERY_STATUS` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'pending'` | `pending`, `sent`, `sendable`。送信失敗時は再試行可能な `sendable` とする。ダミーは実送信せず外部挙動のみ再現 |
| 直前送信日時 | `LAST_SENT_AT` | `TIMESTAMPTZ` / `NOT NULL` | 直前の送信タイムスタンプ（60秒クールダウン判定用） |
| OTP有効期限 | `OTP_EXPIRES_AT` | `TIMESTAMPTZ` / `NOT NULL` | 発行から5分 |
| セッション全体有効期限 | `SESSION_EXPIRES_AT` | `TIMESTAMPTZ` / `NOT NULL` | 初回発行から15分間（パスワードリセットの検証成功時はその時点から15分間に延長。手続き全体の排他維持・失効上限） |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | 初回発行日時 |

> [!NOTE]
> **OTPセッションのパージ方針**
> - 新パスワード設定完了時やアカウント作成確定時等に直ちにDBから物理削除されます。
> - 有効期限切れ（全体最大有効期限 `SESSION_EXPIRES_AT < NOW()`、またはステータスが `active` かつ `OTP_EXPIRES_AT < NOW()` でアクセスされず放置されたレコード）は、Cronジョブ（15分ごと / Cron: `*/15 * * * *` JST）にてDBから一括物理削除されます。
> - ダミーセッションも `SESSION_EXPIRES_AT` 経過後に同じCronで削除します。メール送信が5回連続で失敗した実セッションはCronを待たず、補償処理で直ちに物理削除します。
>
> **ダミーセッションの制約**
> - `OTP_SESSION_ID`、`PURPOSE`、`STATUS`、`IS_DUMMY`、試行・送信回数、配信状態、各期限、作成日時は通常どおり保存します。
> - `PENDING_USERNAME`、`PENDING_EMAIL`、`PENDING_PASSWORD_HASH`、`OTP_HASH` はNULLとします。`EMAIL_CHANGE` のみ所有者認可のため `USER_ID` に認証中ユーザーIDを保存し、`SIGNUP` と `PASSWORD_RESET` はNULLとします。
> - `IS_DUMMY=TRUE` を最終判定とし、他のNULL値からダミーかどうかを推測してはなりません。

---

## 5. ログインレートリミット管理 (LOGIN_IP_RATE_LIMIT)

**Table Name**: `LOGIN_IP_RATE_LIMIT`

パスワードスプレー攻撃等の防止を目的とした、IPアドレス単位のログイン失敗追跡およびアクセス一時遮断用テーブルです。

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| IPアドレス | `IP_ADDRESS` | `VARCHAR(45)` / `PRIMARY KEY` | 対象クライアントIPアドレス (IPv4 / IPv6) |
| 失敗回数 | `FAILED_COUNT` | `INT` / `NOT NULL, DEFAULT 0` | 5分間のインターバルを挟まない累計失敗回数（リクエスト時に `LAST_FAILED_AT` から5分超過で0リセット） |
| 最終失敗日時 | `LAST_FAILED_AT` | `TIMESTAMPTZ` / `NOT NULL` | 最終失敗タイムスタンプ |
| 遮断解除日時 | `BLOCKED_UNTIL` | `TIMESTAMPTZ` | 30回到達時に `NOW() + INTERVAL '15 minutes'` を設定。この時刻まで該当IPからのログインを一律遮断 (HTTP 429) |
| 更新日時 | `UPDATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | レコード更新日時 |

> [!NOTE]
> **保持期間・パージ方針**: 遮断解除日時（`BLOCKED_UNTIL`）を経過し、かつ `LAST_FAILED_AT` から1日（24時間）以上経過した不要レコードは、日次Cronジョブ（毎日 03:00 JST / Cron: `0 3 * * *`）にて物理削除します。

---

## 6. ログ管理 (LOGS)

### 6.1 ログイン情報ログ (LOGIN_LOG)

**Table Name**: `LOGIN_LOG`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| ログID | `LOG_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID |
| ユーザーID | `USER_ID` | `VARCHAR(36)` | 認証成功時または特定可能な場合のUID（未特定時はNULL） |
| メールアドレス | `EMAIL` | `VARCHAR(320)` / `NOT NULL` | ログイン試行対象メールアドレス |
| IPアドレス | `IP_ADDRESS` | `VARCHAR(45)` / `NOT NULL` | アクセス元IPアドレス |
| ログイン成否 | `IS_SUCCESS` | `BOOLEAN` / `NOT NULL` | `TRUE` (成功) / `FALSE` (失敗) |
| セッションの利用 | `IS_SESSION_USED` | `BOOLEAN` / `NOT NULL` | `TRUE` (利用) / `FALSE` (未利用) |
| アクセス日時 | `ACCESS_AT` | `TIMESTAMPTZ` / `NOT NULL` | ログイン試行日時（インデックス対象） |

> [!NOTE]
> **保持期間・パージ方針**: 1年間（365日間）保持し、日次Cronジョブ（毎日 02:00 JST / Cron: `0 2 * * *`）にて経過レコードを物理削除します。

---

### 6.2 APIアクセスログ (ACCESS_LOG)

**Table Name**: `ACCESS_LOG`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| ログID | `LOG_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID |
| ユーザーID | `USER_ID` | `VARCHAR(36)` | 操作元ユーザーID（未ログイン時はNULL） |
| アクセス元IPアドレス | `IP_ADDRESS` | `VARCHAR(45)` / `NOT NULL` | クライアントIPアドレス |
| エンドポイント | `ENDPOINT` | `VARCHAR(1000)` / `NOT NULL` | アクセス先APIパス/メソッド |
| 操作対象リソースID | `RESOURCE_ID` | `VARCHAR(255)` | 閲覧・操作対象リソースID（タスクID等、該当なし時はNULL） |
| アクセス日時 | `ACCESS_AT` | `TIMESTAMPTZ` / `NOT NULL` | アクセス日時（インデックス対象） |

> [!NOTE]
> **保持期間・パージ方針**: 90日間保持し、日次Cronジョブ（毎日 01:00 JST / Cron: `0 1 * * *`）にて経過レコードを物理削除します。

---

### 6.3 メール認証ログ (MAIL_AUTH_LOG)

**Table Name**: `MAIL_AUTH_LOG`

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| ログID | `LOG_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID |
| ユーザーID | `USER_ID` | `VARCHAR(36)` | 対象ユーザーID（未登録・未ログイン時はNULL） |
| 対象メールアドレス | `EMAIL` | `VARCHAR(320)` / `NOT NULL` | 認証対象メールアドレス |
| 認証種別 | `AUTH_TYPE` | `VARCHAR(20)` / `NOT NULL` | `SIGNUP` (新規登録), `PASSWORD_RESET` (パスワードリセット), `EMAIL_CHANGE` (メールアドレス変更) |
| アクセス元IPアドレス | `IP_ADDRESS` | `VARCHAR(45)` / `NOT NULL` | クライアントIPアドレス |
| 処理イベント種別 | `EVENT_TYPE` | `VARCHAR(30)` / `NOT NULL` | `ISSUED` (発行), `VERIFY_SUCCESS` (検証成功), `VERIFY_FAILED` (検証失敗), `RESEND_REQUESTED` (手動再送), `AUTO_RESEND` (5回失敗時自動処理), `EXPIRED` (有効期限切れ), `CANCELLED` (ユーザーキャンセル) |
| 成否 | `IS_SUCCESS` | `BOOLEAN` / `NOT NULL` | 処理の成否 (`TRUE` / `FALSE`) |
| ダミー処理区分 | `IS_DUMMY` | `BOOLEAN` / `NOT NULL` | `TRUE` (ダミー表示) / `FALSE` (実処理) |
| アクセス日時 | `ACCESS_AT` | `TIMESTAMPTZ` / `NOT NULL` | ログ記録日時（インデックス対象） |

> [!NOTE]
> **保持期間・パージ方針**: 1年間（365日間）保持し、日次Cronジョブ（毎日 02:00 JST / Cron: `0 2 * * *`）にて経過レコードを物理削除します。

---

## 7. 推奨インデックス設計 (INDEXES)

要件定義書の性能要件（ターンアラウンドタイム2秒以下）やバッチ・パージ処理を担保するための推奨インデックス定義です。

### 7.1 タスク管理 (`TASK`)
```sql
-- 日本語同一視・部分一致検索用 pg_trgm 拡張および GIN インデックス
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_task_search_text ON TASK USING gin (SEARCH_TEXT gin_trgm_ops);

-- タスク一覧の複合ソート高速化 (ユーザー別、ステータス指定時: ピン留め降順、締切日時昇順 NULLS LAST、作成日時降順)
CREATE INDEX idx_task_user_status_sort ON TASK (USER_ID, STATUS, IS_PINNED DESC, DUE_DATETIME ASC NULLS LAST, CREATED_AT DESC);

-- タスク一覧の複合ソート高速化 (ユーザー別、全ステータス取得時: ピン留め降順、締切日時昇順 NULLS LAST、作成日時降順)
CREATE INDEX idx_task_user_sort ON TASK (USER_ID, IS_PINNED DESC, DUE_DATETIME ASC NULLS LAST, CREATED_AT DESC);

-- カレンダー表示および締切日範囲検索用 (ユーザー別、締切日時)
CREATE INDEX idx_task_user_due_datetime ON TASK (USER_ID, DUE_DATETIME);

-- タスク優先度およびステータスのCHECK制約
ALTER TABLE TASK ADD CONSTRAINT chk_task_priority CHECK (PRIORITY IN ('low', 'medium', 'high'));
ALTER TABLE TASK ADD CONSTRAINT chk_task_status CHECK (STATUS IN ('not_started', 'in_progress', 'completed'));
```

### 7.2 セッション管理 (`LOGIN_SESSION`, `OTP_SESSION`)
```sql
-- ユーザーIDでのセッション照会・失効処理
CREATE INDEX idx_login_session_user ON LOGIN_SESSION (USER_ID);

-- 有効期限切れセッションの日次Cronパージ用
CREATE INDEX idx_login_session_expires ON LOGIN_SESSION (EXPIRES_AT);

-- 同一メールアドレスに対する有効なOTPセッションの重複発行防止（排他制御用部分一意インデックス）
CREATE UNIQUE INDEX uq_otp_session_active_pending_email ON OTP_SESSION (PENDING_EMAIL) WHERE STATUS IN ('active', 'verified');

-- OTPセッションのメールアドレス照会および状態遷移用
CREATE INDEX idx_otp_session_pending_email ON OTP_SESSION (PENDING_EMAIL, STATUS, OTP_EXPIRES_AT);

-- OTPセッションの15分間隔Cronパージ用（全体最大有効期限または失効レコードのクリーンアップ）
CREATE INDEX idx_otp_session_session_expires ON OTP_SESSION (SESSION_EXPIRES_AT);
CREATE INDEX idx_otp_session_active_otp_expires ON OTP_SESSION (STATUS, OTP_EXPIRES_AT) WHERE STATUS = 'active';
```

ダミー対応のため `PENDING_EMAIL` と `OTP_HASH` の無条件NOT NULL制約は外しますが、実セッションまでNULLを許可しないよう次の条件付きCHECK制約を必須とします。

```sql
ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_purpose
CHECK (PURPOSE IN ('SIGNUP', 'PASSWORD_RESET', 'EMAIL_CHANGE'));

ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_status
CHECK (STATUS IN ('active', 'verified'));

ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_real_required_values
CHECK (IS_DUMMY OR (PENDING_EMAIL IS NOT NULL AND OTP_HASH IS NOT NULL));

ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_signup_required_values
CHECK (PURPOSE <> 'SIGNUP' OR IS_DUMMY OR (PENDING_USERNAME IS NOT NULL AND PENDING_PASSWORD_HASH IS NOT NULL));

ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_email_change_owner
CHECK (PURPOSE <> 'EMAIL_CHANGE' OR USER_ID IS NOT NULL);

ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_delivery_status
CHECK (DELIVERY_STATUS IN ('pending', 'sent', 'sendable'));

ALTER TABLE OTP_SESSION ADD CONSTRAINT chk_otp_session_send_failed_count
CHECK (SEND_FAILED_COUNT BETWEEN 0 AND 5);
```

### 7.3 レートリミット管理 (`LOGIN_IP_RATE_LIMIT`)
```sql
-- レートリミットレコードのパージ用
CREATE INDEX idx_login_ip_rate_limit_purge ON LOGIN_IP_RATE_LIMIT (BLOCKED_UNTIL, LAST_FAILED_AT);
```

### 7.4 ログテーブル (`LOGIN_LOG`, `ACCESS_LOG`, `MAIL_AUTH_LOG`)
```sql
-- ログインログのIP/メール別照会および日次パージ
CREATE INDEX idx_login_log_ip ON LOGIN_LOG (IP_ADDRESS, ACCESS_AT DESC);
CREATE INDEX idx_login_log_email ON LOGIN_LOG (EMAIL, ACCESS_AT DESC);
CREATE INDEX idx_login_log_purge ON LOGIN_LOG (ACCESS_AT);

-- アクセスログのパージ用
CREATE INDEX idx_access_log_purge ON ACCESS_LOG (ACCESS_AT);

-- メール認証ログのメール別照会およびパージ用
CREATE INDEX idx_mail_auth_log_email ON MAIL_AUTH_LOG (EMAIL, ACCESS_AT DESC);
CREATE INDEX idx_mail_auth_log_purge ON MAIL_AUTH_LOG (ACCESS_AT);
```