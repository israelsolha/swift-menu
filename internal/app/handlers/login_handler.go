package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"swift-menu-session/config"
	"swift-menu-session/internal/domain/entities"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type AuthCallbackHandlerInterface interface {
	ServeLogin(r *mux.Router)
}

type authCallbackHandler struct {
	oauthConfig *oauth2.Config
	userGateway userCallbackHandlerGateway
	stateString string
	store       config.SessionCookieStore
}

type userCallbackHandlerGateway interface {
	CreateUser(user entities.User) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
}

func NewLoginHandler(oauthConfig *oauth2.Config, userGateway userCallbackHandlerGateway, store config.SessionCookieStore) AuthCallbackHandlerInterface {
	handler := authCallbackHandler{
		oauthConfig: oauthConfig,
		userGateway: userGateway,
		stateString: "swift-menu-state",
		store:       store,
	}

	return &handler
}

func (a *authCallbackHandler) ServeLogin(r *mux.Router) {
	r.HandleFunc("/login", a.handleLogin).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", a.handleCallback).Methods("GET")
}

func (a *authCallbackHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	url := a.oauthConfig.AuthCodeURL(a.stateString, oauth2.SetAuthURLParam("prompt", "consent"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *authCallbackHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	if err := a.validateState(w, r); err != nil {
		return
	}
	token, err := a.exchangeCodeForToken(w, r)
	if err != nil {
		return
	}
	userInfo, err := a.getUserInfo(w, r, token)
	if err != nil {
		return
	}
	email, picture, name, err := a.extractUserInfo(userInfo)
	if err != nil {
		return
	}
	user, err := a.getOrCreateUser(email, picture, name)
	if err != nil {
		http.Error(w, "Failed to create/get user", http.StatusInternalServerError)
		return
	}
	a.setSessionValues(w, r, user)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (a *authCallbackHandler) validateState(w http.ResponseWriter, r *http.Request) error {
	state := r.FormValue("state")
	if state != a.stateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return fmt.Errorf("invalid state")
	}
	return nil
}

func (a *authCallbackHandler) exchangeCodeForToken(w http.ResponseWriter, r *http.Request) (*oauth2.Token, error) {
	code := r.FormValue("code")
	token, err := a.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}
	return token, nil
}

func (a *authCallbackHandler) getUserInfo(w http.ResponseWriter, r *http.Request, token *oauth2.Token) (map[string]interface{}, error) {
	client := a.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return userInfo, nil
}

func (a *authCallbackHandler) extractUserInfo(userInfo map[string]interface{}) (string, string, string, error) {
	email, ok := userInfo["email"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("failed to get email from Google")
	}

	picture, ok := userInfo["picture"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("failed to get picture from Google")
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("failed to get name from Google")
	}

	return email, picture, name, nil
}

func (a *authCallbackHandler) setSessionValues(w http.ResponseWriter, r *http.Request, user entities.User) {
	session, _ := a.store.GetCookie(r)
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Save(r, w)
}

func (a *authCallbackHandler) getOrCreateUser(email string, picture string, name string) (entities.User, error) {
	user, err := a.userGateway.GetUserByEmail(email)
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
	user, err = a.userGateway.CreateUser(user)
	if err != nil {
		return entities.User{}, nil
	}
	return user, nil
}
