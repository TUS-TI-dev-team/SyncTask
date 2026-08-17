# 検索・ソート・レートリミットおよびパージ処理を担保するインデックス設計方針の追記

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 14:00:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書で求められている「ターンアラウンドタイム2秒以下」「タスク一覧の各種複合ソート（ピン留め優先、締切日時昇順、作成日時降順）」「直近5分間に30回失敗のIPレートリミット判定」「セッション・ログの日次パージ処理」について、DB設計書におけるインデックス定義の記載が限定的（`CREATED_AT` のみ言及など）であり、将来的なデータ量増加時の性能劣化リスクがあります。

## 2. 詳細な指摘内容
- `requirements.md` の120〜123行目（タスク一覧並び順）:
  - 1. ピン留め（`IS_PINNED DESC`） 2. 締切日時（`DUE_DATE ASC NULLS LAST`） 3. 作成日時（`CREATED_AT DESC`）
- `requirements.md` の195行目（IPレートリミット）:
  - 同一IPアドレスから直近5分間に累計30回以上のログイン失敗
- `database_design.md` では `TASK` テーブルや `LOGIN_SESSION`、`LOGIN_LOG` テーブルの複合インデックス構成について明記されていません。

## 3. 推奨される修正案
`database_design.md` に各テーブルの推奨インデックス設計節を追加し、以下の方針を明記することを推奨します：
1. **`TASK` テーブル**:
   - `CREATE INDEX idx_task_user_status_sort ON TASK (USER_ID, STATUS, IS_PINNED DESC, DUE_DATE ASC, CREATED_AT DESC);`
2. **`LOGIN_SESSION` テーブル**:
   - `CREATE INDEX idx_login_session_user ON LOGIN_SESSION (USER_ID);`
   - `CREATE INDEX idx_login_session_expires ON LOGIN_SESSION (EXPIRES_AT);`
3. **`LOGIN_LOG` テーブル**:
   - `CREATE INDEX idx_login_log_ip_rate ON LOGIN_LOG (IP_ADDRESS, IS_SUCCESS, CREATED_AT);`
4. **`OTP_SESSION` テーブル**:
   - `CREATE INDEX idx_otp_session_pending_email ON OTP_SESSION (PENDING_EMAIL, STATUS, EXPIRES_AT);`

---

## 修正完了報告

- **Resolved At**: 2026-08-17 14:26:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` に「7. 推奨インデックス設計 (INDEXES)」のセクションを新設しました。
- `TASK`（タスク一覧複合ソート）、`LOGIN_SESSION`（ユーザー別照会・期限切れ日次パージ）、`OTP_SESSION`（排他・有効期限照会・15分間隔パージ）、`LOGIN_IP_RATE_LIMIT`（パージ）、各種ログテーブル（`LOGIN_LOG`, `ACCESS_LOG`, `MAIL_AUTH_LOG`）に対する推奨インデックス定義とSQL構文を網羅的に記載しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
