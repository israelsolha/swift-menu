package handlers_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"swift-menu-session/internal/app/handlers"
	"swift-menu-session/internal/domain/entities"
	"swift-menu-session/mocks"
	"swift-menu-session/testutils"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutWithoutCookies(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)

	// Make the HTTP request to the server
	res, err := http.Get("http://localhost:8077/logout")
	assert.NoError(t, err)

	// Check the status code
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	expectedLocation := "/"
	if loc := res.Request.URL.Path; loc != expectedLocation {
		t.Errorf("handler returned wrong location header: got %v want %v", loc, expectedLocation)
	}

	dummyRequest, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	for _, cookie := range res.Cookies() {
		dummyRequest.AddCookie(cookie)
	}

	session, _ := a.SessionCookieStore.GetCookie(dummyRequest)
	if len(session.Values) != 0 {
		log.Fatal("Expected empty cookie but found values, %w", session.Values)
	}
}

func TestLogoutWithCookies(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)

	// Make the HTTP request to the server
	req, err := http.NewRequest("GET", "http://localhost:8077/logout", nil)
	assert.NoError(t, err)

	user := entities.User{
		Email:          "test@test.com",
		ProfilePicture: "http://profile-picture.com",
		Name:           "Israel Solha",
	}

	res, err := testutils.PerformRequestWithCookie(a.SessionCookieStore, req, user)
	assert.NoError(t, err)

	// Check the status code
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	expectedLocation := "/"
	if loc := res.Request.URL.Path; loc != expectedLocation {
		t.Errorf("handler returned wrong location header: got %v want %v", loc, expectedLocation)
	}

	dummyRequest, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	for _, cookie := range res.Cookies() {
		dummyRequest.AddCookie(cookie)
	}

	session, _ := a.SessionCookieStore.GetCookie(dummyRequest)
	if len(session.Values) != 0 {
		log.Fatal("Expected empty cookie but found values, %w", session.Values)
	}
}

func TestLogoutWithFailingSession(t *testing.T) {
	r := mux.NewRouter()
	mockSession := new(mocks.SessionCookieStore)
	handler := handlers.NewLogoutHandler(mockSession)
	handler.ServeLogout(r)

	req, err := http.NewRequest("GET", "/logout", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()

	mockSession.On("GetCookie", mock.Anything).Return(nil, fmt.Errorf("Error getting session"))

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)

	expectedLocation := "/"
	if loc := rr.Result().Header.Get("Location"); loc != expectedLocation {
		t.Errorf("handler returned wrong location header: got %v want %v", loc, expectedLocation)
	}
}
