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

---

## 3. エンドポイント詳細仕様

### 3.1 認証・アカウント登録 (Auth)

#### 3.1.1 `POST auth/register/request-otp`
新規作成のアカウント情報の入力検証を行い、仮登録セッションおよびOTPを生成してメールを送信します。

- **認証**: 不要

##### Request Body
```json
{
  "username": "exampleUser",
  "email": "user@example.com",
  "password": "Password123!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `username` | string | ○ | 2〜20文字、英数字（大文字小文字可）、前後の空白トリム |
| `email` | string | ○ | 有効なメールアドレス形式、前後の空白トリム、小文字正規化 |
| `password` | string | ○ | 8〜128文字、英大文字/英小文字/数字/記号のうち3種以上を含む。ユーザー名・メールのローカル部（4文字以上の場合）を含まないこと |

##### Response (200 OK)
```json
{
  "otp_session_id": "otp_sess_a1b2c3d4e5",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300
}
```
※既に登録済みのメールアドレス、または**他ユーザーのOTPセッション有効期間中（手続き中）**のメールアドレスが指定された場合も、メールアドレスの登録有無を秘匿するためダミーの `otp_session_id` とマスク文字列を返し、`200 OK`（遅延 1.0s ± 0.1s）を返却します。

##### Errors
- `400 Bad Request`: 入力バリデーション違反（文字数・形式違反等）

---

#### 3.1.2 `POST auth/register/verify-otp`
入力されたOTPを検証し、成功時にアカウントをDBへ本登録して新規ログインセッション（Cookie）およびCSRFトークンCookieを発行します。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_a1b2c3d4e5",
  "otp": "A1B2C3D4"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 発行されたOTPセッションID（ダミーセッションID含む） |
| `otp` | string | ○ | 英数字8桁（大文字・小文字不問） |

##### Response (201 Created)
- **Set-Cookie**: `sync_task_sid=<session_token>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000`
- **Set-Cookie**: `XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/`

```json
{
  "user": {
    "id": "usr_987654321",
    "username": "exampleUser",
    "email": "user@example.com"
  }
}
```
※検証成功後、自動ログイン処理を行います。なおリクエスト時に既存のログインセッションCookie（`sync_task_sid`）が送信された場合は、複数アカウントへの同時重複ログインを防止するため、その旧セッションをDBから物理削除した上で新しいセッションを発行します。

##### Errors
- `400 Bad Request`: OTP不一致（入力試行5回未満。ダミーセッション時も常に本エラーとなり遅延 1.0s ± 0.1s）
- `410 Gone`: OTPセッション有効期限切れ（全体最大15分超過含む）
- `422 Unprocessable Entity`: 5回連続失敗に伴う自動再送実行通知（`"code": "OTP_REISSUED_DUE_TO_FAILURES"`。ダミーセッション時も同一レスポンス）

---

#### 3.1.3 `POST auth/register/resend-otp`
新規登録用OTPを再生成してメールを再送信します（60秒クールダウン制約あり）。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_a1b2c3d4e5"
}
```

##### Response (200 OK)
```json
{
  "message": "OTP has been resent successfully.",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300
}
```
※ダミーセッションの場合も実際のメール送信は行わずに同様の `200 OK`（遅延 1.0s ± 0.1s）を返却します。

##### Errors
- `429 Too Many Requests`: クールダウン期間中（前回の送信から60秒未満）の再送要求
- `410 Gone`: 全体最大有効期限（初回発行から15分）切れ

---

#### 3.1.4 `POST auth/login`
メールアドレスとパスワードでログイン認証を行い、成功時にセッションCookieおよびCSRFトークンCookieを発行します。

- **認証**: 不要

