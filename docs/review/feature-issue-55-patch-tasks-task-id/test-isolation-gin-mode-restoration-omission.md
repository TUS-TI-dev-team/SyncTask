# テスト規約違反: TestPatchTaskHandler における gin.SetMode の状態復元漏れ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-09-05 21:47:00
- **Target Files**:
  - [task_patch_test.go](backend/handler/task_patch_test.go)

## 1. 問題の概要
`backend/TESTING_GUIDE.md` の「4. テストの独立性 (Isolation)」において、「グローバル変数や `gin.SetMode` を変更した場合は、テスト終了時に必ず元の状態に戻してください。」と明記されています。
しかし、`backend/handler/task_patch_test.go` の `TestPatchTaskHandler` 内では `gin.SetMode(gin.TestMode)` を実行しているものの、元のモードへの復元処理（`defer gin.SetMode(previousMode)`）が欠落しています。

## 2. 詳細な指摘内容
1. `backend/TESTING_GUIDE.md` L40-L43 にて以下のように規定されています：
   ```markdown
   ### 4. テストの独立性 (Isolation)
   - 各テストケースは他のテストの実行順序や状態に依存してはなりません。
   - グローバル変数や `gin.SetMode` を変更した場合は、テスト終了時に必ず元の状態に戻してください。
   ```
2. `backend/handler/task_patch_test.go` L19-L20 では以下のように記述されています：
   ```go
   func TestPatchTaskHandler(t *testing.T) {
       gin.SetMode(gin.TestMode)
   ```
   直近に実装された `router_test.go`（L143-L145）などでは以下のように適切に復元処理が行われていますが、本テスト関数では復元が行われていません。
   ```go
   previousMode := gin.Mode()
   gin.SetMode(gin.TestMode)
   defer gin.SetMode(previousMode)
   ```
3. これにより、後続の他のテストスイートが実行される際にグローバルなモード設定が汚染されたままとなり、テストの独立性が損なわれるリスクがあります。

## 3. 推奨される修正案
`backend/handler/task_patch_test.go` の `TestPatchTaskHandler` 冒頭を以下のように修正してください：

```go
func TestPatchTaskHandler(t *testing.T) {
    previousMode := gin.Mode()
    gin.SetMode(gin.TestMode)
    defer gin.SetMode(previousMode)
```

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:51:00
- **Status**: Resolved

### 実施した修正内容
- `backend/handler/task_patch_test.go` の `TestPatchTaskHandler` 冒頭において、`previousMode := gin.Mode()` で現在のモードを退避し、`defer gin.SetMode(previousMode)` でテスト終了時に確実に復元するよう修正しました。
- これにより、`backend/TESTING_GUIDE.md` の「4. テストの独立性 (Isolation)」の規約に適合し、他のテスト実行への影響（モード汚染）を防止しました。

### 変更したファイル
- [task_patch_test.go](backend/handler/task_patch_test.go)
