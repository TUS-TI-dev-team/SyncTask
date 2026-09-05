# POST auth/register/request-otp 開発計画

## 1. 目的と完了条件

`docs/design/` の仕様に基づき、未認証で利用できる `POST /api/auth/register/request-otp` をテスト先行（TDD）で実装する。

完了条件は次のすべてを満たすこととする。

- `backend/TESTING_GUIDE.md` に従った単体テストを作成している（日本語テスト名、Code-as-Docs `@spec`、`require`/`assert` の使い分け）。
- `go test ./...` がすべて成功する。
- `go test -cover ./...` で、今回追加・変更するパッケージ（`handler`, `service`, `repository`, `model`, `util`, `router`）のカバレッジが 80% 以上である。
- `go test -race ./...` が成功する。
- APIレスポンス、DB登録（`OTP_SESSION`, `MAIL_AUTH_LOG`, `ACCESS_LOG`）、メール送信連携、Timing Attack 対策遅延（1.0s ± 0.1s）、アカウント列挙防止が設計仕様と一致する。

## 2. 参照仕様

- `docs/design/api_design/01_overview.md`
  - ベースURL `/api/`、共通エラー形式、列挙防止、レスポンス遅延（1.0s ± 0.1s）、パスワード共通バリデーション仕様（1.6節）。
- `docs/design/api_design/02_auth.md` 3.1.1
  - 入出力、評価順序、HTTPステータス（200 OK, 400 Bad Request, 503 Service Unavailable）。
- `docs/design/process_design/01_account_creation.md` 1.1, 1.2
  - 入力正規化、通常処理・ダミー処理・競合処理・排他制御、ログ記録（`MAIL_AUTH_LOG`, `ACCESS_LOG`）、遅延要件、メール送信失敗時の補償。
- `docs/design/database_design.md`
  - `OTP_SESSION`、`LOGIN_ACCOUNT`、`MAIL_AUTH_LOG`、`ACCESS_LOG` の定義および制約。
- `backend/TESTING_GUIDE.md`
  - テスト配置・命名規約、Code-as-Docs、アサーション規約、独立性。

## 3. 設計判断とアーキテクチャ

既存の `handler`, `service`, `repository`, `model`, `util`, `router` のレイヤードアーキテクチャおよびDI構造を踏襲する。

### 3.1 パッケージ構成
- `backend/model/register_request_otp.go`:
  - リクエスト構造体: `RegisterRequestOtpRequest` (`username`, `email`, `password`)
  - レスポンス構造体: `RegisterRequestOtpResponse` (`otp_session_id`, `masked_email`, `expires_in_seconds`, `cooldown_seconds`)
  - 入力バリデーション関数 `Validate()`:
    - ユーザー名: トリム後 2〜20文字、半角英数字のみ (`^[a-zA-Z0-9]+$`)
    - メールアドレス: トリム & 小文字化、RFC 5322準拠形式、255文字以下
    - パスワード: トリムなし 8〜128文字、英大文字・英小文字・数字・記号（全32種）のうち3種以上、4文字以上のユーザー名またはメールローカル部（大文字小文字不問）を含まない
    - 不備時は `model.NewBadRequestError("BAD_REQUEST", ...)` を生成
- `backend/util/register.go`:
  - `MaskEmail(email string) string`: 先頭4文字（ローカル部が4文字未満の場合は先頭1文字）とドメイン以外を固定10文字の `*`（`**********`）でマスク
  - `GenerateOTPSessionID() (string, error)`: `otp_sess_` + URL-safe ランダム文字列
  - `GenerateOTP() (string, error)`: 8桁の英数字（ランダム）
- `backend/service/mailer.go`:
  - `Mailer` インターフェース: `SendOTP(ctx context.Context, toEmail, otp string) error`
  - テスト容易性のためモック可能にする。
