package config

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionCookieStore interface {
	GetCookie(r *http.Request) (*sessions.Session, error)
}

type sessionCookieStore struct {
	store         *sessions.CookieStore
	sessionString string
}

func NewSessionCookieStore(cookieStoreConfig CookieStore) SessionCookieStore {
	store := sessions.NewCookieStore([]byte(cookieStoreConfig.Secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 3600, // Session expiration time in seconds
		HttpOnly: true,      // Prevents JavaScript access to the cookie
		Secure:   true,      // Ensure the cookie is sent over HTTPS only
	}

	sessionCookieStore := sessionCookieStore{
		store:         store,
		sessionString: "session",
	}

	return &sessionCookieStore
}

func (s *sessionCookieStore) GetCookie(r *http.Request) (*sessions.Session, error) {
	return s.store.Get(r, s.sessionString)
}
