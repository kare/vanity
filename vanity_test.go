package vanity

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	hostname = "kkn.fi"
	config   = map[Path]Package{
		"/gist":              {"/gist", "git", "https://github.com/kare/gist"},
		"/vanity":            {"/vanity", "git", "https://github.com/kare/vanity"},
		"/vanity/cmd":        {"/vanity", "git", "https://github.com/kare/vanity"},
		"/vanity/cmd/vanity": {"/vanity", "git", "https://github.com/kare/vanity"},
		"/foo/bar":           {"/foo", "git", "https://github.com/kare/foo"},
		"/foo/bar/baz":       {"/foo", "git", "https://github.com/kare/foo"},
		"/":                  {"/", "git", "https://github.com/project"},
	}
)

func TestHTTPMethodsSupport(t *testing.T) {
	server := Server{&hostname, config}
	tests := []struct {
		method string
		status int
	}{
		{"GET", 200},
		{"HEAD", 405},
		{"POST", 405},
		{"PUT", 405},
		{"DELETE", 405},
		{"TRACE", 405},
		{"OPTIONS", 405},
	}
	for _, test := range tests {
		req, err := http.NewRequest(test.method, "/gist?go-get=1", nil)
		if err != nil {
			t.Errorf("http request with method %v failed with error: %v", test.method, err)
		}
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		if res.Code != test.status {
			t.Fatalf("Expecting status code %v for method '%v', but got %v", test.status, test.method, res.Code)
		}
	}
}

func TestGoTool(t *testing.T) {
	server := httptest.NewServer(Server{&hostname, config})
	defer server.Close()

	tests := []struct {
		path   string
		result string
	}{
		{"/gist?go-get=1", "kkn.fi/gist git https://github.com/kare/gist"},
		{"/vanity?go-get=1", "kkn.fi/vanity git https://github.com/kare/vanity"},
		{"/foo/bar?go-get=1", "kkn.fi/foo git https://github.com/kare/foo"},
		{"/foo/bar/baz?go-get=1", "kkn.fi/foo git https://github.com/kare/foo"},
	}
	for _, test := range tests {
		url := server.URL + test.path
		res, err := http.Get(url)
		if err != nil {
			t.Errorf("error requesting url %v\n%v", url, err)
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				t.Errorf("error closing response body: %v", err)
			}
		}()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("reading response body failed with error: %v", err)
		}

		expected := `<meta name="go-import" content="` + test.result + `">`
		if !strings.Contains(string(body), expected) {
			log.Fatalf("Expecting body to contain html meta tag: '%v', but got:\n'%v'", expected, string(body))
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
	server := httptest.NewServer(Server{&hostname, config})
	defer server.Close()

	url := server.URL + "/package-not-found?go-get=1"
	res, err := http.Get(url)
	if err != nil {
		t.Errorf("error requesting url %v\n%v", url, err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("reading response body failed with error: %v", err)
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
	server := httptest.NewServer(Server{&hostname, config})
	defer server.Close()

	tests := []struct {
		path   string
		result string
	}{
		{"/gist", "https://godoc.org/kkn.fi/gist"},
		{"/vanity", "https://godoc.org/kkn.fi/vanity"},
		{"/vanity/cmd", "https://godoc.org/kkn.fi/vanity/cmd"},
		{"/vanity/cmd/vanity", "https://godoc.org/kkn.fi/vanity/cmd/vanity"},
		{"/foo/bar", "https://godoc.org/kkn.fi/foo/bar"},
		{"/foo/bar/baz", "https://godoc.org/kkn.fi/foo/bar/baz"},
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
			t.Errorf("error requesting url %v\n%v", url, err)
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
