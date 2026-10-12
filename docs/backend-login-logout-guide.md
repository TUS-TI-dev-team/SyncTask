# login処理の読み方と、新しいエンドポイントを作る手順

既存のlogin実装を題材に、Goの変数・関数の意図と、logoutを実装するときの考え方を整理します。

この文書は学習用の解説です。logoutの完成実装を示すものではありません。仕様は `docs/design/api_design/02_auth.md` の `POST auth/logout` を参照してください。

## 1. 処理の全体像

```text
HTTPリクエスト
    ↓
Router      URLと処理を結びつける
    ↓
Handler     HTTPの入力をGoのデータに変換する
    ↓
Service     ログインに必要な手順を進める
    ↓
Repository  SQLを使ってDBを読み書きする
    ↓
Service     DBの結果をログインの結果にまとめる
    ↓
Handler     Cookie・ステータス・JSONを返す
```

| 場所 | 主な責任 | loginでの具体例 |
| --- | --- | --- |
| router | URL登録、部品の組み立て | `/api/auth/login` の登録 |
| handler | HTTPの読み取り・返却 | JSONやCookieを読み、レスポンスを返す |
| service | 処理の手順・業務上の判断 | 入力検証、トークン生成、失敗時の待機 |
| repository | DB操作 | アカウント検索、セッション登録 |
| model | データの形 | LoginRequest、LoginResponse |
| util | 独立した補助処理 | メール正規化、乱数トークン生成 |

コードを読むときは「何を受け取り、何を判断・変更し、何を返すか」を確認します。新しい処理の配置先は「HTTPに関することか」「処理のルールか」「DBに関することか」で考えます。

## 2. Handlerの読み方

対象: [handler/login.go](../backend/handler/login.go)

```go
func LoginHandler(
    loginService service.LoginService,
    options LoginHandlerOptions,
) gin.HandlerFunc {
    return func(c *gin.Context) {
        // リクエストごとの処理
    }
}
```

ここには二つの実行タイミングがあります。

- `LoginHandler(...)` はRouterの組み立て時に呼ばれます。
- 返された `func(c *gin.Context)` はリクエストが来るたびに呼ばれます。

内側の関数は、外側で受け取ったloginServiceやoptionsを使えます。このような関数をクロージャと呼びます。

| 変数 | 意図 |
| --- | --- |
| loginService | ログイン処理を依頼する相手 |
| options | CookieSecureなどの環境別設定 |
| c | 今回のHTTPリクエストとレスポンスを操作する窓口 |
| req | 入力されたメールアドレスとパスワード |
| oldSessionID | 成功時に削除する既存セッションのID |
| result | Serviceが返すユーザー情報とCookie用情報 |
| err | 処理に失敗した理由 |
| maxAge | Cookieの有効期間を秒にした値 |

### JSONの読み取りとポインタ

```go
var req model.LoginRequest
if err := c.ShouldBindJSON(&req); err != nil {
    // JSONを読み取れなかった場合
}
```

`&req` はreqの置き場所を渡す記法です。ShouldBindJSONがそこへ値を書き込むために必要です。`*model.LoginRequest` は、その置き場所を指すポインタの型です。

`:=` は変数の宣言と代入をまとめた記法です。`if err := ...; err != nil` のerrは、このifとその分岐内で使用できます。

### Serviceへの受け渡し

```go
result, err := loginService.Login(
    c.Request.Context(),
    &req,
    model.LoginMetadata{
        IP:           c.ClientIP(),
        UserAgent:    c.Request.UserAgent(),
        OldSessionID: oldSessionID,
    },
)
```

reqはJSON本文から得た情報、LoginMetadataは本文以外から得た情報です。Contextはリクエストのキャンセルなどを下位の処理へ伝えます。

GinのcをServiceへ直接渡さず、必要な情報を取り出して渡すことで、ServiceをHTTPなしでもテストできます。

`oldSessionID, _ := c.Cookie(...)` の `_` は、その戻り値を使わない記法です。この場面ではCookieがないことも通常のケースとして扱っています。エラーを無視してよいかは場面ごとに判断します。

### 補助関数

