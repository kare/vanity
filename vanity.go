package vanity

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type (
	handler struct {
		log             Logger
		vcs             string
		vcsURL          string
		moduleServerURL string
		indexPage       []byte
	}
	// Option represents a functional option for configuring the vanity middleware.
	Option func(http.Handler)
	// Logger describes functions available for logging purposes.
	Logger interface {
		Printf(format string, v ...interface{})
	}
)

const pkgGoDev = "https://pkg.go.dev/"

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if r.URL.Path == "/" || r.URL.Path == "" {
		http.NotFound(w, r)
		return
	}

	// Respond to Go tool with vcs info meta tag
	if r.FormValue("go-get") == "1" {
		path := r.URL.Path
		const cmd = "/cmd/"
		if strings.HasPrefix(r.URL.Path, cmd) {
			path = path[len(cmd):]
		}
		// redirect github.com/kare/pkg/sub -> github.com/kare/pkg
		vcsroot := h.vcsURL
		f := func(c rune) bool { return c == '/' }
		shortPath := strings.FieldsFunc(path, f)
		if len(shortPath) > 0 {
			vcsroot = h.vcsURL + shortPath[0]
		}

		importRoot := strings.TrimSuffix(r.Host+r.URL.Path, "/")
		metaTag := fmt.Sprintf(`<meta name="go-import" content="%v %v %v">`, importRoot, h.vcs, vcsroot)
		if _, err := w.Write([]byte(metaTag)); err != nil {
			h.log.Printf("vanity: i/o error writing go tool http response: %v", err)
		}
		return
	}

	// Redirect browsers to Go module site.
	url := fmt.Sprintf("%v%v%v", pkgGoDev, r.Host, r.URL.Path)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handler is an HTTP middleware that redirects browsers to pkg.go.dev or Go
// tool to VCS repository. VCS repository is git by default. Configurable Logger
// defaults to os.Stdout.
func Handler(opts ...Option) http.Handler {
	v := &handler{
		log:             log.New(os.Stdout, "", log.LstdFlags),
		vcs:             "git",
		moduleServerURL: pkgGoDev,
	}
	for _, option := range opts {
		option(v)
	}
	return v
}

// VCS sets the version control type.
func VCS(vcs string) Option {
	return func(h http.Handler) {
		v := h.(*handler)
		v.vcs = vcs
	}
}

func addSuffixSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

// VCSURL sets the VCS repository url address.
func VCSURL(vcsURL string) Option {
	return func(h http.Handler) {
		v := h.(*handler)
		v.vcsURL = addSuffixSlash(vcsURL)
	}
}

// SetLogger sets the logger used by vanity package's error logger.
func SetLogger(l Logger) Option {
	return func(h http.Handler) {
		v := h.(*handler)
		v.log = l
	}
}
