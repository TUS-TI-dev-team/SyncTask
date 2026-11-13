# selectLoginUserID における Scan 引数のポインタ渡し漏れと UUID 判定条件の反転

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-11-13 16:09:00
- **Target Files**:
  - [logout.go](backend/repository/logout.go)

## 1. 問題の概要
`backend/repository/logout.go` において、`LOGIN_SESSION` から `USER_ID` を取得する `selectLoginUserID` 関数の `Scan` 引数がポインタになっておらず実行時エラーになる問題、および取得した `uid` の UUID パース判定条件が反転しており正常なセッションでも常に `LogoutStatusNotFound` になってしまう問題があります。

## 2. 詳細な指摘内容

1. **`row.Scan` にポインタではなく値を渡している（65〜71行目）**
   ```go
   func selectLoginUserID(ctx context.Context, tx *sql.Tx, sid string) (uid sql.NullString, err error) {
   	err = tx.QueryRowContext(ctx, selectTargetUserIDQuery, sid).Scan(uid)
   ```
   `database/sql` の `Scan` にはポインタ（`&uid`）を渡す必要があります。値のまま `Scan(uid)` を呼び出すと、レコードが存在した場合に常に `sql: Scan error on column index 0, name "USER_ID": destination not a pointer` エラーとなり、37〜39行目で `LogoutStatusDbFailed` として失敗します。

2. **UUID パース結果の条件判定が逆転している（41〜46行目）**
   ```go
   	if _, parseErr := uuid.Parse(uid.String); !uid.Valid || parseErr == nil {
   		if err = tx.Commit(); err != nil {
   			return nil, err
   		}
   		return &model.LogoutAttemptResult{ Status: model.LogoutStatusNotFound }, nil
   	}
   ```
   `parseErr == nil` は「UUID として**正常にパースできた場合**」に `true` となります。そのため、有効な `USER_ID` が取得できた正常系において必ずこの分岐に入り、セッション削除や再認証失敗回数リセットが行われないまま `LogoutStatusNotFound` が返却されてしまいます。

## 3. 推奨される修正案

- `selectLoginUserID` の `Scan` 引数を `&uid` に修正する。
- 41行目の条件を `parseErr != nil`（または `LOGIN_ACCOUNT.USER_ID` の仕様である UUIDv7 チェックを含める場合は `parsedUUID.Version() != 7`）に修正する。

```go
	if _, parseErr := uuid.Parse(uid.String); !uid.Valid || parseErr != nil {
		if err = tx.Commit(); err != nil {
			return nil, err
		}
		return &model.LogoutAttemptResult{Status: model.LogoutStatusNotFound}, nil
	}
```

```go
func selectLoginUserID(ctx context.Context, tx *sql.Tx, sid string) (uid sql.NullString, err error) {
	err = tx.QueryRowContext(ctx, selectTargetUserIDQuery, sid).Scan(&uid)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullString{}, nil
	}
	return uid, err
}
```
