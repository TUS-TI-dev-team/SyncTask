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
