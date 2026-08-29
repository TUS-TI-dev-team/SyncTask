package main

import (
	"log"

	"synctask/backend/config"
	"synctask/backend/db"
	"synctask/backend/router"
)

// @title SyncTask API
// @version 1.0
// @description SyncTask バックエンド API サーバー
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	if err := db.Migrate(cfg.DB); err != nil {
		log.Fatalf("マイグレーションに失敗しました: %v", err)
	}

	database, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("DB接続に失敗しました: %v", err)
	}
	defer database.Close()

	r := router.SetupRouter(database)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
