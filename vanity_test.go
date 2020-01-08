package vanity_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kkn.fi/vanity"
)

var addr = "https://kkn.fi"

func init() {
	vanity.SetLogger(log.New(ioutil.Discard, "", 0))
}

func TestRedirectFromHttpToHttps(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://kkn.fi", nil)
	srv := vanity.Redirect("git", "kkn.fi", "https://github.com/kare")
	srv.ServeHTTP(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusMovedPermanently {
		t.Errorf("expected response status 301, but got %v", res.StatusCode)
	}
	if res.Header.Get("Location") != addr {
		t.Errorf("expected response location '%v', but got '%v'", addr, res.Header.Get("Location"))
	}
}

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
		req := httptest.NewRequest(test.method, addr+"/gist?go-get=1", nil)
		rec := httptest.NewRecorder()
		srv := vanity.Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != test.status {
			t.Errorf("Expecting status code %v for method '%v', but got %v", test.status, test.method, res.StatusCode)
		}
	}
}

func TestIndexPageNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", addr, nil)
	srv := vanity.Redirect("git", "kkn.fi", "https://github.com/kare")
	srv.ServeHTTP(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected response status 404, but got %v", res.StatusCode)
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
		{"/pkg/subpkg?go-get=1", "kkn.fi/pkg/subpkg git https://github.com/kare/pkg"},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", addr+test.path, nil)
		srv := vanity.Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(rec, req)

		res := rec.Result()
		body, _ := ioutil.ReadAll(res.Body)
		expected := fmt.Sprintf(`<meta name="go-import" content="%v">`, test.result)
		if !strings.Contains(string(body), expected) {
			t.Errorf("Expecting url '%v' body to contain html meta tag: '%v', but got:\n'%v'", test.path, expected, string(body))
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected response status 200, but got %v", res.StatusCode)
		}
	}
}

func TestBrowserGoDoc(t *testing.T) {
	tests := []struct {
		path   string
		result string
	}{
		{"/gist", "https://pkg.go.dev/kkn.fi/gist"},
		{"/set", "https://pkg.go.dev/kkn.fi/set"},
		{"/cmd/vanity", "https://pkg.go.dev/kkn.fi/cmd/vanity"},
		{"/cmd/tcpproxy", "https://pkg.go.dev/kkn.fi/cmd/tcpproxy"},
		{"/pkgabc/sub/foo", "https://pkg.go.dev/kkn.fi/pkgabc/sub"},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", addr+test.path, nil)
		srv := vanity.Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Errorf("Expected response status %v, but got %v", http.StatusTemporaryRedirect, res.StatusCode)
		}
		body, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(body), test.result) {
			t.Errorf("Expecting '%v' be contained in '%v'", test.result, string(body))
		}
	}
}
