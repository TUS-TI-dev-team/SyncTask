# PATCH /tasks/{task_id} エンドポイント開発計画書

本ドキュメントでは、SyncTask バックエンドにおけるタスク部分更新 API（`PATCH /api/tasks/:task_id`）の実装計画を定義します。
セッションユーザーに基づく所有権検証（IDOR/BOLA防止）、部分更新（指定フィールドのみの変更、空ボディ `{}` 時の既存値返却、非Null許容フィールドへの明示的 null に対する 400 Bad Request、コメント/締切日時の null クリア、検索用正規化文字列 `SEARCH_TEXT` の再生成）、ならびに `backend/TESTING_GUIDE.md` に準拠した単体テスト先行（TDD）の作成・検証手順を具体的に定めます。

---

## 1. 概要・要件定義

### 1.1 エンドポイント仕様概要
- **パス / メソッド**: `PATCH /api/tasks/:task_id`（ベースURL `/api/` 配下）
- **認証**: 必須（Cookie セッション）
- **機能概要**:
  - 指定された `task_id` のタスク情報を部分更新します。
  - リクエストボディに含まれるフィールドのみを更新対象とします。
  - 更新対象フィールドが一つも指定されていない空リクエストボディ `{}` が送信された場合、変更を行わずに `200 OK` と現在のタスク情報を返却します。
  - システムの読み取り専用フィールド（`id`, `user_id`, `created_at`, `updated_at`）が含まれている場合は、エラーとせず更新対象外として単に無視します。
  - `title` または `comment` が更新された場合、バックエンドで検索用正規化文字列（小文字化＋NFKC正規化＋ひらがな→カタカナ変換）を再生成し、`TASK.SEARCH_TEXT` カラムを自動更新します。
  - 存在しないタスク、または他ユーザーが所有するタスクが指定された場合は、情報漏洩を防ぐため同一の `404 NOT_FOUND` を返却します（IDOR/BOLA防止）。

### 1.2 リクエスト仕様
- **パスパラメータ**:
  - `task_id` (string): 更新対象タスクID（UUID形式）
- **リクエストボディ**:
  - JSON形式。各フィールドは省略可能。

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `title` | string | × | 前後の空白文字（半角・全角スペース、タブ、改行）を除去（トリム）した上で1〜100文字必須。空白文字のみの入力および文字列内部への改行・タブ等の制御文字入力は 400 Bad Request（code: `"BAD_REQUEST"`）。明示的な `null` 指定は不可（400 Bad Request）。 |
| `comment` | string / null | × | 0〜1000文字（トリム後）。改行は `\n` に正規化。空文字 `""` または `null` 指定でコメントをクリア（空文字 `""` に更新）。 |
| `priority` | string | × | `high`, `medium`, `low` のいずれか。明示的な `null` または定義外列挙値指定は不可（400 Bad Request）。 |
| `status` | string | × | `not_started`, `in_progress`, `completed` のいずれか。明示的な `null` または定義外列挙値指定は不可（400 Bad Request）。 |
| `due_datetime` | string / null | × | ISO 8601 日時文字列、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を設定。※タイムゾーンオフセットを含まない場合は JST (`+09:00`) と解釈し、UTC (`Z`) や他タイムゾーン指定時は JST (`+09:00`) に変換・正規化）。`null` 指定で締切日時をクリア（`null` に更新）。形式不正時は 400 Bad Request。 |
| `is_pinned` | boolean | × | ピン留め状態（`true` / `false`）。明示的な `null` または非 boolean 型指定は不可（400 Bad Request）。 |

### 1.3 レスポンス仕様 (200 OK)
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

### 1.4 エラーレスポンス仕様
- **`400 Bad Request`**: バリデーション不正、非Null許容フィールドへの明示的 null 指定、型不正、JSON構文エラー等
- **`401 Unauthorized`**: 未ログインまたはセッション無効（code: `"UNAUTHORIZED"`）
- **`404 Not Found`**: 存在しないタスクまたは他ユーザー所有タスク（code: `"NOT_FOUND"`）
- **`500 Internal Server Error`**: サーバー内部エラー（code: `"INTERNAL_SERVER_ERROR"`）

---

## 2. アーキテクチャ・設計方針

### 2.1 パッケージ構成
```
backend/
├── model/
│   ├── task_patch.go        # PatchTaskRequest / PatchTaskResponse / UnmarshalJSON / Validate()
│   └── task_patch_test.go   # PatchTaskRequest バリデーション単体テスト
├── repository/
│   ├── task.go              # TaskRepository インターフェースに UpdateTask を追加
│   ├── task_patch.go        # UpdateTask 実装
│   └── task_patch_test.go   # UpdateTask 単体テスト（sqlmock）
├── service/
│   ├── task.go              # TaskService インターフェースに PatchTask を追加
│   ├── task_patch.go        # PatchTask 実装
│   └── task_patch_test.go   # PatchTask 単体テスト
├── handler/
│   ├── task_patch.go        # PatchTaskHandler 実装
│   └── task_patch_test.go   # PatchTaskHandler 単体テスト
└── router/
    ├── router.go            # PATCH /api/tasks/:task_id のルーティング追加
    └── router_test.go       # ルーター単体テスト
```

