# Step 6: テストを書く

[← Step 5](step-5.md) | [README に戻る](README.md) | [Step 7 →](step-7.md)

## やること(概要)

Step 2〜5 で先取りしたテストを整理・拡充し、認証ミドルウェアのテストを完成させます。DBを使わず、service を fake に差し替えてミドルウェア単体を検証する方式です。

## やること(詳細)

1. `backend/TESTING_GUIDE.md` を読み、規約(日本語テスト名 `<分類>: <期待される結果>`、`@spec`、require/assert の使い分け)を確認する
2. service の interface を満たす fake をテスト内に定義し、ミドルウェア関数に注入する
   - Go ではモックライブラリ無しで struct を1つ書けば fake になる(これが interface 学習の締めくくり)
3. テストケース案(自分で網羅性を考えて増やすこと):
   - Cookieが無い → 401 で abort
   - Cookieがあるがセッション不正 → 401
   - セッション有効 → 次のhandlerが実行され、context に user_id が入っている
4. `require` と `assert` をどこで使い分けるか意識して書く(中断すべき場面と、値検証を続けたい場面)

## 調べるキーワード

- `httptest.NewRecorder`、`gin.CreateTestContext` またはテスト用ルーターの構築
- `testify` の require / assert

## 参考リンク

- 既存リポジトリ文書: `backend/TESTING_GUIDE.md`(最優先で読む)
- httptest パッケージ — https://pkg.go.dev/net/http/httptest
- testify — https://pkg.go.dev/github.com/stretchr/testify
- 既存コード: `backend/router/router_test.go`、`backend/handler/login_test.go`(テストの型の写経元)

## ステップ末尾のテスト課題

- (このステップ自体がテストの完成なので、テスト課題は上記のテストケース案を参照)

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 6 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- `backend/TESTING_GUIDE.md` を最初に読む。テスト名・`@spec`・require/assert の規約が全部書いてある
- `backend/router/router_test.go` と `backend/handler/login_test.go`: リクエストの組み立て方・実行・検証の一式を真似る
- pkg.go.dev の httptest: `NewRecorder` の役割(レスポンスの記録係)を確認する

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- fake は「interface と同じシグネチャのメソッドを持つ struct」をテストファイル内に書くだけ。メソッド内で呼び出し引数を記録したり、フィールドに事前設定しておいた戻り値/エラーを返したりする
- ミドルウェア単体の検証とルート全体の検証は道具が違う:
  - `gin.CreateTestContext`: ミドルウェア1つに絞って `c.Next()` の前後の挙動を見る
  - テスト用ルーター(ルートを登録して `router.ServeHTTP(recorder, req)`): URL・ステータスまで含めて見る
- abort されたかの判定は「後続 handler を模したフラグ関数が呼ばれたか」で見るのが確実。`c.IsAborted()` も使える
- require は「ここで失敗したら以降の検証は無意味」な場所(例: レスポンス生成がそもそも失敗したら本文検証は不要)。使いすぎると1実行で1つの失敗しか分からない
- テスト名の `<分類>` は `middleware` など既存規約に合わせる。日本語で期待結果を書く点を忘れない

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- まず「セッション有効 → 次のhandlerが実行される」の1ケースを fake で作り、緑にしてから異常系を足す。fake の形が固まれば残りは応用
- 各ステップ末尾で先取りしたテストを集め、重複を整理して一つのテストファイルにまとめ上げる
- テストが直らないときは、fake の戻り値の型(特に `SessionInfo` のゼロ値)が実装とずれていないかを最初に疑う

</details>
</details>
