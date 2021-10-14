package vanity_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"kkn.fi/vanity"
)

var addr = "https://kkn.fi"

func TestRedirectFromHttpToHttps(t *testing.T) {
	tests := []struct {
		url      string
		location string
	}{
		{
			url:      "http://kkn.fi",
			location: "https://kkn.fi",
		},
		{
			url:      "http://kkn.fi/",
			location: "https://kkn.fi/",
		},
		{
			url:      "http://kkn.fi/pkg/sub/foo",
			location: "https://kkn.fi/pkg/sub/foo",
		},
		{
			url:      "http://kkn.fi/vanity",
			location: "https://kkn.fi/vanity",
		},
		{
			url:      "http://kkn.fi/vanity?go-get=1",
			location: "https://kkn.fi/vanity?go-get=1",
		},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, test.url, nil)
		srv := vanity.Handler(
			vanity.Log(log.New(ioutil.Discard, "", 0)),
		)
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusMovedPermanently {
			t.Errorf("expected response status 301, but got %v", res.StatusCode)
		}
		if test.location != res.Header.Get("Location") {
			t.Errorf("expected response location '%v', but got '%v'", test.location, res.Header.Get("Location"))
		}
	}
}

func TestHTTPMethodsSupport(t *testing.T) {
	tests := []struct {
		method string
		status int
	}{
		{
			method: http.MethodGet,
			status: http.StatusOK,
		},
		{
			method: http.MethodHead,
			status: http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodPut,
			status: http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodDelete,
			status: http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodTrace,
			status: http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodOptions,
			status: http.StatusMethodNotAllowed,
		},
	}
	for _, test := range tests {
		req := httptest.NewRequest(test.method, addr+"/gist?go-get=1", nil)
		rec := httptest.NewRecorder()
		srv := vanity.Handler(
			vanity.Log(log.New(ioutil.Discard, "", 0)),
		)
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != test.status {
			t.Errorf("Expecting status code %v for method '%v', but got %v", test.status, test.method, res.StatusCode)
		}
	}
}

func TestIndexPageNotFound(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "with trailing slash",
			url:  "https://kkn.fi/",
		},
		{
			name: "without trailing slash",
			url:  "https://kkn.fi",
		},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, test.url, nil)
		srv := vanity.Handler(
			vanity.VCSURL("https://github.com/kare"),
			vanity.Log(log.New(ioutil.Discard, "", 0)),
			vanity.StaticDir("/not-found", "/.static/"),
		)
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("%v: expected response status 404, but got %v\n", test.name, res.StatusCode)
		}
	}
}

func TestBrowserGoDoc(t *testing.T) {
	tests := []struct {
		path         string
		moduleServer string
		result       string
	}{
		{
			path:         "/gist",
			moduleServer: "https://pkg.go.dev/",
			result:       "https://pkg.go.dev/kkn.fi/gist",
		},
		{
			path:         "/gist/",
			moduleServer: "https://pkg.go.dev",
			result:       "https://pkg.go.dev/kkn.fi/gist",
		},
		{
			path:         "/set",
			moduleServer: "https://pkg.go.dev",
			result:       "https://pkg.go.dev/kkn.fi/set",
		},
		{
			path:         "/cmd/kkn.fi-srv",
			moduleServer: "https://pkg.go.dev",
			result:       "https://pkg.go.dev/kkn.fi/cmd/kkn.fi-srv",
		},
		{
			path:         "/cmd/tcpproxy/",
			moduleServer: "https://pkg.go.dev",
			result:       "https://pkg.go.dev/kkn.fi/cmd/tcpproxy",
		},
		{
			path:         "/pkgabc/sub/foo",
			moduleServer: "https://pkg.go.dev",
			result:       "https://pkg.go.dev/kkn.fi/pkgabc/sub",
		},
		{
			path:         "/pkgabc/sub/foo",
			moduleServer: "https://pkg.go.dev",
			result:       "https://pkg.go.dev/kkn.fi/pkgabc/sub",
		},
		{
			path:         "/vanity",
			moduleServer: "",
			result:       "https://pkg.go.dev/kkn.fi/vanity",
		},
		{
			path:         "/cmd/healthcheck",
			moduleServer: "https://github.com/kare",
			result:       "https://github.com/kare/healthcheck",
		},
		{
			path:         "/vanity",
			moduleServer: "https://github.com/kare/",
			result:       "https://github.com/kare/vanity",
		},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, addr+test.path, nil)
		srv := vanity.Handler(
			vanity.ModuleServerURL(test.moduleServer),
			vanity.Log(log.New(ioutil.Discard, "", 0)),
		)
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Errorf("expected response status %v, but got %v", http.StatusTemporaryRedirect, res.StatusCode)
		}
		body, _ := ioutil.ReadAll(res.Body)
		if !strings.Contains(string(body), test.result) {
			t.Errorf("expecting\n%v be contained in\n%v", test.result, string(body))
		}
	}
}

