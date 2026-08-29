package config

import (
	"fmt"
	"os"
)

// Config はアプリケーション全体の設定を保持する構造体です。
type Config struct {
	GinMode     string
	FrontendURL string
	DB          DBConfig
}

// DBConfig はデータベース接続設定を保持する構造体です。
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// Load は環境変数から設定を読み込み Config 構造体を返します。
func Load() *Config {
	return &Config{
		GinMode:     getEnvOrDefault("GIN_MODE", "debug"),
		FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		DB: DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "synctask"),
			Password: getEnvOrDefault("DB_PASSWORD", "synctask_pass"),
			Name:     getEnvOrDefault("DB_NAME", "synctask_dev"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
	}
}

// DSN は PostgreSQL 接続文字列を生成して返します。
func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
