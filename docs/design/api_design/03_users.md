### 3.2 ユーザー管理 (Users)

#### 3.2.1 `GET users/{user_id}`
ログインユーザーのプロフィール情報（ユーザー名、メールアドレス等）を取得します。

- **認証**: 必須（Cookie）
- **Path Parameters**:
  - `user_id` (string): 取得対象のユーザーID（セッションと一致必須。不一致時は 404）

##### リクエスト評価順序
1. **認証検証 (`401 Unauthorized`)**:
   ログインセッションの有効性を確認（未ログインまたはセッション無効時は 401 `UNAUTHORIZED`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証（不一致または存在しない場合は 404 `NOT_FOUND`）。

##### Response (200 OK)

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `user` | object | ○ | ユーザー情報オブジェクト |
| `user.id` | string | ○ | ユーザーID（例: `usr_987654321`） |
| `user.username` | string | ○ | ユーザー名 |
| `user.email` | string | ○ | 登録メールアドレス |
| `user.created_at` | string | ○ | アカウント作成日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |
| `user.updated_at` | string | ○ | アカウント最終更新日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |

```json
{
  "user": {
    "id": "usr_987654321",
    "username": "exampleUser",
    "email": "user@example.com",
    "created_at": "2026-08-01T10:00:00+09:00",
    "updated_at": "2026-08-17T16:00:00+09:00"
  }
}
```

##### Errors
- `401 Unauthorized`: 未ログインまたはセッション無効（code: `"UNAUTHORIZED"`）
- `404 Not Found`: ユーザーが存在しない、または他ユーザーの `user_id` を指定した場合（code: `"NOT_FOUND"`）

---

#### 3.2.2 `PUT users/{user_id}`
プロフィール情報（ユーザー名）を更新します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須フィールド `username` の有無および型（`null` 指定不可）、トリム後の文字数（2〜20文字）・使用可能文字（英数字）を検証（不備時は 400 `BAD_REQUEST` を即座に返却）。なお、読み取り専用・変更不可フィールド（`id`, `email`, `created_at`, `updated_at` 等）がボディに含まれている場合は、エラーとせず更新対象外として単に無視します。
3. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
4. **ビジネスルール検証 (`422 Unprocessable Entity`)**:
   トリム後の `username` が現在のユーザー名と同一か検証（同一の場合は 422 `SAME_AS_CURRENT_USERNAME`）。

##### Request Body
```json
{
  "username": "newUsername"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `username` | string | ○ | 必須。null 不可。前後の空白を自動トリム後、2〜20文字かつ英数字（大文字小文字可）であること。要件違反・未入力時は 400 エラー。トリム後のユーザー名が現在のユーザー名と同一の場合は 422 エラー |

##### Response (200 OK)

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `user` | object | ○ | ユーザー情報オブジェクト |
| `user.id` | string | ○ | ユーザーID（例: `usr_987654321`） |
| `user.username` | string | ○ | 更新後のユーザー名 |
| `user.email` | string | ○ | 登録メールアドレス |
| `user.created_at` | string | ○ | アカウント作成日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |
| `user.updated_at` | string | ○ | アカウント最終更新日時（ISO 8601 / `YYYY-MM-DDTHH:mm:ss+09:00`） |

```json
{
  "user": {
    "id": "usr_987654321",
    "username": "newUsername",
    "email": "user@example.com",
    "created_at": "2026-08-01T10:00:00+09:00",
    "updated_at": "2026-08-17T16:35:00+09:00"
  }
}
```

##### Errors
- `400 Bad Request`: リクエスト構文違反・必須パラメータ欠落/null指定、またはトリム後のユーザー名要件違反（文字数2〜20文字範囲外・使用可能文字不正・空白のみ送信等、code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 未ログインまたはセッション無効（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークンヘッダーの欠落または不一致（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー（他ユーザーID指定または存在しないユーザー、code: `"NOT_FOUND"`）
- `422 Unprocessable Entity`: 現在のユーザー名と同一（code: `"SAME_AS_CURRENT_USERNAME"`）

---

#### 3.2.3 `DELETE users/{user_id}`
パスワード再認証を行い、アカウントを論理削除（`IS_DELETED=true`, `DELETED_AT=NOW()`）します。同メールアドレスでの将来の再登録を可能とするため、`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムを退避フォーマット（`deleted_<USER_ID>_<EMAIL>`）へ更新します。なお、所有タスクデータおよび該当ユーザーの全ログインセッション（`LOGIN_SESSION`）とアクティブなOTPセッション（`OTP_SESSION`）は即座に物理削除されます。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須フィールドの有無を検証（未入力時は即座に 400 を返し、遅延・失敗カウンター加算なし）。
3. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証（不一致時は即座に 404 を返し、遅延・失敗カウンター加算なし）。
4. **パスワード再認証検証 (`401 REAUTH_FAILED` / `SESSION_DESTROYED`)**:
   現在のパスワードのハッシュ照合を実施。不一致時は一律 `1.0s ± 0.1s` の遅延を適用し、失敗回数をカウントアップ。5回達成分は操作中の該当ログインセッションのみを強制破棄および Cookie 削除ヘッダーを付与。**照合成功時は即座に `REAUTH_FAILED_COUNT` を 0 にリセットしてアカウント削除処理（論理削除・セッション削除等）を実行する。**

##### Request Body
```json
{
  "password": "Password123!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `password` | string | ○ | 再認証用パスワード。未入力時は 400 エラー、パスワード不一致時は 401 エラー |

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`

