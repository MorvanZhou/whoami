package gormdb

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open opens a SQLite database and runs AutoMigrate for whoami tables.
func Open(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("databaseURL is empty")
	}

	// Ensure parent directory exists for file-backed DBs.
	if databaseURL != ":memory:" && databaseURL != "file::memory:?cache=shared" {
		if dir := filepath.Dir(databaseURL); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("create db dir: %w", err)
			}
		}
	}

	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1) // SQLite-friendly
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}
	// WAL helps file-backed SQLite concurrency; ignore failures (e.g. :memory:).
	_, _ = sqlDB.Exec("PRAGMA journal_mode = WAL")

	if err := db.AutoMigrate(&userModel{}, &apiKeyModel{}, &storeTokenModel{}); err != nil {
		return nil, fmt.Errorf("automigrate: %w", err)
	}

	return db, nil
}
