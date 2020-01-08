package vanity

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	// log is error log.
	log Logger
)

// Logger describes functions available for logging purposes.
type Logger interface {
	Printf(format string, v ...interface{})
}

// SetLogger sets the logger used by vanity package's error log.
func SetLogger(l Logger) {
	log = l
}

// Redirect is an HTTP middleware that redirects browsers to pkg.go.dev or
// Go tool to VCS repository.
func Redirect(vcs, importPath, repoRoot string) http.Handler {
	redirect := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme == "http" {
			r.URL.Scheme = "https"
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}
		if r.Method != http.MethodGet {
			status := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(status), status)
			return
		}

		if !strings.HasPrefix(strings.TrimSuffix(r.Host+r.URL.Path, "/"), importPath+"/") {
			http.NotFound(w, r)
			return
		}

		// Redirect browsers to Go module site.
		// Such as pkg.go.dev or something similar
		if r.FormValue("go-get") != "1" {
			goProxyHostname := "pkg.go.dev"
			url := "https://" + goProxyHostname + "/" + r.Host + r.URL.Path
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		var path string
		if strings.HasPrefix(r.URL.Path, "/cmd/") {
			path = r.URL.Path[4:]
		} else {
			path = r.URL.Path
		}

		// redirect github.com/kare/pkg/sub -> github.com/kare/pkg
		vcsroot := repoRoot
		f := func(c rune) bool { return c == '/' }
		shortPath := strings.FieldsFunc(path, f)
		if len(shortPath) > 0 {
			vcsroot = repoRoot + "/" + shortPath[0]
		}

		// Respond to Go tool with vcs info meta tag
		importRoot := r.Host + r.URL.Path
		meta := `<meta name="go-import" content="%v %v %v">`
		s := fmt.Sprintf(meta, importRoot, vcs, vcsroot)
		if _, err := w.Write([]byte(s)); err != nil {
			log.Printf("vanity: i/o error: %v", err)
		}
	}
	return http.HandlerFunc(redirect)
}
