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

#### 3.3.4 `PATCH tasks/{task_id}`
タスク情報を部分更新します。リクエストボディに含まれるフィールドのみが更新対象となります（ステータス変更 `status` やピン留め `is_pinned` の単体更新含む）。

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
