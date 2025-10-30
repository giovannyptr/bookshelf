package platform

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type Migrator interface {
	AutoMigrate(...interface{}) error
}

// OpenGorm creates a GORM connection using env vars.
func OpenGorm() (*gorm.DB, error) {
	host := env("DB_HOST", "127.0.0.1")
	port := env("DB_PORT", "5433")
	user := env("DB_USER", "bookshelf")
	pass := env("DB_PASSWORD", "")
	name := env("DB_NAME", "bookshelf")
	ssl := env("DB_SSLMODE", "disable")

	auth := user
	if pass != "" {
		auth = fmt.Sprintf("%s:%s", user, pass)
	}
	dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s", auth, host, port, name, ssl)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
