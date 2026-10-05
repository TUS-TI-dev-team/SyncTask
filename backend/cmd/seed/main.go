// seed コマンドは開発時の動作テスト用シードデータをデータベースに投入します。
//
// 使い方:
//
//	go run ./cmd/seed          # 既存データは残し、不足分のみ投入 (冪等)
//	go run ./cmd/seed --force  # TASK / LOGIN_SESSION / LOGIN_ACCOUNT を全件削除して再投入
//
// DB 接続先は backend 本体と同様に環境変数 (DB_HOST / DB_PORT / DB_USER /
// DB_PASSWORD / DB_NAME / DB_SSLMODE) で指定します。
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"synctask/backend/config"
	"synctask/backend/db"
	"synctask/backend/db/seed"
)

func main() {
	force := flag.Bool("force", false, "既存のシード対象データを全件削除してから再投入する")
	flag.Parse()

	if err := run(*force); err != nil {
		log.Fatalf("seed command failed: %v", err)
	}
}

func run(force bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	// スキーマが未適用の場合に備え、まずマイグレーションを実行する。
	if err := db.Migrate(cfg.DB); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	database, err := db.Connect(cfg.DB)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer database.Close()

	if err := seed.Run(ctx, database, force); err != nil {
		return fmt.Errorf("seed failed: %w", err)
	}

	if force {
		fmt.Println("シードデータの全削除&再投入が完了しました。")
	} else {
		fmt.Println("シードデータの投入が完了しました。")
	}
	return nil
}
