package model

// GetTaskResponse はタスク詳細取得成功時のレスポンスを表します。
type GetTaskResponse struct {
	// Task は詳細取得されたタスクオブジェクトです
	Task *Task `json:"task"`
}
