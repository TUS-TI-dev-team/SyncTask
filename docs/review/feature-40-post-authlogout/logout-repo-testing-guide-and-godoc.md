# TESTING_GUIDE.md 規約違反（GoDoc @spec 欠落および logout_test.go のビルド失敗）

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-11-13 16:09:00
- **Target Files**:
  - [logout.go](backend/repository/logout.go)
  - [logout_test.go](backend/repository/logout_test.go)
  - [TESTING_GUIDE.md](backend/TESTING_GUIDE.md)

## 1. 問題の概要
`backend/repository/logout.go` において、`BeginTx` / `Commit` エラー時のみ第1戻り値が `nil` となっており `LogoutStatusDbFailed` が返らない点、および `backend/repository/logout_test.go` の単体テストが未実装でコンパイルエラーとなる点が残っています。

## 2. 詳細な指摘内容

1. **`BeginTx` および `Commit` 失敗時の戻り値の統一（`repository/logout.go` 33, 49, 67行目）**
   - クエリ実行失敗時（44, 56, 62行目）は `&model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err` を返している一方、`r.db.BeginTx`（33行目）や `tx.Commit()`（49, 67行目）が失敗した時だけ `return nil, err` になっています。
   - この状態で `BeginTx` や `Commit` が失敗すると、呼び出し元の `service/logout.go`（50〜52行目 `if result == nil { return fmt.Errorf("logout result is nil") }`）に入り、DB フェールセーフ（`LogoutStatusDbFailed`）にならず 500 エラーになってしまいます。
   - 33, 49, 67行目も `return &model.LogoutAttemptResult{ Status: model.LogoutStatusDbFailed }, err` に統一すると、すべての DB 障害に対して一貫してフェールセーフが機能します。

2. **`backend/repository/logout_test.go` の単体テスト実装**
   - `logout_test.go` 内の未使用 import および未使用変数 `attempt` により、`go test ./repository/...` がビルド失敗します。`sqlmock` を用いた単体テストの実装が必要です。

## 3. 推奨される修正案

- `LogoutRepository`, `NewLogoutRepository`, `AttemptLogout` に GoDoc コメントと `@spec` を追記する。
- `repository/logout_test.go` に `go-sqlmock` を用いた単体テストを実装し、`go test ./...` がパスする状態にする。
