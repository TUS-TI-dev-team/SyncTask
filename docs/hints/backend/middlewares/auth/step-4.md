# Step 4: 検証結果を後続に渡す

[← Step 3](step-3.md) | [README に戻る](README.md) | [Step 5 →](step-5.md)

## やること(概要)

認証が通ったリクエストの `USER_ID` を `*gin.Context` に載せ、後続の handler から参照できるようにします。「ミドルウェアが後続処理にコンテキストを引き渡す」設計の定石です。

## やること(詳細)

1. コンテキストキーを定数化する(マジックストリングを散らさない)
2. `c.Set(key, value)` で user_id を注入する
3. 参照用のヘルパーを用意する(handler 側の書き心地を良くするため)
   - 骨組み:
     ```go
     // UserID はコンテキストから認証済みユーザーIDを取り出します。
     // 型アサーションに失敗した場合の挙動を自分で設計すること。
     ```
4. 既存の task 系 handler(`handler/task_get.go` など)を一度読み、将来「user_id で自分のタスクだけを返す」改造を想像しながら、どこでこの値を使うかメモを取る(handlerの改造自体は本ロードマップのスコープ外)

## 調べるキーワード

- `gin.Context.Set` / `Get` / `GetString` 系
- Go の型アサーション(`v, ok := x.(T)`)

## 参考リンク

- `gin.Context.Set` — https://pkg.go.dev/github.com/gin-gonic/gin#Context.Set
- Go公式ツアー: 型アサーション(復習) — https://go.dev/tour/methods/15

## ステップ末尾のテスト課題

- 「`c.Set` で入れた値をヘルパーで取り出せる」「値が無い場合の挙動が設計通り」のテスト

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 4 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- pkg.go.dev の `gin.Context.Set` / `Get` のシグネチャ。値は `any` で保持される点を確認する
- Go公式ツアーの型アサーションの節。`v, ok := x.(T)` の2戻り値形式を復習する
- 既存 handler で共通の値を取り出している箇所があれば、その書き心地を観察する

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- `c.Get` の戻り値は `any`。使う側で型アサーションが必ず必要で、**1戻り値形式(`x.(T)`)は値が型違い/未設定だと panic する**。ミドルウェアのヘルパーでは2戻り値形式(`v, ok :=`)が安全
- キーを文字列リテラルで各所に書くと、handler 側での打ち間違えがコンパイルエラーにならず発見が遅れる。**定数化 + ヘルパー関数経由**にすれば文字列は1箇所に閉じる
- `gin.Context` と `context.Context` は**別物**。`c.Set` の相手は前者。Step 3 で使った `ctx` の方ではない
- `GetString` 系の便利メソッドもあるが、user_id の型(文字列？ 数値？)を決めるのは設計判断。既存の `USER_ID` の型を migrations で確認してから決める

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- 作る順番は「キー定数 → `c.Set` → ヘルパー関数」。ヘルパーは `UserID(c *gin.Context) (型, bool)` のような形を想定
- ヘルパーの型アサーション失敗時の挙動(ゼロ値を返す / false を返す / panic させる)は**自分で決めてコメントに書く**。ここがこのステップの学習ポイント
- Step 7 の振り返りに向け、「キー名を何にしたか、なぜか」もメモしておくと後で楽

</details>
</details>