func TestGoTool(t *testing.T) {
	tests := []struct {
		path   string
		vcs    string
		vcsURL string
		result string
	}{
		{
			path:   "/gist?go-get=1",
			vcs:    "hg",
			vcsURL: "https://bitbucket.org/kare/",
			result: "kkn.fi/gist hg https://bitbucket.org/kare/gist",
		},
		{
			path:   "/set/?go-get=1",
			vcs:    "hg",
			vcsURL: "https://bitbucket.org/kare",
			result: "kkn.fi/set hg https://bitbucket.org/kare/set",
		},
		{
			path:   "/cmd/kkn.fi-srv?go-get=1",
			vcs:    "git",
			vcsURL: "https://github.com/kare/",
			result: "kkn.fi/cmd/kkn.fi-srv git https://github.com/kare/kkn.fi-srv",
		},
		{
			path:   "/cmd/kkn.fi-srv/?go-get=1",
			vcs:    "git",
			vcsURL: "https://github.com/kare",
			result: "kkn.fi/cmd/kkn.fi-srv git https://github.com/kare/kkn.fi-srv",
		},
		{
			path:   "/pkgabc/sub/foo?go-get=1",
			vcs:    "git",
			vcsURL: "https://github.com/kare",
			result: "kkn.fi/pkgabc/sub/foo git https://github.com/kare/pkgabc",
		},
		{
			path:   "/pkgabc/sub/foo/?go-get=1",
			vcs:    "git",
			vcsURL: "https://github.com/kare",
			result: "kkn.fi/pkgabc/sub/foo git https://github.com/kare/pkgabc",
		},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, addr+test.path, nil)
		srv := vanity.Handler(
			vanity.VCS(test.vcs),
			vanity.VCSURL(test.vcsURL),
			vanity.Log(log.New(ioutil.Discard, "", 0)),
		)
		srv.ServeHTTP(rec, req)

		res := rec.Result()
		body, _ := ioutil.ReadAll(res.Body)
		expected := fmt.Sprintf(`<meta name="go-import" content="%v">`, test.result)
		if !strings.Contains(string(body), expected) {
			t.Errorf("expecting url '%v' body to contain html meta tag:\n%v, but got:\n%v", test.path, expected, string(body))
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected response status 200, but got %v", res.StatusCode)
		}
	}
}

func TestStaticDir(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "with trailing slash",
			url:  "https://kkn.fi/",
		},
		{
			name: "without trailing slash",
			url:  "https://kkn.fi",
		},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, test.url, nil)
		srv := vanity.Handler(
			vanity.StaticDir("testdata", "dir"),
			vanity.Log(log.New(ioutil.Discard, "", 0)),
		)
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("%v: expected response status 200, but got %v", test.name, res.StatusCode)
		}

		body, _ := ioutil.ReadAll(res.Body)
		expected := "<html>homepage</html>\n"
		if string(body) != expected {
			t.Errorf("%v: expecting body to match:\n'%v', but got:\n'%s'", test.name, expected, body)
		}
	}
}

func TestRobotsTxt(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "GET /robots.txt",
			url:  "https://kkn.fi/robots.txt",
		},
	}
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, test.url, nil)
		srv := vanity.Handler(
			vanity.RobotsTxt(""),
			//vanity.Log(log.New(ioutil.Discard, "", 0)),
		)
		srv.ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("%v: expected response status 200, but got %v", test.name, res.StatusCode)
		}

		body, _ := ioutil.ReadAll(res.Body)
		expected := vanity.DefaultRobotsTxt
		if string(body) != expected {
			t.Errorf("%v: expecting body to match:\n'%v', but got:\n'%s'", test.name, expected, body)
		}
	}
}

func ExampleHandler() {
	errorLog := log.New(os.Stderr, "vanity: ", log.Ldate|log.Ltime|log.LUTC)
	srv := vanity.Handler(
		vanity.ModuleServerURL("https://pkg.go.dev"),
		vanity.Log(errorLog),
		vanity.VCSURL("https://github.com/kare"),
		vanity.VCS("git"),
		vanity.StaticDir("testdata", "/.static/"),
		vanity.IndexPageHandler(vanity.DefaultIndexPageHandler("testdata/index.html")),
	)
	http.Handle("/", srv)
	// Output:
}
