package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"swift-menu-session/testutils"
	"testing"
)

func TestHomeWithoutBeingLogged(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)
	expectedBody := `<html><body><a href="/login">Google Log In</a></body></html>`

	// Make the HTTP request to the server
	res, err := http.Get("http://localhost:8077")
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code
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

	w := httptest.NewRecorder()
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8077", nil)
	if err != nil {
		t.Fatal(err)
	}

	session, _ := a.SessionCookieStore.GetCookie(req)
	session.Values["authenticated"] = true
	session.Values["email"] = "test@test.com"
	session.Save(req, w)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 500 {
		t.Fatalf("unexpected status code %d", res.StatusCode)
	}

}
