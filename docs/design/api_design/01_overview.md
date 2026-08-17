# API Design (API設計書)

## 1. 概要・共通仕様

- **ベースURL**: `https://<domain>/api/`
- **通信形式**: JSON (HTTP REST API)
- **文字コード**: UTF-8
- **タイムゾーン / 日時フォーマット**:
  - 日時文字列: ISO 8601 拡張形式 / 日本標準時（例: `2026-08-17T12:00:00+09:00` または `YYYY-MM-DDTHH:mm:ss+09:00`）
  - 日付文字列: `YYYY-MM-DD`（カレンダー・締切日絞り込み時など）

### 1.1 セッション管理 & 認証方式
- **ログインセッション**:
  - トークンをリクエスト本文で送受信するのではなく、`HttpOnly`, `Secure`, `SameSite=Lax` 属性が付与されたセッションCookie（名称: `sync_task_sid`）によって管理します。
  - 認証が必要なAPIリクエストでは、ブラウザにより自動送信される Cookie からセッションを検証します。
  - セッション有効期限は 43,200分（1ヶ月）であり、APIアクセスごとに自動延長（Sliding Expiration）されます。
- **OTPセッション**:
  - アカウント新規作成、パスワードリセット、メールアドレス変更の手続き中は、手続きごとの `otp_session_id` をリクエストボディで送受信（または一時Cookie管理）します。
  - OTP有効期限は発行から5分（手続き全体の最大有効期限は15分）です。

### 1.2 セキュリティ & CSRF・アカウント列挙対策
- **CSRF対策**:
  - Cookieベースの認証を行うため、状態を変更するすべてのHTTPメソッド（`POST`, `PUT`, `PATCH`, `DELETE`）において CSRFトークンの検証を必須とします。
  - CSRFトークンは **Double Submit Cookie 方式** にて管理します。ログイン成功（`auth/login`）およびアカウント新規登録完了（`auth/register/verify-otp`）時に、レスポンスヘッダーで `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/`（JavaScriptから読み取り可能な `HttpOnly=false`）を発行します。
  - クライアントは JavaScript で Cookie から CSRF トークンを取得し、状態変更リクエストの `X-CSRF-Token` ヘッダーに付与して送信します。
- **アカウント列挙防止 (User Enumeration 対策)**:
  - 新規登録（`auth/register/request-otp`）、パスワードリセット（`auth/password-reset/request-otp`）、メールアドレス変更（`auth/change-email/request-otp`）において、指定されたメールアドレスの登録有無、他ユーザーとの重複、現在と同一メールアドレスの指定、または**他ユーザーの有効なOTPセッション期間中（手続き中）**の指定にかかわらず、**一貫してダミーOTPセッションを発行して `200 OK`（応答遅延 1.0s ± 0.1s）を返却**します。これにより、エラーコードや応答差分からメールアドレスの登録状況が推測されることを完全に防止します。
  - ダミーOTPセッションに対する後続の検証（`verify-otp`）や再送（`resend-otp`）に対しても、実セッションと全く同一のエラーコード（400, 410, 422, 429）および応答遅延（1.0s ± 0.1s）を適用します。
  - ログイン失敗時は、メールアドレス不一致・パスワード不一致・論理削除済みアカウントのいずれも一律で `401 Unauthorized`（code: `UNAUTHORIZED`、遅延 1.0s ± 0.1s）を返却します。
- **認可制御 (IDOR / BOLA 対策)**:
  - ユーザー情報（`/api/users/{user_id}`）およびタスク情報（`/api/tasks/{task_id}`）へのアクセス・変更・削除時は、セッション内のログインユーザーIDとリソースの所有ユーザーIDの一致を厳格に検証します。
  - 他ユーザー所有のリソースまたは存在しないリソースへのアクセスに対しては、リソースの存在有無を秘匿するため一律 `404 Not Found` を返却します。
- **遅延制御 (Timing Attack 対策)**:
  - ログイン失敗、OTP検証失敗、アカウント存在有無のダミー処理時は、一律 `1.0s ± 0.1s` のレスポンス遅延を適用します。

### 1.3 共通エラーレスポンス構造
すべてのエラー応答は以下のJSONフォーマットで返却されます（HTTP ステータスコード: `4xx` または `5xx`）。

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