### 2.2 部分更新・バリデーション制御設計 (PatchTaskRequest)
1. **未指定 vs 明示的 null vs 値の識別**:
   - `PatchTaskRequest` に `UnmarshalJSON(data []byte) error` を実装し、`map[string]json.RawMessage` を経由して各フィールドのキー存在と値を判定します。
   - `title`, `priority`, `status`, `is_pinned` にキーが存在し値が `null` の場合は、即座にバリデーションエラー（400 Bad Request）として記録。
   - `comment` にキーが存在し値が `null` の場合は、空文字 `""` への更新として処理。
   - `due_datetime` にキーが存在し値が `null` の場合は、締切クリア（`nil`）として処理。
   - 各フィールドの型違反（例: `is_pinned` に文字列や数値が渡された場合など）も 400 Bad Request として記録。
   - `id`, `user_id`, `created_at`, `updated_at` が含まれていても無視。
2. **更新ロジック (Service層)**:
   - `repo.GetTaskByID(ctx, userID, taskID)` で既存タスクを取得。
   - レコードが存在しない場合は `model.NewNotFoundError("指定されたタスクが見つかりません。")` を返却（IDOR/BOLA防止）。
   - 更新フィールドが一つも指定されていない場合（`req.HasChanges() == false`）、DB更新を行わず既存タスクをそのまま `200 OK` で返却。
   - 指定されたフィールドのみ既存タスク構造体に上書き反映。
   - `title` または `comment` が更新された場合、`util.NormalizeSearchText(task.Title, task.Comment)` で `SearchText` を再計算。
   - `task.UpdatedAt = time.Now().In(jst)` を設定。
   - `repo.UpdateTask(ctx, task)` を実行し、更新後のタスクを返却。

### 2.3 リポジトリ層設計 (UpdateTask)
```sql
UPDATE TASK
SET
    TITLE = $1,
    PRIORITY = $2,
    DUE_DATETIME = $3,
    STATUS = $4,
    IS_PINNED = $5,
    COMMENT = $6,
    SEARCH_TEXT = $7,
    UPDATED_AT = $8
WHERE TASK_ID = $9 AND USER_ID = $10
RETURNING
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
```

---

## 3. 開発手順・ステップ

### Step 1: テストデータ・単体テストプログラムの作成 (TDD先行作成)
`backend/TESTING_GUIDE.md` の規約（ファイル命名、日本語 `t.Run`、`require`/`assert` の使い分け、`@spec` 連携）に従い、以下の単体テストを実装します。

#### 1. `backend/model/task_patch_test.go`
- `正常系: 全フィールド指定時に正常にデコードおよびバリデーションを通過すること`
- `正常系: 単一フィールド（statusのみ、is_pinnedのみ等）指定時に正常にデコードされること`
- `正常系: 空ボディ {} 指定時に HasChanges が false となりバリデーションを通過すること`
- `正常系: 読み取り専用フィールド（id, user_id 等）が含まれていても無視されて通過すること`
- `正常系: comment に null または空文字が指定された場合クリア対象となること`
- `正常系: due_datetime に null が指定された場合クリア対象となること`
- `正常系: due_datetime に YYYY-MM-DD 形式が指定された場合 JST 23:59:00 として解釈されること`
- `異常系: title に null が指定された場合に 400 エラーを返すこと`
- `異常系: title が空文字または空白のみの場合に 400 エラーを返すこと`
- `異常系: title に改行やタブが含まれる場合に 400 エラーを返すこと`
- `異常系: title が 100 文字を超える場合に 400 エラーを返すこと`
- `異常系: comment が 1000 文字を超える場合に 400 エラーを返すこと`
- `異常系: priority に null または無効な値が指定された場合に 400 エラーを返すこと`
- `異常系: status に null または無効な値が指定された場合に 400 エラーを返すこと`
- `異常系: is_pinned に null または非 boolean 値が指定された場合に 400 エラーを返すこと`
- `異常系: due_datetime の形式が不正な場合に 400 エラーを返すこと`
- `異常系: JSON の形式自体が不正な場合に 400 エラーを返すこと`

#### 2. `backend/repository/task_patch_test.go` (`sqlmock` 利用)
- `正常系: タスクが正常に更新され、更新後のモデルが返却されること`
- `準正常系: 更新対象レコードが存在しない（または他者所有）場合に nil が返却されること`
- `異常系: DBクエリエラー発生時にエラーが返却されること`

#### 3. `backend/service/task_patch_test.go`
- `正常系: 指定フィールドが正常に更新され PatchTaskResponse が返却されること`
- `正常系: title または comment 更新時に SearchText が再生成されること`
- `正常系: 空リクエストボディ {} の場合に DB 更新をスキップして既存タスクがそのまま返却されること`
- `正常系: comment に null 指定時にコメントが空文字にクリアされること`
- `正常系: due_datetime に null 指定時に締切日時が nil にクリアされること`
- `異常系: リクエストバリデーションエラー時に 400 BAD_REQUEST が返却されること`
- `異常系: 対象タスクが存在しない（または他者所有）場合に 404 NOT_FOUND が返却されること`
- `異常系: リポジトリ層でエラーが発生した場合にエラーがそのまま返却されること`

