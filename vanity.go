package vanity

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type (
	handler struct {
		log              Logger
		vcs              string
		vcsURL           string
		host             string
		moduleServerURL  string
		static           *staticDir
		indexPageHandler http.Handler
		robotsTxt        string
	}
	staticDir struct {
		uRLPath string
		path    string
		fs      http.Handler
	}
	// Option represents a functional option for configuring the vanity middleware.
	Option func(http.Handler) error
	// Logger describes functions available for logging purposes.
	Logger interface {
		Printf(format string, v ...interface{})
	}
)

const (
	// mPkgGoDev is the default module server by Google.
	mPkgGoDev = "https://pkg.go.dev/"
	// mGitHub is not an actual module server, but a source code repository.
	mGitHub = "https://github.com/"
)

// DefaultIndexPageHandler serves given indexFilePath over HTTP via http.ServeFile(w, r, name).
func DefaultIndexPageHandler(indexFilePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexFilePath)
	})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	if h.static != nil && strings.HasPrefix(r.URL.Path, h.static.uRLPath) {
		h.static.fs.ServeHTTP(w, r)
		return
	}

	if r.URL.Path == "/" || r.URL.Path == "" {
		if h.indexPageHandler != nil {
			h.indexPageHandler.ServeHTTP(w, r)
			return
		}
		DefaultIndexPageHandler(h.static.path+"/index.html").ServeHTTP(w, r)
	}

	if r.URL.Path == "/robots.txt" {
		_, _ = fmt.Fprintf(w, "%s", h.robotsTxt)
		return
	}

	host := r.Host
	if h.host != "" {
		host = h.host
	}
	// Respond to Go tool with vcs info meta tag
	if r.FormValue("go-get") == "1" {
		path := r.URL.Path
		const cmd = "/cmd/"
		if strings.HasPrefix(r.URL.Path, cmd) {
			path = path[len(cmd):]
		}
		vcsroot := h.vcsURL
		shortPath := pathComponents(path)
		stripSubPackagesFromPath := len(shortPath) > 0
		if stripSubPackagesFromPath {
			vcsroot = h.vcsURL + shortPath[0]
		}

		importRoot := strings.TrimSuffix(host+r.URL.Path, "/")
		metaTag := fmt.Sprintf(`<meta name="go-import" content="%v %v %v">`, importRoot, h.vcs, vcsroot)
		if _, err := w.Write([]byte(metaTag)); err != nil {
			h.log.Printf("vanity: i/o error writing go tool http response: %v", err)
		}
		return
	}

	// Redirect browsers to Go module site.
	url := h.browserURL(host, r.URL.Path)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func pathComponents(path string) []string {
	f := func(c rune) bool {
		return c == '/'
	}
	return strings.FieldsFunc(path, f)
}

func (h *handler) browserURL(host, path string) string {
	// host = kkn.fi
	// path = /foo/bar

	if strings.HasPrefix(h.moduleServerURL, mGitHub) {
		pkg := path
		components := pathComponents(pkg)
		return stripSuffixSlash(h.moduleServerURL) + "/" + components[len(components)-1]
	}
	switch h.moduleServerURL {
	case mPkgGoDev:
		fallthrough
	default:
		pkg := path
		return fmt.Sprintf("%v%v%v", mPkgGoDev, host, pkg)
	}
}

func stripSuffixSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s[0 : len(s)-1]
	}
	return s
}

// Handler is an HTTP middleware that redirects browsers to Go module server
// (pkg.go.dev or similar) or Go cmd line tool to VCS repository. Handler can
// be configured by providing options. VCS repository is git by default. VCS
// can be set with VCS(). Configurable Logger defaults to os.Stderr. Logger can
// be configured with SetLogger(). Module server URL is https://pkg.go.dev/ and
// it can be configured via ModuleServerURL() func.  VCSURL() func must be used
// to set VCS repository URL (such as https://github.com/kare/).
func NewHandlerWithOptions(opts ...Option) (http.Handler, error) {
	v := &handler{
		log:             log.New(os.Stderr, "", log.LstdFlags),
		vcs:             "git",
		moduleServerURL: mPkgGoDev,
	}
	for _, option := range opts {
		if err := option(v); err != nil {
			return nil, err
		}
	}
	return v, nil
}

// VCS sets the version control type.
func VCS(vcs string) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		v.vcs = vcs
		return nil
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
	return func(h http.Handler) error {
		v := h.(*handler)
		v.vcsURL = addSuffixSlash(vcsURL)
		return nil
	}
}

// Host sets the hostname of the vanity server. Host defaults to HTTP request
// hostname.
func Host(host string) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		v.host = host
		return nil
	}
}

// Log sets the logger used by vanity package's error logger.
func Log(l Logger) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		v.log = l
		return nil
	}
}

// ModuleServerURL sets Go module server address for browser redirect.
func ModuleServerURL(moduleServerURL string) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		v.moduleServerURL = addSuffixSlash(moduleServerURL)
		return nil
	}
}

var ErrNotReadable = errors.New("vanity: static dir path directory is not readable")

// StaticDir serves a file system directory over HTTP. Given path is the local
// file system path to directory. Given urlPath is the path portion of the URL for the server.
func StaticDir(path, URLPath string) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("vanity: static dir path stat error: %w", err)
		}
		if !info.IsDir() {
			return errors.New("vanity: static dir path is not a directory")
		}
		const (
			read = 4
			// write = 2
			// exec = 1
		)
		const (
			readUser  = read << 6
			readGroup = read << 3
			readOther = read << 0
		)
		if info.Mode().Perm()&readOther == 0 {
			return ErrNotReadable
		}
		if info.Mode().Perm()&readGroup == 0 {
			return ErrNotReadable
		}
		if info.Mode().Perm()&readUser == 0 {
			return ErrNotReadable
		}
		dir := http.Dir(path)
		server := http.FileServer(dir)
		v.static = &staticDir{
			path:    path,
			uRLPath: URLPath,
			fs:      http.StripPrefix(URLPath, server),
		}
		return nil
	}
}

// IndexPageHandler sets a handler for index.html page.
func IndexPageHandler(index http.Handler) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		v.indexPageHandler = index
		return nil
	}
}

// DefaultRobotsTxt is the default value for /robots.txt file.
var DefaultRobotsTxt = `user-agent: *
Allow: /$
Allow: /.static/*$
Disallow: /`

// RobotsTxt takes in a value for robots.txt. If value is empty, the value of
// `DefaultRobotsTxt` is used
func RobotsTxt(robotsTxt string) Option {
	return func(h http.Handler) error {
		v := h.(*handler)
		if robotsTxt != "" {
			v.robotsTxt = robotsTxt
		} else {
			v.robotsTxt = DefaultRobotsTxt
		}
		return nil
	}
}
