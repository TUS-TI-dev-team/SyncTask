# 認証ミドルウェア自作ステップ(学習用ドキュメント)

このドキュメントは、**自分の手**で SyncTask バックエンドにセッション認証ミドルウェアを実装するための学習ロードマップです。
答えを書くのではなく、「各ステップで何をやるか」「何を調べるか」「詰まったらどこを見るか」を示します。コード例は最小限の骨組みのみです。

## 前提知識

| 項目 | 内容 |
| --- | --- |
| 技術レベル想定 | Go は A Tour of Go 完了レベル。Web開発は Hono で2ヶ月程度の経験(記憶は薄れ気味) |
| 既存実装(活用するもの) | `backend/handler/login.go` がログイン成功時に **`sync_task_sid`**(HttpOnly セッションCookie)と **`XSRF-TOKEN`**(CSRF Cookie)を発行済み |
| 既存テーブル | `LOGIN_SESSION`(`SESSION_ID`, `USER_ID`, `EXPIRES_AT`, `USER_AGENT`, `CREATED_AT`)。`backend/db/migrations/000001_init.up.sql` 参照 |
| このドキュメントのスコープ | **セッション認証のみ**(CSRF検証は次の学習テーマとして切り出し) |

## 全体工程図

```
[Step 1] Ginミドルウェアの仕組み理解(ロガー自作)
    ↓
[Step 2] Cookieを読む(c.Cookie)
    ↓
[Step 3] セッション検証の層を作る(repository / service)
    ↓
[Step 4] 検証結果を後続に渡す(c.Set / c.Get)
    ↓
[Step 5] 401応答の設計とルート保護(AppError + tasksグループ適用)
    ↓
[Step 6] テストを書く(httptest + fake service)
    ↓
[Step 7] まとめと次の一歩
```

各ステップの目安は **1日** です。Step 2〜5 は順番に依存しますが、Step 6 は Step 3 以降の各ステップで少しずつ先取りして構いません。

## 実装方針(先に決めておくこと)

以下は本ドキュメント策定時に決定済みの方針です。迷ったらこれに従ってください。

- **配置**: 新規パッケージ `backend/middleware/`(`middleware/auth.go`)
- **レイヤー経路**: `repository`(新規 LoginSessionRepository)→ `service`(新規 SessionService)→ `middleware`。既存の `router.go` のDIパターン(`router.go` の 57〜65 行目付近)に倣う
- **保護対象**: `POST/GET/PATCH /api/tasks` のみ。`/api/auth/*` と `/health-check` は公開のまま
- **401応答**: 既存の `model.AppError` / `model.ErrorResponse` パターンで統一(code: `UNAUTHORIZED`)
- **テスト**: `httptest` + service の fake(interface 差し替え)。DBや testcontainers は使わない

---

## Step 1: Ginミドルウェアの仕組みを理解する

### やること(概要)

認証の前に、まず「ミドルウェアとは何か」を小さなロガーを自作して掴みます。Hono でも `app.use()` で似た概念に触れていますが、Go/Gin では「**`*gin.Context` を引数に取り、次の処理を `c.Next()` で明示的に呼ぶ関数**」という素朴な構造です。

### やること(詳細)

1. `backend/middleware/` パッケージを新規作成する
2. リクエストのメソッド・パス・処理時間を `log` に出力するだけのロガーミドルウェアを書く
   - 骨組み(この程度のヒントまで。中身は自分で書く):
     ```go
     package middleware

     import "github.com/gin-gonic/gin"

     // Logger はリクエストの概要をログ出力するミドルウェアを返します。
     func Logger() gin.HandlerFunc {
         return func(c *gin.Context) {
             // before: リクエスト情報を記録
             c.Next()
             // after: ステータスコードと処理時間を記録
         }
     }
     ```
3. `router.SetupRouter` 内で `r.Use(middleware.Logger())` を呼び、全ルートで動くことを確認する
4. 時間があれば「`r.Use` と `group.Use` の違い」「ミドルウェア内で `c.Abort()` すると何が起きるか」を実験する

