package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config はアプリケーション全体の設定を保持する構造体です。
type Config struct {
	GinMode        string
	FrontendURL    string
	CookieSecure   bool
	TrustedProxies []string
	AutoSeed       bool
	DB             DBConfig
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

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getCSVEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

// Load は環境変数から設定を読み込み Config 構造体を返します。
//
// @spec 設定: GIN_MODE が release 以外かつ AUTO_SEED が明示的に false でない場合、AutoSeed を true に設定すること
// @spec 設定: GIN_MODE が release の場合、AUTO_SEED の指定に関わらず AutoSeed を強制的に false に設定すること
func Load() *Config {
	ginMode := getEnvOrDefault("GIN_MODE", "debug")

	var autoSeed bool
	if ginMode != "release" {
		autoSeed = getBoolEnvOrDefault("AUTO_SEED", true)
	}

	return &Config{
		GinMode:        ginMode,
		FrontendURL:    getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		CookieSecure:   getBoolEnvOrDefault("COOKIE_SECURE", false),
		TrustedProxies: getCSVEnv("TRUSTED_PROXIES"),
		AutoSeed:       autoSeed,
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
