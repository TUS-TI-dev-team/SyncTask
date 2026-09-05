# タスク部分更新 API におけるリクエスト評価順序（認可チェックと入力バリデーション）の仕様乖離

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-09-05 21:47:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [task_patch.go](backend/handler/task_patch.go)
  - [task_patch.go](backend/service/task_patch.go)

## 1. 問題の概要
API 設計書 `docs/design/api_design/04_tasks.md` の 3.3.4 節 (`PATCH tasks/{task_id}`) における「リクエスト評価順序」では、認可チェック（IDOR/BOLA 検証・404）が入力バリデーション（400）より先に評価されるよう規定されています。
しかし、現在の実装コードではハンドラーおよびサービス層において、認可チェック（DB 存在・所有権確認）の前にリクエストボディの JSON パースおよびバリデーション（`req.Validate()`）が実行されており、仕様書と実装で評価順序の逆転（乖離）が発生しています。

## 2. 詳細な指摘内容
1. **仕様書の規定 (`docs/design/api_design/04_tasks.md` L251-L258)**:
   ```markdown
   ##### リクエスト評価順序
   1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
      ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
   2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
      パスパラメータ `task_id` の存在およびセッションユーザーの所有タスクかを検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
   3. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
      リクエストボディの JSON 形式、非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）への `null` 指定有無、文字数・列挙値制約を検証。不備がある場合は 400 `BAD_REQUEST` を返却。
   ```
2. **実装コード (`backend/service/task_patch.go` L21-L32)**:
   ```go
   func (s *taskService) PatchTask(ctx context.Context, userID, taskID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
       if err := req.Validate(); err != nil {
           return nil, err
       }

       existing, err := s.repo.GetTaskByID(ctx, userID, taskID)
       if err != nil {
           return nil, err
       }
       if existing == nil {
           return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
       }
   ```
3. **具体的な問題点**:
   - 存在しない `task_id`、または他ユーザーが所有する `task_id` を指定したリクエストにおいて、不正なボディ（例: `{"title": ""}` や `{"priority": "invalid"}`）が送信された場合、仕様書上は Step 2 で `404 Not Found` が返却されるべきですが、実装上は `400 Bad Request` が返却されます。
   - 攻撃者が不正なリクエストボディをプローブとして送信することで、該当 `task_id` が実在するかどうかをエラーコードの違い（400 vs 404）から推測できるリソース列挙（Resource Enumeration）リスクを低減するために仕様書で認可チェックが先行定義されていると考えられます。

## 3. 推奨される修正案
以下のいずれかの方針で整合性を確保することを推奨します：

- **方針 A（推奨: 仕様書に実装を合わせる）**:
  `service.task_patch.go` において、`existing, err := s.repo.GetTaskByID(ctx, userID, taskID)` による認可チェック（存在確認・所有権確認）を `req.Validate()` より先に実行します。
  ```go
  func (s *taskService) PatchTask(ctx context.Context, userID, taskID string, req *model.PatchTaskRequest) (*model.PatchTaskResponse, error) {
      existing, err := s.repo.GetTaskByID(ctx, userID, taskID)
      if err != nil {
          return nil, err
      }
      if existing == nil {
          return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
      }

      if err := req.Validate(); err != nil {
          return nil, err
      }
      ...
  ```
  ※ なお、JSON 自体の構文エラー（未整形の JSON 文字列等）については HTTP リクエストパースの観点からハンドラー層で 400 を返却し、業務バリデーション（`req.Validate()`）を認可チェック後に行う形とすることで、無用な DB 負荷とリソース列挙リスクのバランスをとることが可能です。

- **方針 B（実装に合わせて仕様書を更新する場合）**:
  一般的な Web API としてリクエストバリデーション（400）を DB アクセスより常に先行させることがチーム方針である場合は、`04_tasks.md` 3.3.4 の「リクエスト評価順序」を「2. リクエスト構文・入力バリデーション (`400 Bad Request`) → 3. 認可チェック・IDOR/BOLA検証 (`404 Not Found`)」へ改訂してください。

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:51:00
- **Status**: Resolved

### 実施した修正内容
- 推奨方針 A を採用し、`backend/service/task_patch.go` の `PatchTask` において、`existing, err := s.repo.GetTaskByID(ctx, userID, taskID)` による認可・存在チェックを `req.Validate()` より先に実行するように処理順序を変更しました。
- GoDoc の `@spec` コメントの順序も実装と合わせて更新しました。
- `backend/service/task_patch_test.go` において、既存タスクが存在する場合にのみバリデーションエラーとなるようモックを修正し、さらに「異常系: 対象タスクが存在しない場合、不正なリクエストボディであっても 400 ではなく先に 404 NOT_FOUND が返却されること」のテストケースを追加して認可チェックが先行することを検証しました。

### 変更したファイル
- [task_patch.go](backend/service/task_patch.go)
- [task_patch_test.go](backend/service/task_patch_test.go)
