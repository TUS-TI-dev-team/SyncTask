# Database Design (データベース設計)

## 概要

システムのデータベーステーブル構造およびカラム定義です。

---

## 1. アカウント管理 (LOGIN_ACCOUNT)

**Table Name**: `LOGIN_ACCOUNT`

| 項目名 | カラム名 | データ型 / 制約 | 備考 | 実装 |
| --- | --- | --- | --- | --- |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID | ✅ |
| ユーザー名 | `USERNAME` | `VARCHAR(20)` / `UNIQUE, NOT NULL` | 2〜20文字、英大小数字 | ✅ |
| メールアドレス | `EMAIL` | `VARCHAR(255)` / `UNIQUE, NOT NULL` | 認証用メール | ✅ |
| パスワードハッシュ | `PASSWORD_HASH` | `VARCHAR(255)` / `NOT NULL` | ハッシュ化保存 | ✅ |
| 削除フラグ (論理削除) | `IS_DELETED` | `BOOLEAN` / `NOT NULL, DEFAULT FALSE` | アカウント削除時は論理削除 | ✅ |
| アカウント削除日時 | `DELETED_AT` | `TIMESTAMPTZ` | 削除処理タイムスタンプ | ✅ |
| ログイン失敗回数 | `LOGIN_FAILED_COUNT` | `INT` / `DEFAULT 0` | 5分ごとにリセット (判定時動的リセット) | ✅ |
| ログイン最終失敗日時 | `LOGIN_LAST_FAILED_AT` | `TIMESTAMPTZ` | 最終失敗タイムスタンプ | ✅ |
| ロック解除日時 | `LOGIN_LOCK_UNTIL` | `TIMESTAMPTZ` | 5回失敗で30分間ロック | ✅ |
| パスワード変更失敗回数 | `CHPASS_FAILED_COUNT` | `INT` / `DEFAULT 0` | 5分ごとにリセット (判定時動的リセット) | ✅ |
| パスワード変更最終失敗日時 | `CHPASS_LAST_FAILED_AT` | `TIMESTAMPTZ` | 最終失敗タイムスタンプ | ✅ |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | | ✅ |
| 更新日時 | `UPDATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | | ✅ |

> [!NOTE]
> **削除方針のまとめ**
> - **アカウント (`LOGIN_ACCOUNT`)**: 退会・アカウント削除時は論理削除 (`IS_DELETED = TRUE`, `DELETED_AT = NOW()`) を行います。
> - **セッション (`LOGIN_SESSION`, `OTP_SESSION`)**: ログアウト・アカウント削除時および期限切れ時は物理削除 (`DELETE`) されます。

---

## 2. タスク管理 (TASK)

**Table Name**: `TASK`

| 項目名 | カラム名 | データ型 / 制約 | 備考 | 実装 |
| --- | --- | --- | --- | --- |
| タスクID | `TASK_ID` | `VARCHAR(36)` / `PRIMARY KEY` | UUID | ✅ |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` | 所有ユーザー | ✅ |
| タスク名 | `TITLE` | `VARCHAR(255)` / `NOT NULL` | | ✅ |
| 優先度 | `PRIORITY` | `VARCHAR(20)` / `NOT NULL` | LOW, MEDIUM, HIGH など | ✅ |
| 締切日時 | `DUE_DATE` | `TIMESTAMPTZ` | 72時間前判定等 | ✅ |
| タスクステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL` | IN_PROGRESS, COMPLETED | ✅ |
| ピン止めフラグ | `IS_PINNED` | `BOOLEAN` / `DEFAULT FALSE` | | ✅ |
| コメント | `COMMENT` | `TEXT` | 補足メモ | ✅ |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | | ✅ |
| 更新日時 | `UPDATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | | ✅ |

---

## 3. ログインセッション管理 (LOGIN_SESSION)

**Table Name**: `LOGIN_SESSION`

| 項目名 | カラム名 | データ型 / 制約 | 備考 | 実装 |
| --- | --- | --- | --- | --- |
| セッションID | `SESSION_ID` | `VARCHAR(64)` / `PRIMARY KEY` | ランダムトークン / Cookie保存 | ✅ |
| ユーザーID | `USER_ID` | `VARCHAR(36)` / `FOREIGN KEY (LOGIN_ACCOUNT.USER_ID)` | | ✅ |
| 有効期限 | `EXPIRES_AT` | `TIMESTAMPTZ` / `NOT NULL` | 1週間(10080分) / CRONで日次削除 | ✅ |
| IPアドレス | `IP_ADDRESS` | `VARCHAR(45)` | アクセスログ用 | ✅ |
| User-Agent | `USER_AGENT` | `TEXT` | クライアント情報 | ✅ |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | | ✅ |

---

## 4. OTPセッション管理 (OTP_SESSION)

**Table Name**: `OTP_SESSION`

| 項目名 | カラム名 | データ型 / 制約 | 備考 | 実装 |
| --- | --- | --- | --- | --- |
| OTPセッションID | `OTP_SESSION_ID` | `VARCHAR(64)` / `PRIMARY KEY` | Cookie/パラメータ保存 | ✅ |
| 登録予定ユーザー名 | `PENDING_USERNAME` | `VARCHAR(20)` | アカウント作成時は登録予定値 | ✅ |
| 登録予定メールアドレス | `PENDING_EMAIL` | `VARCHAR(255)` | メール変更時 / 新規作成時 | ✅ |
| 登録予定パスワードハッシュ | `PENDING_PASSWORD_HASH` | `VARCHAR(255)` | メール変更時はNULL | ✅ |
| OTPハッシュ | `OTP_HASH` | `VARCHAR(255)` / `NOT NULL` | 8桁英数字のハッシュ | ✅ |
| ステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL` | `active`, `verified`, `expired`, `locked`, `completed` | ✅ |
| 試行失敗回数 | `ATTEMPT_COUNT` | `INT` / `DEFAULT 0` | 最大5回でロック | ✅ |
| 再送回数 | `SEND_COUNT` | `INT` / `DEFAULT 0` | | ✅ |
| 有効期限 | `EXPIRES_AT` | `TIMESTAMPTZ` / `NOT NULL` | 5分 (検証成功後も5分間有効) | ✅ |
| 作成日時 | `CREATED_AT` | `TIMESTAMPTZ` / `NOT NULL` | | ✅ |

---

## 参考リンク
- [見習いエンジニアがログイン機能のデータベースを構築してみた｜k](https://note.com/glossy_lemur2953/n/n92e01e11481c)
- [MENTA ログインDB構築に関する解説](https://menta.work/post/detail/3024/PnXbtPMBnWVHUPAoMqby)
