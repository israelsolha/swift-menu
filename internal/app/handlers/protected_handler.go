package handlers

import (
	"encoding/json"
	"net/http"
	"swift-menu-session/config"
	"swift-menu-session/internal/domain/entities"
	"time"

	"github.com/gorilla/mux"
)

type ProtectedHandlerInterface interface {
	ServeProtected(r *mux.Router)
}

type protectedHandler struct {
	store       config.SessionCookieStore
	userGateway userProtectedHandlerGateway
}

type userProtectedHandlerGateway interface {
	GetUserByEmail(email string) (entities.User, error)
}

func NewProtectedHandler(userGateway userProtectedHandlerGateway, store config.SessionCookieStore) ProtectedHandlerInterface {
	handler := protectedHandler{
		store:       store,
		userGateway: userGateway,
	}
	return &handler
}

func (p *protectedHandler) ServeProtected(r *mux.Router) {
	r.Handle("/protected", p.authMiddleware(http.HandlerFunc(p.handleProtected))).Methods("GET")
}

func (p *protectedHandler) handleProtected(w http.ResponseWriter, r *http.Request) {
	session, _ := p.store.GetCookie(r)
	email, ok := session.Values["email"].(string)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	user, err := p.userGateway.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "This is a protected endpoint",
		"user":    user,
		"time":    time.Now(),
	})
}

func (p *protectedHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := p.store.GetCookie(r)
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
