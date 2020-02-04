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
	}
	// Option represents a functional option for configuring the vanity middleware.
	Option func(http.Handler)
	// Logger describes functions available for logging purposes.
	Logger interface {
		Printf(format string, v ...interface{})
	}
)

const (
	// pkgGoDev is the default module server by Google.
	pkgGoDev = "https://pkg.go.dev/"
	// searchGocenterIo is a module server by JFrog.
	searchGocenterIo = "https://search.gocenter.io/"
)

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
	url := h.browserURL(r.Host, r.URL.Path)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) browserURL(host, path string) string {
	switch h.moduleServerURL {
	case searchGocenterIo:
		pkg := strings.ReplaceAll(path, "/", "~2F")
		return fmt.Sprintf("%v%v%v/info", searchGocenterIo, host, pkg)
	case pkgGoDev:
		fallthrough
	default:
		pkg := path
		return fmt.Sprintf("%v%v%v", pkgGoDev, host, pkg)
	}
}

// Handler is an HTTP middleware that redirects browsers to Go module server
// (pkg.go.dev or similar) or Go tool to VCS repository. VCS repository is git
// by default. VCS can be set with VCS(). Configurable Logger defaults to
// os.Stdout. Logger can be configured with SetLogger(). Module server URL is
// https://pkg.go.dev/ and it can be configured via ModuleServerURL() func.
// VCSURL() func must be used to set VCS repository URL (such as
// https://github.com/kare/).
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

// ModuleServerURL sets Go module server address for browser redirect.
func ModuleServerURL(moduleServerURL string) Option {
	return func(h http.Handler) {
		v := h.(*handler)
		v.moduleServerURL = addSuffixSlash(moduleServerURL)
	}
}
