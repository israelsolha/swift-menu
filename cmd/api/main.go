package main

import (
	"log"
	"net/http"
	"os"
	"swift-menu-session/config"
	"swift-menu-session/internal/app/handlers"
	"swift-menu-session/internal/gateways/mysql"

	"github.com/gorilla/mux"
)

func main() {
	env := os.Getenv("env")
	conf, err := config.LoadConfig(env)
	if err != nil {
		panic(err)
	}

	db, err := config.NewDb(conf)
	if err != nil {
		panic(err)
	}

	oauthConfig := config.NewOauth2Config(conf)

	userGateway := mysql.NewUserGateway(db)

	query := `
    CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		profile_picture VARCHAR(255),
		name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_email (email)
	);`

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	sessionHandler := handlers.NewSessionHandler(oauthConfig, userGateway, conf.CookieStore.Secret)
	sessionHandler.HandleSession(r)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
