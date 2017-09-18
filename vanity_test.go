package vanity

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPMethodsSupport(t *testing.T) {
	tests := []struct {
		method string
		status int
	}{
		{http.MethodGet, http.StatusOK},
		{http.MethodHead, http.StatusMethodNotAllowed},
		{http.MethodPost, http.StatusMethodNotAllowed},
		{http.MethodPut, http.StatusMethodNotAllowed},
		{http.MethodDelete, http.StatusMethodNotAllowed},
		{http.MethodTrace, http.StatusMethodNotAllowed},
		{http.MethodOptions, http.StatusMethodNotAllowed},
	}
	for _, test := range tests {
		req, err := http.NewRequest(test.method, "http://kkn.fi"+"/gist?go-get=1", nil)
		if err != nil {
			t.Skipf("http request with method %v failed with error: %v", test.method, err)
		}
		res := httptest.NewRecorder()
		srv := Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(res, req)
		if res.Code != test.status {
			t.Fatalf("Expecting status code %v for method '%v', but got %v", test.status, test.method, res.Code)
		}
	}
}

func TestGoTool(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"/gist?go-get=1", "kkn.fi/gist git https://github.com/kare/gist"},
		{"/set?go-get=1", "kkn.fi/set git https://github.com/kare/set"},
		{"/cmd/vanity?go-get=1", "kkn.fi/cmd/vanity git https://github.com/kare/vanity"},
		{"/cmd/tcpproxy?go-get=1", "kkn.fi/cmd/tcpproxy git https://github.com/kare/tcpproxy"},
	}
	for _, test := range tests {
		res := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://kkn.fi"+test.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		srv := Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(res, req)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("reading response body failed with error: %v", err)
		}

		expected := `<meta name="go-import" content="` + test.result + `">`
		if !strings.Contains(string(body), expected) {
			t.Fatalf("Expecting url '%v' body to contain html meta tag: '%v', but got:\n'%v'", test.path, expected, string(body))
		}

		expected = "text/html; charset=utf-8"
		if res.HeaderMap.Get("content-type") != expected {
			t.Fatalf("Expecting content type '%v', but got '%v'", expected, res.HeaderMap.Get("content-type"))
		}

		if res.Code != http.StatusOK {
			t.Fatalf("Expected response status 200, but got %v", res.Code)
		}
	}
}

func TestBrowserGoDoc(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"/gist", "https://godoc.org/kkn.fi/gist"},
		{"/set", "https://godoc.org/kkn.fi/set"},
		{"/cmd/vanity", "https://godoc.org/kkn.fi/cmd/vanity"},
		{"/cmd/tcpproxy", "https://godoc.org/kkn.fi/cmd/tcpproxy"},
	}
	for _, test := range tests {
		res := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://kkn.fi"+test.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		srv := Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(res, req)

		if res.Code != http.StatusTemporaryRedirect {
			t.Fatalf("Expected response status %v, but got %v", http.StatusOK, res.Code)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("reading response body failed with error: %v", err)
		}
		if !strings.Contains(string(body), test.result) {
			t.Fatalf("Expecting '%v' be contained in '%v'", test.result, string(body))
		}
	}
}
