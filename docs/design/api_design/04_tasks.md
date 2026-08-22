### 3.3 タスク管理 (Tasks)

#### 3.3.1 `GET tasks`
タスク一覧を取得します。クエリパラメータにより、通常一覧、優先タスク・締切間近・ピン留めビュー、検索絞り込み、カレンダー表示用期間取得に対応します。
※本APIは条件に合致するタスク全件を一括返却します。画面（UI）側でのページネーション（1ページ20件表示）は、取得したタスク配列をもとにクライアントサイドで分割制御を行います。ホーム画面においては、フロントエンドが選択中のタブに応じて `view_type=high_priority`、`view_type=near_deadline`、`view_type=pinned` のAPIリクエストを発行し、20件単位のページ分割制御およびページネーションUI操作を行います。

- **認証**: 必須（Cookie）

##### Query Parameters

| パラメータ名 | 型 | 必須 | デフォルト | 説明 |
| :--- | :--- | :---: | :--- | :--- |
| `view_type` | string | × | - | ビュー指定: `high_priority`（優先高）, `near_deadline`（72時間以内/期限超過。`include_completed` の指定に関わらず常に完了タスクは除外。※`status=completed` との同時指定はパラメータ競合として 400 Bad Request）, `pinned`（ピン留めのみ）。定義外の値指定時は 400 Bad Request（code: `"BAD_REQUEST"`） |
| `include_completed`| boolean | × | 通常/due_date: `false`<br>カレンダー: `true` | 完了タスクを含めるか（`true` / `false`）。`true` または `false` 以外の文字列・数値が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却。通常一覧時および `due_date` 指定時はデフォルト `false`、`start_date` / `end_date` 指定（カレンダー表示）時はデフォルト `true`。なお `status` が明示指定された場合は本パラメータは無視されます |
| `keyword` | string | × | - | タスク名およびコメントの部分一致検索（日本語同一視: 英大文字/小文字、日本語全角/半角（NFKC正規化）、ひらがな/カタカナを同一視。アプリケーション層で入力キーワードを小文字化・NFKC正規化・ひらがな→カタカナ変換した上で、`TASK.SEARCH_TEXT` に対する部分一致検索（`ILIKE` / `pg_trgm`）を実行。前後の空白文字（半角・全角スペース、タブ、改行）をトリム処理。トリム後のキーワードが100文字を超える場合は 400 Bad Request（code: `"BAD_REQUEST"`）。SQLワイルドカード特殊文字 % や _ や \ はリテラルエスケープ。未入力またはトリム後空文字時は絞り込みを行わない） |
| `priority` | string | × | - | 優先度絞り込み: `high`, `medium`, `low`。定義外の値指定時は 400 Bad Request（code: `"BAD_REQUEST"`） |
| `status` | string | × | - | ステータス絞り込み: `not_started`, `in_progress`, `completed`。明示指定時は `include_completed` の指定を無視して指定ステータスを最優先適用（※`view_type=near_deadline` との同時指定は 400 Bad Request）。定義外の値指定時は 400 Bad Request（code: `"BAD_REQUEST"`） |
| `due_date` | string | × | - | 締切日絞り込み（`YYYY-MM-DD`。抽出ロジック: `include_completed=true` の場合は指定日当日の全タスク［未完了＋完了］および指定日より過去の未完了タスク［`status != 'completed'`］を抽出対象とし、過去の完了済みタスクは常に除外。`include_completed=false`［デフォルト］の場合は指定日当日の未完了タスクおよび指定日より過去の未完了タスクを抽出対象とし、当日・過去を問わずすべての完了タスクを除外。締切日時未設定 `null` のタスクは除外。未指定時は絞り込みを行わない。※`start_date` / `end_date` との同時指定は不可） |
| `start_date` | string | × | - | カレンダー表示用: グリッド取得開始日（`YYYY-MM-DD`）。`end_date` とペアで指定必須。最大許容期間幅は 42日間（6週間）。※`due_date` との同時指定は不可（同時指定時は 400 Bad Request） |
| `end_date` | string | × | - | カレンダー表示用: グリッド取得終了日（`YYYY-MM-DD`）。`start_date` とペアで指定必須（`start_date <= end_date`）。※`due_date` との同時指定は不可（同時指定時は 400 Bad Request） |
| `sort_by` | string | × | `default` | ソート種別。指定可能値: `default`（ピン留め優先→締切昇順→作成日時降順）、`due_date_asc`（締切昇順）、`due_date_desc`（締切降順）、`created_at_desc`（作成日時降順）、`priority_desc`（優先度降順: `high` → `medium` → `low`）。定義外の値指定時は 400 Bad Request（code: `"BAD_REQUEST"`）。※`due_date_asc` および `due_date_desc` 指定時、締切日時未設定（`null`）のタスクは常に末尾に配置されます（`NULLS LAST`）。なお、すべてのソート指定において第一ソート条件で同一値となるタスクについては、確定的なソート順序を保証するためタイブレーク条件として `created_at DESC` → `id DESC` が適用されます。なお `view_type` 指定時に `sort_by` が明示的に指定された場合は、指定された `sort_by` のソート条件が最優先で適用されます（`sort_by` 省略時は `view_type` ごとのデフォルトソートが適用されます）。 |

