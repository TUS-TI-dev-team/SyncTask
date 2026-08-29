package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	// 一時的に環境変数をクリア
	os.Unsetenv("GIN_MODE")
	os.Unsetenv("FRONTEND_URL")
	os.Unsetenv("COOKIE_SECURE")
	os.Unsetenv("TRUSTED_PROXIES")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")

	cfg := Load()

	assert.Equal(t, "debug", cfg.GinMode)
	assert.Equal(t, "http://localhost:3000", cfg.FrontendURL)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "synctask", cfg.DB.User)
	assert.False(t, cfg.CookieSecure)
	assert.Empty(t, cfg.TrustedProxies)
	assert.Equal(t, "synctask_pass", cfg.DB.Password)
	assert.Equal(t, "synctask_dev", cfg.DB.Name)
	assert.Equal(t, "disable", cfg.DB.SSLMode)
	assert.Equal(t, "host=localhost port=5432 user=synctask password=synctask_pass dbname=synctask_dev sslmode=disable", cfg.DB.DSN())
}

func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("GIN_MODE", "release")
	t.Setenv("FRONTEND_URL", "https://synctask.app")
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "custom_user")
	t.Setenv("DB_PASSWORD", "custom_pass")
	t.Setenv("DB_NAME", "custom_db")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("TRUSTED_PROXIES", "10.0.0.0/8,192.0.2.10")

	cfg := Load()

	assert.Equal(t, "release", cfg.GinMode)
	assert.Equal(t, "https://synctask.app", cfg.FrontendURL)
	assert.True(t, cfg.CookieSecure)
	assert.Equal(t, []string{"10.0.0.0/8", "192.0.2.10"}, cfg.TrustedProxies)
	assert.Equal(t, "db", cfg.DB.Host)
	assert.Equal(t, "5433", cfg.DB.Port)
	assert.Equal(t, "custom_user", cfg.DB.User)
	assert.Equal(t, "custom_pass", cfg.DB.Password)
	assert.Equal(t, "custom_db", cfg.DB.Name)
	assert.Equal(t, "require", cfg.DB.SSLMode)
	assert.Equal(t, "host=db port=5433 user=custom_user password=custom_pass dbname=custom_db sslmode=require", cfg.DB.DSN())
}
