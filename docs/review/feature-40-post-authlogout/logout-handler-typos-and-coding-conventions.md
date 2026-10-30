# Swaggerアノテーション（@Failure 433）・コメントのタイポおよびコーディング規約

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-10-30 19:27:00
- **Target Files**:
  - [logout.go](backend/handler/logout.go)
  - [logout_test.go](backend/handler/logout_test.go)

## 1. 問題の概要
前回の修正により主要なタイポやテストケースは解消されましたが、`backend/handler/logout.go` の Swagger アノテーションで `@Failure 403` が `@Failure 433`（存在しないステータスコード）になっている点やコメント内のタイポ、`@spec` 記述、`gofmt` 未適用箇所が残っています。

## 2. 詳細な指摘内容
1. **Swagger アノテーションおよびコメントのタイポ (`backend/handler/logout.go`)**:
   - L16: `// LogoutHandler は POST /api/auth/logout のHandker を返す。` → 後半の `Handker` がタイポのまま残っています（`Handler`）。
   - L24: `// @Failure 433 {object} model.ErrorResponse` → `403` が `433` に誤記されています。Swagger ドキュメント生成時に不正なステータスコードになるため `403` に修正が必要です。
   - L11: `// LogoutHandlerOptionsはログインCookieの環境別設定です．` → 句点が全角ピリオド `．` になっています（`。` に統一）。
2. **Code-as-Docs (`@spec`) とテストでの空ボディ確認**:
   - `TESTING_GUIDE.md` の規約に従い、`LogoutHandler` の GoDoc コメントに `@spec` を記載することが推奨されます。
   - `backend/handler/logout_test.go` (L39, L67) において、成功時にレスポンスボディが空（`assert.Empty(t, w.Body.String())`）であることをアサートしておくと、`c.Status(http.StatusOK)` の挙動がテストでも保証されます。
   - また `assert.Equal` の引数順序は `assert.Equal(t, expected, actual)` が標準のため、L31 は `assert.Equal(t, "4c7bcdd6-898b-483f-bd0d-7a4191a0384b", input.SessionID)`、L61 は `assert.Equal(t, "", input.SessionID)` とすると失敗時の diff 表示が正確になります。
3. **フォーマット (`gofmt`) と不要なコメントアウト**:
   - `backend/handler/logout.go`: L4（スペースインデント混在）、L34（タブ位置）、L65（`){` のスペース不足）に `gofmt` が未適用です。
   - `backend/handler/logout_test.go`: L14（`//"github.com/stretchr/testify/require"` のコメントアウト残存）および L18, L21, L30, L60, L105 の戻り値 `(error)` の不要な括弧があります。

## 3. 推奨される修正案
1. `backend/handler/logout.go` の L16 (`Handker` → `Handler`) および L24 (`@Failure 433` → `@Failure 403`) を修正し、`@spec` コメントを追記してください。
2. `backend/handler/logout_test.go` の正常系2ケースに `assert.Empty(t, w.Body.String())` を追加し、不要なコメントアウト行（L14）を削除してください。
3. `gofmt -w logout.go logout_test.go` を実行してフォーマットを整えてください。
