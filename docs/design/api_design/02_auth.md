### 3.1 認証・アカウント登録 (Auth)

#### OTP API共通規則

- クライアントが同じ用途の有効な `otp_session_id` を保持している場合はOTP入力画面へ復帰し、`request-otp` を呼び出しません。それでも `request-otp` が呼ばれ、対象メールアドレスまたは認証ユーザーに有効なOTPセッションが存在する場合、既存セッションを更新せず、新しいダミーセッションを作成して同一形式の `200 OK` を返します。
- ダミーセッションは `OTP_SESSION.IS_DUMMY=true` を最終判定とします。`SIGNUP` と `PASSWORD_RESET` のダミーは `USER_ID` をNULL、`EMAIL_CHANGE` のダミーは所有者認可のため認証ユーザーIDを保持します。`PENDING_USERNAME`、`PENDING_EMAIL`、`PENDING_PASSWORD_HASH`、`OTP_HASH` はNULLとし、OTP照合を成功させません（メールアドレス表示には `MASKED_EMAIL` を使用）。
- 実メール送信に失敗した場合は `DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT+=1` とし、`503 Service Unavailable`（code: `OTP_DELIVERY_FAILED`）を返します（`error.details` は空配列 `[]`）。失敗送信には60秒クールダウンを適用せず、OTP再送時は同一の `otp_session_id` を維持して既存レコードを直接更新（`otp_session_id` の変更は行わない）して再送処理を行います。
- 送信成功時は `DELIVERY_STATUS='sent'`、`SEND_FAILED_COUNT=0` とします。初回送信、手動再送、5回照合失敗による自動再送を通じて5回連続で送信に失敗した場合、対象OTPセッションを物理削除し、`410 Gone`（code: `OTP_SESSION_INVALIDATED`）を返します。

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
| `password` | string | ○ | 8〜128文字、英大文字/英小文字/数字/記号（全32種）のうち3種以上を含む（01_overview.md 1.4節準拠）。ユーザー名・メールのローカル部（4文字以上の場合、大文字小文字を区別せず比較）を含まないこと |

##### Response (200 OK)
```json
{
  "otp_session_id": "otp_sess_a1b2c3d4e5",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300,
  "cooldown_seconds": 60
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 生成されたOTPセッションID（例: `otp_sess_a1b2c3d4e5`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス（例: `user**********@example.com`） |
| `expires_in_seconds` | integer | ○ | OTPの有効期限（秒、デフォルト: 300） |
| `cooldown_seconds` | integer | ○ | 再送可能になるまでのクールダウン秒数（60秒） |

※アカウント列挙防止および Timing Attack 対策として、正常成功時（実メール送信時）およびダミー発行時（既に登録済みのメールアドレス、他ユーザーの有効なOTPセッション期間中等の指定時）を一貫して区別せず、一律でレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`username`, `email`, `password`）の有無、文字数・形式制約を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし）。
2. **既存OTPセッション検証・列挙対策 (`200 OK`)**:
   指定メールに有効なOTPセッションがある場合は、経過時間にかかわらず既存レコードを更新せず、ダミーセッションを新規作成します。未登録かつ排他なしの場合だけ実OTPを発行し、通常・登録済み・予約中・重複要求のすべてで `1.0s ± 0.1s` 後に同一形式の `200 OK` を返します。

##### Errors
- `400 Bad Request`: 入力バリデーション違反（文字数・形式違反等、code: `"BAD_REQUEST"`）
- `503 Service Unavailable`: 実メール送信失敗（code: `"OTP_DELIVERY_FAILED"`）

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
- **Set-Cookie**: `XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000`

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "exampleUser",
    "email": "user@example.com",
    "created_at": "2026-08-17T12:00:00+09:00",
    "updated_at": "2026-08-17T12:00:00+09:00"
  }
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `user` | object | ○ | 本登録完了したユーザーオブジェクト |
| `user.id` | string | ○ | ユーザーID（UUID形式、例: `550e8400-e29b-41d4-a716-446655440000`） |
| `user.username` | string | ○ | ユーザー名（例: `exampleUser`） |
| `user.email` | string | ○ | メールアドレス（例: `user@example.com`） |
| `user.created_at` | string | ○ | アカウント登録日時（ISO 8601 JST 形式） |
| `user.updated_at` | string | ○ | アカウント更新日時（ISO 8601 JST 形式） |

※検証成功後、自動ログイン処理を行います。本登録完了と同時に、使用済みの手続き用OTPセッション（`OTP_SESSION`）をDBから直ちに物理削除します。なおリクエスト時に既存のログインセッションCookie（`sync_task_sid`）が送信された場合は、複数アカウントへの同時重複ログインを防止するため、その旧セッション（`LOGIN_SESSION`）もDBから物理削除した上で新しいログインセッションを発行します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`, `otp`）の有無、および `otp` の形式（英数字8桁）を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし、試行回数 `ATTEMPT_COUNT` 加算なし）。
2. **OTPセッション状態・目的・期限検証 (`400 Bad Request` / `410 Gone`)**:
   指定された `otp_session_id` の存在、用途 `PURPOSE` が新規登録（`SIGNUP`）であること、ステータス（`active` であること）、および有効期限（`OTP_EXPIRES_AT` / `SESSION_EXPIRES_AT`）を検証します。セッション不在・`PURPOSE`不一致・失効・既に検証済み等の非 `active` ステータスの場合は、Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"`）を返却します（ダミーセッション時含む）。