##### リクエスト評価順序
1. **認証検証 (`401 Unauthorized`)**:
   ログインセッションの有効性を確認（未ログインまたはセッション無効時は 401 `UNAUTHORIZED`）。
2. **リクエスト構文・クエリパラメータバリデーション (`400 Bad Request`)**:
   クエリパラメータの型・形式（`include_completed` への非 boolean 指定、`keyword` トリム後100文字超過、`priority`, `status`, `view_type`, `sort_by` パラメータへの定義外の無効な値指定、`view_type=near_deadline` と `status=completed` の同時指定、日付フォーマット `YYYY-MM-DD` 違反、`start_date` / `end_date` 片側のみ指定、`start_date` / `end_date` と `due_date` の同時指定、`start_date > end_date`、期間幅42日超過等）を検証。不備がある場合は 400 `BAD_REQUEST` を返却。

##### Response (200 OK)
```json
{
  "tasks": [
    {
      "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "課題レポート提出",
      "comment": "第5章の要約を含むこと\n参考文献を記載",
      "priority": "high",
      "status": "in_progress",
      "due_datetime": "2026-08-20T23:59:00+09:00",
      "is_pinned": true,
      "created_at": "2026-08-17T10:00:00+09:00",
      "updated_at": "2026-08-17T11:30:00+09:00"
    }
  ]
}
```

##### Response Body フィールド定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `tasks` | array[object] | ○ | タスクオブジェクトの配列（検索・絞り込みに該当するタスクが0件の場合は空配列 `[]`） |
| `tasks[].id` | string | ○ | タスクID（UUID形式、例: `7c9e6679-7425-40de-944b-e07fc1f90ae7`） |
| `tasks[].user_id` | string | ○ | 所有ユーザーID（UUID形式、例: `550e8400-e29b-41d4-a716-446655440000`） |
| `tasks[].title` | string | ○ | タスクタイトル |
| `tasks[].comment` | string | ○ | タスクコメント（未入力時は空文字 `""`） |
| `tasks[].priority` | string | ○ | 優先度（`high`, `medium`, `low`） |
| `tasks[].status` | string | ○ | ステータス（`not_started`, `in_progress`, `completed`） |
| `tasks[].due_datetime` | string / null | ○ | 締切日時（ISO 8601 JST 形式。締切未設定時は `null`） |
| `tasks[].is_pinned` | boolean | ○ | ピン留めフラグ |
| `tasks[].created_at` | string | ○ | 作成日時（ISO 8601 JST 形式） |
| `tasks[].updated_at` | string | ○ | 更新日時（ISO 8601 JST 形式） |

