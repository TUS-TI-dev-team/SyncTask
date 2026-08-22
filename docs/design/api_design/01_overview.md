# API 共通仕様書 (API Overview)

本ドキュメントでは、SyncTask のすべての REST API エンドポイントに共通する通信仕様、認証・認可アーキテクチャ、セキュリティ制御、共通ヘッダー、エラーハンドリング、パスワードバリデーション、およびエンドポイント一覧を定義します。

---

## 1. 共通通信仕様

### 1.1 プロトコル & データフォーマット
- **ベースURL**: `/api/`
- **プロトコル**: `HTTPS` (本番環境) / `HTTP` (ローカル開発環境)
- **データ形式**: `application/json` (リクエストボディおよびレスポンスボディ)
- **文字コード**: `UTF-8`
- **日時フォーマット**: ISO 8601 拡張形式・JSTタイムゾーン固定 (`YYYY-MM-DDTHH:mm:ss+09:00`)
- **タイムアウト**: 30秒
- **Cookie & セッション方式**:
  - セッション管理には Cookie ベースセッションを採用します（Cookie名: `sync_task_sid`）。
  - セッションCookie属性: `HttpOnly: true`, `SameSite: Lax`, `Path: /`, `Max-Age: 2592000` (30日間)。
  - `Secure`: 本番環境（HTTPS）では常時 `Secure: true`。ローカル開発環境（HTTP）では環境変数（`COOKIE_SECURE=false`）により切替可能とします。
  - OTP有効期限は発行から5分（300秒）、手続き全体の最大有効期限は15分です。
  - **OTPセッション破棄（戻る・キャンセル時）**: OTP入力画面からの離脱や「戻る」ボタン押下時は、クライアント側の `otp_session_id` を破棄するとともに、セッション破棄API（`POST auth/otp-session/cancel`）を呼び出してサーバー側 `OTP_SESSION` を即座に物理削除（無効化）します。

### 1.2 CORS (Cross-Origin Resource Sharing) & プリフライト仕様
フロントエンドとバックエンドが異なるオリジンで稼働する場合のクロスオリジン通信制御は以下の通り規定します：
- **許可オリジン (`Access-Control-Allow-Origin`)**: 環境変数 `FRONTEND_URL` で指定された単一オリジンを明示的に許可します（例: `http://localhost:3000`, `https://synctask.app`。ワイルドカード `*` は使用不可）。
- **認証情報許可 (`Access-Control-Allow-Credentials`)**: `true`（セッションCookieおよびCSRF Cookieの送受信を許可）。
- **許可メソッド (`Access-Control-Allow-Methods`)**: `GET, POST, PUT, PATCH, DELETE, OPTIONS`
- **許可ヘッダー (`Access-Control-Allow-Headers`)**: `Content-Type, X-CSRF-Token, Authorization`
- **公開ヘッダー (`Access-Control-Expose-Headers`)**: `Retry-After`（フロントエンドからの再試行待機時間読み取り用）
- **プリフライトキャッシュ (`Access-Control-Max-Age`)**: `86400` (24時間)
- **プリフライト（OPTIONS）レスポンス**: 条件に合致する OPTIONS リクエストに対しては `204 No Content` を返却します。

### 1.3 共通レスポンスヘッダー & セキュリティヘッダー
全APIレスポンスにおいて、以下のセキュリティヘッダーを出力します：
- `Content-Type: application/json; charset=utf-8`
- `X-Content-Type-Options: nosniff`（MIMEタイプスニッフィング防止）
- `X-Frame-Options: DENY`（クリックジャッキング防止）
- `Referrer-Policy: strict-origin-when-cross-origin`（リファラ情報の漏洩防止）
- `Strict-Transport-Security: max-age=31536000; includeSubDomains`（本番HTTPS環境のみ出力。HSTS有効化）