### 調べるキーワード

- `gin.HandlerFunc`、`gin.Context.Next()`、`gin.Context.Abort()`
- `time.Since`

### 参考リンク

- Gin公式: Middleware の書き方 — https://gin-gonic.com/docs/examples/custom-middleware/
- Gin公式: RouterGroup(`Use` のドキュメント) — https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.Use
- Go標準 `net/http` の理解(ページ単位で目を通せばOK) — https://pkg.go.dev/net/http
- Go ミドルウェアの考え方(解説記事、公式ほど信頼できない点に注意) — https://echo.labstack.com/docs/middleware(Echoのものだが概念共通)

### ステップ末尾のテスト課題

- ロガーミドルウェアのテストは Step 6 まで先行しなくてOK。ただし `httptest.NewRecorder()` + `gin.CreateTestContext` の存在だけ調べておくこと

---

## Step 2: Cookieを読む

### やること(概要)

ブラウザが送ってくる `sync_task_sid` Cookie をミドルウェアで取り出すところを作ります。ここが認証の入り口です。

### やること(詳細)

1. `middleware/auth.go` に認証ミドルウェアの関数シグネチャを定義する(まだ検証はしない)
   - 後続のStepで service を注入するため、**関数型のミドルウェア(`func Auth(...) gin.HandlerFunc`)にする**設計にしておく
2. `c.Cookie("sync_task_sid")` で Cookie 値を取得する
3. Cookie が存在しない/空の場合の分岐を作る(この時点ではとりあえず `c.AbortWithStatus(401)` でOK。401の応答設計は Step 5)
4. `curl -i` で Cookie 無し/有りを送り、ログに取得した値が出ることを確認する
   - Cookie付きリクエストの例: `curl -i http://localhost:8080/api/tasks -H "Cookie: sync_task_sid=dummy"`

### 調べるキーワード

- `gin.Context.Cookie()`
- HTTP Cookie の仕組み(`HttpOnly`, `Secure`, `SameSite` — login.go で既に使われているので読み返す)

### 参考リンク

- `gin.Context.Cookie` — https://pkg.go.dev/github.com/gin-gonic/gin#Context.Cookie
- MDN: HTTP Cookie(仕組みの復習) — https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/Cookie
- Go標準 `http.Cookie` — https://pkg.go.dev/net/http#Cookie
- 既存コード: `backend/handler/login.go` の `http.SetCookie` 部分(Cookie発行側との対応を確認)

### ステップ末尾のテスト課題

- 「Cookieが無いと abort される」テストを書いてみる(規約: `middleware: Cookieが無い場合リクエストを中断する` のような日本語テスト名)

---

## Step 3: セッション検証の層を作る

### やること(概要)

Cookie から取り出したセッションIDが `LOGIN_SESSION` テーブルに存在し、期限切れでないかを検証する repository / service を作ります。ここがこのロードマップで一番の山場です。

### やること(詳細)

1. **repository**: `repository/login_session.go` に新規作成
   - 既存の `repository/` 配下(例: `login` 関連)のコードを読み、pgx を使ったクエリの書き方・戻り値の型を合わせる
   - 必要なメソッド案: `FindBySessionID(ctx, sessionID)`(存在確認+有効期限取得)
   - SQLは「セッションIDで1件取得」のみ。期限判定は service 側で行うか SQL 側で行うか、**自分で設計して理由をコメントに書く**
2. **service**: `service/session.go` に新規作成
   - インターフェースを切る(Step 6 のfake差し替えに必要)
     ```go
     type SessionService interface {
         ValidateSession(ctx context.Context, sessionID string) (SessionInfo, error)
     }
     ```
   - `SessionInfo` は `UserID` + 有効期限など、ミドルウェアが後続に渡したい情報を持たせる
   - 既存 `service.NewLoginService` のDIパターン(`service.LoginDependencies` を使う構造)に倣う
3. **middleware 側**: service を受け取り、Cookie値で `ValidateSession` を呼び、エラーなら abort する分岐まで実装
4. 手動確認: 実際にログインして得た `sync_task_sid` を curl に付け、検証が通ること。存在しない値を付けて落ちること

