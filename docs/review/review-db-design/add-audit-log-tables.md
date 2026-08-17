# ログ管理（ログイン情報・アクセスログ・メール認証ログ）テーブル定義の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書において詳細に規定されている「ログイン情報」「DBへのアクセスログ」「メール認証ログ」の3種類のログおよび日次Cronパージ方針に対応するテーブル設計が、データベース設計書に記載されていません。

## 2. 詳細な指摘内容
- `docs/req-def/requirements.md` の232〜241行目（ログ管理）において、以下のログデータ保持および日次Cronパージが要件として明記されています：
  - **ログイン情報**: 日時、ユーザーID (UID)、メールアドレス、IPアドレス、ログイン成否（保持期間: 1年間）
  - **DBへのアクセスログ**: 日時、ユーザーID (UID)、アクセス元IPアドレス、エンドポイント、閲覧・操作対象リソースID（保持期間: 90日間）
  - **メール認証ログ**: 日時、ユーザーID (UID)、対象メールアドレス、認証種別、アクセス元IPアドレス、処理イベント種別、成否、ダミー処理区分（保持期間: 1年間）
- しかし、`docs/design/database_design.md` には `LOGIN_ACCOUNT`, `TASK`, `LOGIN_SESSION`, `OTP_SESSION` の4テーブルしか存在せず、これらログを永続化・管理・パージするためのテーブルスキーマ定義が完全に欠落しています。

## 3. 推奨される修正案
要件定義書の「ログ管理」要件を満たすテーブル定義（例: `LOGIN_LOG`, `ACCESS_LOG`, `MAIL_AUTH_LOG` 等）および各カラム・型・インデックス・保持期間/パージ仕様をDB設計書に追加してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
要件定義書に完全準拠し、データベース設計書に以下の3つのログ管理テーブル定義および保持期間・日次Cronパージ方針を追加しました：
1. `LOGIN_LOG`: ログイン情報ログ（保持期間: 1年間、Cron: `0 2 * * *`）
2. `ACCESS_LOG`: DBアクセスログ（保持期間: 90日間、Cron: `0 1 * * *`）
3. `MAIL_AUTH_LOG`: メール認証ログ（保持期間: 1年間、Cron: `0 2 * * *`）

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