- `backend/repository/register_request_otp.go`:
  - インターフェース `RegisterRequestOtpRepository`:
    - `FindActiveUserByEmail(ctx context.Context, email string) (bool, error)`: 有効アカウント（`IS_DELETED=false`）存在確認
    - `FindActiveOtpSessionByEmail(ctx context.Context, email string) (bool, error)`: 有効OTPセッション（`STATUS IN ('active', 'verified')` かつ `SESSION_EXPIRES_AT > NOW()`）確認
    - `CreateOtpSession(ctx context.Context, session *model.OtpSessionRecord) error`: `OTP_SESSION` 挿入
    - `UpdateOtpSessionDeliveryStatus(ctx context.Context, sessionID, status string, sendFailedCount int) error`: 送信状態・失敗回数の更新
    - `RecordMailAuthLog(ctx context.Context, log *model.MailAuthLogRecord) error`: `MAIL_AUTH_LOG` 記録
    - `RecordAccessLog(ctx context.Context, log *model.AccessLogRecord) error`: `ACCESS_LOG` 記録
- `backend/service/register_request_otp.go`:
  - `RegisterRequestOtpService`:
    - 依存: Repository, Mailer, Clock, Sleeper, PasswordHasher, TokenGenerator
    - 処理フロー:
      1. 入力正規化・バリデーション（違反時は遅延なし 400 Bad Request）
      2. 既存アカウント・既存有効OTPセッションの照会
      3. 通常処理 vs ダミー処理 の分岐
         - 通常処理: OTP生成、セッションID生成、bcryptハッシュ生成、`OTP_SESSION` (IS_DUMMY=false) 作成、ログ記録、コミット
         - コミット後メール送信: 失敗時は `DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT=1` に更新し 503 `OTP_DELIVERY_FAILED` 返却
         - ダミー処理: ダミーセッションID生成、`OTP_SESSION` (IS_DUMMY=true, PENDING系はNULL) 作成、ログ記録、メール送信スキップ
      4. タイミング攻撃対策遅延（1.0s ± 0.1s）を適用（通常成功時・ダミー時）
      5. レスポンス返却
- `backend/handler/register_request_otp.go`:
  - JSONバインド、ClientIP取得、Service呼び出し
  - レスポンスヘッダー設定（`Content-Type`, `Cache-Control: no-store, no-cache, must-revalidate`, `Pragma: no-cache`）
  - 共通エラー形式（`model.AppError`）へのマッピング
- `backend/router/router.go`:
  - `api.POST("/auth/register/request-otp", ...)` の追加と DI 構成

---

## 4. ステップ別作業計画

### Step 1: テストデータ・単体テストコード作成
- `backend/model/register_request_otp_test.go`:
  - 入力検証（正常系、各種バリデーションエラー境界値）
- `backend/util/register_test.go`:
  - メールマスク処理、OTPセッションID生成、OTP生成の単体テスト
- `backend/repository/register_request_otp_test.go`:
  - `go-sqlmock` を用いたDB操作テスト（通常セッション登録、ダミー登録、競合ロールバック、ログ記録）
- `backend/service/register_request_otp_test.go`:
  - 手動モックを用いたServiceテスト（通常成功・メール送信、ダミー成功、メール送信失敗503、遅延呼出確認）
- `backend/handler/register_request_otp_test.go`:
  - Gin TestModeを用いたHandlerテスト（200 OK, 400 Bad Request, 503 Service Unavailable, ヘッダー検証）
- `backend/router/router_test.go`:
  - エンドポイント登録テスト

### Step 2: プログラム本体実装
- `backend/model/register_request_otp.go`
- `backend/util/register.go`
- `backend/service/mailer.go`
- `backend/repository/register_request_otp.go`
- `backend/service/register_request_otp.go`
- `backend/handler/register_request_otp.go`
- `backend/router/router.go`

### Step 3: 単体テスト実行とカバレッジ確認
- `go test -v ./...` で全テスト通過を確認
- `go test -cover ./...` で今回関係する各パッケージ 80% 以上を確認
- `go test -race ./...` でデータ競合のないことを確認

### Step 4: プログラム修正（失敗時）
- テスト失敗やカバレッジ不足があれば修正し、Step 3 を再実行
