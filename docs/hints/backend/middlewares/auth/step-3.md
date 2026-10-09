# Step 3: セッション検証の層を作る

[← Step 2](step-2.md) | [README に戻る](README.md) | [Step 4 →](step-4.md)

## やること(概要)

Cookie から取り出したセッションIDが `LOGIN_SESSION` テーブルに存在し、期限切れでないかを検証する repository / service を作ります。ここがこのロードマップで一番の山場です。

## やること(詳細)

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

## 調べるキーワード

- Go の interface(暗黙的実装。Javaと違い `implements` と書かない)
- `context.Context`(既存コードで既に使われているので読み返す)
- pgx の QueryRow / Scan、`sql.ErrNoRows` 相当の扱い

## 参考リンク

- Go公式: Effective Go — Interfaces — https://go.dev/doc/effective_go#interfaces
- Go公式ブログ: Goとcontext — https://go.dev/blog/context
- pgx v5 ドキュメント — https://pkg.go.dev/github.com/jackc/pgx/v5
- 既存コード: `backend/repository/` 配下(特に login 系)、`backend/service/login.go`(DI構造の写経元)
- セッション管理の一般論(設計の背景理解) — https://developer.mozilla.org/ja/docs/Web/HTTP/Guides/Authentication

## ステップ末尾のテスト課題

- service層をrepositoryのfakeでテスト:「有効なセッションIDは検証成功」「期限切れセッションIDはエラー」「存在しないIDはエラー」の3ケース

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 3 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- `backend/repository/` 配下の login 系: pgx のクエリ実行・`Scan` の型合わせ・戻り値の作り方の写経元。まず60行ほど読む
- `backend/service/login.go` の `LoginDependencies` 構造: service 層の DI パターンの写経元
- Go公式ブログの context 記事: 「第1引数に `ctx context.Context` を置く」慣習の背景が分かる

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- Go の interface は `implements` 宣言が不要。**メソッドを持つ具体型は、その形の interface に勝手に代入できる**。service の依存を interface 型で受けておくと、テストで別の実装を差し込める
- pgx で「0行だった」ことの検出方法は、既存 repository がどう扱っているかをそのまま倣う。`pgx.ErrNoRows` を初めて見る型でも、既存コードのパターンを写せば問題ない
- 期限判定の場所の判断材料:
  - SQL 側(`WHERE expires_at > now()`): 「有効な行だけ」が取れる。ただし「存在しない」と「期限切れ」を区別したくなったときに情報が減る
  - service 側(取得後に Go で比較): 行が取れるので「存在しない/期限切れ」を区別できる。ただし service が時刻比較の責務を持つ
  - どちらを選んでも構わない。**選んだ理由をコメントに書くこと**がこのステップの狙い
- レイヤーの役割: repository はテーブル構造を知る層、middleware は HTTP(Cookie・ステータスコード)を知る層。service はその間を繋ぐ。互いの専門用語を持ち込まないのが清潔

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- 実装の順番は repository → service → middleware。**repository だけを先に動かして SQL が正しいか手動確認**(ログインして得た実物の `SESSION_ID` で引き、1行取れるか)してから service に繋ぐと切り分けが楽
- service は「有効なら `SessionInfo`、無効なら error」だけを返す薄い層に保つ。Cookie やステータスコードの知識を service に持ち込まない
- middleware 側は「service の error を abort に変換する」だけの分岐にする。ここで応答の作り込み(401のJSON)をしないのは Step 5 でやるため

</details>
</details>
