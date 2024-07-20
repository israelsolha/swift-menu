package config

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"swift-menu-session/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewDb(databaseConfig Database) *sql.DB {
	db := createDb(databaseConfig)
	driver := createDriver(db)
	migrationManager := createMigrationManager(driver)
	runMigrations(migrationManager)
	return db
}

func createDb(databaseConfig Database) *sql.DB {
	maxRetries := 15
	retryInterval := 2 * time.Second

	var lastErr error

	for retries := 0; retries < maxRetries; retries++ {
		connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			databaseConfig.Username,
			databaseConfig.Password,
			databaseConfig.Host,
			databaseConfig.Port,
			databaseConfig.Schema)
		var err error
		db, err := sql.Open("mysql", connectionStr)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to database (attempt %d/%d): %w", retries+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		err = db.Ping()
		if err != nil {
			lastErr = fmt.Errorf("failed to ping database (attempt %d/%d): %w", retries+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		log.Println("Successfully connected to the database!")
		return db
	}

	log.Panicf("Failed to connect to database after %d retries: %v", maxRetries, lastErr)
	return nil
}

func createDriver(db *sql.DB) database.Driver {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	return driver
}

func createMigrationManager(driver database.Driver) *migrate.Migrate {
	rootPath := utils.FindProjectRoot()
	migrationFiles := "file://" + path.Join(rootPath, "migrations")
	m, err := migrate.NewWithDatabaseInstance(
		migrationFiles, // Path to migration files
		"mysql",        // Database type
		driver,         // Database driver
	)
	if err != nil {
		log.Panicf("Migration initialization failed: %v", err)
	}
	return m
}

func runMigrations(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database: %v", err)
	}
}
