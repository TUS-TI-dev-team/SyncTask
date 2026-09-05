# Service 層における UpdateTask 失敗時および更新対象消失時の単体テスト欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-09-05 21:47:00
- **Target Files**:
  - [task_patch.go](backend/service/task_patch.go)
  - [task_patch_test.go](backend/service/task_patch_test.go)

## 1. 問題の概要
`backend/service/task_patch.go` の `PatchTask` メソッドでは、リポジトリ層の `UpdateTask` 呼び出しにおける DB エラー返却（L84-L87）および `updated == nil`（L88-L90: 取得後に並行処理等で対象タスクが消失した場合の 404 NOT_FOUND 返却）のエラーハンドリング分岐が実装されています。
しかし、`backend/service/task_patch_test.go` においては `GetTaskByID` が失敗するケースしか検証されておらず、`UpdateTask` の呼び出し失敗および `nil` 返却時を検証する単体テストケースが欠落しています。

## 2. 詳細な指摘内容
1. `backend/service/task_patch.go` L84-L90:
   ```go
   updated, err := s.repo.UpdateTask(ctx, existing)
   if err != nil {
       return nil, err
   }
   if updated == nil {
       return nil, model.NewNotFoundError("指定されたタスクが見つかりません。")
   }
   ```
2. 現在の `backend/service/task_patch_test.go`（L223-L240）では、`TestTaskService_PatchTask` の「異常系: リポジトリ層でエラーが発生した場合にエラーがそのまま返却されること」において、`getTaskByIDFunc` が `dbErr` を返すケースのみがテストされています。
3. `UpdateTask` 実行時に DB エラーが発生した場合にエラーが伝播するか、および並行削除等で `UpdateTask` が `nil, nil` を返却した際に `model.AppError`（404 `NOT_FOUND`）が適切に返却されるかの分岐が単体テストで網羅されていません。

## 3. 推奨される修正案
`backend/service/task_patch_test.go` に以下の2つのサブテストケースを追加してください：

```go
t.Run("異常系: UpdateTask 実行時にリポジトリがエラーを返却した場合にエラーが伝播すること", func(t *testing.T) {
    var req model.PatchTaskRequest
    err := json.Unmarshal([]byte(`{"title": "更新後タイトル"}`), &req)
    require.NoError(t, err)

    taskCopy := *baseTask
    dbErr := errors.New("update failed")
    repo := &mockTaskRepository{
        getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
            return &taskCopy, nil
        },
        updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
            return nil, dbErr
        },
    }
    svc := NewTaskService(repo)

    res, err := svc.PatchTask(context.Background(), userID, taskID, &req)
    require.Error(t, err)
    assert.Nil(t, res)
    assert.Equal(t, dbErr, err)
    assert.Equal(t, 1, repo.getTaskByIDCalls)
    assert.Equal(t, 1, repo.updateTaskCalls)
})

t.Run("異常系: UpdateTask の戻り値が nil の場合に 404 NOT_FOUND が返却されること", func(t *testing.T) {
    var req model.PatchTaskRequest
    err := json.Unmarshal([]byte(`{"title": "更新後タイトル"}`), &req)
    require.NoError(t, err)

    taskCopy := *baseTask
    repo := &mockTaskRepository{
        getTaskByIDFunc: func(ctx context.Context, uID, tID string) (*model.Task, error) {
            return &taskCopy, nil
        },
        updateTaskFunc: func(ctx context.Context, task *model.Task) (*model.Task, error) {
            return nil, nil // 並行削除等でレコードが存在しなくなった場合
        },
    }
    svc := NewTaskService(repo)

    res, err := svc.PatchTask(context.Background(), userID, taskID, &req)
    require.Error(t, err)
    assert.Nil(t, res)

    var appErr *model.AppError
    require.True(t, errors.As(err, &appErr))
    assert.Equal(t, 404, appErr.StatusCode)
    assert.Equal(t, "NOT_FOUND", appErr.Code)
    assert.Equal(t, "指定されたタスクが見つかりません。", appErr.Message)
    assert.Equal(t, 1, repo.getTaskByIDCalls)
    assert.Equal(t, 1, repo.updateTaskCalls)
})
```

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:51:00
- **Status**: Resolved

### 実施した修正内容
- `backend/service/task_patch_test.go` に推奨修正案通りの2つのサブテストケースを追加しました：
  1. `異常系: UpdateTask 実行時にリポジトリがエラーを返却した場合にエラーが伝播すること`（DBエラーの伝播を検証）
  2. `異常系: UpdateTask の戻り値が nil の場合に 404 NOT_FOUND が返却されること`（並行更新等でタスクが消失した場合の 404 エラー返却を検証）
- TESTING_GUIDE.md の規約に従い、日本語でのテスト名プレフィックス（`異常系:`）およびアサーションの使い分け（`require` / `assert`）を遵守しています。

### 変更したファイル
- [task_patch_test.go](backend/service/task_patch_test.go)
