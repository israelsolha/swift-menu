package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"swift-menu-session/internal/domain/entities"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type SessionHandlerInterface interface {
	HandleSession(r *mux.Router)
}

type sessionHandler struct {
	oauthConfig *oauth2.Config
	userGateway userGateway
	stateString string
	store       sessions.Store
}

type userGateway interface {
	CreateUser(user entities.User) (entities.User, error)
	GetUserByID(id int) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
}

func NewSessionHandler(oauthConfig *oauth2.Config, userGateway userGateway, cookieSecret string) SessionHandlerInterface {
	handler := sessionHandler{
		oauthConfig: oauthConfig,
		userGateway: userGateway,
		stateString: "swift-menu-state",
		store:       sessions.NewCookieStore([]byte(cookieSecret)),
	}

	return &handler
}

func (s *sessionHandler) HandleSession(r *mux.Router) {
	r.HandleFunc("/auth/{provider}/callback", s.handleCallback).Methods("GET")
	r.HandleFunc("/login", s.handleLogin).Methods("GET")
	r.HandleFunc("/logout", s.handleLogout).Methods("GET")
	r.HandleFunc("/", s.handleHome)

}

func (s *sessionHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != s.stateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		http.Error(w, "Failed to get email from Google", http.StatusInternalServerError)
		return
	}

	picture, ok := userInfo["picture"].(string)
	if !ok {
		http.Error(w, "Failed to get picture from Google", http.StatusInternalServerError)
		return
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		http.Error(w, "Failed to get name from Google", http.StatusInternalServerError)
		return
	}

	user, err := s.getOrCreateUser(email, picture, name)
	if err != nil {
		http.Error(w, "Failed to create/get user", http.StatusInternalServerError)
		return
	}

	session, _ := s.store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Options.MaxAge = 15
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (s *sessionHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	url := s.oauthConfig.AuthCodeURL(s.stateString, oauth2.SetAuthURLParam("prompt", "consent"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *sessionHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")

	// Clear session values
	session.Values["authenticated"] = false
	session.Values["email"] = nil
	session.Values["token"] = nil // Clear stored token

	// Set MaxAge to -1 to delete the session cookie
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *sessionHandler) getOrCreateUser(email string, picture string, name string) (entities.User, error) {
	user, err := s.userGateway.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return entities.User{}, err
	}
	if err == nil {
		return user, nil
	}

	user = entities.User{
		ProfilePicture: picture,
		Name:           name,
		Email:          email,
	}
	user, err = s.userGateway.CreateUser(user)
	if err != nil {
		return entities.User{}, nil
	}
	return user, nil
}

func (s *sessionHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		fmt.Fprint(w, `<html><body><a href="/login">Google Log In</a></body></html>`)
		return
	}
	user, _ := s.userGateway.GetUserByEmail(session.Values["email"].(string))
	html := fmt.Sprintf(`<html><body>
	Welcome %s! You are logged in.<br>
	<img src="%s" alt="Description of Image" /><br>
	%s<br>
	<a href="/protected">Protected Endpoint</a><br>
	<a href="/logout">Logout</a>
</body></html>`, user.Name, user.ProfilePicture, user.Email)
	w.Write([]byte(html))
}

func (s *sessionHandler) HandleProtected(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	email, ok := session.Values["email"].(string)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	user, err := s.userGateway.GetUserByEmail(email)
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

func (s *sessionHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "session")
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
