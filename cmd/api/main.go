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
	conf := config.LoadConfig(env)
	db := config.NewDb(conf)
	defer db.Close()

	oauthConfig := config.NewOauth2Config(conf)
	userGateway := mysql.NewUserGateway(db)

	sessionCookieStore := config.NewSessionCookieStore(conf.CookieStore.Secret)

	r := mux.NewRouter()

	loginHandler := handlers.NewLoginHandler(oauthConfig, userGateway, sessionCookieStore)
	loginHandler.ServeLogin(r)

	logoutHandler := handlers.NewLogoutHandler(sessionCookieStore)
	logoutHandler.ServeLogout(r)

	homeHandler := handlers.NewHomeHandler(sessionCookieStore, userGateway)
	homeHandler.ServeHome(r)

	protectedHandler := handlers.NewProtectedHandler(userGateway, sessionCookieStore)
	protectedHandler.ServeProtected(r)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