### 調べるキーワード

- Go の interface(暗黙的実装。Javaと違い `implements` と書かない)
- `context.Context`(既存コードで既に使われているので読み返す)
- pgx の QueryRow / Scan、`sql.ErrNoRows` 相当の扱い

### 参考リンク

- Go公式: Effective Go — Interfaces — https://go.dev/doc/effective_go#interfaces
- Go公式ブログ: Goとcontext — https://go.dev/blog/context
- pgx v5 ドキュメント — https://pkg.go.dev/github.com/jackc/pgx/v5
- 既存コード: `backend/repository/` 配下(特に login 系)、`backend/service/login.go`(DI構造の写経元)
- セッション管理の一般論(設計の背景理解) — https://developer.mozilla.org/ja/docs/Web/HTTP/Authentication

### ステップ末尾のテスト課題

- service層をrepositoryのfakeでテスト:「有効なセッションIDは検証成功」「期限切れセッションIDはエラー」「存在しないIDはエラー」の3ケース

---

## Step 4: 検証結果を後続に渡す

### やること(概要)

認証が通ったリクエストの `USER_ID` を `*gin.Context` に載せ、後続の handler から参照できるようにします。「ミドルウェアが後続処理にコンテキストを引き渡す」設計の定石です。

### やること(詳細)

1. コンテキストキーを定数化する(マジックストリングを散らさない)
2. `c.Set(key, value)` で user_id を注入する
3. 参照用のヘルパーを用意する(handler 側の書き心地を良くするため)
   - 骨組み:
     ```go
     // UserID はコンテキストから認証済みユーザーIDを取り出します。
     // 型アサーションに失敗した場合の挙動を自分で設計すること。
     ```
4. 既存の task 系 handler(`handler/task_get.go` など)を一度読み、将来「user_id で自分のタスクだけを返す」改造を想像しながら、どこでこの値を使うかメモを取る(handlerの改造自体は本ロードマップのスコープ外)

### 調べるキーワード

- `gin.Context.Set` / `Get` / `GetString` 系
- Go の型アサーション(`v, ok := x.(T)`)

### 参考リンク

- `gin.Context.Set` — https://pkg.go.dev/github.com/gin-gonic/gin#Context.Set
- Go公式ツアー: 型アサーション(復習) — https://go.dev/tour/methods/15

### ステップ末尾のテスト課題

- 「`c.Set` で入れた値をヘルパーで取り出せる」「値が無い場合の挙動が設計通り」のテスト

---

## Step 5: 401応答の設計とルート保護

### やること(概要)

認証失敗時のレスポンスを既存のエラー体系に統合し、実際に `/api/tasks` を保護します。これで認証ミドルウェアとして完成です。

### やること(詳細)

1. `model.AppError` / `model.ErrorResponse` / `handler/writeLoginError` を読み、既存のエラー応答パターンを把握する
2. 認証ミドルウェアから `401 + code: UNAUTHORIZED` の `ErrorResponse` を返す実装にする
   - **設計課題(必ず自分で考え、本文にメモを残すこと)**: 「Cookieが無い」「セッションが存在しない」「期限切れ」をすべて同じ401メッセージにするか、区別するか。ヒント: 攻撃者に情報を与えないという観点
3. `router.SetupRouter` で tasks ルートを認証必須のグループに移す
   ```go
   authorized := api.Group("/", middleware.Auth(sessionService))
   {
       authorized.POST("/tasks", ...)
       authorized.GET("/tasks/:task_id", ...)
       authorized.PATCH("/tasks/:task_id", ...)
   }
   ```
4. curl で総合確認: Cookie無し→401、期限切れ→401、正常ログイン後のCookie→200
5. `router_test.go` を読み、既存テストがルート変更で壊れていないか確認する

### 調べるキーワード

- HTTP 401 と 403 の違い(どちらを返すべきかの判断材料)
- `gin.RouterGroup` のネストとパスの合成

