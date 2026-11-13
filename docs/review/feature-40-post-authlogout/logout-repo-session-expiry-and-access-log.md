# セッション有効期限（EXPIRES_AT）判定および ACCESS_LOG 記録の仕様未反映

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-11-13 16:09:00
- **Target Files**:
  - [logout.go](backend/repository/logout.go)
  - [05_logout.md](docs/design/process_design/05_logout.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`backend/repository/logout.go` に追加された `EXPIRES_AT` 判定周りにおいて、コンパイルエラーとなる箇所が3点（47行目の `||` 残存、60行目の `time.Time` 比較演算子および条件式構造、74行目の戻り値不足）あり、また `service/logout.go` の 50行目で `result == nil` 時に nil ポインタパニックが発生する問題があります。

## 2. 詳細な指摘内容

1. **47行目: 末尾の `||` による構文エラー（`repository/logout.go`）**
   ```go
   if _, parseErr := uuid.Parse(uid.String); !uid.Valid || parseErr != nil || {
   ```
   `parseErr != nil` の後ろに `||` が残っておりコンパイルエラーになります。

2. **74行目: `selectLoginUserID` の戻り値不足によるコンパイルエラー（`repository/logout.go`）**
   ```go
   func selectLoginUserID(ctx context.Context, tx *sql.Tx, sid string) (uid sql.NullString, expires_at sql.NullTime, err error) {
   	err = tx.QueryRowContext(ctx, selectTargetUserIDQuery, sid).Scan(&uid, &expires_at)
   	if errors.Is(err, sql.ErrNoRows) {
   		return sql.NullString{}, nil
   	}
   ```
   戻り値が3つ（`uid`, `expires_at`, `err`）定義されていますが、74行目で2つしか返していません。`return sql.NullString{}, sql.NullTime{}, nil` に修正が必要です。

3. **60行目: `time.Time` の `<=` 比較エラーおよび無条件クエリ実行バグ（`repository/logout.go`）**
   ```go
   if _, err = tx.ExecContext(ctx, resetReauthFailedCountQuery, attempt.Now, uid.String); expires_at.Time <= attempt.Now && err != nil {
   	return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err
   }
   ```
   - Go の `time.Time` 型は `<=` 演算子で比較できないためコンパイルエラーになります（`expires_at.Time.After(attempt.Now)` 等を使用します）。
   - `if 初期化文; 条件` の形になっているため、**有効期限切れかどうかに関わらず `tx.ExecContext`（再認証失敗回数リセット）が必ず実行されてしまいます**。
   - さらに、有効期限内（`expires_at.Time > attempt.Now`）のときに DB エラー（`err != nil`）が起きると、条件式が `false` になってエラーハンドリングを素通りしてしまいます。

4. **`service/logout.go` 50行目: `result == nil` 時の nil ポインタパニック**
   ```go
   result, err := s.repo.AttemptLogout(ctx, attempt)
   if err != nil && result.Status != model.LogoutStatusDbFailed{
   	return err
   }
   ```
   `BeginTx` や `Commit` 失敗時（またはモックが `(nil, err)` を返した時）に `result` が `nil` になるため、50行目の `result.Status` アクセスでパニック（nil pointer dereference）が発生します。

## 3. 推奨される修正案

- **`repository/logout.go`（60行目付近）**:
  有効期限内（`expiresAt.Valid && expiresAt.Time.After(attempt.Now)`）の場合にのみ `resetReauthFailedCountQuery` を実行するよう `if` 文を分けてください。
  ```go
  // reset reauth count (有効期限内のセッションの場合のみ実行)
  if expiresAt.Valid && expiresAt.Time.After(attempt.Now) {
  	if _, err = tx.ExecContext(ctx, resetReauthFailedCountQuery, attempt.Now, uid.String); err != nil {
  		return &model.LogoutAttemptResult{Status: model.LogoutStatusDbFailed}, err
  	}
  }
  ```
- **`service/logout.go`（50行目付近）**:
  `result == nil` チェックを先に行うか、`err != nil` 内で `result != nil` を確認してください。
  ```go
  result, err := s.repo.AttemptLogout(ctx, attempt)
  if err != nil {
  	if result != nil && result.Status == model.LogoutStatusDbFailed {
  		return nil
  	}
  	return err
  }
  ```
