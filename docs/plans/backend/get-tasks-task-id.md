# GET /tasks/{task_id} エンドポイント開発計画書

本ドキュメントでは、SyncTask バックエンドにおけるタスク詳細取得 API（`GET /api/tasks/:task_id`）の実装計画を定義します。
セッションユーザーに基づく所有権検証（IDOR/BOLA防止）、404 Not Found 制御、ならびに `backend/TESTING_GUIDE.md` に準拠した単体テスト先行（TDD）の作成・検証手順を具体的に定めます。

---

## 1. 概要・要件定義

### 1.1 エンドポイント仕様概要
- **パス / メソッド**: `GET /api/tasks/:task_id`（ベースURL `/api/` 配下）
- **認証**: 必須（Cookie セッション）
- **機能概要**:
  - 指定された `task_id` に該当するタスクの詳細情報を取得して返却します。
  - セッションユーザーが所有するタスクのみ取得可能です。
  - 存在しないタスク、または他ユーザーが所有するタスクが指定された場合は、情報漏洩を防ぐため同一の `404 NOT_FOUND` を返却します（IDOR/BOLA防止）。

### 1.2 リクエスト仕様
- **パスパラメータ**:
  - `task_id` (string): 取得対象のタスクID（UUID形式）
- **リクエストボディ / クエリパラメータ**: なし

### 1.3 レスポンス仕様 (200 OK)
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

#### レスポンスボディ フィールド定義
| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `task` | object | ○ | 詳細取得されたタスクオブジェクト |
| `task.id` | string | ○ | タスクID（UUID形式） |
| `task.user_id` | string | ○ | 所有ユーザーID（UUID形式） |
| `task.title` | string | ○ | タスクタイトル |
| `task.comment` | string | ○ | タスクコメント（未入力時は空文字 `""`） |
| `task.priority` | string | ○ | 優先度（`high`, `medium`, `low`） |
| `task.status` | string | ○ | ステータス（`not_started`, `in_progress`, `completed`） |
| `task.due_datetime` | string / null | ○ | 締切日時（ISO 8601 JST 形式。締切未設定時は `null`） |
| `task.is_pinned` | boolean | ○ | ピン留めフラグ |
| `task.created_at` | string | ○ | 作成日時（ISO 8601 JST 形式） |
| `task.updated_at` | string | ○ | 更新日時（ISO 8601 JST 形式） |

### 1.4 エラーレスポンス仕様
#### 401 Unauthorized
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "認証が必要です。",
    "details": []
  }
}
```

#### 404 Not Found
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "指定されたタスクが見つかりません。",
    "details": []
  }
}
```

#### 500 Internal Server Error
```json
{
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "サーバー内部でエラーが発生しました。",
    "details": []
  }
}
```

---

## 2. アーキテクチャ・設計方針

### 2.1 パッケージ構成と責務
各レイヤーの肥大化を防ぎ、保守性・可読性を向上させるため、同一パッケージ内で操作単位（`task_<action>.go` / `task_<action>_test.go`）にファイルを分割します。

```
backend/
├── model/
│   ├── error.go             # NewNotFoundError(message string) の追加
│   ├── error_test.go        # エラーモデル単体テスト
│   ├── task.go              # Task 共通モデル（エンティティ・定数など）
│   ├── task_create.go       # CreateTaskRequest / CreateTaskResponse / Validate()（既存から分離）
│   ├── task_create_test.go  # CreateTaskRequest バリデーション単体テスト（既存から分離）
│   └── task_get.go          # GetTaskResponse 構造体の追加（新規）
├── repository/
│   ├── task.go              # TaskRepository インターフェース定義、NewTaskRepository
│   ├── task_create.go       # CreateTask, CreateTasks（既存から分離）
│   ├── task_create_test.go  # CreateTask, CreateTasks 単体テスト（既存から分離）
│   ├── task_get.go          # GetTaskByID(ctx, userID, taskID) の追加（新規）
│   └── task_get_test.go     # リポジトリ層詳細取得単体テスト（sqlmock）（新規）
├── service/
│   ├── task.go              # TaskService インターフェース定義、NewTaskService、共通定義
│   ├── task_create.go       # CreateTask 実装（既存から分離）
│   ├── task_create_test.go  # CreateTask 単体テスト（既存から分離）
│   ├── task_get.go          # GetTask(ctx, userID, taskID) の追加（新規）
│   └── task_get_test.go     # サービス層詳細取得単体テスト（新規）
├── handler/
│   ├── task_create.go       # CreateTaskHandler（既存から分離）
│   ├── task_create_test.go  # CreateTaskHandler 単体テスト（既存から分離）
│   ├── task_get.go          # GetTaskHandler の追加（新規）
│   └── task_get_test.go     # ハンドラー層詳細取得単体テスト（新規）
└── router/
    ├── router.go            # GET /api/tasks/:task_id のルーティング追加
    └── router_test.go       # ルーター層単体テスト
```