3. **OTP照合検証 (`400 BAD_REQUEST` / `422 OTP_REISSUED_DUE_TO_FAILURES`)**:
   入力された `otp` のハッシュ照合を実施します。
   - 不一致（試行1〜4回目）: 失敗回数（`ATTEMPT_COUNT`）を+1加算し、`400 Bad Request`（code: `"BAD_REQUEST"`、遅延 1.0s ± 0.1s）を返却します。
   - 不一致（試行5回達成）: 失敗回数をリセットし、新OTPを自動再発行・送信します。
     - 実メール送信に成功した場合（ダミーセッション含む）: `422 Unprocessable Entity`（code: `"OTP_REISSUED_DUE_TO_FAILURES"`、遅延 1.0s ± 0.1s）を返却します。
     - 自動再送における実メール送信に失敗した場合（1〜4回目の送信失敗）: `DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT+=1` とし、`503 Service Unavailable`（code: `"OTP_DELIVERY_FAILED"`）を返却します。
     - 自動再送を含めて5回連続送信失敗となった場合: 対象セッションを物理削除し、`410 Gone`（code: `"OTP_SESSION_INVALIDATED"`）を返却します。

##### Errors
- `400 Bad Request`: 入力形式違反（即時返却、遅延なし）または無効なセッション/PURPOSE不一致/STATUS非active指定またはOTP照合不一致（試行1〜4回目。ダミーセッション時も一律遅延 1.0s ± 0.1s、code: `"BAD_REQUEST"`）
- `410 Gone`: OTPセッション有効期限切れ（全体最大15分超過含む、code: `"GONE"`）、またはメール送信5回連続失敗に伴うセッション失効（code: `"OTP_SESSION_INVALIDATED"`）
- `422 Unprocessable Entity`: 5回連続失敗に伴う自動再送実行通知（応答遅延 1.0s ± 0.1s、code: `"OTP_REISSUED_DUE_TO_FAILURES"`。ダミーセッション時も実際のメール再送を行わずに全く同一のレスポンスを返却）
- `503 Service Unavailable`: 5回目不一致時の自動再送処理におけるメール送信失敗（再送可能状態を維持。code: `"OTP_DELIVERY_FAILED"`）

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

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 発行されたOTPセッションID |

##### Response (200 OK)
```json
{
  "message": "OTP has been resent successfully.",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300,
  "cooldown_seconds": 60
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理結果メッセージ（例: `"OTP has been resent successfully."`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス（例: `user**********@example.com`） |
| `expires_in_seconds` | integer | ○ | 再発行されたOTPの有効期限（秒、デフォルト: 300） |
| `cooldown_seconds` | integer | ○ | 再送可能になるまでのクールダウン秒数（60秒） |

※Timing Attack 対策として、正常再送時（実メール送信時）およびダミーセッション再送時（実際のメール送信を行わない場合）を一貫して区別せず、一律でレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却します。また、リクエスト対象の `OTP_SESSION` の `PURPOSE` が `SIGNUP` かつステータスが `active` であることを必須で検証します。再送処理成功時、対象の `OTP_SESSION` レコードにおいて新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、試行失敗回数（`ATTEMPT_COUNT`）を 0 にリセット、送信回数（`SEND_COUNT`）を +1 加算、直前送信日時（`LAST_SENT_AT`）を更新するとともに、有効期限（`OTP_EXPIRES_AT`）を再送信時点から5分間（全体最大有効期限 `SESSION_EXPIRES_AT` の範囲内）へ更新延長します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`）の有無を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし）。
2. **OTPセッション状態・目的・期限検証 (`400 Bad Request` / `410 Gone`)**:
   指定された `otp_session_id` の存在、用途 `PURPOSE` が新規登録（`SIGNUP`）であること、ステータス（`active` であること）、および全体最大有効期限（`SESSION_EXPIRES_AT` 15分）を検証します。セッション不在・`PURPOSE`不一致・無効なセッションまたは失効時の場合は、Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"`）を返却します（ダミーセッション時含む）。
3. **クールダウン検証 (`429 Too Many Requests`)**:
   前回の送信（`LAST_SENT_AT`）から60秒未満である場合は `429 Too Many Requests`（code: `"OTP_RESEND_COOLDOWN"`、遅延なし）を返却します。
4. **OTP再発行・Timing Attack 対策処理 (`200 OK`)**:
   新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、`ATTEMPT_COUNT` を 0 にリセット、`SEND_COUNT` を +1 加算、`LAST_SENT_AT` および `OTP_EXPIRES_AT` を更新するとともに、一律 `1.0s ± 0.1s` のレスポンス遅延を適用した上で `200 OK` を返却します（ダミーセッション時含む）。

