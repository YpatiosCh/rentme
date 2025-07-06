package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_PATH")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              false,
		CreateBatchSize:                          1000,
		Logger:                                   logger.Default.LogMode(logger.Info),

		NowFunc: func() time.Time { return time.Now().UTC() }, // optional but recommended

		SkipDefaultTransaction:   false,
		DryRun:                   false,
		DisableNestedTransaction: true,
		AllowGlobalUpdate:        false,
		NamingStrategy:           nil,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Enable foreign key constraints
	db.Exec("PRAGMA foreign_keys = ON")

	// Enable WAL mode for better concurrency
	db.Exec("PRAGMA journal_mode = WAL")

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(5)

	return db, nil
}