### 2.2 認可・データアクセス制御 (IDOR防止)
- Repository 層のクエリにおいて、必ず `WHERE TASK_ID = $1 AND USER_ID = $2` を指定します。
- これにより、指定した `task_id` が存在しない場合、または他ユーザー所有のタスクである場合の両方で `sql.ErrNoRows` となり、Service 層で `404 NOT_FOUND`（`model.NewNotFoundError`）へマッピングされます。

---

## 3. 開発手順・ステップ

### Step 0: 既存タスク処理のファイル分割・リファクタリング
新機能実装に先立ち、既存の作成関連コードとテストを操作別ファイルへ切り出します。リファクタリング後、既存テストが全件 PASS することを確認します。

1. **`backend/model/`**:
   - `model/task.go`: `Task` コアモデル定義を残す。
   - `model/task_create.go`: `CreateTaskRequest`, `CreateTaskResponse`, `Validate()` を移動。
   - `model/task_create_test.go`: `model/task_test.go` から作成用バリデーションテストを移動。
2. **`backend/repository/`**:
   - `repository/task.go`: `TaskRepository` インターフェース定義、`NewTaskRepository` を残す。
   - `repository/task_create.go`: `CreateTask`, `CreateTasks` 実装を移動。
   - `repository/task_create_test.go`: `repository/task_test.go` から作成用単体テストを移動。
3. **`backend/service/`**:
   - `service/task.go`: `TaskService` インターフェース定義、`NewTaskService`、共通変数（`jst`, `weekdayMap`）を残す。
   - `service/task_create.go`: `CreateTask` 実装を移動。
   - `service/task_create_test.go`: `service/task_test.go` から作成用単体テストを移動。
4. **`backend/handler/`**:
   - `handler/task_create.go`: `CreateTaskHandler` 実装を移動。
   - `handler/task_create_test.go`: `handler/task_test.go` から作成用単体テストを移動。

---

### Step 1: テストデータ・単体テストプログラムの作成 (TDD先行作成)
`backend/TESTING_GUIDE.md` の規約（ファイル命名、日本語 `t.Run`、`require`/`assert` の使い分け、`@spec` 連携）に従い、以下の単体テストを実装します。

#### 1. `backend/model/error_test.go`
- `正常系: NewNotFoundError で 404 ステータスと NOT_FOUND コードを持つ AppError が生成されること`

#### 2. `backend/repository/task_get_test.go` (`sqlmock` 利用)
- `正常系: 指定した task_id と user_id に一致するタスクが取得できること`
- `準正常系: 該当レコードが存在しない場合（sql.ErrNoRows）に nil が返却されること`
- `異常系: DBクエリエラー発生時にエラーが返却されること`

#### 3. `backend/service/task_get_test.go`
- `正常系: 指定した task_id のタスクが正常に取得され GetTaskResponse が返却されること`
- `異常系: リポジトリが nil を返した場合（該当タスクなし）に 404 NOT_FOUND エラーが返却されること`
- `異常系: リポジトリ層で予期せぬエラーが発生した場合にエラーがそのまま返却されること`

#### 4. `backend/handler/task_get_test.go`
- `正常系: 有効な task_id 指定時に 200 OK とタスク詳細JSONが返却されること`
- `異常系: 未ログイン（Context に userID なし）の場合に 401 UNAUTHORIZED を返すこと`
- `異常系: タスクが存在しない（または他者所有）場合に 404 NOT_FOUND を返すこと`
- `異常系: サーバー内部エラー発生時に 500 INTERNAL_SERVER_ERROR を返すこと`

#### 5. `backend/router/router_test.go`
- `正常系: GET /api/tasks/:task_id が登録されておりリクエストをルーティングできること`

---

### Step 2: プログラムの実装

#### 1. モデル層の実装
- `backend/model/error.go`:
  - `NewNotFoundError(message string) *AppError` を実装
- `backend/model/task_get.go`:
  - `GetTaskResponse` 構造体を定義（Swaggerタグ・example含む）
    ```go
    type GetTaskResponse struct {
        Task *Task `json:"task"`
    }
    ```

#### 2. リポジトリ層の実装
- `backend/repository/task.go`:
  - `TaskRepository` インターフェースに `GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error)` を追加
- `backend/repository/task_get.go`:
  - `selectTaskByIDQuery` 定義および `GetTaskByID` メソッドを実装:
    ```sql
    SELECT
        TASK_ID,
        USER_ID,
        TITLE,
        PRIORITY,
        DUE_DATETIME,
        STATUS,
        IS_PINNED,
        COMMENT,
        SEARCH_TEXT,
        CREATED_AT,
        UPDATED_AT
    FROM TASK
    WHERE TASK_ID = $1 AND USER_ID = $2
    ```
  - `sql.ErrNoRows` の場合は `nil, nil` を返却