##### Errors
- `400 Bad Request`: リクエストボディ不正・必須パラメータ欠落、または無効なセッション/STATUS非active指定（code: `"BAD_REQUEST"`）
- `410 Gone`: 全体最大有効期限（初回発行から15分）切れ（code: `"GONE"`）、またはメール送信5回連続失敗に伴うセッション失効（code: `"OTP_SESSION_INVALIDATED"`）
- `429 Too Many Requests`: クールダウン期間中（前回の送信から60秒未満）の再送要求（code: `"OTP_RESEND_COOLDOWN"`）
- `503 Service Unavailable`: 実メール送信失敗（code: `"OTP_DELIVERY_FAILED"`）

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
- **Set-Cookie**: `XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000`

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "exampleUser",
    "email": "user@example.com",
    "created_at": "2026-08-17T12:00:00+09:00",
    "updated_at": "2026-08-17T12:00:00+09:00"
  }
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `user` | object | ○ | 認証成功したユーザーオブジェクト |
| `user.id` | string | ○ | ユーザーID（UUID形式、例: `550e8400-e29b-41d4-a716-446655440000`） |
| `user.username` | string | ○ | ユーザー名（例: `exampleUser`） |
| `user.email` | string | ○ | メールアドレス（例: `user@example.com`） |
| `user.created_at` | string | ○ | アカウント登録日時（ISO 8601 JST 形式） |
| `user.updated_at` | string | ○ | アカウント更新日時（ISO 8601 JST 形式） |

※なおリクエスト時に既存のログインセッションCookie（`sync_task_sid`）が送信された場合は、複数アカウントへの同時重複ログインを防止するため、その旧セッション（`LOGIN_SESSION`）もDBから物理削除した上で新しいログインセッションを発行します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`email`, `password`）の有無を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし、失敗カウンター加算なし）。
2. **IPレートリミット・アカウントロックアウト検証 (`429 Too Many Requests` / `401 Unauthorized`)**:
   該当IPからの連続失敗遮断（`BLOCKED_UNTIL`）を検証し、遮断中の場合は `429 Too Many Requests`（code: `"RATE_LIMIT_EXCEEDED"`、`Retry-After: <残秒数>`、遅延 1.0s ± 0.1s）を返却します。また、該当メールアドレスのアカウントロックアウト（`LOGIN_LOCK_UNTIL`）中である場合は、アカウントの登録有無を秘匿するため直ちにはエラーを返さず、ステップ3の認証照合において正しいパスワードが入力された場合であっても認証拒否とし、一貫して `401 Unauthorized`（code: `"UNAUTHORIZED"`、遅延 1.0s ± 0.1s）を返却します。
3. **認証照合処理 (`401 Unauthorized` / `200 OK`)**:
   メールアドレスの存在確認およびパスワードハッシュ照合を実施します（論理削除アカウント含む）。
   - 認証失敗: 失敗カウンターを加算し、`401 Unauthorized`（code: `"UNAUTHORIZED"`、遅延 1.0s ± 0.1s）を返却します。
   - 認証成功: 失敗カウンターを0にリセットし、セッションを発行して `200 OK` を返却します。

##### Errors
- `400 Bad Request`: リクエストボディ不正・必須パラメータ欠落（code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 認証失敗（メールアドレス未登録、パスワード不一致、論理削除済みアカウントのいずれも本エラーで一律返却。遅延 1.0s ± 0.1s、code: `"UNAUTHORIZED"`）
- `429 Too Many Requests`: IPレートリミット超過（直近5分間に30回失敗で15分遮断）。遅延 1.0s ± 0.1s（code: `"RATE_LIMIT_EXCEEDED"`）

---

#### 3.1.5 `POST auth/logout`
現在操作中の端末のログインセッションをDBから物理削除し、対象ユーザーの再認証失敗カウンター（`REAUTH_FAILED_COUNT`）を 0 にリセット、`REAUTH_LAST_FAILED_AT` を NULL に更新した上でCookieを消去します。
※未ログインまたはセッションが既に無効化・期限切れの状態で呼び出された場合もエラーとせず、`200 OK` を返却してCookieを確実に消去します（ログアウト操作の冪等性確保）。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### リクエスト評価順序
1. **未ログイン・セッション無効時の挙動 (`200 OK`)**:
   ログインセッションの有効性を確認し、未ログインまたはセッションが無効・期限切れの場合は CSRF 検証の有無にかかわらず `200 OK` を返却し、Cookie 削除ヘッダー（`Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`, `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0`）を返します。
2. **有効ログインセッション存在時の CSRF検証 (`403 Forbidden`)**:
   ログインセッションが有効な場合のみ CSRF トークン（`X-CSRF-Token`）を検証し、欠落・不一致時は `403 Forbidden`（code: `"FORBIDDEN"`）を返却します。
