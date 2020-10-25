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
		log              Logger
		vcs              string
		vcsURL           string
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
	Option func(http.Handler)
	// Logger describes functions available for logging purposes.
	Logger interface {
		Printf(format string, v ...interface{})
	}
)

const (
	// mPkgGoDev is the default module server by Google.
	mPkgGoDev = "https://pkg.go.dev/"
	// mSearchGocenterIo is a module server by JFrog.
	mSearchGocenterIo = "https://search.gocenter.io/"
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
		fmt.Fprintf(w, "%s", h.robotsTxt)
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
		shortPath := pathComponents(path)
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
	case mSearchGocenterIo:
		pkg := strings.ReplaceAll(path, "/", "~2F")
		return fmt.Sprintf("%v%v%v/info", mSearchGocenterIo, host, pkg)
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
		moduleServerURL: mPkgGoDev,
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

// Log sets the logger used by vanity package's error logger.
func Log(l Logger) Option {
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

// StaticDir serves a file system directory over HTTP. Given path is the local
// file system path to directory. Given urlPath is the path portition of the URL for the server.
func StaticDir(path, URLPath string) Option {
	return func(h http.Handler) {
		// TODO: path must be a readable directory or fail
		v := h.(*handler)
		dir := http.Dir(path)
		server := http.FileServer(dir)
		v.static = &staticDir{
			path:    path,
			uRLPath: URLPath,
			fs:      http.StripPrefix(URLPath, server),
		}
	}
}

// IndexPageHandler sets a handler for index.html page.
func IndexPageHandler(index http.Handler) Option {
	return func(h http.Handler) {
		v := h.(*handler)
		v.indexPageHandler = index
	}
}

// DefaultRobotsTxt is the default value for /robots.txt file.
var DefaultRobotsTxt = `user-agent: *
Allow: /$
Allow: /.static/*$
Disallow: /`

// RobotsTxt takes in option robotsTxt value. If value is empty, the value of DefaultRobotsTxt is used
func RobotsTxt(robotsTxt string) Option {
	return func(h http.Handler) {
		v := h.(*handler)
		if robotsTxt != "" {
			v.robotsTxt = robotsTxt
		} else {
			v.robotsTxt = DefaultRobotsTxt
		}
	}
}
