# OTPセッションクリーンアップクエリにおける主キーカラム名の不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
  - [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)

## 1. 問題の概要
`docs/design/job_design.md` の「3-1. OTP セッションクリーンアップ (`CLEANUP_OTP_SESSIONS`)」に記載されているクリーンアップ SQL において、`OTP_SESSION` テーブルの主キーカラム名が `SESSION_ID` と誤って記述されています。
`docs/design/database_design.md` で定義されている主キーカラム名は `OTP_SESSION_ID` であり、このまま実装すると SQL 実行時エラー（`column "session_id" does not exist`）が発生し、ジョブが失敗します。

## 2. 詳細な指摘内容
- `docs/design/database_design.md`「4. OTPセッション管理 (OTP_SESSION)」では、テーブル主キーが以下のように定義されています。
  - カラム名: `OTP_SESSION_ID` (`VARCHAR(64)` / `PRIMARY KEY`)
- しかし、`docs/design/job_design.md`「3-1. OTP セッションクリーンアップ (`CLEANUP_OTP_SESSIONS`)」のクエリ（L68-78）では `SESSION_ID` が使用されています：
  ```sql
  WITH target_rows AS (
      SELECT SESSION_ID
      FROM OTP_SESSION
      WHERE MAX_EXPIRES_AT < NOW()
         OR (STATUS IN ('expired', 'locked', 'completed') AND EXPIRES_AT < NOW())
      LIMIT :batch_size
  )
  DELETE FROM OTP_SESSION
  WHERE SESSION_ID IN (SELECT SESSION_ID FROM target_rows);
  ```
- `LOGIN_SESSION` テーブルの主キー（`SESSION_ID`）と混同して記述されているため、テーブルスキーマとの不整合が発生しています。

## 3. 推奨される修正案
`docs/design/job_design.md` の該当クエリ内の `SESSION_ID` を、スキーマ定義通りの `OTP_SESSION_ID` に修正してください。

```sql
WITH target_rows AS (
    SELECT OTP_SESSION_ID
    FROM OTP_SESSION
    WHERE MAX_EXPIRES_AT < NOW()
       OR (STATUS IN ('expired', 'locked', 'completed') AND EXPIRES_AT < NOW())
    LIMIT :batch_size
)
DELETE FROM OTP_SESSION
WHERE OTP_SESSION_ID IN (SELECT OTP_SESSION_ID FROM target_rows);
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:15:30
- **Status**: Resolved

### 実施した修正内容
`docs/design/job_design.md` の「3-1. OTP セッションクリーンアップ (`CLEANUP_OTP_SESSIONS`)」におけるクリーンアップ SQL の主キーカラム指定について、`docs/design/database_design.md` のテーブル定義に従い `SESSION_ID` から `OTP_SESSION_ID` へ修正しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)

