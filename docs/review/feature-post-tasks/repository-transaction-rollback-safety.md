# 複数タスク一括登録（CreateTasks）におけるトランザクションの Rollback 安全性

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-29 14:55:00
- **Target Files**:
  - [backend/repository/task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/repository/task.go#L60-L90)

## 1. 問題の概要
`backend/repository/task.go` の `CreateTasks` メソッドにおいて、トランザクション開始後にループ内エラー時のみ明示的に `tx.Rollback()` を呼び出していますが、`defer tx.Rollback()` が設定されていません。

## 2. 詳細な指摘内容
1. `CreateTasks` 内で `tx, err := r.db.BeginTx(ctx, nil)` を行った後、ループ内でのクエリ失敗時にのみ `_ = tx.Rollback()` が呼び出されています。
2. もしループ処理中やその前後にパニックが発生した場合、または `tx.Commit()` 自身がエラーを返した場合に、トランザクションがロールバックされず接続リソースがリークまたはロックされたまま残る可能性があります。
3. Go言語のデータベース処理における標準的なプラクティスとして、`BeginTx` の直後に `defer tx.Rollback()` を呼び出す（正常に `Commit()` された後は `Rollback()` は `sql.ErrTxDone` となり安全に無視される）パターンが推奨されます。

## 3. 推奨される修正案
`CreateTasks` において、`BeginTx` 直後に `defer tx.Rollback()` を配置し、ループ内での個別ロールバック呼び出しを簡素化・安全化します。

```go
func (r *taskRepository) CreateTasks(ctx context.Context, tasks []*model.Task) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, task := range tasks {
		if _, err := tx.ExecContext(
			ctx,
			insertTaskQuery,
			task.ID,
			task.UserID,
			task.Title,
			task.Priority,
			task.DueDatetime,
			task.Status,
			task.IsPinned,
			task.Comment,
			task.SearchText,
			task.CreatedAt,
			task.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
```

---

## 修正完了報告

- **Resolved At**: 2026-08-29 15:13:00
- **Status**: Resolved

### 実施した修正内容
- `backend/repository/task.go` の `CreateTasks` メソッドにおいて、`tx, err := r.db.BeginTx(ctx, nil)` の直後に `defer tx.Rollback()` を配置しました。
- ループ内の個別エラーハンドリングにおける手動 `tx.Rollback()` 呼び出しを整理し、エラー時やパニック発生時、コミット失敗時にも確実にトランザクションが安全にロールバックされるようにしました。

### 変更したファイル
- [task.go](file:///C:/Users/kazuh/Programming/repos/SyncTask/backend/repository/task.go)