### 1.4 セキュリティ & CSRF・アカウント列挙対策
- **CSRF対策**:
  - Cookieベースの認証を行うため、**認証を必要とするすべての状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）** において CSRFトークン（`X-CSRF-Token` ヘッダー）の検証を必須とします（未認証のログイン・会員登録・パスワードリセット等のリクエスト、および `auth/otp-session/cancel` 等の未認証可能エンドポイントを除く）。
  - CSRFトークンは **Double Submit Cookie 方式** にて管理します。ログイン成功時（`auth/login`）およびアカウント新規登録完了時（`auth/register/verify-otp`）に、レスポンスヘッダーで `Set-Cookie: XSRF-TOKEN=<csrf_token>; SameSite=Lax; Path=/; Max-Age=2592000`（JavaScriptから読み取り可能な `HttpOnly=false`）を発行します。本番HTTPS環境では `Secure` 属性を付与、開発・テスト環境では付与しないものとします。
  - クライアント側 JavaScript で Cookie から CSRF トークンを取得し、状態変更リクエストの `X-CSRF-Token` ヘッダーに付与して送信します。
- **アカウント列挙防止 (User Enumeration 対策)**:
  - 新規登録（`auth/register/request-otp`）、パスワードリセット（`auth/password-reset/request-otp`）、メールアドレス変更（`auth/change-email/request-otp`）および各 OTP 再送（`resend-otp`）において、指定されたメールアドレスの登録有無、他ユーザーとの重複、または**他ユーザーの有効なOTPセッション期間中（手続き中）** の持続にかかわらず、**正常成功時・ダミー発行時を一貫して区別せずレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却**します。これにより、エラーコード・レスポンス構造・応答時間のあらゆる差異からメールアドレスの登録状況が推測されることを完全に防止します。
  - ダミーOTPセッションに対する後続の検証（`verify-otp`）や再送（`resend-otp`）に対しても、実セッションと全く同一のエラーコード（400, 410, 422, 429）および応答遅延（1.0s ± 0.1s）を適用します。
  - ログイン失敗時は、メールアドレス不一致・パスワード不一致・論理削除済みアカウント・アカウントロックアウト（同一メールアドレスへの連続失敗によるロック）のいずれも一律で `401 Unauthorized`（code: `UNAUTHORIZED`、遅延 1.0s ± 0.1s）を返却します（アカウントの登録有無およびロックアウト状況を攻撃者に推測させないため）。
- **メールアドレスマスク表示仕様**:
  - OTP送信先画面およびDB保存（`MASKED_EMAIL`）におけるメールアドレスマスク処理ルール：
    - ローカル部（`@` より前）が **4文字以上** の場合: 先頭4文字 ＋ 固定10文字のアスタリスク「`*`」 ＋ 「`@`」 ＋ ドメイン名（例: `user1234@example.com` → `user**********@example.com`）
    - ローカル部が **4文字未満（1〜3文字）** の場合: 先頭1文字 ＋ 固定10文字のアスタリスク「`*`」 ＋ 「`@`」 ＋ ドメイン名（例: `a@example.com` → `a**********@example.com`, `ab@example.com` → `a**********@example.com`, `abc@example.com` → `a**********@example.com`）。文字数の推測を完全に防止します。
- **認可制御 (IDOR / BOLA 対策)**:
  - ユーザー情報（`users/{user_id}`）およびタスク情報（`tasks/{task_id}`）へのアクセス・変更・削除時は、セッション内のログインユーザーIDとリソースの所有ユーザーIDの一致を厳格に検証します。
  - 他ユーザー所有のリソースまたは存在しないリソースへのアクセスに対しては、リソースの存在有無を秘匿するため一律 `404 Not Found` を返却します。
- **遅延制御 (Timing Attack 対策)**:
  - ログイン失敗（アカウントロックアウト中やIPレートリミット超過時429を含む）、OTP発行・再送処理（正常成功時およびアカウント存在有無のダミー処理を含む一括）、OTP検証失敗、パスワード再認証失敗時は、一律 `1.0s ± 0.1s` のレスポンス遅延を適用します。
- **レートリミット待機通知 (`Retry-After` ヘッダー)**:
  - `429 Too Many Requests`（IPレートリミット超過、OTP再送クールダウン中）返却時には、クライアントが再試行可能になるまでの待機時間を通知するため、標準の `Retry-After: <秒数>` レスポンスヘッダー（例: `Retry-After: 60`, `Retry-After: 900`）を必須で付与します。