#### 代表的なエラーコード一覧

| HTTP Status | エラーコード (`code`) | 説明 |
| :--- | :--- | :--- |
| 400 | `BAD_REQUEST` | リクエスト形式またはバリデーション不正、OTP不一致 |
| 401 | `UNAUTHORIZED` | 未ログイン、セッション無効・期限切れ、ログイン認証失敗、または再認証失敗（5回失敗によるセッション破棄含む） |
| 403 | `FORBIDDEN` | CSRFトークン不正または権限不足、未検証OTPセッションでの更新試行 |
| 404 | `NOT_FOUND` | 指定されたリソース（または他者所有リソース）が存在しない |
| 409 | `CONFLICT` | リソースの競合（※アカウント列挙防止のためメールアドレス重複等には使用しません） |
| 410 | `GONE` | OTPセッションの有効期限切れ（全体最大15分経過含む） |
| 422 | `UNPROCESSABLE_ENTITY` | ビジネスルール違反（パスワード変更時の同一パスワード再利用 `SAME_AS_CURRENT_PASSWORD`、同一ユーザー名への変更 `SAME_AS_CURRENT_USERNAME`、OTP 5回連続失敗による自動再送 `OTP_REISSUED_DUE_TO_FAILURES` 等） |
| 429 | `RATE_LIMIT_EXCEEDED` | 連続ログイン試行失敗（アカウントロック）、IPレートリミット超過、またはOTP再送クールダウン期間中 |
| 500 | `INTERNAL_SERVER_ERROR` | サーバー内部エラー |

---

## 2. エンドポイント一覧

| カテゴリ | メソッド | URI | 役割・機能 | 認証要否 |
| :--- | :--- | :--- | :--- | :---: |
| **認証 (Auth)** | `POST` | `auth/register/request-otp` | 新規登録情報のバリデーション・OTP発行・メール送信 | 不要 |
| | `POST` | `auth/register/verify-otp` | 新規登録OTP検証・アカウント本登録・セッション発行 | 不要 |
| | `POST` | `auth/register/resend-otp` | 新規登録OTPの再送信 | 不要 |
| | `POST` | `auth/login` | メールアドレス・パスワードによるログイン認証 | 不要 |
| | `POST` | `auth/logout` | ログインセッションの破棄・ログアウト | 必須 |
| | `POST` | `auth/password-reset/request-otp` | パスワードリセット用OTP発行・メール送信 | 不要 |
| | `POST` | `auth/password-reset/verify-otp` | パスワードリセット用OTP検証 | 不要 |
| | `POST` | `auth/password-reset/resend-otp` | パスワードリセット用OTPの再送信 | 不要 |
| | `POST` | `auth/password-reset/reset` | 新パスワードの設定完了処理 | 不要 |
| | `POST` | `auth/change-email/request-otp` | メールアドレス変更用OTP作成・送信 | 必須 |
| | `POST` | `auth/change-email/verify-otp` | メールアドレス変更用OTP検証・変更確定 | 必須 |
| | `POST` | `auth/change-email/resend-otp` | メールアドレス変更用OTPの再送信 | 必須 |
| **ユーザー (Users)** | `GET` | `users/{user_id}` | ログインユーザーのプロフィール情報取得 | 必須 |
| | `PUT` | `users/{user_id}` | プロフィール情報（ユーザー名等）の更新 | 必須 |
| | `DELETE` | `users/{user_id}` | アカウント論理削除 | 必須 |
| | `PATCH` | `users/{user_id}/password` | ログイン状態でのパスワード変更 | 必須 |
| **タスク (Tasks)** | `GET` | `tasks` | タスク一覧取得（検索・絞り込み・カレンダー期間取得・ページネーション） | 必須 |
| | `POST` | `tasks` | 新規タスク作成（単一作成 / 毎週繰り返し一括作成） | 必須 |
| | `GET` | `tasks/{task_id}` | 単一タスクの詳細取得 | 必須 |
| | `PATCH` | `tasks/{task_id}` | タスク情報の部分更新（特定フィールドのみの更新） | 必須 |
| | `DELETE` | `tasks/{task_id}` | タスクの物理削除 | 必須 |