##### Request Body
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `email` | string | ○ | トリム・小文字正規化 |
| `password` | string | ○ | 8〜128文字 |

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=<session_token>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000`
- **Set-Cookie**: `XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/`

```json
{
  "user": {
    "id": "usr_987654321",
    "username": "exampleUser",
    "email": "user@example.com"
  }
}
```

##### Errors
- `401 Unauthorized`: 認証失敗（メールアドレス未登録、パスワード不一致、論理削除済みアカウントのいずれも本エラーで一律返却。遅延 1.0s ± 0.1s）
- `429 Too Many Requests`: メールアドレス単位ロックアウト（直近15分間に5回連続失敗で30分ロック）またはIPレートリミット超過（直近5分間に30回失敗で15分遮断）。遅延 1.0s ± 0.1s

---

#### 3.1.5 `POST auth/logout`
現在操作中の端末のログインセッションをDBから物理削除し、Cookieを消去します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Max-Age=0`

##### Errors
- `401 Unauthorized`: 未ログイン

---

#### 3.1.6 `POST auth/password-reset/request-otp`
パスワードリセット用OTPを発行します。

- **認証**: 不要

##### Request Body
```json
{
  "email": "user@example.com"
}
```

##### Response (200 OK)
```json
{
  "otp_session_id": "otp_sess_reset_12345",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300
}
```

##### Errors
- `400 Bad Request`: 入力形式違反

---

#### 3.1.7 `POST auth/password-reset/verify-otp`
パスワードリセット用OTPを検証します。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_reset_12345",
  "otp": "A1B2C3D4"
}
```

##### Response (200 OK)
```json
{
  "message": "OTP verified successfully."
}
```

##### Errors
- `400 Bad Request`: OTP不一致
- `410 Gone`: 有効期限切れ

---

#### 3.1.8 `POST auth/password-reset/resend-otp`
パスワードリセット用OTPを再送します。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_reset_12345"
}
```

##### Response (200 OK)
```json
{
  "message": "OTP has been resent successfully."
}
```

---