| 関数 | 役割 |
| --- | --- |
| setLoginResponseHeaders | キャッシュやセキュリティ関連のヘッダー設定 |
| writeLoginError | エラーをHTTPステータスとJSONへ変換 |

## 3. Serviceの型と依存関係

対象: [service/login.go](../backend/service/login.go)

```go
type LoginService interface {
    Login(context.Context, *model.LoginRequest, model.LoginMetadata) (*model.LoginServiceResult, error)
    RecordInvalidRequest(context.Context, string) error
}

type loginService struct {
    repo repository.LoginRepository
    deps LoginDependencies
}
```

| 型・関数 | 意図 |
| --- | --- |
| LoginService | 呼び出す側に見せる「できることの一覧」 |
| loginService | 実際に処理する部品 |
| NewLoginService | 必要な部品を受け取り、使える状態に組み立てる関数 |
| LoginDependencies | 時刻、乱数、待機、パスワード照合を差し替える設定 |

interfaceは必要なメソッドの約束です。Goでは、そのメソッドを持つ型が約束を満たします。Handlerのテストでは、この約束を満たす偽物のServiceを渡せます。

先頭が大文字の名前は別のパッケージから参照できます。小文字のloginServiceはパッケージ内部の実装です。

```go
func (s *loginService) Login(...) (...)
```

`(s *loginService)` はメソッドのレシーバーです。この関数がloginServiceに属することを示します。sは変数名で、`s.repo` や `s.deps` から部品の持つ値を使えます。

### 関数をフィールドに持たせる理由

```go
Now func() time.Time
```

これは時刻を返す関数を入れられるフィールドです。本番ではtime.Now、テストでは固定日時を返す関数を入れられます。時刻や乱数を固定すると、テストで結果を予測できます。

NewLoginServiceは、指定されていない依存に本番用の処理を設定します。

### Loginの処理順序

1. 開始時刻を記録する。
2. 入力を検証する。
3. セッション用とCSRF用のトークンを生成する。
4. Repositoryへ渡す情報をまとめる。
5. Repositoryへログイン試行を依頼する。
6. 成功・認証失敗・回数制限に応じた結果を返す。

| 変数 | 意図 |
| --- | --- |
| startedAt | 処理時間の計測と日時の基準 |
| sessionID | セッションを識別するランダムな値 |
| csrfToken | CSRF検証に使う値 |
| attempt | DB処理に必要な入力・日時・トークンのまとまり |
| result.Status | Repositoryで確定したログイン結果 |

### 関数を引数として渡す

```go
result, err := s.repo.AttemptLogin(ctx, attempt, func(hash string) bool {
    return s.deps.ComparePassword(hash, req.Password)
})
```

RepositoryはDBからハッシュを取得し、Serviceは照合方法を持っています。そのため、取得したハッシュを渡せば照合できる関数をRepositoryへ渡しています。この実装では照合とDB更新を一つのトランザクション内で進めます。

waitForFailureDelayは失敗時の応答時間を一定範囲にそろえる処理です。既にかかった時間を引き、足りない分だけ待機します。randomLoginFailureDelayは目標の待機時間をランダムに決めます。これらはlogin特有の要件です。

RecordInvalidRequestはJSON構文不正などの要求をアクセスログへ記録する処理です。

## 4. Repositoryの読み方

対象: [repository/login.go](../backend/repository/login.go)

| 変数 | 意図 |
| --- | --- |
| r | Repository自身 |
| r.db | DBへの接続を管理する部品 |
| tx | 一連のDB操作をまとめるトランザクション |
| account | 取得したアカウント情報 |
| accountExists | アカウントが見つかったか |
| locked | 現在アカウントがロック中か |
| hash | 照合に使うパスワードハッシュ |
| passwordMatches | パスワードが一致したか |
| blockedUntil | IP制限の解除日時 |
| retry | 再試行までの残り秒数 |

| 関数 | 役割 |
| --- | --- |
| NewLoginRepository | DBを受け取りRepositoryを組み立てる |
| AttemptLogin | 制限確認、照合、状態更新、セッション作成、ログ記録 |
| selectIPLimit | IP単位の制限を取得 |
| selectLoginAccount | アカウントを取得 |
| nextAccountFailureCount | 次の失敗回数を計算 |
| upsertIPFailure | IP単位の失敗記録を追加・更新 |
| insertAttemptLogs | ログイン・アクセスの記録を保存 |
| RecordInvalidRequest | 入力不正のアクセスログを保存 |