3. **セッション破棄および再認証カウンターリセット (`200 OK`)**:
   当該端末のログインセッション（`LOGIN_SESSION`）を物理削除し、対象ユーザーの再認証失敗カウンター（`LOGIN_ACCOUNT.REAUTH_FAILED_COUNT = 0`）および最終失敗日時（`REAUTH_LAST_FAILED_AT = NULL`）を初期化リセットした上で、Cookie 削除ヘッダーを返却します。

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`

レスポンスボディなし（JSONオブジェクトなし）。

##### Errors
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）

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

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `email` | string | ○ | 有効なメールアドレス形式、前後の空白トリム、小文字正規化 |

##### Response (200 OK)
```json
{
  "otp_session_id": "otp_sess_reset_12345",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300,
  "cooldown_seconds": 60
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 生成されたパスワードリセット用OTPセッションID（例: `otp_sess_reset_12345`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス（例: `user**********@example.com`） |
| `expires_in_seconds` | integer | ○ | OTPの有効期限（秒、デフォルト: 300） |
| `cooldown_seconds` | integer | ○ | 再送可能になるまでのクールダウン秒数（60秒） |

※アカウント列挙防止および Timing Attack 対策として、正常成功時（実メール送信時）およびダミー発行時（未登録のメールアドレス、他ユーザーの有効なOTPセッション期間中等の指定時）を一貫して区別せず、一律でレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`email`）の有無・形式制約を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし）。
2. **既存OTPセッション検証・列挙対策 (`200 OK`)**:
   指定メールに有効なOTPセッションがある場合は、経過時間にかかわらず既存レコードを更新せず、ダミーセッションを新規作成します。有効アカウントが存在し排他なしの場合だけ実OTPを発行し、通常・未登録・削除済み・予約中・重複要求のすべてで `1.0s ± 0.1s` 後に同一形式の `200 OK` を返します。

##### Errors
- `400 Bad Request`: 入力形式違反（code: `"BAD_REQUEST"`）
- `503 Service Unavailable`: 実メール送信失敗（code: `"OTP_DELIVERY_FAILED"`）

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

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 発行されたパスワードリセット用OTPセッションID |
| `otp` | string | ○ | 英数字8桁（大文字・小文字不問） |

##### Response (200 OK)
```json
{
  "message": "OTP verified successfully."
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理結果メッセージ（例: `"OTP verified successfully."`） |

※検証成功時、当該OTPセッション（`OTP_SESSION`）のステータスを `verified` に変更し、全体有効期限（`SESSION_EXPIRES_AT`）を検証成功時点から15分間に延長します（この検証済みOTPセッションは後続の `POST auth/password-reset/reset` エンドポイントでのみ使用可能となります）。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`, `otp`）の有無、および `otp` の形式（英数字8桁）を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし、試行回数 `ATTEMPT_COUNT` 加算なし）。
2. **OTPセッション状態・目的・期限検証 (`400 Bad Request` / `410 Gone`)**:
   指定された `otp_session_id` の存在、用途 `PURPOSE` がパスワードリセット（`PASSWORD_RESET`）であること、ステータス（`active` であること）、および有効期限（`OTP_EXPIRES_AT` / `SESSION_EXPIRES_AT`）を検証します。セッション不在・`PURPOSE`不一致・失効・既に検証済み等の非 `active` ステータスの場合は、Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"`）を返却します（ダミーセッション時含む）。
3. **OTP照合検証 (`400 BAD_REQUEST` / `422 OTP_REISSUED_DUE_TO_FAILURES`)**:
   入力された `otp` のハッシュ照合を実施します。
   - 不一致（試行1〜4回目）: 失敗回数（`ATTEMPT_COUNT`）を+1加算し、`400 Bad Request`（code: `"BAD_REQUEST"`、遅延 1.0s ± 0.1s）を返却します。
   - 不一致（試行5回達成）: 失敗回数をリセットし、新OTPを自動再発行・送信します。
     - 実メール送信に成功した場合（ダミーセッション含む）: `422 Unprocessable Entity`（code: `"OTP_REISSUED_DUE_TO_FAILURES"`、遅延 1.0s ± 0.1s）を返却します。
     - 自動再送における実メール送信に失敗した場合（1〜4回目の送信失敗）: `DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT+=1` とし、`503 Service Unavailable`（code: `"OTP_DELIVERY_FAILED"`）を返却します。
     - 自動再送を含めて5回連続送信失敗となった場合: 対象セッションを物理削除し、`410 Gone`（code: `"OTP_SESSION_INVALIDATED"`）を返却します。

##### Errors
- `400 Bad Request`: 入力形式違反（即時返却、遅延なし）または無効なセッション/PURPOSE不一致/STATUS非active指定またはOTP照合不一致（試行1〜4回目。ダミーセッション時も一律遅延 1.0s ± 0.1s、code: `"BAD_REQUEST"`）
- `410 Gone`: 有効期限切れ（全体最大15分超過含む、code: `"GONE"`）、またはメール送信5回連続失敗に伴うセッション失効（code: `"OTP_SESSION_INVALIDATED"`）
- `422 Unprocessable Entity`: 5回連続失敗に伴う自動再送実行通知（応答遅延 1.0s ± 0.1s、code: `"OTP_REISSUED_DUE_TO_FAILURES"`。ダミーセッション時も実際のメール再送を行わずに全く同一のレスポンスを返却）
- `503 Service Unavailable`: 5回目不一致時の自動再送処理におけるメール送信失敗（再送可能状態を維持。code: `"OTP_DELIVERY_FAILED"`）

---

#### 3.1.8 `POST auth/password-reset/resend-otp`
パスワードリセット用OTPを再送します（60秒クールダウン制約あり）。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_reset_12345"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 発行されたパスワードリセット用OTPセッションID |

##### Response (200 OK)
```json
{
  "message": "OTP has been resent successfully.",
  "masked_email": "user**********@example.com",
  "expires_in_seconds": 300,
  "cooldown_seconds": 60
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理結果メッセージ（例: `"OTP has been resent successfully."`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス（例: `user**********@example.com`） |
| `expires_in_seconds` | integer | ○ | 再発行されたOTPの有効期限（秒、デフォルト: 300） |
| `cooldown_seconds` | integer | ○ | 再送可能になるまでのクールダウン秒数（60秒） |

※Timing Attack 対策として、正常再送時（実メール送信時）およびダミーセッション再送時（実際のメール送信を行わない場合）を一貫して区別せず、一律でレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却します。また、リクエスト対象の `OTP_SESSION` の `PURPOSE` が `PASSWORD_RESET` かつステータスが `active` であることを必須で検証します。再送処理成功時、対象の `OTP_SESSION` レコードにおいて新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、試行失敗回数（`ATTEMPT_COUNT`）を 0 にリセット、送信回数（`SEND_COUNT`）を +1 加算、直前送信日時（`LAST_SENT_AT`）を更新するとともに、有効期限（`OTP_EXPIRES_AT`）を再送信時点から5分間（全体最大有効期限 `SESSION_EXPIRES_AT` の範囲内）へ更新延長します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`）の有無を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし）。
2. **OTPセッション状態・目的・期限検証 (`400 Bad Request` / `410 Gone`)**:
   指定された `otp_session_id` の存在、用途 `PURPOSE` がパスワードリセット（`PASSWORD_RESET`）であること、ステータス（`active` であること）、および全体最大有効期限（`SESSION_EXPIRES_AT` 15分）を検証します。セッション不在・`PURPOSE`不一致・無効なセッションまたは失効時の場合は、Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"`）を返却します（ダミーセッション時含む）。
3. **クールダウン検証 (`429 Too Many Requests`)**:
   前回の送信（`LAST_SENT_AT`）から60秒未満である場合は `429 Too Many Requests`（code: `"OTP_RESEND_COOLDOWN"`、遅延なし）を返却します。
4. **OTP再発行・Timing Attack 対策処理 (`200 OK`)**:
   新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、`ATTEMPT_COUNT` を 0 にリセット、`SEND_COUNT` を +1 加算、`LAST_SENT_AT` および `OTP_EXPIRES_AT` を更新するとともに、一律 `1.0s ± 0.1s` のレスポンス遅延を適用した上で `200 OK` を返却します（ダミーセッション時含む）。

##### Errors
- `400 Bad Request`: リクエストボディ不正・必須パラメータ欠落、または無効なセッション/PURPOSE不一致/STATUS非active指定（code: `"BAD_REQUEST"`）
- `410 Gone`: 全体最大有効期限（初回発行から15分）切れ（code: `"GONE"`）、またはメール送信5回連続失敗に伴うセッション失効（code: `"OTP_SESSION_INVALIDATED"`）
- `429 Too Many Requests`: クールダウン期間中（60秒未満）の再送要求（code: `"OTP_RESEND_COOLDOWN"`）
- `503 Service Unavailable`: 実メール送信失敗（code: `"OTP_DELIVERY_FAILED"`）

---

#### 3.1.9 `POST auth/password-reset/reset`
OTP検証完了後のセッションに基づきパスワードをリセットします。

- **認証**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_reset_12345",
  "new_password": "NewPassword123!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 検証成功済み（ステータス `verified`）のOTPセッションID |
| `new_password` | string | ○ | 8〜128文字、英大文字/英小文字/数字/記号（全32種）のうち3種以上を含む（01_overview.md 1.4節準拠）。検証済み `OTP_SESSION` 経由で取得したユーザーのユーザー名・メールのローカル部（4文字以上の場合、大文字小文字を区別せず比較）を含まないこと。現在のパスワードと同一の場合は 422 エラー |

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`

```json
{
  "message": "Password has been reset successfully. Please log in with your new password."
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理結果メッセージ（例: `"Password has been reset successfully. Please log in with your new password."`） |

※`OTP_SESSION` の `PURPOSE` が `PASSWORD_RESET` かつステータスが `verified` かつ有効期限内であることを確認し、検証済み `OTP_SESSION` 経由で取得した対象ユーザーの属性検証および新パスワードが現在のパスワードと同一でないことを検証します（現在のパスワードと同一の場合やユーザー名/メールローカル部が含まれる場合は 422 エラー）。パスワード更新成功後に当該OTPセッション（`OTP_SESSION`）および該当ユーザーに紐づくすべてのアクティブな `OTP_SESSION`、ならびに該当ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）をDBから直ちに一括物理削除し、Cookieを消去して再ログインを要求します。

##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`, `new_password`）の有無、`new_password` の文字数（8〜128文字）・文字種要件（3種以上）、および `otp_session_id` の形式チェックを検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし）。
2. **OTPセッション状態・目的・期限検証 (`403 Forbidden` / `410 Gone`)**:
   指定された `otp_session_id` の存在、用途 `PURPOSE` がパスワードリセット（`PASSWORD_RESET`）であること、ステータスが `verified` であること、および仮セッションの有効期限（検証成功後15分）内であることを検証します。未検証・`PURPOSE`不一致・無効時は `403 Forbidden`（code: `"FORBIDDEN"`）、期限切れ時は `410 Gone`（code: `"GONE"`）を返却します。
3. **ビジネスルール・ユーザー属性検証 (`422 Unprocessable Entity`)**:
   検証成功した `OTP_SESSION` 経由で取得した対象ユーザーの `username` および `email` のローカル部（4文字以上の場合、大文字小文字を区別せず比較）が `new_password` 内に含まれていないこと、および新パスワードが現在のパスワードと同一でないことを検証します。違反時は `422 Unprocessable Entity`（code: `"SAME_AS_CURRENT_PASSWORD"` または `"INVALID_PASSWORD_CONTENT"`）を返却します。
4. **パスワード更新・セッション破棄 (`200 OK`)**:
   パスワードハッシュを更新し、該当 `OTP_SESSION` および該当ユーザーに紐づくすべてのアクティブな `OTP_SESSION`、ならびに該当ユーザーの全ログインセッション（`LOGIN_SESSION`）を一括物理削除し、Cookieを消去して `200 OK` を返却します。

##### Errors
- `400 Bad Request`: リクエスト構文違反・必須パラメータ欠落、単体パスワード要件違反（文字数・文字種不足）、または形式不正な `otp_session_id` 指定（code: `"BAD_REQUEST"`）
- `403 Forbidden`: 未検証のOTPセッション（`verified` でない場合）でのリセット試行（code: `"FORBIDDEN"`）
- `410 Gone`: OTP検証完了後の仮セッション有効期限切れ（検証成功後15分経過、code: `"GONE"`）
- `422 Unprocessable Entity`: 現在のパスワードと同一のパスワード（code: `"SAME_AS_CURRENT_PASSWORD"`）、または検証済みユーザーのユーザー名/メールローカル部含有（code: `"INVALID_PASSWORD_CONTENT"`）

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

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `new_email` | string | ○ | 有効なメールアドレス形式、前後の空白トリム、小文字正規化。現在のメールアドレスと同一の場合は 422 エラー |

##### Response (200 OK)
```json
{
  "otp_session_id": "otp_sess_chg_998877",
  "masked_email": "new_**********@example.com",
  "expires_in_seconds": 300,
  "cooldown_seconds": 60
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 生成されたメールアドレス変更用OTPセッションID（例: `otp_sess_chg_998877`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス（例: `new_**********@example.com`、短ローカル部時は01_overview.md準拠） |
| `expires_in_seconds` | integer | ○ | OTPの有効期限（秒、デフォルト 300） |
| `cooldown_seconds` | integer | ○ | 再送可能になるまでのクールダウン秒数（60秒） |

※アカウント列挙防止および Timing Attack 対策として、正常成功時（実メール送信時）およびダミー発行時（登録済みの他ユーザーのメールアドレス、他ユーザーの有効なOTPセッション期間中等の指定時、および直前送信から60秒以内の再要求時）を一貫して区別せず、一律でレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却します（60秒以内の再要求時は実際のメール送信を行わず、残クールダウン秒数を返却）。

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request` / `422 SAME_AS_CURRENT_EMAIL`)**:
   リクエストボディの `new_email` 形式を検証（不備時は 400 `BAD_REQUEST`）。現在のメールアドレスと同一かを検証（同一時は 422 `SAME_AS_CURRENT_EMAIL`）。
3. **既存OTPセッション検証・列挙対策 (`200 OK`)**:
   認証ユーザーに有効なメール変更OTPセッションがある場合、前回の送信から60秒以内であればメールを再送せず `200 OK`（残 `cooldown_seconds` 返却）とします。変更先が未使用かつ排他なしの場合だけ実OTPを発行し、正常・登録済み・予約中・重複要求のすべてで `1.0s ± 0.1s` 後に同一形式の `200 OK` を返します。

##### Errors
- `400 Bad Request`: メールアドレス形式不正（code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `422 Unprocessable Entity`: 現在のメールアドレスと同一（code: `"SAME_AS_CURRENT_EMAIL"`）
- `503 Service Unavailable`: 実メール送信失敗（code: `"OTP_DELIVERY_FAILED"`）

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

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 発行されたメールアドレス変更用OTPセッションID |
| `otp` | string | ○ | 英数字8桁（大文字・小文字不問） |

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`

```json
{
  "message": "Email address updated successfully.",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "exampleUser",
    "email": "new_user@example.com",
    "created_at": "2026-08-17T12:00:00+09:00",
    "updated_at": "2026-08-17T12:30:00+09:00"
  }
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理結果メッセージ |
| `user` | object | ○ | 更新後のユーザーオブジェクト |
| `user.id` | string | ○ | ユーザーID（UUID形式） |
| `user.username` | string | ○ | ユーザー名 |
| `user.email` | string | ○ | 変更後の新メールアドレス |
| `user.created_at` | string | ○ | アカウント作成日時（ISO 8601 JST 形式） |
| `user.updated_at` | string | ○ | メールアドレス更新日時（ISO 8601 JST 形式） |

※指定された `otp_session_id` に紐づくユーザーID（`OTP_SESSION.USER_ID`）が現在認証中のログインユーザーIDと一致していることを検証します。検証成功後、アカウントのメールアドレスを更新し、旧メールアドレス宛てに変更完了通知メールを送信（非同期処理）します。同時に、使用済みの手続き用OTPセッション（`OTP_SESSION`）および当該ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）をDBから直ちに物理削除してCookieを消去し、新メールアドレスでの再ログインを要求します。

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`, `otp`）の有無、および `otp` の形式（英数字8桁）を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし、試行回数 `ATTEMPT_COUNT` 加算なし）。
3. **認可・OTPセッション状態・目的・期限検証 (`400 Bad Request` / `403 Forbidden` / `410 Gone`)**:
   指定された `otp_session_id` の存在および現在ログイン中のユーザーに紐づくセッションであることを検証（他者所有の場合は 403 `FORBIDDEN`）。また、用途 `PURPOSE` がメールアドレス変更（`EMAIL_CHANGE`）であること、ステータスが `active` であること、および有効期限（`OTP_EXPIRES_AT` / `SESSION_EXPIRES_AT`）を検証します。セッション不在・`PURPOSE`不一致・失効・既に検証済み等の非 `active` ステータスの場合、Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"` / `"OTP_SESSION_INVALIDATED"`）を返却します（ダミーセッション時含む）。
4. **OTP照合検証 (`400 BAD_REQUEST` / `422 OTP_REISSUED_DUE_TO_FAILURES`)**:
   入力された `otp` のハッシュ照合を実施します。
   - 不一致（試行1〜4回目）: 失敗回数（`ATTEMPT_COUNT`）を+1加算し、`400 Bad Request`（code: `"BAD_REQUEST"`、遅延 1.0s ± 0.1s）を返却します。
   - 不一致（試行5回達成）: 失敗回数をリセットし、新OTPを自動再発行・送信します。
     - 実メール送信に成功した場合（ダミーセッション含む）: `422 Unprocessable Entity`（code: `"OTP_REISSUED_DUE_TO_FAILURES"`、遅延 1.0s ± 0.1s）を返却します。
     - 自動再送における実メール送信に失敗した場合（1〜4回目の送信失敗）: `DELIVERY_STATUS='sendable'`、`SEND_FAILED_COUNT+=1` とし、`503 Service Unavailable`（code: `"OTP_DELIVERY_FAILED"`）を返却します。
     - 自動再送を含めて5回連続送信失敗となった場合: 対象セッションを物理削除し、`410 Gone`（code: `"OTP_SESSION_INVALIDATED"`）を返却します。
5. **メールアドレス更新確定・全ログインセッション物理削除 (`200 OK`)**:
   OTP照合成功時、ユーザーのメールアドレスを新アドレスへ更新し、セキュリティ要件として**当該ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）および使用済み `OTP_SESSION` をDBから直ちに物理削除**します。また、コミット後に更新前の旧メールアドレス宛へ変更完了通知メールを非同期送信します（通知送信失敗時も確定済みメールアドレス更新はロールバックせず、エラーログを記録）。レスポンスヘッダーに Cookie 削除ヘッダー（`Max-Age=0`）を付与して返却し、フロントエンド側でログイン画面へリダイレクト（新メールアドレスでの再ログイン要求）させます。

##### Errors
- `400 Bad Request`: 入力形式違反（即時返却、遅延なし）または無効なセッション/PURPOSE不一致/STATUS非active指定またはOTP照合不一致（試行1〜4回目。ダミーセッション時も一律遅延 1.0s ± 0.1s、code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: 他ユーザー所有の `otp_session_id` 指定（認可不一致）または CSRFトークン不正（code: `"FORBIDDEN"`）
- `410 Gone`: 有効期限切れ（全体最大15分超過含む、code: `"GONE"`）、またはメール送信5回連続失敗に伴うセッション失効（code: `"OTP_SESSION_INVALIDATED"`）
- `422 Unprocessable Entity`: 5回連続失敗に伴う自動再送実行通知（応答遅延 1.0s ± 0.1s、code: `"OTP_REISSUED_DUE_TO_FAILURES"`。ダミーセッション時も実際のメール再送を行わずに全く同一のレスポンスを返却）
- `503 Service Unavailable`: 5回目不一致時の自動再送処理におけるメール送信失敗（再送可能状態を維持。code: `"OTP_DELIVERY_FAILED"`）

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

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 発行されたメールアドレス変更用OTPセッションID |

##### Response (200 OK)
```json
{
  "message": "OTP has been resent successfully.",
  "masked_email": "new_**********@example.com",
  "expires_in_seconds": 300,
  "cooldown_seconds": 60
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理結果メッセージ（例: `"OTP has been resent successfully."`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス（例: `new_**********@example.com`） |
| `expires_in_seconds` | integer | ○ | 再発行されたOTPの有効期限（秒、デフォルト 300） |
| `cooldown_seconds` | integer | ○ | 再送可能になるまでのクールダウン秒数（60秒） |

※Timing Attack 対策として、正常再送時（実メール送信時）およびダミーセッション再送時（実際のメール送信を行わない場合）を一貫して区別せず、一律でレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却します。また、リクエスト対象の `OTP_SESSION` の `PURPOSE` が `EMAIL_CHANGE` かつステータスが `active` であることを必須で検証します。再送処理成功時、対象の `OTP_SESSION` レコードにおいて新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、試行失敗回数（`ATTEMPT_COUNT`）を 0 にリセット、送信回数（`SEND_COUNT`）を +1 加算、直前送信日時（`LAST_SENT_AT`）を更新するとともに、有効期限（`OTP_EXPIRES_AT`）を再送信時点から5分後（全体最大有効期限 `SESSION_EXPIRES_AT` の範囲内）へ更新延長します。

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`otp_session_id`）の有無を検証します。不備がある場合は即座に `400 Bad Request`（code: `"BAD_REQUEST"`）を返却します（遅延なし）。
3. **認可・OTPセッション状態・目的・期限検証 (`400 Bad Request` / `403 Forbidden` / `410 Gone`)**:
   指定された `otp_session_id` の存在および現在ログイン中のユーザーに紐づくセッションであることを検証（他者所有の場合は 403 `FORBIDDEN`）。用途 `PURPOSE` が `EMAIL_CHANGE` であること、ステータスが `active` であること、および全体最大有効期限（`SESSION_EXPIRES_AT` 15分）を検証します。セッション不在・`PURPOSE`不一致・失効時の場合は Timing Attack 対策として一律 `1.0s ± 0.1s` の遅延を適用し `400 Bad Request`（code: `"BAD_REQUEST"`）または `410 Gone`（code: `"GONE"` / `"OTP_SESSION_INVALIDATED"`）を返却します（ダミーセッション時含む）。
4. **クールダウン検証 (`429 Too Many Requests`)**:
   前回の送信（`LAST_SENT_AT`）から60秒未満である場合は `429 Too Many Requests`（code: `"OTP_RESEND_COOLDOWN"`、`Retry-After: <残秒数>` ヘッダー付与、遅延なし）を返却します。
5. **OTP再発行・Timing Attack 対策処理 (`200 OK`)**:
   新たな8桁OTPコード（`OTP_HASH`）を発行・保存し、`ATTEMPT_COUNT` を 0 にリセット、`SEND_COUNT` を +1 加算、`LAST_SENT_AT` および `OTP_EXPIRES_AT` を更新するとともに、一律 `1.0s ± 0.1s` のレスポンス遅延を適用した上で `200 OK` を返却します（ダミーセッション時含む）。

##### Errors
- `400 Bad Request`: リクエストボディ不正・必須パラメータ欠落、または無効なセッション/PURPOSE不一致/STATUS非active指定（code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: 他ユーザー所有の `otp_session_id` 指定（認可不一致）または CSRFトークン不正（code: `"FORBIDDEN"`）
- `410 Gone`: 初回発行から15分経過（code: `"GONE"`）、またはメール送信5回連続失敗に伴うセッション失効（code: `"OTP_SESSION_INVALIDATED"`）
- `429 Too Many Requests`: クールダウン期間中（60秒未満、code: `"OTP_RESEND_COOLDOWN"`、`Retry-After` ヘッダー付与）
- `503 Service Unavailable`: 実メール送信失敗（code: `"OTP_DELIVERY_FAILED"`）

---

#### 3.1.13 `POST auth/otp-session/cancel`
ユーザー操作（OTP入力画面での「戻る」ボタン押下、画面離脱等）により、進行中のOTPセッションをサーバー側で即座に物理削除（無効化）します。

- **認証**: 不要（会員登録・パスワードリセット・メール変更の全種別で利用可能）
- **Headers**: 不要

##### Request Body
```json
{
  "otp_session_id": "otp_sess_chg_998877"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 無効化対象のOTPセッションID |

##### Response (200 OK)
```json
{
  "message": "OTP session cancelled successfully."
}
```

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 処理完了メッセージ（例: `"OTP session cancelled successfully."`） |

※指定された `otp_session_id` が存在する場合、メール認証ログテーブル（`MAIL_AUTH_LOG`）に `EVENT_TYPE='CANCELLED'` を記録した上で、該当レコードを DB から直ちに物理削除（`DELETE FROM OTP_SESSION`）します。指定されたセッションが存在しない場合やダミーセッションの場合も、アカウント列挙防止のため区別せず同一の `200 OK` を返却します。

##### Errors
- `400 Bad Request`: リクエストボディ不正・`otp_session_id` 欠落（code: `"BAD_REQUEST"`）