※締切日時が未設定のタスクの場合、`due_datetime` は `null` として返却されます。またコメントが未入力の場合、`comment` は空文字 `""` として返却されます。
※ `sort_by` 省略時に `view_type` ごとに適用されるデフォルトソート順序は以下の通りです：
- `high_priority`: 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）※ピン留めによる並び替えは行いません
- `near_deadline`: 締切日時昇順（`due_datetime ASC`） → 作成日時降順（`created_at DESC`）
- `pinned`: 締切日時昇順（`due_datetime ASC NULLS LAST`） → 作成日時降順（`created_at DESC`）

※ `start_date` / `end_date` を指定したカレンダー期間取得時は、`due_datetime` が設定されているタスクのうち `start_date 00:00:00+09:00 <= due_datetime <= end_date 23:59:59+09:00` の範囲に該当するタスクのみが抽出返却されます（`due_datetime` が `null` のタスクは除外されます）。なお、`due_date` パラメータとの同時指定は不可となり、同時に指定された場合は 400 `BAD_REQUEST` を返却します。`start_date` / `end_date` 指定時に `view_type`, `priority`, `status`, `keyword` 等の絞り込みパラメータが併用された場合は、指定された期間内で該当する絞り込み条件を満たすタスクのみが一括返却されます。※ `start_date` / `end_date` 指定（カレンダー期間取得）時に `sort_by` が省略された場合のデフォルトソート順序は、`is_pinned DESC`（ピン留め優先） → `due_datetime ASC`（締切日時昇順） → `created_at DESC`（作成日時降順） → `id DESC`（タイブレーク）となります。

##### Errors
- `400 Bad Request`: クエリパラメータ不正（include_completedの非boolean指定、priority/status/view_type/sort_byの不正値指定、view_type=near_deadlineとstatus=completedの同時指定、日付フォーマット違反、start_date/end_dateとdue_dateの同時指定、片側期間指定、start_date > end_date、期間幅超過(42日超)等、code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 未ログインまたはセッション無効（code: `"UNAUTHORIZED"`）

---

#### 3.3.2 `POST tasks`
新規タスクを作成します。単一タスク作成に加え、期間と曜日を指定した毎週タスクの即時一括生成（最大100件）に対応します。タスク作成時、タイトルおよびコメントから検索用正規化文字列（小文字化＋NFKC正規化＋ひらがな→カタカナ変換）を生成し、`TASK.SEARCH_TEXT` カラムに自動保存します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須パラメータ（`title` 等）の有無、`title` 文字数・制御文字制約、`priority` への `null` や定義外列挙値（`high`, `medium`, `low` 以外）指定有無、`is_pinned` への `null` や非 boolean 型（数値・文字列等）指定有無、`recurring_rule` の各フィールド値・形式・期間制約（生成件数1〜100件範囲内）等を検証。不備時は 400 `BAD_REQUEST` を返却。

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
| `title` | string | ○ | 前後の空白文字（半角・全角スペース、タブ、改行）を除去（トリム）した上で1〜100文字必須。空白文字のみの入力および文字列内部への改行・タブ等の制御文字入力は 400 Bad Request（code: `"BAD_REQUEST"`）として拒否 |
| `comment` | string / null | × | 0〜1000文字（トリム後）。改行は `\n` に正規化。未入力時は空文字 `""` として登録 |
| `priority` | string | × | `high`, `medium`, `low`（デフォルト: `medium`）。明示的に `null` または定義外の列挙値（`high`, `medium`, `low` 以外）が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却 |
| `due_datetime` | string / null | × | ISO 8601 日時文字列（例: `2026-08-20T23:59:00+09:00`）、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を設定）。※タイムゾーンオフセットを含まない ISO 8601 文字列が指定された場合はデフォルトで JST (`+09:00`) と解釈し、UTC (`Z`) や他タイムゾーン指定時は JST (`+09:00`) に変換・正規化して登録・返却します。省略時または `null` 指定時は締切日時未設定（`null`）として作成。※`is_recurring: true` 時は無視されます |
| `is_pinned` | boolean | × | ピン留めフラグ（`true` / `false`）。省略時はデフォルト `false` として作成。明示的に `null` または非 boolean 型（数値・文字列等）が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却 |
| `is_recurring` | boolean | × | 繰り返し一括作成フラグ（デフォルト: `false`） |
| `recurring_rule` | object | △ | `is_recurring: true` 時のみ必須（`false` 時は無視） |
| `recurring_rule.start_date` | string | ○ | 開始日（`YYYY-MM-DD`）。`start_date <= end_date`。なお、業務要件に基づき `start_date` には過去日付の指定も可能であり、`start_date <= end_date` かつ期間が1年以内（生成件数1〜100件）の条件を満たしていれば、過去日に該当するタスクも正常に一括生成されます |
| `recurring_rule.end_date` | string | ○ | 終了日（`YYYY-MM-DD`）。最大1年間（52週以内） |
| `recurring_rule.days_of_week` | array[string] | ○ | 1つ以上の要素が必須（空配列 `[]` 不可）。`["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]` より選択。文字列は小文字へ正規化され重複はデデュプリケーション処理 |
| `recurring_rule.due_time` | string | × | 締切時刻 24時間表記の `HH:mm`（`00:00`〜`23:59`、省略時は `23:59`）。`HH:mm` 形式でない場合、または `00:00`〜`23:59` の範囲外の数値や秒数を含む形式が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却 |

