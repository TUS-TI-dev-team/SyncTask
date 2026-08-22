# `CLEANUP_OTP_SESSIONS` の OR 条件クエリに対するインデックス効率とパージ性能の整合性

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

## 1. 問題の概要
`job_design.md` 3-1 における `CLEANUP_OTP_SESSIONS` のクリーンアップ SQL では、`WHERE SESSION_EXPIRES_AT < NOW() OR (STATUS = 'active' AND OTP_EXPIRES_AT < NOW())` という `OR` 条件で抽出を行っています。
一方、`database_design.md` 7.2 で定義されているパージ用インデックスは `CREATE INDEX idx_otp_session_purge ON OTP_SESSION (SESSION_EXPIRES_AT, STATUS, OTP_EXPIRES_AT);` の1本のみです。
PostgreSQL において、複合インデックスの先頭列が `SESSION_EXPIRES_AT` の場合、`OR` 条件の後半 `(STATUS = 'active' AND OTP_EXPIRES_AT < NOW())` に対してインデックスレンジスキャンが効率的に効かず、パージ処理時の検索負荷が高まる懸念があります。

## 2. 詳細な指摘内容
1. **インデックススキャンの挙動**:
   - `WHERE A < NOW() OR (B = 'active' AND C < NOW())` の条件に対して、インデックス `(A, B, C)` は `A` の条件には有効ですが、`OR` 結合された `(B, C)` の条件を満たす行を探すためにはインデックス全体のスキャンまたはテーブルフルスキャンが必要になります。
2. **インデックス設計の乖離**:
   - `database_design.md` 7.2 では `idx_otp_session_pending_email (PENDING_EMAIL, STATUS, OTP_EXPIRES_AT)` も定義されていますが、クリーンアップクエリに最適化されたインデックス構成になっていません。

## 3. 推奨される修正案
以下のいずれかの対応を実施してください：

- **案1（インデックス側の最適化）**:
  `database_design.md` 7.2 のインデックス定義を見直し、以下の2本に分割または調整する：
  1. `CREATE INDEX idx_otp_session_session_expires ON OTP_SESSION (SESSION_EXPIRES_AT);`
  2. `CREATE INDEX idx_otp_session_active_otp_expires ON OTP_SESSION (STATUS, OTP_EXPIRES_AT) WHERE STATUS = 'active';` (部分インデックス)
- **案2（クエリ側の分割実行）**:
  `job_design.md` において、全体期限切れパージと放置OTPパージを2つのステップ/クエリに分割して実行する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
`OTP_SESSION` のクリーンアップ OR 条件クエリに対応するため、パージ用インデックスを `idx_otp_session_session_expires (SESSION_EXPIRES_AT)` および部分インデックス `idx_otp_session_active_otp_expires (STATUS, OTP_EXPIRES_AT) WHERE STATUS = 'active'` に最適化・分割しました。

### 変更したファイル
- [database_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\database_design.md)
