package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewDb(cfg Config) *sql.DB {
	db := createDb(cfg)
	driver := createDriver(db)
	migrationManager := createMigrationManager(driver)
	runMigrations(migrationManager)
	return db
}

func createDb(cfg Config) *sql.DB {
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Schema)
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createDriver(db *sql.DB) database.Driver {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return driver
}

func createMigrationManager(driver database.Driver) *migrate.Migrate {
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations", // Path to migration files
		"mysql",                   // Database type
		driver,                    // Database driver
	)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}
	return m
}

func runMigrations(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database: %v", err)
	}
}