※`is_recurring: true` の場合、生成件数が1〜100件の範囲内で即時一括生成されます（0件または101件以上の場合はエラーとなり作成されません）。繰り返し一括生成される各タスクの `due_datetime` には、該当日に `due_time`（省略時は `23:59`）と JST オフセット `+09:00` が結合された ISO 8601 日時文字列（例: `"2026-08-22T18:00:00+09:00"`）が自動的に設定されて登録・返却されます。なお、`is_recurring: true` 時に `is_pinned` フィールドが指定された場合、一括生成されるすべてのタスクにそのピン留め設定が適用されます。

##### Response (201 Created)
```json
{
  "created_count": 10,
  "tasks": [
    {
      "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
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

##### Response Body フィールド定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `created_count` | integer | ○ | 作成されたタスク件数（単一作成時は `1`、繰り返し一括作成時は作成数） |
| `tasks` | array[object] | ○ | 作成されたタスクオブジェクトの配列 |
| `tasks[].id` | string | ○ | タスクID（UUID形式、例: `7c9e6679-7425-40de-944b-e07fc1f90ae7`） |
| `tasks[].user_id` | string | ○ | 所有ユーザーID（UUID形式、例: `550e8400-e29b-41d4-a716-446655440000`） |
| `tasks[].title` | string | ○ | タスクタイトル |
| `tasks[].comment` | string | ○ | タスクコメント（未入力時は空文字 `""`） |
| `tasks[].priority` | string | ○ | 優先度（`high`, `medium`, `low`） |
| `tasks[].status` | string | ○ | ステータス（初期値: `not_started`） |
| `tasks[].due_datetime` | string / null | ○ | 締切日時（ISO 8601 JST 形式。締切未設定時は `null`） |
| `tasks[].is_pinned` | boolean | ○ | ピン留めフラグ |
| `tasks[].created_at` | string | ○ | 作成日時（ISO 8601 JST 形式） |
| `tasks[].updated_at` | string | ○ | 更新日時（ISO 8601 JST 形式） |

※単一タスク作成時も `created_count: 1` および要素数1の `tasks` 配列を返却します。繰り返し一括生成時に返却される `tasks` 配列は `due_datetime` の昇順（時系列順）でソートされて返却されます。

##### Errors
- `400 Bad Request`: リクエスト形式またはバリデーション不正（code: `"BAD_REQUEST"`）
  - `is_recurring: true` 時の `recurring_rule` 欠落・null・非オブジェクト指定時: `error.details: [{ "field": "recurring_rule", "message": "is_recurringがtrueの場合、recurring_ruleオブジェクトの指定は必須です" }]`
  - 生成件数0件時: `error.details: [{ "field": "recurring_rule", "message": "指定された期間内に該当する曜日が存在しません" }]`
  - 生成件数100件超過時: `error.details: [{ "field": "recurring_rule", "message": "生成件数が上限（100件）を超えています" }]`
  - `due_time` 形式不正時: `error.details: [{ "field": "recurring_rule.due_time", "message": "締切時刻の形式が不正です（HH:mm形式で指定してください）" }]`
  - タイトル文字数違反、日付範囲不整合（`start_date > end_date` または 1年超）、無効な曜日指定、`priority` への `null` / 定義外列挙値指定、`is_pinned` への `null` / 非 boolean 指定等
  - ※ `recurring_rule` オブジェクト内のネストされた個別フィールドでバリデーションエラーが発生した場合（例: `start_date`, `end_date`, `days_of_week`, `due_time`）、`error.details[].field` にはドット記法（例: `"recurring_rule.start_date"`, `"recurring_rule.due_time"`, `"recurring_rule.days_of_week"` 等）で該当キー名が設定されます。
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）

---

#### 3.3.3 `GET tasks/{task_id}`
指定されたタスクの詳細情報を取得します。タスク一覧画面やカレンダー日付詳細ポップアップからの「タスク編集ポップアップ」展開時、および直接アクセス時に最新のタスク詳細情報を取得するために利用されます。

- **認証**: 必須（Cookie）
- **Path Parameters**:
  - `task_id` (string): タスクID

##### リクエスト評価順序
1. **認証検証 (`401 Unauthorized`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `task_id` の存在およびセッションユーザーの所有タスクかを検証（不一致または存在しない場合は 404 `NOT_FOUND`）。

##### Response (200 OK)
```json
{
  "task": {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
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

##### Response Body フィールド定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `task` | object | ○ | 詳細取得されたタスクオブジェクト |
| `task.id` | string | ○ | タスクID（UUID形式、例: `7c9e6679-7425-40de-944b-e07fc1f90ae7`） |
| `task.user_id` | string | ○ | 所有ユーザーID（UUID形式、例: `550e8400-e29b-41d4-a716-446655440000`） |
| `task.title` | string | ○ | タスクタイトル |
| `task.comment` | string | ○ | タスクコメント（未入力時は空文字 `""`） |
| `task.priority` | string | ○ | 優先度（`high`, `medium`, `low`） |
| `task.status` | string | ○ | ステータス（`not_started`, `in_progress`, `completed`） |
| `task.due_datetime` | string / null | ○ | 締切日時（ISO 8601 JST 形式。締切未設定時は `null`） |
| `task.is_pinned` | boolean | ○ | ピン留めフラグ |
| `task.created_at` | string | ○ | 作成日時（ISO 8601 JST 形式） |
| `task.updated_at` | string | ○ | 更新日時（ISO 8601 JST 形式） |

※締切日時が未設定のタスクの場合、`due_datetime` は `null` として返却されます。またコメントが未入力の場合、`comment` は空文字 `""` として返却されます。

##### Errors
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `404 Not Found`: 存在しないタスクまたは他ユーザー所有タスク（code: `"NOT_FOUND"`）

---

#### 3.3.4 `PATCH tasks/{task_id}`
タスク情報を部分更新します。リクエストボディに含まれるフィールドのみが更新対象となります（ステータス変更 `status` やピン留め `is_pinned` の単体更新含む）。タイトル（`title`）またはコメント（`comment`）が更新された場合は、バックエンドで検索用正規化文字列（小文字化＋NFKC正規化＋ひらがな→カタカナ変換）を再生成し、`TASK.SEARCH_TEXT` カラムを自動更新します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `task_id` (string): 更新対象タスクID

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `task_id` の存在およびセッションユーザーの所有タスクかを検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
3. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）への `null` 指定有無、文字数・列挙値制約を検証。不備がある場合は 400 `BAD_REQUEST` を返却。なお、システムの読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）が含まれている場合は、エラーとせず更新対象外として単に無視します。

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
※ステータスのみを変更する場合は `{"status": "completed"}`、ピン留めのみを変更する場合は `{"is_pinned": true}` のように、更新対象のフィールドのみを指定して送信可能です。なお、更新対象フィールドが一つも指定されていない空リクエストボディ `{}` が送信された場合、変更を行わずに `200 OK` と現在のタスク情報を返却します。また、非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）に対して明示的に `null`（例: `{"title": null}`, `{"status": null}`）が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却します。リクエストボディに読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）が含まれる場合は更新対象外として無視されます。なお、`title` または `comment` が更新された場合、バックエンドで検索用正規化文字列（小文字化＋NFKC正規化＋ひらがな→カタカナ変換）を再生成し、`TASK.SEARCH_TEXT` カラムを自動更新します。

##### Request Body フィールド定義

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `title` | string | × | 前後の空白文字（半角・全角スペース、タブ、改行）を除去（トリム）した上で1〜100文字必須。空白文字のみの入力および文字列内部への改行・タブ等の制御文字入力は 400 Bad Request（code: `"BAD_REQUEST"`）として拒否 |
| `comment` | string / null | × | 0〜1000文字（トリム後）。改行は `\n` に正規化。空文字 `""` または `null` 指定でコメントをクリア（削除） |
| `priority` | string | × | `high`, `medium`, `low` |
| `status` | string | × | `not_started`, `in_progress`, `completed` |
| `due_datetime` | string / null | × | ISO 8601 日時文字列、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を設定。※タイムゾーンオフセットを含まない ISO 8601 文字列が指定された場合はデフォルトで JST (`+09:00`) と解釈し、UTC (`Z`) や他タイムゾーン指定時は JST (`+09:00`) に変換・正規化して登録・返却。`null` 指定で締切解除） |
| `is_pinned` | boolean | × | ピン留め状態（`true` / `false`） |

##### Response (200 OK)
```json
{
  "task": {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
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

##### Response Body フィールド定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `task` | object | ○ | 更新後のタスクオブジェクト |
| `task.id` | string | ○ | タスクID（UUID形式、例: `7c9e6679-7425-40de-944b-e07fc1f90ae7`） |
| `task.user_id` | string | ○ | 所有ユーザーID（UUID形式、例: `550e8400-e29b-41d4-a716-446655440000`） |
| `task.title` | string | ○ | タスクタイトル |
| `task.comment` | string | ○ | タスクコメント（未入力時は空文字 `""`） |
| `task.priority` | string | ○ | 優先度（`high`, `medium`, `low`） |
| `task.status` | string | ○ | ステータス（`not_started`, `in_progress`, `completed`） |
| `task.due_datetime` | string / null | ○ | 締切日時（ISO 8601 JST 形式。締切未設定時は `null`） |
| `task.is_pinned` | boolean | ○ | ピン留めフラグ |
| `task.created_at` | string | ○ | 作成日時（ISO 8601 JST 形式） |
| `task.updated_at` | string | ○ | 更新日時（ISO 8601 JST 形式） |

##### Errors
- `400 Bad Request`: バリデーション不正（文字数違反、ステータス不正値等、code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー（存在しないタスクまたは他者所有タスク、code: `"NOT_FOUND"`）

---

#### 3.3.5 `DELETE tasks/{task_id}`
タスクをDBから物理削除します。

- **認証**: 必須（Cookie）
- **Headers**: `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `task_id` (string): 削除対象タスクID

##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `task_id` の存在およびセッションユーザーの所有タスクかを検証（不一致または存在しない場合は 404 `NOT_FOUND`）。

##### Response (200 OK)
```json
{
  "message": "Task has been deleted successfully."
}
```

##### Response Body フィールド定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `message` | string | ○ | 削除完了メッセージ（例: `"Task has been deleted successfully."`） |

##### Errors
- `401 Unauthorized`: 未ログイン（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー（存在しないタスクまたは他者所有タスク、code: `"NOT_FOUND"`）

---
