package handlers

import (
	"net/http"
	"swift-menu-session/config"

	"github.com/gorilla/mux"
)

type LogoutHandlerInterface interface {
	ServeLogout(r *mux.Router)
}

type loginCallbackHandler struct {
	store config.SessionCookieStore
}

func NewLogoutHandler(store config.SessionCookieStore) LogoutHandlerInterface {
	handler := loginCallbackHandler{
		store: store,
	}
	return &handler
}

func (l *loginCallbackHandler) ServeLogout(r *mux.Router) {
	r.HandleFunc("/logout", l.handleLogout).Methods("GET")
}

func (l *loginCallbackHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := l.store.GetCookie(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
