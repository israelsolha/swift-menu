package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"swift-menu-session/internal/domain/entities"
	"swift-menu-session/testutils"
	"testing"
)

func TestHomeWithoutBeingLogged(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)
	expectedBody := `<html><body><a href="/login">Google Log In</a></body></html>`

	res, err := http.Get("http://localhost:8077")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://localhost:8077", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := testutils.PerformRequestWithCookie(a.SessionCookieStore, req, user)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("unexpected status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != expectedBody {
		t.Fatalf("unexpected body returned: %s", body)
	}

}

func TestHomeBeingLoggedWithoutUser(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)

	req, err := http.NewRequest("GET", "http://localhost:8077", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := testutils.PerformRequestWithCookie(a.SessionCookieStore, req, entities.User{})
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 500 {
		t.Fatalf("unexpected status code %d", res.StatusCode)
	}
}
