# Step 2: Cookieを読む

[← Step 1](step-1.md) | [README に戻る](README.md) | [Step 3 →](step-3.md)

## やること(概要)

ブラウザが送ってくる `sync_task_sid` Cookie をミドルウェアで取り出すところを作ります。ここが認証の入り口です。

## やること(詳細)

1. `middleware/auth.go` に認証ミドルウェアの関数シグネチャを定義する(まだ検証はしない)
   - 後続のStepで service を注入するため、**関数型のミドルウェア(`func Auth(...) gin.HandlerFunc`)にする**設計にしておく
2. `c.Cookie("sync_task_sid")` で Cookie 値を取得する
3. Cookie が存在しない/空の場合の分岐を作る(この時点ではとりあえず `c.AbortWithStatus(401)` でOK。401の応答設計は Step 5)
4. `curl -i` で Cookie 無し/有りを送り、ログに取得した値が出ることを確認する
   - Cookie付きリクエストの例: `curl -i http://localhost:8080/api/tasks -H "Cookie: sync_task_sid=dummy"`

## 調べるキーワード

- `gin.Context.Cookie()`
- HTTP Cookie の仕組み(`HttpOnly`, `Secure`, `SameSite` — login.go で既に使われているので読み返す)

## 参考リンク

- `gin.Context.Cookie` — https://pkg.go.dev/github.com/gin-gonic/gin#Context.Cookie
- MDN: HTTP Cookie(仕組みの復習) — https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Headers/Cookie
- Go標準 `http.Cookie` — https://pkg.go.dev/net/http#Cookie
- 既存コード: `backend/handler/login.go` の `http.SetCookie` 部分(Cookie発行側との対応を確認)

## ステップ末尾のテスト課題

- 「Cookieが無いと abort される」テストを書いてみる(規約: `middleware: Cookieが無い場合リクエストを中断する` のような日本語テスト名)

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 2 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- pkg.go.dev の `gin.Context.Cookie` のドキュメントで、引数と**戻り値の個数・型**を確認する
- `backend/handler/login.go` の Cookie発行側(`http.SetCookie`)を読み、「発行された Cookie の名前・属性」と「読み出し側」を対応づける
- ブラウザの開発者ツール(Network タブ)で、実際のリクエストに付く `Cookie:` ヘッダーの形を見てみる

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- `c.Cookie(name)` は戻り値が2つ(値と error)。Cookie が存在しないことは **error で判別**する。値の空文字チェックだけでは不十分
- `HttpOnly` な Cookie は JavaScript から読めないが、**サーバー側(ミドルウェア)からは普通に読める**。混乱しやすい点
- curl での Cookie 送信は `-H "Cookie: name=value"` の形式。set-cookie された値をそのまま手で貼ってもよいし、curl の Cookie 機能を使ってもよい
- 「Cookie が無い」「Cookie はあるが空」の2ケースを分けて考えると、Step 5 の401設計で考えやすくなる

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- まず「Cookie を取得してログに出す」だけのミドルウェアを作り、curl で動作確認する
- その後、Cookie 無し/空の分岐を追加。401の応答作り込みは Step 5 の担当なので、ここでは `c.AbortWithStatus(401)` の仮実装で先に進んでOK

</details>
</details>
