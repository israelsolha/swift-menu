package handlers

import (
	"fmt"
	"net/http"
	"swift-menu-session/config"
	"swift-menu-session/internal/domain/entities"

	"github.com/gorilla/mux"
)

type HomeHandlerInterface interface {
	ServeHome(r *mux.Router)
}

type homeHandler struct {
	store       config.SessionCookieStore
	userGateway userHomeHandlerGateway
}

type userHomeHandlerGateway interface {
	GetUserByEmail(email string) (entities.User, error)
}

func NewHomeHandler(store config.SessionCookieStore, userGateway userHomeHandlerGateway) HomeHandlerInterface {
	handler := homeHandler{
		store:       store,
		userGateway: userGateway,
	}
	return &handler
}

func (l *homeHandler) ServeHome(r *mux.Router) {
	r.HandleFunc("/", l.handleHome).Methods("GET")
}

func (h *homeHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.GetCookie(r)
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		fmt.Fprint(w, `<html><body><a href="/login">Google Log In</a></body></html>`)
		return
	}
	user, err := h.userGateway.GetUserByEmail(session.Values["email"].(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	html := fmt.Sprintf(`<html><body>
		Welcome %s! You are logged in.<br>
		<img src="%s" alt="Description of Image" /><br>
		%s<br>
		<a href="/protected">Protected Endpoint</a><br>
		<a href="/logout">Logout</a>
		</body></html>`,
		user.Name, user.ProfilePicture, user.Email)
	w.Write([]byte(html))
}