### 1.5 共通エラーレスポンス構造
すべてのエラー応答は以下のJSONフォーマットで返却されます（HTTP ステータスコード `4xx` または `5xx`）：

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "入力内容に不備があります。",
    "details": [
      {
        "field": "email",
        "message": "メールアドレスの形式が正しくありません。"
      }
    ]
  }
}
```

#### エラーレスポンス スキーマ定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `error.code` | string | ○ | エラーコード（例: `BAD_REQUEST`, `UNAUTHORIZED`, `SAME_AS_CURRENT_PASSWORD` 等） |
| `error.message` | string | ○ | ユーザー向け汎用エラーメッセージ |
| `error.details` | array | ○ | フィールド単位のバリデーション詳細情報リスト。対象フィールドが存在しないエラー応答の場合は空配列 `[]` を返却（`null` やキー省略は不可） |
| `error.details[].field` | string | ○ | エラー対象のフィールド名またはクエリパラメータ名。リクエストデータがネストされたオブジェクト構造を持つ場合は `recurring_rule.due_time` のようにドット記法（`親オブジェクト.子フィールド`）で指定し、`GET` リクエスト等のクエリパラメータバリデーション違反時は対象のクエリパラメータ名（例: `priority`, `due_date`, `start_date`）を指定 |
| `error.details[].message` | string | ○ | フィールド固有のエラーメッセージ |

#### 代表的なエラーコード一覧

| HTTP Status | エラーコード (`code`) | 説明 |
| :--- | :--- | :--- |
| 400 | `BAD_REQUEST` | リクエスト形式またはバリデーション不正、OTP不一致、期間指定不正 |
| 401 | `UNAUTHORIZED` | 未ログイン、セッション無効・期限切れ、ログイン認証失敗、アカウントロックアウト |
| 401 | `REAUTH_FAILED` | アカウント削除/パスワード変更時の再認証失敗（1〜4回目、遅延 1.0s ± 0.1s） |
| 401 | `SESSION_DESTROYED` | アカウント削除/パスワード変更時の再認証失敗5回連続達成に伴うセッション強制破棄（遅延 1.0s ± 0.1s） |
| 403 | `FORBIDDEN` | CSRFトークン不正または権限不足、未検証OTPセッションでの更新試行、異ユーザー所有OTPセッション指定 |
| 404 | `NOT_FOUND` | 指定されたリソース（または他者所有リソース）が存在しない |
| 410 | `GONE` | OTPセッションの有効期限切れ（全体最大15分経過含む） |
| 410 | `OTP_SESSION_INVALIDATED` | メール送信5回連続失敗またはユーザーによるキャンセル・失効に伴うOTPセッション無効化・物理削除 |
| 422 | `UNPROCESSABLE_ENTITY` | 汎用ビジネスルール違反 |
| 422 | `INVALID_PASSWORD_CONTENT` | パスワード変更/リセット時のユーザー名・メールアドレスローカル部含有違反 |
| 422 | `SAME_AS_CURRENT_PASSWORD` | パスワード変更/リセット時の同一パスワード指定 |
| 422 | `SAME_AS_CURRENT_USERNAME` | ユーザー名変更時の同一ユーザー名指定 |
| 422 | `SAME_AS_CURRENT_EMAIL` | メールアドレス変更時の同一メールアドレス指定 |
| 422 | `OTP_REISSUED_DUE_TO_FAILURES` | OTP検証 5回連続失敗による自動再送・カウンターリセット |
| 429 | `RATE_LIMIT_EXCEEDED` | IPレートリミット超過（`Retry-After` ヘッダー付与） |
| 429 | `OTP_RESEND_COOLDOWN` | OTP再送クールダウン期間（60秒）中（`Retry-After` ヘッダー付与） |
| 500 | `INTERNAL_SERVER_ERROR` | サーバー内部エラー |
| 503 | `OTP_DELIVERY_FAILED` | 認証コードメール送信失敗（再送信可能状態の維持） |

### 1.6 パスワード共通バリデーション仕様
新規登録（`auth/register/request-otp`）、パスワードリセット（`auth/password-reset/reset`）、パスワード変更（`users/{user_id}/password`）におけるパスワードの複雑性要件は以下の通り共通で適用されます：
- **文字数**: 8文字以上128文字以下
- **文字種**: 以下の4種類中、**3種類以上** を含むこと
  1. 英大文字（`A-Z`）
  2. 英小文字（`a-z`）
  3. 数字（`0-9`）
  4. 許可記号（ASCII印字可能半角記号全32種類: ``! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~`` / 正規表現文字クラス: `[\x21-\x2f\x3a-\x40\x5b-\x60\x7b-\x7e]`）
- **禁止パターン**: ユーザー名およびメールアドレスのローカル部（`@` より前の文字列）が**4文字以上の場合**、大文字小文字を区別せず、それらと完全一致または部分文字列として含んではならない（違反時: `422 INVALID_PASSWORD_CONTENT`）。4文字未満の場合は本禁止パターンの判定対象外とします。

---

## 2. エンドポイント一覧

| カテゴリ | メソッド | URI | 役割・機能 | 認証要否 | CSRFヘッダー |
| :--- | :--- | :--- | :--- | :--- | :---: |
| **認証 (Auth)** | `POST` | `auth/register/request-otp` | 新規登録情報のバリデーション・OTP発行・メール送信 | 不要 | 不要 |
| | `POST` | `auth/register/verify-otp` | 新規登録OTP検証・アカウント本登録・セッション発行 | 不要 | 不要 |
| | `POST` | `auth/register/resend-otp` | 新規登録OTPの再送信 | 不要 | 不要 |
| | `POST` | `auth/login` | メールアドレス・パスワードによるログイン認証 | 不要 | 不要 |
| | `POST` | `auth/logout` | ログインセッションの破棄・ログアウト | 必須 | 必須 |
| | `POST` | `auth/password-reset/request-otp` | パスワードリセット用OTP発行・メール送信 | 不要 | 不要 |
| | `POST` | `auth/password-reset/verify-otp` | パスワードリセット用OTP検証 | 不要 | 不要 |
| | `POST` | `auth/password-reset/resend-otp` | パスワードリセット用OTPの再送信 | 不要 | 不要 |
| | `POST` | `auth/password-reset/reset` | 新パスワードの設定完了処理 | 不要 | 不要 |
| | `POST` | `auth/change-email/request-otp` | メールアドレス変更用OTP作成・送信 | 必須 | 必須 |
| | `POST` | `auth/change-email/verify-otp` | メールアドレス変更用OTP検証・変更確定 | 必須 | 必須 |
| | `POST` | `auth/change-email/resend-otp` | メールアドレス変更用OTPの再送信 | 必須 | 必須 |
| | `POST` | `auth/otp-session/cancel` | ユーザー操作（戻る/離脱）によるOTPセッションの即時破棄・無効化 | 不要 | 不要 |
| **ユーザー (Users)** | `GET` | `users/{user_id}` | ログインユーザーのプロフィール情報取得 | 必須 | 不要 |
| | `PUT` | `users/{user_id}` | プロフィール情報（ユーザー名等）の更新 | 必須 | 必須 |
| | `DELETE` | `users/{user_id}` | アカウント論理削除 | 必須 | 必須 |
| | `PATCH` | `users/{user_id}/password` | ログイン状態でのパスワード変更 | 必須 | 必須 |
| **タスク (Tasks)** | `GET` | `tasks` | タスク一覧取得（検索・絞り込み・カレンダー期間取得） | 必須 | 不要 |
| | `POST` | `tasks` | 新規タスク作成（単一作成 / 毎週繰り返し一括作成） | 必須 | 必須 |
| | `GET` | `tasks/{task_id}` | 単一タスクの詳細取得 | 必須 | 不要 |
| | `PATCH` | `tasks/{task_id}` | タスク情報の部分更新（特定フィールドのみの更新） | 必須 | 必須 |
| | `DELETE` | `tasks/{task_id}` | タスクの物理削除 | 必須 | 必須 |