### 参考リンク

- MDN: 401 Unauthorized — https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/401
- MDN: 403 Forbidden — https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/403
- 既存コード: `backend/model/`(AppError / ErrorResponse の定義)、`backend/handler/login.go` の `writeLoginError`
- Gin公式: ルーティングのグループ化 — https://gin-gonic.com/docs/examples/grouping-routes/

### ステップ末尾のテスト課題

- router レベルのテスト:「認証Cookie無しの `/api/tasks` リクエストは401」「有効なセッションなら200」

---

## Step 6: テストを書く

### やること(概要)

Step 2〜5 で先取りしたテストを整理・拡充し、認証ミドルウェアのテストを完成させます。DBを使わず、service を fake に差し替えてミドルウェア単体を検証する方式です。

### やること(詳細)

1. `backend/TESTING_GUIDE.md` を読み、規約(日本語テスト名 `<分類>: <期待される結果>`、`@spec`、require/assert の使い分け)を確認する
2. service の interface を満たす fake をテスト内に定義し、ミドルウェア関数に注入する
   - Go ではモックライブラリ無しで struct を1つ書けば fake になる(これが interface 学習の締めくくり)
3. テストケース案(自分で網羅性を考えて増やすこと):
   - Cookieが無い → 401 で abort
   - Cookieがあるがセッション不正 → 401
   - セッション有効 → 次のhandlerが実行され、context に user_id が入っている
4. `require` と `assert` をどこで使い分けるか意識して書く(中断すべき場面と、値検証を続けたい場面)

### 調べるキーワード

- `httptest.NewRecorder`、`gin.CreateTestContext` またはテスト用ルーターの構築
- `testify` の require / assert

### 参考リンク

- 既存リポジトリ文書: `backend/TESTING_GUIDE.md`(最優先で読む)
- httptest パッケージ — https://pkg.go.dev/net/http/httptest
- testify — https://pkg.go.dev/github.com/stretchr/testify
- 既存コード: `backend/router/router_test.go`、`backend/handler/login_test.go`(テストの型の写経元)

---

## Step 7: まとめと次の一歩

### やること(概要)

実装の振り返りと、今回あえてスコープから外した項目の接続先を整理して、学習を次に繋げます。

### やること(詳細)

1. 自分の書いた `middleware/auth.go` を上から読み直し、各判断(期限判定の場所、401の応答内容、コンテキストキー名など)に理由コメントが書かれているか確認する
2. 振り返りメモをこのファイルの下に追記する(下にテンプレあり)
3. 次の学習テーマの候補を把握する:
   - **CSRF検証**: 発行済みの `XSRF-TOKEN` をダブルサブミットCookie方式で検証するミドルウェア
   - **ログアウト**: `LOGIN_SESSION` の削除+Cookie失効
   - **user_id によるタスクの所有者検証**: Step 4 で注入した user_id を handler/service に渡し、「他人のタスクを見れない」ようにする

### 参考リンク

- OWASP Session Management Cheat Sheet(セッション設計のチェックリスト。自分の実装と突き合わせる) — https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
- OWASP Authentication Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html

---

## 振り返りメモ(テンプレート)

```markdown
## 振り返り(Step 7 完了時に追記)

- うまくいったこと:
- 詰まったことと解決方法:
- 自分が下した設計判断とその理由:
  - 期限判定の場所(SQL側 or service側):
  - 401メッセージの粒度:
  - コンテキストキーの設計:
- 次に学びたいこと:
```

## 進め方のヒント

- 各ステップの「調べるキーワード」で公式ドキュメント(pkg.go.dev / gin-gonic.com / go.dev)を**必ず一次情報として**確認する。Qiitaやブログ記事はあくまで補助
- Hono の経験と対応づけると理解が早い: `app.use()` ≒ `r.Use()`、`c.set()` ≒ `c.Set()`、`c.env` ≒ 依存性注入(関数引数)
- 詰まったら、まず既存コード(`service/login.go` や `router_test.go`)に同じパターンがないか探す。このリポジトリはパターンの見本が揃っている