#### 3. サービス層の実装
- `backend/service/task.go`:
  - `TaskService` インターフェースに `GetTask(ctx context.Context, userID, taskID string) (*model.GetTaskResponse, error)` を追加
- `backend/service/task_get.go`:
  - `GetTask` メソッドを実装。`repo.GetTaskByID` を呼び出し、タスクが取得できなかった場合は `model.NewNotFoundError("指定されたタスクが見つかりません。")` を返却

#### 4. ハンドラー層の実装
- `backend/handler/task_get.go`:
  - `GetTaskHandler(service service.TaskService) gin.HandlerFunc` を実装
  - `c.GetString("userID")` の認証チェック（未認証時 401）
  - `c.Param("task_id")` からパスパラメータを取得
  - `service.GetTask` を実行し、結果を `200 OK` で返却
  - `model.AppError` を判定して適切な HTTP ステータスコード（404/500等）とエラーレスポンスを返却

#### 5. ルーターへの登録
- `backend/router/router.go`:
  - `api.GET("/tasks/:task_id", handler.GetTaskHandler(taskService))` を追加

---

### Step 3: テスト実行・検証
`backend/TESTING_GUIDE.md` に従い、単体テストを実行・検証します。

```bash
# バックエンドディレクトリで実行
cd backend

# 全単体テストの実行
go test ./...

# 詳細出力（Verbose）での確認
go test -v ./...

# カバレッジの計測と確認
go test -cover ./...
```

**成功基準**:
1. すべてのパッケージ（`model`, `service`, `repository`, `handler`, `router`）の単体テストが PASS すること。
2. 正常系、認可エラー/非存在（404）、未認証（401）、サーバーエラー（500）のシナリオがすべて網羅されていること。

---

### Step 4: プログラムの修正・リファクタリング（テスト失敗時のサイクル）
1. テスト失敗が発生した場合は、エラーログおよびアサーション差分を確認。
2. 実装コード（またはテストコード）の不整合を修正。
3. 再度 `go test -v ./...` を実行し、全テスト通過を確認。
4. GoDoc コメント（`@spec` 記法）を追記・整備。

---

## 4. 変更対象ファイル一覧

| 操作 | ファイルパス | 説明 |
|---|---|---|
| 変更 | `backend/model/error.go` | `NewNotFoundError` の追加 |
| 変更 | `backend/model/error_test.go` | `NewNotFoundError` の単体テスト追加 |
| 変更/分割 | `backend/model/task.go` | コアモデル定義の整理 |
| 新規/移動 | `backend/model/task_create.go` | `CreateTaskRequest` / `CreateTaskResponse` の切り出し |
| 新規/移動 | `backend/model/task_create_test.go` | 作成バリデーション単体テストの切り出し |
| 新規 | `backend/model/task_get.go` | `GetTaskResponse` 構造体の追加 |
| 変更/分割 | `backend/repository/task.go` | `TaskRepository` IF に `GetTaskByID` 追加、共通部整理 |
| 新規/移動 | `backend/repository/task_create.go` | `CreateTask`, `CreateTasks` の切り出し |
| 新規/移動 | `backend/repository/task_create_test.go` | 作成用単体テストの切り出し（sqlmock） |
| 新規 | `backend/repository/task_get.go` | `GetTaskByID` 実装の追加 |
| 新規 | `backend/repository/task_get_test.go` | `GetTaskByID` 単体テストの追加（sqlmock） |
| 変更/分割 | `backend/service/task.go` | `TaskService` IF に `GetTask` 追加、共通部整理 |
| 新規/移動 | `backend/service/task_create.go` | `CreateTask` 実装の切り出し |
| 新規/移動 | `backend/service/task_create_test.go` | `CreateTask` 単体テストの切り出し |
| 新規 | `backend/service/task_get.go` | `GetTask` 実装の追加 |
| 新規 | `backend/service/task_get_test.go` | `GetTask` 単体テストの追加 |
| 新規/移動 | `backend/handler/task_create.go` | `CreateTaskHandler` の切り出し |
| 新規/移動 | `backend/handler/task_create_test.go` | `CreateTaskHandler` 単体テストの切り出し |
| 新規 | `backend/handler/task_get.go` | `GetTaskHandler` の追加 |
| 新規 | `backend/handler/task_get_test.go` | `GetTaskHandler` 単体テストの追加 |
| 変更 | `backend/router/router.go` | ルート登録（`GET /api/tasks/:task_id`） |
| 変更 | `backend/router/router_test.go` | ルーター単体テスト更新 |
