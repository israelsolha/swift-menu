package testutils

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"swift-menu-session/app"
	"time"
)

func PreapreIntegrationTest() *app.App {
	// Set environment variable for the test
	os.Setenv("env", "test")

	// Start the API server
	a := app.SetupAPP()
	app.SetupHandlers(a)

	return a
}

func TearDown(a *app.App) {
	dropAllTables(a.Db)
	shutDownServer(a.Srv)
}

func dropAllTables(db *sql.DB) {
	// Query for all tables in the database
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterate through the results and drop each table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}

		// Drop the table
		_, err := db.Exec("DROP TABLE " + tableName)
		if err != nil {
			panic(err)
		}
		log.Printf("Dropped table: %s", tableName)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func shutDownServer(srv *http.Server) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}
