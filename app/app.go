package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"swift-menu-session/config"
	"swift-menu-session/internal/app/handlers"
	"swift-menu-session/internal/gateways/mysql"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type App struct {
	Env                string
	Config             config.Config
	Db                 *sql.DB
	OauthConfig        *oauth2.Config
	UserGateway        *mysql.UserGateway
	SessionCookieStore config.SessionCookieStore
	R                  *mux.Router
	Srv                *http.Server
}

func SetupAPP() *App {

	env := os.Getenv("env")
	config.SetupDocker(env)
	conf := config.LoadConfig(env)
	db := config.NewDb(conf.Database)

	oauthConfig := config.NewOauth2Config(conf.Oauth2)
	userGateway := mysql.NewUserGateway(db)

	sessionCookieStore := config.NewSessionCookieStore(conf.CookieStore)

	r := mux.NewRouter()

	port := conf.Api.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: r,
	}

	return &App{
		Env:                env,
		Config:             conf,
		Db:                 db,
		OauthConfig:        oauthConfig,
		UserGateway:        userGateway,
		SessionCookieStore: sessionCookieStore,
		Srv:                srv,
		R:                  r,
	}
}

func SetupHandlers(app *App) {
	loginHandler := handlers.NewLoginHandler(app.OauthConfig, app.UserGateway, app.SessionCookieStore)
	loginHandler.ServeLogin(app.R)

	logoutHandler := handlers.NewLogoutHandler(app.SessionCookieStore)
	logoutHandler.ServeLogout(app.R)

	homeHandler := handlers.NewHomeHandler(app.SessionCookieStore, app.UserGateway)
	homeHandler.ServeHome(app.R)

	protectedHandler := handlers.NewProtectedHandler(app.UserGateway, app.SessionCookieStore)
	protectedHandler.ServeProtected(app.R)

	go func() {
		log.Printf("Starting server on :%s", app.Srv.Addr)
		if err := app.Srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (a *App) TearDown() {
	config.TearDown(a.Env, a.Config.Docker)
}
