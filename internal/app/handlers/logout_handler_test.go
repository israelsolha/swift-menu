package handlers_test

import (
	"log"
	"net/http"
	"swift-menu-session/testutils"
	"testing"
)

func TestLogoutWithoutCookies(t *testing.T) {
	a := testutils.PreapreIntegrationTest()
	defer testutils.TearDown(a)

	// Make the HTTP request to the server
	res, err := http.Get("http://localhost:8077/logout")
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}

	expectedLocation := "/"
	if loc := res.Request.URL.Path; loc != expectedLocation {
		t.Errorf("handler returned wrong location header: got %v want %v", loc, expectedLocation)
	}

	dummyRequest, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, cookie := range res.Cookies() {
		dummyRequest.AddCookie(cookie)
	}

	session, _ := a.SessionCookieStore.GetCookie(dummyRequest)
	if len(session.Values) != 0 {
		log.Fatal("Expected empty cookie but found values, %w", session.Values)
	}
}
