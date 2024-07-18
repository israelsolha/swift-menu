package testutils

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"swift-menu-session/app"
	"swift-menu-session/config"
	"swift-menu-session/internal/domain/entities"
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

func PerformRequestWithCookie(store config.SessionCookieStore, r *http.Request, user entities.User) (*http.Response, error) {
	w := httptest.NewRecorder()
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8077", nil)
	if err != nil {
		log.Fatal(err)
	}

	session, _ := store.GetCookie(req)
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Save(req, w)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return client.Do(req)
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