#### 4. `backend/handler/task_patch_test.go`
- `正常系: 有効なリクエストで 200 OK と更新後のタスク詳細が返却されること`
- `正常系: 空ボディ {} でも 200 OK と現在のタスク詳細が返却されること`
- `異常系: 未ログイン（Context に userID なし）の場合に 401 UNAUTHORIZED を返すこと`
- `異常系: JSON 構文不正またはバリデーションエラー時に 400 BAD_REQUEST を返すこと`
- `異常系: 対象タスクが存在しない（または他者所有）場合に 404 NOT_FOUND を返すこと`
- `異常系: サーバー内部エラー発生時に 500 INTERNAL_SERVER_ERROR を返すこと`

#### 5. `backend/router/router_test.go`
- `正常系: PATCH /api/tasks/:task_id が登録されておりリクエストをルーティングできること`

---

### Step 2: プログラム本体の実装

#### 1. モデル層の実装
- `backend/model/task_patch.go`:
  - `PatchTaskRequest`:
    - フィールド: `Title`, `Comment`, `Priority`, `Status`, `DueDatetime`, `IsPinned`
    - フィールド指定フラグおよびクリアフラグ
    - `UnmarshalJSON(data []byte) error` 実装
    - `Validate() error` 実装
    - `HasChanges() bool` 実装
    - 適用ヘルパー `ApplyTo(task *Task)` 実装
  - `PatchTaskResponse`:
    ```go
    type PatchTaskResponse struct {
        Task *Task `json:"task"`
    }
    ```

#### 2. リポジトリ層の実装
- `backend/repository/task.go`:
  - `TaskRepository` インターフェースに `UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error)` を追加
- `backend/repository/task_patch.go`:
  - `UpdateTask` メソッドを実装

#### 3. サービス層の実装
- `backend/service/task.go`:
  - `TaskService` インターフェースに `PatchTask(ctx context.Context, userID, taskID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error)` を追加
- `backend/service/task_patch.go`:
  - `PatchTask` メソッドを実装

#### 4. ハンドラー層の実装
- `backend/handler/task_patch.go`:
  - `PatchTaskHandler(service service.TaskService) gin.HandlerFunc` を実装
  - Swagger アノテーション（`@Summary`, `@Tags`, `@Router /api/tasks/{task_id} [patch]` 等）を付与

#### 5. ルーターへの登録
- `backend/router/router.go`:
  - `api.PATCH("/tasks/:task_id", handler.PatchTaskHandler(taskService))` を登録

---

### Step 3: 単体テスト実行・検証
```bash
cd backend
go test ./...
go test -v ./...
go test -cover ./...
```
**成功基準**:
1. 全パッケージの単体テストが PASS すること。
2. 仕様書（`docs/design/api_design/04_tasks.md`）に記載された全ケース（部分更新、空ボディ、null禁止フィールド、nullクリアフィールド、IDOR検証、400/401/404/500）が検証されていること。

---

### Step 4: プログラムの修正・リファクタリング（テスト失敗時のサイクル）
テスト失敗が発生した場合、ログを確認してコード修正を行い、再度 Step 3 を実行する。
Code-as-Docs 原則に基づき、GoDoc コメントの `@spec` 記述を整備する。

---

## 4. 変更対象ファイル一覧

| 操作 | ファイルパス | 説明 |
|---|---|---|
| 新規 | `backend/model/task_patch.go` | `PatchTaskRequest`, `PatchTaskResponse`, `UnmarshalJSON`, `Validate` |
| 新規 | `backend/model/task_patch_test.go` | リクエストモデルのバリデーション・デコード単体テスト |
| 変更 | `backend/repository/task.go` | `TaskRepository` IF に `UpdateTask` 追加 |
| 新規 | `backend/repository/task_patch.go` | `UpdateTask` メソッド実装 |
| 新規 | `backend/repository/task_patch_test.go` | `UpdateTask` 単体テスト（sqlmock） |
| 変更 | `backend/service/task.go` | `TaskService` IF に `PatchTask` 追加 |
| 新規 | `backend/service/task_patch.go` | `PatchTask` ビジネスロジック実装 |
| 新規 | `backend/service/task_patch_test.go` | `PatchTask` 単体テスト |
| 新規 | `backend/handler/task_patch.go` | `PatchTaskHandler` 実装 |
| 新規 | `backend/handler/task_patch_test.go` | `PatchTaskHandler` 単体テスト |
| 変更 | `backend/router/router.go` | ルート登録（`PATCH /api/tasks/:task_id`） |
| 変更 | `backend/router/router_test.go` | ルーター単体テスト（PATCH ルート検証） |
