package internal

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	hostname = "kkn.fi"
	config   = map[string]*Package{
		"/gist":         NewPackage("kkn.fi/gist", "git", "https://github.com/kare/gist"),
		"/set":          NewPackage("kkn.fi/set", "git", "https://github.com/kare/set"),
		"/cmd/vanity":   NewPackage("kkn.fi/cmd/vanity", "git", "https://github.com/kare/vanity"),
		"/cmd/tcpproxy": NewPackage("kkn.fi/cmd/tcpproxy", "git", "https://github.com/kare/tcpproxy"),
	}
)

func TestPackage(t *testing.T) {
	p := NewPackage("kkn.fi/gist", "git", "https://github.com/kare/gist")
	if p.name() != "gist" {
		t.Errorf("expected 'gist', got %v", p.name())
	}
	p = NewPackage("kkn.fi/cmd/tcpproxy", "git", "https://github.com/kare/tcpproxy")
	if p.name() != "cmd/tcpproxy" {
		t.Errorf("expected 'cmd/tcpproxy', got %v", p.name())
	}
}

func TestHTTPMethodsSupport(t *testing.T) {
	server := NewServer(hostname, config)
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
		req, err := http.NewRequest(test.method, "/gist?go-get=1", nil)
		if err != nil {
			t.Skipf("http request with method %v failed with error: %v", test.method, err)
		}
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		if res.Code != test.status {
			t.Fatalf("Expecting status code %v for method '%v', but got %v", test.status, test.method, res.Code)
		}
	}
}

func TestGoTool(t *testing.T) {
	server := httptest.NewServer(NewServer(hostname, config))
	defer server.Close()

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
		url := server.URL + test.path
		res, err := http.Get(url)
		if err != nil {
			t.Skipf("error requesting url %v\n%v", url, err)
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				t.Errorf("error closing response body: %v", err)
			}
		}()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("reading response body failed with error: %v", err)
		}

		expected := `<meta name="go-import" content="` + test.result + `">`
		if !strings.Contains(string(body), expected) {
			t.Fatalf("Expecting url '%v' body to contain html meta tag: '%v', but got:\n'%v'", url, expected, string(body))
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

func TestGoToolPackageNotFound(t *testing.T) {
	server := httptest.NewServer(NewServer(hostname, config))
	defer server.Close()

	url := server.URL + "/package-not-found?go-get=1"
	res, err := http.Get(url)
	if err != nil {
		t.Skipf("error requesting url %v\n%v", url, err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading response body failed with error: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected response status 404, but got %v", res.StatusCode)
	}
	expected := "404 page not found\n"
	if string(body) != expected {
		t.Fatalf("Expecting '%v', but got '%v'", expected, string(body))
	}
}

func TestBrowserGoDoc(t *testing.T) {
	server := httptest.NewServer(NewServer(hostname, config))
	defer server.Close()

	tests := []struct {
		path   string
		result string
	}{
		{"/gist", "https://godoc.org/kkn.fi/gist"},
		{"/set", "https://godoc.org/kkn.fi/set"},
		{"/cmd/vanity", "https://godoc.org/kkn.fi/cmd/vanity"},
		{"/cmd/tcpproxy", "https://godoc.org/kkn.fi/cmd/tcpproxy"},
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	for _, test := range tests {
		url := server.URL + test.path
		res, err := client.Get(url)
		if err != nil {
			t.Skipf("error requesting url %v\n%v", url, err)
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				t.Errorf("error closing response body: %v", err)
			}
		}()

		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Fatalf("Expected response status %v, but got %v", http.StatusTemporaryRedirect, res.StatusCode)
		}

		location := res.Header.Get("location")
		if location != test.result {
			t.Fatalf("Expecting location header to match '%v', but got '%v'", test.result, location)
		}

		expected := "text/html; charset=utf-8"
		contentType := res.Header.Get("content-type")
		if contentType != expected {
			t.Fatalf("Expecting content type '%v', but got '%v'", expected, contentType)
		}
	}
}
