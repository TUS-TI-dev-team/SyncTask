# Step 5: 401応答の設計とルート保護

[← Step 4](step-4.md) | [README に戻る](README.md) | [Step 6 →](step-6.md)

## やること(概要)

認証失敗時のレスポンスを既存のエラー体系に統合し、実際に `/api/tasks` を保護します。これで認証ミドルウェアとして完成です。

## やること(詳細)

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

## 調べるキーワード

- HTTP 401 と 403 の違い(どちらを返すべきかの判断材料)
- `gin.RouterGroup` のネストとパスの合成

## 参考リンク

- MDN: 401 Unauthorized — https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/401
- MDN: 403 Forbidden — https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/403
- 既存コード: `backend/model/`(AppError / ErrorResponse の定義)、`backend/handler/login.go` の `writeLoginError`
- Gin公式: ルートのグループ化 — https://gin-gonic.com/ja/docs/routing/grouping-routes/

## ステップ末尾のテスト課題

- router レベルのテスト:「認証Cookie無しの `/api/tasks` リクエストは401」「有効なセッションなら200」

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 5 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- `backend/model/` の `AppError` / `ErrorResponse` と、`backend/handler/login.go` の `writeLoginError`: エラー応答 JSON の形(キー名・code の付け方)の写経元
- MDN の 401 と 403 の両ページ: 「どちらを返すべきか」の判断材料になる
- `backend/router/router.go` の既存のグループ化と `router_test.go`: どのルートを authorized グループへ移すか、壊れそうなテストはどれか

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- 401 は「未認証(誰だかわからない)」、403 は「認証済みだが権限がない」。セッションが無い/切れたリクエストは前者に当たる
- 応答内容の粒度: 「Cookieが無い」「存在しないセッション」「期限切れ」を区別して返すと、**攻撃者に「そのセッションIDは実在した」という情報を与える**。原則は同一の401に寄せる
- グループ移動時のパス合成: グループのパスと個別ルートのパスの**つなぎ目**で `/api/tasks` が `/api/api/tasks` や `/api/task` になることがある。既存テストが壊れたらまず URL がずれていないか疑う
- JSON 応答を書くときは Content-Type や HTTP ステータスの設定も既存の `writeLoginError` と同じ形に合わせる

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- 応答まわりから着手: 既存の `writeLoginError` と同じ JSON 形を作る関数を middleware 側に用意し、curl の `-i` でステータスと本文を確認する
- 3ケース(Cookie無し/存在しない/期限切れ)の全てで**同じ応答**になることを curl で確認する
- グループ化は「既存ルート定義をなるべく動かさず authorized グループの下に移す」ことを優先し、`GET /api/tasks` 系の curl でパスが変わっていないか確認してから次へ

</details>
</details>
