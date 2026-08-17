# OTP_SESSIONテーブルにおけるPENDING_EMAILのカラム長不整合および小文字正規化仕様の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 15:20:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`OTP_SESSION` テーブルの `PENDING_EMAIL` カラムのデータ型が `VARCHAR(255)` と定義されていますが、`LOGIN_ACCOUNT`、`LOGIN_LOG`、`MAIL_AUTH_LOG` テーブルの `EMAIL` カラム（`VARCHAR(320)`）と長さが不整合になっています。また、要件定義書で規定されている「システム内部での小文字正規化保持」についての注記が欠落しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` において、他テーブルのメールアドレスカラムは以下のように定義されています：
  - `LOGIN_ACCOUNT.EMAIL`: `VARCHAR(320)`（論理削除時の退避プレフィックス対応およびRFC 5321準拠）
  - `LOGIN_LOG.EMAIL`: `VARCHAR(320)`
  - `MAIL_AUTH_LOG.EMAIL`: `VARCHAR(320)`
- しかし、`OTP_SESSION.PENDING_EMAIL` のみ `VARCHAR(255)` となっており、最大長や型の一貫性が損なわれています。長大なメールアドレスでの新規登録やパスワードリセット要求時に桁あふれ等のリスクがあります。
- さらに、`docs/req-def/requirements.md` の259行目において「メールアドレスはシステム内部（登録・ログイン・重複判定・認証処理等）で一律小文字（`toLowerCase()`）に正規化して保持・比較し、Case-Insensitive な一意性を保証する」と規定されていますが、`LOGIN_ACCOUNT.EMAIL` には小文字正規化の記載があるものの、`OTP_SESSION.PENDING_EMAIL` には小文字正規化に関する備考が明記されていません。

## 3. 推奨される修正案
1. `OTP_SESSION` テーブルの `PENDING_EMAIL` カラムのデータ型を `VARCHAR(320)` に変更してください。
2. `PENDING_EMAIL` の備考欄に「登録・更新・認証要求時に一律小文字へ正規化して保存」する旨を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:22:00
- **Status**: Resolved

### 実施した修正内容
1. `docs/design/database_design.md` の第4章「OTPセッション管理 (OTP_SESSION)」において、`PENDING_EMAIL` カラムのデータ型を `VARCHAR(320)` / `NOT NULL` に変更し、他テーブル（`LOGIN_ACCOUNT`、`LOGIN_LOG` 等）とのカラム長整合性を確保しました。
2. `PENDING_EMAIL` の備考欄に「登録・更新・認証要求時に一律小文字 `toLowerCase()` へ正規化して保存」する仕様を明記しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