DummyPasswordHashは、未登録アカウントでもハッシュ照合を行い、応答時間から登録の有無が推測されにくくするための値です。

### トランザクション

```text
BeginTx → 複数のDB操作 → Commit
                ↓ 失敗
             Rollback
```

Commitは変更の確定、Rollbackは変更の取り消しです。旧セッション削除と新セッション登録などをまとめ、途中までしか変更されない状態を避けます。

AttemptLoginの戻り値には名前があります。

```go
(_ *model.LoginAttemptResult, err error)
```

名前付き戻り値のerrを、関数終了時に実行されるdefer内でも確認し、エラーがあればRollbackしています。

sql.NullTimeはDBの日時がNULLかもしれない場合の型です。Validは値の有無、Timeは日時を表します。sql.ErrNoRowsは検索結果がないことを表し、DB接続失敗などとは区別します。

SQLの `$1` などは値を渡すためのプレースホルダーです。入力をSQL文字列へ直接連結する代わりに、ExecContextなどの引数から渡します。

## 5. modelとutilの役割

対象: [model/login.go](../backend/model/login.go)、[util/login.go](../backend/util/login.go)

| 型 | どこで使うか |
| --- | --- |
| LoginRequest | クライアントから受け取るJSON本文 |
| LoginMetadata | HTTP由来の補足情報 |
| LoginAttempt | ServiceからRepositoryへ渡す情報 |
| LoginAttemptResult | RepositoryからServiceへ返す結果 |
| LoginServiceResult | ServiceからHandlerへ返す結果とCookie情報 |
| LoginResponse | クライアントへ返すJSON本文 |
| LoginUser | 成功時に返すユーザー情報 |
| LoginStatus | 成功・認証失敗・回数制限を表す値 |

内部で必要な情報と、外部へ返す情報を分ける考え方です。HandlerがJSON本文として返すのはresult.Responseです。

Validateは必須項目とパスワード長を検証し、メールを正規化します。r.Email自体を書き換えるので入力の加工も行います。パスワードは加工しません。

NormalizeLoginEmailはメールの前後空白除去と小文字化、GenerateSecureTokenは乱数バイト列をURLで使える文字列へ変換する処理です。

## 6. logoutを作る前に仕様を分岐表にする

参照: [API仕様](design/api_design/02_auth.md)、[共通API仕様](design/api_design/01_overview.md)、[DB設計](design/database_design.md)

| 状態 | 処理 | 応答 |
| --- | --- | --- |
| セッションCookieなし | Cookie消去 | 本文なしの200 |
| セッション不明・期限切れ | Cookie消去 | 本文なしの200 |
| 有効セッションあり、CSRF欠落・不一致 | エラーを返す | 403 |
| 有効セッションあり、CSRF正常 | セッション削除、再認証失敗状態リセット、Cookie消去 | 本文なしの200 |
| DB操作失敗 | サーバーエラーを返す | 500 |

### 判断の順序

```text
セッションの有効性を確認
    ├─ 無効・なし → Cookie消去 → 200
    └─ 有効 → CSRF検証
                  ├─ 不正 → 403
                  └─ 正常 → DB更新 → Cookie消去 → 200
```

有効性を確認する前に一律のCSRF検証を行うと、未ログイン時も200を返す仕様を満たせません。未ログイン時に401で中断する認証ミドルウェアも、そのまま適用すると仕様に合いません。

DB検索失敗を「セッションなし」と扱ってはいけません。正常に検索して見つからない場合と、検索できなかった場合を分けます。

有効セッションの場合は、セッションからユーザーIDを取得します。クライアント指定のユーザーIDをそのまま更新対象にしません。

### DBでまとめて行う変更

- 現在の端末のLOGIN_SESSIONを削除する。
- 対象ユーザーのREAUTH_FAILED_COUNTを0にする。
- REAUTH_LAST_FAILED_ATをNULLにする。

