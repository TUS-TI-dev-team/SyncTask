# Step 1: Ginミドルウェアの仕組みを理解する

[← README に戻る](README.md) | [Step 2 →](step-2.md)

## やること(概要)

認証の前に、まず「ミドルウェアとは何か」を小さなロガーを自作して掴みます。Hono でも `app.use()` で似た概念に触れていますが、Go/Gin では「**`*gin.Context` を引数に取り、次の処理を `c.Next()` で明示的に呼ぶ関数**」という素朴な構造です。

## やること(詳細)

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

## 調べるキーワード

- `gin.HandlerFunc`、`gin.Context.Next()`、`gin.Context.Abort()`
- `time.Since`

## 参考リンク

- Gin公式: カスタムミドルウェアの書き方 — https://gin-gonic.com/ja/docs/middleware/custom-middleware/
- Gin公式: RouterGroup(`Use` のドキュメント) — https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.Use
- Go標準 `net/http` の理解(ページ単位で目を通せばOK) — https://pkg.go.dev/net/http
- Go ミドルウェアの考え方(解説記事、公式ほど信頼できない点に注意) — https://echo.labstack.com/middleware/(Echoのものだが概念共通)

## ステップ末尾のテスト課題

- ロガーミドルウェアのテストは Step 6 まで先行しなくてOK。ただし `httptest.NewRecorder()` + `gin.CreateTestContext` の存在だけ調べておくこと

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 1 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- 本ステップの参考リンクにある Gin公式のカスタムミドルウェア解説を通読し、「ミドルウェア = `gin.HandlerFunc` を返す関数」という構造を確認する
- 既存の `backend/router/router.go` で `r.Use(...)` がどう使われているか(何を登録しているか)を探す
- Hono の `app.use()` の記憶と対応づける。「引数に渡した関数が全リクエストで動く」という点は同じ

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- `c.Next()` を呼び忘れると、後続の handler が一切実行されない(ログだけ出て終わる)
- 「before」処理と「after」処理の境界は `c.Next()` の**書いた位置**。処理時間は「before の先頭で時刻を取る → after で差分を取る」のが定石
- ステータスコードは `c.Writer.Status()` で取れる。パスとメソッドは `c.Request` の方に入っている
- `r.Use()` を呼ぶ位置で適用範囲が変わる。`SetupRouter` のルート定義より前に置くのが「全ルート」の意味

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- まず「開始時刻を取る → `c.Next()` → 差分とステータスを1行で出力する」形から始める。フォーマットを凝るのは動いてから
- 動いたら実験課題: ミドルウェア内で `c.Abort()` を呼び、後続の handler が止まることを curl で確認する。これができればミドルウェアの仕組みは掴めている

</details>
</details>
