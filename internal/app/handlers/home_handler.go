package handlers

import (
	"fmt"
	"net/http"
	"swift-menu-session/config"
	"swift-menu-session/internal/domain/entities"

	"github.com/gorilla/mux"
)

var (
	loggedOutBody = `<html><body><a href="/login">Google Log In</a></body></html>`
	loggedInBody  = `<html><body>
		Welcome %s! You are logged in.<br>
		<img src="%s" alt="Description of Image" /><br>
		%s<br>
		<a href="/protected">Protected Endpoint</a><br>
		<a href="/logout">Logout</a>
		</body></html>`
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
	session, err := h.store.GetCookie(r)
	if err != nil {
		fmt.Fprint(w, loggedOutBody)
		return
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprint(w, loggedOutBody)
		return
	}
	user, err := h.userGateway.GetUserByEmail(session.Values["email"].(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	html := fmt.Sprintf(loggedInBody, user.Name, user.ProfilePicture, user.Email)
	w.Write([]byte(html))
}