これらを同じトランザクションで行うのが自然です。loginで使うLOGIN_FAILED_COUNTと、再認証用のREAUTH_FAILED_COUNTは別の項目です。

セッション確認と更新の間に状態が変わる可能性も考えます。トランザクションや必要なロックを使い、確認した対象と更新対象の整合性を保つ設計にします。

### CSRFとCookie

共通仕様はDouble Submit Cookie方式です。XSRF-TOKEN CookieとX-CSRF-Tokenヘッダーを読み、欠落や不一致を検証します。有効セッションがある場合だけ検証します。

Cookie削除にはGo特有の点があります。http.CookieのMaxAge: 0は削除指定になりません。削除にはMaxAge: -1を指定し、Goに削除用のMax-Age=0ヘッダーを出力させます。Cookie名とPathなどの設定は発行時に合わせます。

成功レスポンスは本文なしの200です。空のJSONオブジェクトを返す仕様ではありません。

## 7. 他のエンドポイントにも使える開発手順

### 手順1: 入出力と分岐を書き出す

URL、HTTPメソッド、認証条件、入力、成功応答、失敗応答、DB変更を整理します。最初にコードをコピーするのではなく、仕様から分岐表を作ります。

### 手順2: 判断の順番を文章で書く

logoutなら「セッション確認 → 無効なら200 → 有効ならCSRF確認 → DB更新 → Cookie削除」です。チェック項目だけでなく、順番も仕様の一部です。

### 手順3: 各処理の担当を決める

| logoutの処理 | 主な担当 |
| --- | --- |
| Cookieとヘッダーの読み取り | Handler |
| セッション状態とCSRFによる分岐 | Service |
| セッション検索とDB更新 | Repository |
| Cookie削除とHTTP応答 | Handler |

### 手順4: 層の間の入出力を決める

各関数が何を受け取り、何を返すかを先に決めます。logoutはJSON本文がないため、loginと同じ形のRequest型を無理に作る必要はありません。

出発点として、次のように責任を文章にします。

> Logout ServiceはセッションID・CSRF用Cookie・CSRFヘッダーを受け取り、セッションの状態に応じてDB更新を依頼し、成功またはエラーを返す。

その後、セッション検索に必要なRepositoryの関数と戻り値を考えます。時刻による期限判定があるなら、テストで時刻を固定できる設計も検討します。

### 手順5: 分岐表からテストケースを作る

| テスト対象 | 確認すること |
| --- | --- |
| Service | 無効セッション時の成功、有効時のCSRF検証、DB更新の呼び出し、エラー伝播 |
| Repository | 対象セッション検索、対象だけの削除、カウンター初期化、CommitとRollback |
| Handler | ステータス、Cookie削除ヘッダー、成功本文が空であること |
| Router | URLとHTTPメソッドが登録されていること |

[TESTING_GUIDE.md](../backend/TESTING_GUIDE.md)とリポジトリのAGENTS.mdに従います。GoDocの@specと日本語のテストケースを対応させ、即時中断が必要な前提確認にはrequire、値の比較にはassertを使います。

### 手順6: 実装し、Routerで組み立てる

[router/router.go](../backend/router/router.go)を参考にRepository → Service → Handlerを接続し、POST /api/auth/logoutを登録します。必要なSwaggerコメントや仕様との対応も確認します。

依存する部品を外から渡すことを依存性注入（DI）と呼びます。これにより各層を独立してテストできます。

### 手順7: 検証する

backendディレクトリで実行します。

```bash
gofmt -w <変更したGoファイル>
go test ./...
```

単にテストが通ることだけでなく、最初に作った分岐表の各行を確認できているかも見ます。

## 8. コードを読むときの確認項目

- この関数の入力と戻り値は何か。
- 値を読むだけか、引数やDBを変更するか。
- エラーはどこで判断し、どこでHTTP応答に変換するか。
- この変数は誰から受け取り、誰へ渡すか。
- この処理は共通の仕組みか、このAPI固有のルールか。
- 複数の変更をまとめて成功・失敗させる必要があるか。

変数名だけで意図が分からない場合は、値を設定する場所と使用する場所の両方を追います。loginを参考にするときも、役割や接続方法を参考にし、logoutに必要な処理はlogoutの仕様から決めます。
