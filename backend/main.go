package main

import (
	"log"

	"synctask/backend/router"
)

// @title SyncTask API
// @version 1.0
// @description SyncTask バックエンド API サーバー
// @host localhost:8080
// @BasePath /
func main() {
	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
