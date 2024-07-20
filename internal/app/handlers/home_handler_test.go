package handlers_test

import (
	"fmt"
	"io"
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

func TestHomeWithoutBeingLogged(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)
	expectedBody := `<html><body><a href="/login">Google Log In</a></body></html>`

	res, err := http.Get("http://localhost:8077")
	assert.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	if string(body) != expectedBody {
		t.Fatalf("unexpected body returned: %s", body)
	}
}

func TestHomeBeingLogged(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)

	email := "test@test.com"

	user := entities.User{
		Email:          email,
		ProfilePicture: "http://mypicture.com",
		Name:           "Israel Solha",
	}

	expectedBody := fmt.Sprintf(`<html><body>
		Welcome %s! You are logged in.<br>
		<img src="%s" alt="Description of Image" /><br>
		%s<br>
		<a href="/protected">Protected Endpoint</a><br>
		<a href="/logout">Logout</a>
		</body></html>`,
		user.Name, user.ProfilePicture, user.Email)

	_, err := a.UserGateway.CreateUser(user)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "http://localhost:8077", nil)
	assert.NoError(t, err)

	res, err := testutils.PerformRequestWithCookie(a.SessionCookieStore, req, user)
	assert.NoError(t, err)

	if res.StatusCode != 200 {
		t.Fatalf("unexpected status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	if string(body) != expectedBody {
		t.Fatalf("unexpected body returned: %s", body)
	}

}

func TestHomeBeingLoggedWithoutUser(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)

	req, err := http.NewRequest("GET", "http://localhost:8077", nil)
	assert.NoError(t, err)

	res, err := testutils.PerformRequestWithCookie(a.SessionCookieStore, req, entities.User{})
	assert.NoError(t, err)

	if res.StatusCode != 500 {
		t.Fatalf("unexpected status code %d", res.StatusCode)
	}
}

func TestHomeFailingSession(t *testing.T) {
	r := mux.NewRouter()
	mockSession := new(mocks.SessionCookieStore)
	mockUserGateway := new(mocks.UserCallbackHandlerGateway)
	handler := handlers.NewHomeHandler(mockSession, mockUserGateway)
	handler.ServeHome(r)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()

	expectedBody := `<html><body><a href="/login">Google Log In</a></body></html>`

	mockSession.On("GetCookie", mock.Anything).Return(nil, fmt.Errorf("Error getting session"))

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	body, err := io.ReadAll(rr.Body)
	assert.NoError(t, err)

	assert.Equal(t, expectedBody, string(body))
}
