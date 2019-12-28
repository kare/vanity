package vanity

import (
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var addr = "https://kkn.fi"

func init() {
	SetLogger(stdlog.New(ioutil.Discard, "", 0))
}

func TestRedirectFromHttpToHttps(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://kkn.fi", nil)
	srv := Redirect("git", "kkn.fi", "https://github.com/kare")
	srv.ServeHTTP(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("expected response status 301, but got %v", res.StatusCode)
	}
	if res.Header.Get("Location") != addr {
		t.Fatalf("expected response location '%v', but got '%v'", addr, res.Header.Get("Location"))
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
		srv := Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != test.status {
			t.Fatalf("Expecting status code %v for method '%v', but got %v", test.status, test.method, res.StatusCode)
		}
	}
}

func TestIndexPageNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", addr, nil)
	srv := Redirect("git", "kkn.fi", "https://github.com/kare")
	srv.ServeHTTP(rec, req)
	res := rec.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected response status 404, but got %v", res.StatusCode)
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
		srv := Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(rec, req)

		res := rec.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("reading response body failed with error: %v", err)
		}

		expected := `<meta name="go-import" content="` + test.result + `">`
		if !strings.Contains(string(body), expected) {
			t.Fatalf("Expecting url '%v' body to contain html meta tag: '%v', but got:\n'%v'", test.path, expected, string(body))
		}

		expected = "text/html; charset=utf-8"
		if res.Header.Get("content-type") != expected {
			t.Fatalf("Expecting content type '%v', but got '%v'", expected, res.Header.Get("content-type"))
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected response status 200, but got %v", res.StatusCode)
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
		srv := Redirect("git", "kkn.fi", "https://github.com/kare")
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Fatalf("Expected response status %v, but got %v", http.StatusTemporaryRedirect, res.StatusCode)
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