```json
{
  "message": "Account has been deleted successfully."
}
```
※ パスワード再認証失敗時（1〜4回目および5回連続達成分）は、Timing Attack 対策として一律 `1.0s ± 0.1s` のレスポンス遅延を適用してエラーを返却します。なお、パスワード再認証失敗が5回連続に達した場合（`SESSION_DESTROYED`）、セキュリティ保護のため操作中の該当ログインセッションのみを直ちに物理削除し、Cookieを消去します。また、再認証失敗カウンター（`REAUTH_FAILED_COUNT`）は、再認証成功時、5回失敗によるセッション強制破棄時、およびログアウト（`auth/logout`）時に 0 にリセットされます。

##### Errors
- `400 Bad Request`: リクエストボディ不正・パスワード未入力（code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 
  - 未ログイン・セッション無効（code: `"UNAUTHORIZED"`）
  - パスワード再認証失敗 1〜4回目（code: `"REAUTH_FAILED"`、遅延 1.0s ± 0.1s）
  - パスワード再認証失敗 5回連続達成分（セッション強制破棄・Cookie消去、code: `"SESSION_DESTROYED"`、遅延 1.0s ± 0.1s）
    - **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
    - **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー（他ユーザーID指定または存在しないユーザー、code: `"NOT_FOUND"`）

---

#### 3.2.4 `PATCH users/{user_id}/password`
現在のパスワードを検証した上で、新しいパスワードへ変更し、該当ユーザーの全ログインセッション（`LOGIN_SESSION`）およびアクティブなOTPセッション（`OTP_SESSION`）を一括物理削除します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `user_id` (string): 対象ユーザーID（セッションと一致必須）

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須フィールド（`current_password`, `new_password`）の有無、および `new_password` の文字数（8〜128文字）・使用可能文字種（英大文字/小文字/数字/記号のうち3種以上）を検証。不備がある場合は即座に 400 エラーを返却（遅延・失敗カウンター加算なし）。
3. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証。不一致時は即座に 404 エラーを返却（遅延・失敗カウンター加算なし）。
4. **パスワード再認証検証 (`401 REAUTH_FAILED` / `SESSION_DESTROYED`)**:
   `current_password` のハッシュ照合を実施。不一致時は一律 `1.0s ± 0.1s` の遅延を適用し、失敗回数をカウントアップ。5回達成分は操作中の該当ログインセッションのみを強制破棄および Cookie 削除ヘッダーを付与。**照合成功時は即座に `REAUTH_FAILED_COUNT` を 0 にリセットしてステップ5に進む。**
5. **ビジネスルール検証 (`422 Unprocessable Entity`)**:
   再認証成功後、`new_password` に対象ユーザーのユーザー名・メールアドレスのローカル部（4文字以上の場合、大文字・小文字を区別せず Case-Insensitive）が含まれていないか検証（含有時は 422 `INVALID_PASSWORD_CONTENT` を返却）、および照合済みの `current_password` と同一でないか検証（同一の場合は 422 `SAME_AS_CURRENT_PASSWORD` を返却）。

##### Request Body
```json
{
  "current_password": "Password123!",
  "new_password": "NewSecurePassword456!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `current_password` | string | ○ | 現在のパスワード。不一致時は 401 (`REAUTH_FAILED`) エラー |
| `new_password` | string | ○ | 8〜128文字。英大文字/英小文字/数字/記号のうち3種以上を含む。ログインセッション/DBから取得した対象ユーザーのユーザー名・メールアドレスのローカル部（4文字以上の場合）を、大文字・小文字を区別せず（Case-Insensitive）新パスワード内に含まないこと。現在のパスワードと同一の場合は 422 エラー |

##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`

```json
{
  "message": "Password has been updated successfully. Please log in again."
}
```
※ パスワード再認証失敗時（1〜4回目および5回連続達成分）は、Timing Attack 対策として一律 `1.0s ± 0.1s` のレスポンス遅延を適用してエラーを返却します。なお、パスワード再認証失敗が5回連続に達した場合（`SESSION_DESTROYED`）、セキュリティ保護のため操作中の該当ログインセッションのみを直ちに物理削除し、Cookieを消去します。また、再認証失敗カウンター（`REAUTH_FAILED_COUNT`）は、再認証成功時、5回失敗によるセッション強制破棄時、およびログアウト（`auth/logout`）時に 0 にリセットされます。

##### Errors
- `400 Bad Request`: リクエスト構文違反・必須パラメータ欠落、新パスワード単体要件違反（文字数・文字種不足、code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 
  - 未ログイン・セッション無効（code: `"UNAUTHORIZED"`）
  - 現在のパスワード不一致 1〜4回目（code: `"REAUTH_FAILED"`、遅延 1.0s ± 0.1s）
  - 現在のパスワード不一致 5回連続達成分（セッション強制破棄・Cookie消去、code: `"SESSION_DESTROYED"`、遅延 1.0s ± 0.1s）
    - **Set-Cookie**: `sync_task_sid=; Path=/; Max-Age=0`
    - **Set-Cookie**: `XSRF-TOKEN=; Path=/; Max-Age=0`

- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー（code: `"NOT_FOUND"`）
- `422 Unprocessable Entity`: 新パスワードにユーザー名/メールローカル部含有（code: `"INVALID_PASSWORD_CONTENT"`）、または新パスワードが現在のパスワードと同一（code: `"SAME_AS_CURRENT_PASSWORD"`）

---
