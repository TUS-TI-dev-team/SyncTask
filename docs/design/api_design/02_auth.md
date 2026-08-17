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