#### 3.1.9 `POST auth/password-reset/reset`
パスワードをリセットします。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_reset_12345",
  "otp": "A1B2C3D4",
  "new_password": "NewPassword123!"
}
```

##### Response (200 OK)
```json
{
  "message": "Password has been reset successfully."
}
```

---

#### 3.1.10 `POST auth/change-email/request-otp`
メールアドレス変更のためのOTPを送信します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### Request Body
```json
{
  "new_email": "new_user@example.com"
}
```

##### Response (200 OK)
```json
{
  "otp_session_id": "otp_sess_chg_998877",
  "expires_in_seconds": 300
}
```

---

#### 3.1.11 `POST auth/change-email/verify-otp`
メールアドレス変更を確定させます。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### Request Body
```json
{
  "otp_session_id": "otp_sess_chg_998877",
  "otp": "A1B2C3D4"
}
```

##### Response (200 OK)
```json
{
  "message": "Email address has been updated successfully."
}
```

---

#### 3.1.12 `POST auth/change-email/resend-otp`
メールアドレス変更用OTPを再送信します（60秒クールダウン）。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### Request Body
```json
{
  "otp_session_id": "otp_sess_chg_998877"
}
```

##### Response (200 OK)
```json
{
  "message": "OTP has been resent successfully.",
  "masked_email": "new_**********@example.com",
  "expires_in_seconds": 300
}
```
※ダミーセッションの場合も同様に `200 OK`（遅延 1.0s ± 0.1s）を返却します。

##### Errors
- `429 Too Many Requests`: クールダウン期間中（60秒未満）
- `410 Gone`: 初回発行から15分経過

---

### 3.2 ユーザー管理 (Users)

#### 3.2.1 `GET users/{user_id}`
ログインユーザーのプロフィール情報（ユーザー名、メールアドレス等）を取得します。

- **認証**: 必須（Cookie）
- **Path Parameters**:
  - `user_id` (string): 取得対象のユーザーID（セッションと一致必須。不一致時は 404）

##### Response (200 OK)
```json
{
  "user": {
    "id": "usr_987654321",
    "username": "exampleUser",
    "email": "user@example.com",
    "created_at": "2026-08-01T10:00:00+09:00"
  }
}
```

##### Errors
- `401 Unauthorized`: 未ログイン
- `404 Not Found`: ユーザーが存在しない、または他ユーザーの `user_id` を指定した場合

---

#### 3.2.2 `PUT users/{user_id}`
プロフィール情報（ユーザー名）を更新します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### Request Body
```json
{
  "username": "newUsername"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `username` | string | ○ | 2〜20文字、英数字。現在のユーザー名と同一の場合は 422 エラー |

##### Response (200 OK)
```json
{
  "user": {
    "id": "usr_987654321",
    "username": "newUsername",
    "email": "user@example.com"
  }
}
```

##### Errors
- `400 Bad Request`: ユーザー名要件違反（文字数・使用可能文字不正）
- `422 Unprocessable Entity`: 現在のユーザー名と同一（`"code": "SAME_AS_CURRENT_USERNAME"`）
- `404 Not Found`: 認可エラー

---

#### 3.2.3 `DELETE users/{user_id}`
パスワード再認証を行い、アカウントを論理削除（`IS_DELETED=true`）します。所有タスクデータおよび全セッションは物理削除されます。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### Request Body
```json
{
  "password": "Password123!"
}
```

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Max-Age=0`

```json
{
  "message": "Account has been deleted successfully."
}
```

##### Errors
- `400 Bad Request`: パスワード再認証失敗（5回連続失敗時はセッション強制破棄・401 Unauthorized）
- `404 Not Found`: 認可エラー

---

#### 3.2.4 `PATCH users/{user_id}/password`
現在のパスワードを検証した上で、新しいパスワードへ変更し、全セッションを一括物理削除します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### Request Body
```json
{
  "current_password": "Password123!",
  "new_password": "NewSecurePassword456!"
}
```

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Max-Age=0`（再ログイン要求）

```json
{
  "message": "Password has been updated successfully. Please log in again."
}
```

##### Errors
- `400 Bad Request`: 新パスワード要件違反、または現在のパスワード不一致（5回連続失敗でセッション破棄・401）
- `422 Unprocessable Entity`: 新パスワードが現在のパスワードと同一（`"code": "SAME_AS_CURRENT_PASSWORD"`）
- `404 Not Found`: 認可エラー

---

### 3.3 タスク管理 (Tasks)

#### 3.3.1 `GET tasks`
タスク一覧を取得します。クエリパラメータにより、通常一覧、優先タスク・締切間近・ピン留めビュー、検索絞り込み、カレンダー表示用期間取得、ページネーションに対応します。

- **認証**: 必須（Cookie）

##### Query Parameters

| パラメータ名 | 型 | 必須 | デフォルト | 説明 |
| :--- | :--- | :---: | :--- | :--- |
| `page` | integer | × | `1` | ページ番号（1始まり） |
| `limit` | integer | × | `20` | 1ページあたりの取得件数（最大100件） |
| `view_type` | string | × | - | ビュー指定: `high_priority`（優先高）, `near_deadline`（72時間以内/期限超過）, `pinned`（ピン留めのみ） |
| `include_completed`| boolean | × | `false` | 完了タスクを含めるか（`true` / `false`） |
| `keyword` | string | × | - | タスク名およびコメントの部分一致検索（Case-Insensitive、トリム処理） |
| `priority` | string | × | - | 優先度絞り込み: `high`, `medium`, `low` |
| `status` | string | × | - | ステータス絞り込み: `not_started`, `in_progress`, `completed` |
| `due_date` | string | × | - | 締切日絞り込み（`YYYY-MM-DD`。指定日 23:59:59 までのタスクを検索） |
| `start_date` | string | × | - | カレンダー表示用: グリッド取得開始日（`YYYY-MM-DD`） |
| `end_date` | string | × | - | カレンダー表示用: グリッド取得終了日（`YYYY-MM-DD`） |
| `sort_by` | string | × | `default` | ソート種別（`default`: ピン留め優先→締切昇順→作成日時降順） |

##### Response (200 OK)
```json
{
  "items": [
    {
      "id": "tsk_1001",
      "user_id": "usr_987654321",
      "title": "課題レポート提出",
      "comment": "第5章の要約を含むこと\n参考文献を記載",
      "priority": "high",
      "status": "in_progress",
      "due_datetime": "2026-08-20T23:59:00+09:00",
      "is_pinned": true,
      "created_at": "2026-08-17T10:00:00+09:00",
      "updated_at": "2026-08-17T11:30:00+09:00"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_count": 45,
    "total_pages": 3
  }
}
```

##### Errors
- `400 Bad Request`: クエリパラメータ不正（日付フォーマット違反、limit超過等）
- `401 Unauthorized`: 未ログイン

---

#### 3.3.2 `POST tasks`
新規タスクを作成します。単一タスク作成に加え、期間と曜日を指定した毎週タスクの即時一括生成（最大100件）に対応します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### Request Body (1: 単一タスク作成時)
```json
{
  "title": "課題レポート提出",
  "comment": "第5章の要約を含むこと",
  "priority": "high",
  "due_datetime": "2026-08-20T23:59:00+09:00"
}
```

##### Request Body (2: 毎週繰り返し一括作成時)
```json
{
  "title": "週次ゼミ発表準備",
  "comment": "進捗スライド作成",
  "priority": "medium",
  "is_recurring": true,
  "recurring_rule": {
    "start_date": "2026-08-22",
    "end_date": "2026-10-31",
    "days_of_week": ["saturday"],
    "due_time": "18:00"
  }
}
```

##### Request Body フィールド定義

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `title` | string | ○ | 1〜100文字（トリム後）。改行・タブ等の制御文字禁止 |
| `comment` | string | × | 0〜1000文字（トリム後）。改行は `\n` に正規化 |
| `priority` | string | × | `high`, `medium`, `low`（デフォルト: `medium`） |
| `due_datetime` | string | × | ISO 8601 日時文字列。単一作成時用（省略時は当日 `23:59:00+09:00`）。※`is_recurring: true` 時は指定されていても無視されます |
| `is_recurring` | boolean | × | 繰り返し一括作成フラグ（デフォルト: `false`） |
| `recurring_rule` | object | △ | `is_recurring: true` 時のみ必須（`false` 時は無視） |
| `recurring_rule.start_date` | string | ○ | 開始日（`YYYY-MM-DD`）。`start_date <= end_date` |
| `recurring_rule.end_date` | string | ○ | 終了日（`YYYY-MM-DD`）。最大1年間（52週以内） |
| `recurring_rule.days_of_week` | array[string] | ○ | `["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]` より1つ以上選択 |
| `recurring_rule.due_time` | string | × | 締切時刻 `HH:mm`（省略時は `23:59`） |

※`is_recurring: true` の場合、生成件数が1〜100件の範囲内で即時一括生成されます（0件または101件以上の場合はエラーとなり作成されません）。

##### Response (201 Created)
```json
{
  "created_count": 10,
  "tasks": [
    {
      "id": "tsk_2001",
      "user_id": "usr_987654321",
      "title": "週次ゼミ発表準備",
      "comment": "進捗スライド作成",
      "priority": "medium",
      "status": "not_started",
      "due_datetime": "2026-08-22T18:00:00+09:00",
      "is_pinned": false,
      "created_at": "2026-08-17T12:00:00+09:00",
      "updated_at": "2026-08-17T12:00:00+09:00"
    }
  ]
}
```

##### Errors
- `400 Bad Request`: バリデーション不正（文字数違反、期間・曜日不整合等）
- `401 Unauthorized`: 未ログイン
- `403 Forbidden`: CSRFトークン不正

---

#### 3.3.3 `GET tasks/{task_id}`
指定されたタスクの詳細情報を取得します。

- **認証**: 必須（Cookie）
- **Path Parameters**:
  - `task_id` (string): タスクID

##### Response (200 OK)
```json
{
  "task": {
    "id": "tsk_1001",
    "user_id": "usr_987654321",
    "title": "課題レポート提出",
    "comment": "第5章の要約を含むこと",
    "priority": "high",
    "status": "in_progress",
    "due_datetime": "2026-08-20T23:59:00+09:00",
    "is_pinned": true,
    "created_at": "2026-08-17T10:00:00+09:00",
    "updated_at": "2026-08-17T11:30:00+09:00"
  }
}
```

##### Errors
- `401 Unauthorized`: 未ログイン
- `404 Not Found`: 存在しないタスクまたは他ユーザー所有タスク

---

#### 3.3.4 `PATCH tasks/{task_id}`
タスク情報を部分更新します。リクエストボディに含まれるフィールドのみが更新対象となります。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `task_id` (string): 更新対象タスクID

##### Request Body
```json
{
  "title": "課題レポート提出（修正版）",
  "comment": "参考文献の追記完了",
  "priority": "high",
  "status": "completed",
  "due_datetime": "2026-08-21T23:59:00+09:00",
  "is_pinned": true
}
```

##### Response (200 OK)
```json
{
  "task": {
    "id": "tsk_1001",
    "user_id": "usr_987654321",
    "title": "課題レポート提出（修正版）",
    "comment": "参考文献の追記完了",
    "priority": "high",
    "status": "completed",
    "due_datetime": "2026-08-21T23:59:00+09:00",
    "is_pinned": true,
    "created_at": "2026-08-17T10:00:00+09:00",
    "updated_at": "2026-08-17T13:00:00+09:00"
  }
}
```

##### Errors
- `400 Bad Request`: バリデーション不正（文字数違反、ステータス不正値等）
- `401 Unauthorized`: 未ログイン
- `403 Forbidden`: CSRFトークン不正
  "username": "newUsername"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `username` | string | ○ | 2〜20文字、英数字。現在のユーザー名と同一の場合は 422 エラー |

##### Response (200 OK)
```json
{
  "user": {
    "id": "usr_987654321",
    "username": "newUsername",
    "email": "user@example.com"
  }
}
```

##### Errors
- `400 Bad Request`: ユーザー名要件違反（文字数・使用可能文字不正）
- `422 Unprocessable Entity`: 現在のユーザー名と同一（`"code": "SAME_AS_CURRENT_USERNAME"`）
- `404 Not Found`: 認可エラー

---

#### 3.2.3 `DELETE users/{user_id}`
パスワード再認証を行い、アカウントを論理削除（`IS_DELETED=true`）します。所有タスクデータおよび全セッションは物理削除されます。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### Request Body
```json
{
  "password": "Password123!"
}
```

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Max-Age=0`

```json
{
  "message": "Account has been deleted successfully."
}
```

##### Errors
- `400 Bad Request`: パスワード再認証失敗（5回連続失敗時はセッション強制破棄・401 Unauthorized）
- `404 Not Found`: 認可エラー

---

#### 3.2.4 `PATCH users/{user_id}/password`
現在のパスワードを検証した上で、新しいパスワードへ変更し、全セッションを一括物理削除します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### Request Body
```json
{
  "current_password": "Password123!",
  "new_password": "NewSecurePassword456!"
}
```

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Max-Age=0`（再ログイン要求）

```json
{
  "message": "Password has been updated successfully. Please log in again."
}
```

##### Errors
- `400 Bad Request`: 新パスワード要件違反、または現在のパスワード不一致（5回連続失敗でセッション破棄・401）
- `422 Unprocessable Entity`: 新パスワードが現在のパスワードと同一（`"code": "SAME_AS_CURRENT_PASSWORD"`）
- `404 Not Found`: 認可エラー

---

### 3.3 タスク管理 (Tasks)

#### 3.3.1 `GET tasks`
タスク一覧を取得します。クエリパラメータにより、通常一覧、優先タスク・締切間近・ピン留めビュー、検索絞り込み、カレンダー表示用期間取得、ページネーションに対応します。

- **認証**: 必須（Cookie）

##### Query Parameters

| パラメータ名 | 型 | 必須 | デフォルト | 説明 |
| :--- | :--- | :---: | :--- | :--- |
| `page` | integer | × | `1` | ページ番号（1始まり） |
| `limit` | integer | × | `20` | 1ページあたりの取得件数（最大100件） |
| `view_type` | string | × | - | ビュー指定: `high_priority`（優先高）, `near_deadline`（72時間以内/期限超過）, `pinned`（ピン留めのみ） |
| `include_completed`| boolean | × | `false` | 完了タスクを含めるか（`true` / `false`） |
| `keyword` | string | × | - | タスク名およびコメントの部分一致検索（Case-Insensitive、トリム処理） |
| `priority` | string | × | - | 優先度絞り込み: `high`, `medium`, `low` |
| `status` | string | × | - | ステータス絞り込み: `not_started`, `in_progress`, `completed` |
| `due_date` | string | × | - | 締切日絞り込み（`YYYY-MM-DD`。指定日 23:59:59 までのタスクを検索） |
| `start_date` | string | × | - | カレンダー表示用: グリッド取得開始日（`YYYY-MM-DD`） |
| `end_date` | string | × | - | カレンダー表示用: グリッド取得終了日（`YYYY-MM-DD`） |
| `sort_by` | string | × | `default` | ソート種別（`default`: ピン留め優先→締切昇順→作成日時降順） |

##### Response (200 OK)
```json
{
  "items": [
    {
      "id": "tsk_1001",
      "user_id": "usr_987654321",
      "title": "課題レポート提出",
      "comment": "第5章の要約を含むこと\n参考文献を記載",
      "priority": "high",
      "status": "in_progress",
      "due_datetime": "2026-08-20T23:59:00+09:00",
      "is_pinned": true,
      "created_at": "2026-08-17T10:00:00+09:00",
      "updated_at": "2026-08-17T11:30:00+09:00"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_count": 45,
    "total_pages": 3
  }
}
```

##### Errors
- `400 Bad Request`: クエリパラメータ不正（日付フォーマット違反、limit超過等）
- `401 Unauthorized`: 未ログイン

---

#### 3.3.2 `POST tasks`
新規タスクを作成します。単一タスク作成に加え、期間と曜日を指定した毎週タスクの即時一括生成（最大100件）に対応します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### Request Body (1: 単一タスク作成時)
```json
{
  "title": "課題レポート提出",
  "comment": "第5章の要約を含むこと",
  "priority": "high",
  "due_datetime": "2026-08-20T23:59:00+09:00"
}
```

##### Request Body (2: 毎週繰り返し一括作成時)
```json
{
  "title": "週次ゼミ発表準備",
  "comment": "進捗スライド作成",
  "priority": "medium",
  "is_recurring": true,
  "recurring_rule": {
    "start_date": "2026-08-22",
    "end_date": "2026-10-31",
    "days_of_week": ["saturday"],
    "due_time": "18:00"
  }
}
```

##### Request Body フィールド定義

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `title` | string | ○ | 1〜100文字（トリム後）。改行・タブ等の制御文字禁止 |
| `comment` | string | × | 0〜1000文字（トリム後）。改行は `\n` に正規化 |
| `priority` | string | × | `high`, `medium`, `low`（デフォルト: `medium`） |
| `due_datetime` | string | × | ISO 8601 日時文字列。単一作成時用（省略時は当日 `23:59:00+09:00`）。※`is_recurring: true` 時は指定されていても無視されます |
| `is_recurring` | boolean | × | 繰り返し一括作成フラグ（デフォルト: `false`） |
| `recurring_rule` | object | △ | `is_recurring: true` 時のみ必須（`false` 時は無視） |
| `recurring_rule.start_date` | string | ○ | 開始日（`YYYY-MM-DD`）。`start_date <= end_date` |
| `recurring_rule.end_date` | string | ○ | 終了日（`YYYY-MM-DD`）。最大1年間（52週以内） |
| `recurring_rule.days_of_week` | array[string] | ○ | `["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]` より1つ以上選択 |
| `recurring_rule.due_time` | string | × | 締切時刻 `HH:mm`（省略時は `23:59`） |

※`is_recurring: true` の場合、生成件数が1〜100件の範囲内で即時一括生成されます（0件または101件以上の場合はエラーとなり作成されません）。

##### Response (201 Created)
```json
{
  "created_count": 10,
  "tasks": [
    {
      "id": "tsk_2001",
      "user_id": "usr_987654321",
      "title": "週次ゼミ発表準備",
      "comment": "進捗スライド作成",
      "priority": "medium",
      "status": "not_started",
      "due_datetime": "2026-08-22T18:00:00+09:00",
      "is_pinned": false,
      "created_at": "2026-08-17T12:00:00+09:00",
      "updated_at": "2026-08-17T12:00:00+09:00"
    }
  ]
}
```
※単一タスク作成時も `created_count: 1` および要素数1の `tasks` 配列を返却します。

##### Errors
- `400 Bad Request`: タイトル文字数違反、期間・曜日不整合、生成件数超過（0件または101件以上）等
- `401 Unauthorized`: 未ログイン
- `403 Forbidden`: CSRFトークン不正

---

#### 3.3.3 `GET tasks/{task_id}`
指定されたタスクの詳細情報を取得します。

- **認証**: 必須（Cookie）
- **Path Parameters**:
  - `task_id` (string): タスクID

##### Response (200 OK)
```json
{
  "task": {
    "id": "tsk_1001",
    "user_id": "usr_987654321",
    "title": "課題レポート提出",
    "comment": "第5章の要約を含むこと",
    "priority": "high",
    "status": "in_progress",
    "due_datetime": "2026-08-20T23:59:00+09:00",
    "is_pinned": true,
    "created_at": "2026-08-17T10:00:00+09:00",
    "updated_at": "2026-08-17T11:30:00+09:00"
  }
}
```

##### Errors
- `401 Unauthorized`: 未ログイン
- `404 Not Found`: 存在しないタスクまたは他ユーザー所有タスク

---

#### 3.3.4 `PUT tasks/{task_id}`
タスク情報を更新します。全項目の一括更新に加え、リクエストボディに含まれるフィールドのみを部分更新（ステータス変更 `status` やピン留め `is_pinned` の単体更新含む）する実質的な部分更新兼用仕様とします。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `task_id` (string): 更新対象タスクID

##### Request Body
```json
{
  "title": "課題レポート提出（修正版）",
  "comment": "参考文献の追記完了",
  "priority": "high",
  "status": "completed",
  "due_datetime": "2026-08-21T23:59:00+09:00",
  "is_pinned": true
}
```
※ステータスのみを変更する場合は `{"status": "completed"}`、ピン留めのみを変更する場合は `{"is_pinned": true}` のように、更新対象のフィールドのみを指定して送信可能です。

##### Request Body フィールド定義

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `title` | string | × | 1〜100文字（トリム後）。改行等の制御文字禁止 |
| `comment` | string | × | 0〜1000文字（トリム後）。改行は `\n` に正規化 |
| `priority` | string | × | `high`, `medium`, `low` |
| `status` | string | × | `not_started`, `in_progress`, `completed` |
| `due_datetime` | string / null | × | ISO 8601 日時文字列（`null` 指定で締切解除） |
| `is_pinned` | boolean | × | ピン留め状態（`true` / `false`） |

##### Response (200 OK)
```json
{
  "task": {
    "id": "tsk_1001",
    "user_id": "usr_987654321",
    "title": "課題レポート提出（修正版）",
    "comment": "参考文献の追記完了",
    "priority": "high",
    "status": "completed",
    "due_datetime": "2026-08-21T23:59:00+09:00",
    "is_pinned": true,
    "created_at": "2026-08-17T10:00:00+09:00",
    "updated_at": "2026-08-17T13:00:00+09:00"
  }
}
```

##### Errors
- `400 Bad Request`: バリデーション不正（文字数違反、ステータス不正値等）
- `401 Unauthorized`: 未ログイン
- `403 Forbidden`: CSRFトークン不正
- `404 Not Found`: 認可エラー（存在しないタスクまたは他者所有タスク）

---

#### 3.3.5 `DELETE tasks/{task_id}`
タスクをDBから物理削除します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `task_id` (string): 削除対象タスクID

##### Response (200 OK)
```json
{
  "message": "Task has been deleted successfully."
}
```

##### Errors
- `401 Unauthorized`: 未ログイン
- `403 Forbidden`: CSRFトークン不正
- `404 Not Found`: 認可エラー（存在しないタスクまたは他者所有タスク）
