# 再認証失敗回数リセット処理における誤ったクエリ参照・SQL構文エラー・更新カラム不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-11-13 16:09:00
- **Target Files**:
  - [logout.go](backend/repository/logout.go)
  - [05_logout.md](docs/design/process_design/05_logout.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`backend/repository/logout.go` の再認証失敗回数リセット処理において、ログイン失敗回数リセット用のクエリ定数 `resetAccountFailureQuery` が誤って使用されているほか、定義されている `resetReauthFailedCountQuert` に定数名タイポ・SQL 構文エラー・仕様で定められた更新カラム（`REAUTH_LAST_FAILED_AT = NULL` 等）の欠落があります。

## 2. 詳細な指摘内容

1. **誤ったクエリ定数 `resetAccountFailureQuery` の使用（53〜56行目）**
   ```go
   	// reset reauth count
   	if _, err = tx.ExecContext(ctx, resetAccountFailureQuery, uid); err != nil {
   		return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
   	}
   ```
   `resetAccountFailureQuery` は同一パッケージ内の `repository/login.go` で定義されている「ログイン失敗回数 (`LOGIN_FAILED_COUNT`, `LOGIN_LAST_FAILED_AT`, `LOGIN_LOCK_UNTIL`)」のリセット用クエリです。ログアウト時にリセットすべき対象は「再認証失敗回数 (`REAUTH_FAILED_COUNT`, `REAUTH_LAST_FAILED_AT`)」であり、仕様と異なるカラムを更新しています。
   また、引数に `sql.NullString` 型の `uid` をそのまま渡すよりも、検証済みの文字列 `uid.String` を渡す方が明確です。

2. **`resetReauthFailedCountQuert` のタイポ・SQL構文エラー・カラム不足（81〜85行目）**
   ```go
   const resetReauthFailedCountQuert = `
   UPDATE LOGIN_ACCOUNT
   SET REAUTH_FAILED_COUNT = 0,
   WHERE USER_ID = $1
   `
   ```
   - 定数名が `resetReauthFailedCountQuert`（末尾が `Quert`）とタイポしており、未使用になっています。
   - `SET REAUTH_FAILED_COUNT = 0,` の末尾にカンマ `,` が残っており、`WHERE` の直前で PostgreSQL の SQL 構文エラーになります。
   - `docs/design/process_design/05_logout.md`（3.2節: 「セッションの `USER_ID` に対応する `LOGIN_ACCOUNT.REAUTH_FAILED_COUNT` を `0`、`REAUTH_LAST_FAILED_AT` を `NULL` に更新する」）および `docs/design/api_design/02_auth.md`（3.1.5節）の仕様通り、`REAUTH_LAST_FAILED_AT = NULL`（および `UPDATED_AT = CURRENT_TIMESTAMP`）の更新が必要です。

## 3. 推奨される修正案

定数名と SQL 文を修正し、`AttemptLogout` 内でその定数および `uid.String` を使用するように修正してください。

```go
	// reset reauth count
	if _, err = tx.ExecContext(ctx, resetReauthFailedCountQuery, uid.String); err != nil {
		return &model.LogoutAttemptResult{Status: model.LogoutStatusDbFailed}, err
	}
```

```go
const resetReauthFailedCountQuery = `
UPDATE LOGIN_ACCOUNT
SET REAUTH_FAILED_COUNT = 0,
    REAUTH_LAST_FAILED_AT = NULL,
    UPDATED_AT = CURRENT_TIMESTAMP
WHERE USER_ID = $1
`
```
