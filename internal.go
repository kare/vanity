package main // import "kkn.fi/cmd/vanity"

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	// packageConfig defines Go package that has vanity import defined by Path,
	// VCS system type and VCS URL.
	packageConfig struct {
		// Name is the name of the Go package.
		Name string
		// VCS is version control system used by the project.
		VCS string
		// URL of the git repository
		URL string
	}
	// vanityServer is the actual HTTP server for Go vanity domains.
	vanityServer struct {
		// Domain is the vanity domain.
		Domain string
		// Packages contains settings for vanity Packages.
		Packages map[string]*packageConfig
	}
)

// newPackage returns a new Package given a path and VCS.
func newPackage(name, vcs, url string) *packageConfig {
	return &packageConfig{
		Name: name,
		VCS:  vcs,
		URL:  url,
	}
}

// newServer returns a new Vanity Server given domain name and
// vanity package configuration.
func newServer(domain string, config map[string]*packageConfig) *vanityServer {
	return &vanityServer{
		Domain:   domain,
		Packages: config,
	}
}

// path returns the path portition of the package name
func (p packageConfig) path() string {
	i := strings.Index(p.Name, "/")
	return p.Name[i+1:]
}

// goDocURL returns the HTTP URL to godoc.org.
func (p packageConfig) goDocURL(domain string) string {
	return fmt.Sprintf("https://godoc.org/%v%v", domain, p.Name)
}

// goImportMeta creates the <meta/> HTML tag containing name and content attributes.
func (p packageConfig) goImportMeta(domain string) string {
	s := `<meta name="go-import" content="%v/%v %v %v">`
	return fmt.Sprintf(s, domain, p.path(), p.VCS, p.URL)
}

// ServeHTTP is an HTTP Handler for Go vanity domain.
func (s vanityServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	pack, ok := s.Packages[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("go-get") != "1" {
		url := pack.goDocURL(s.Domain)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprint(w, pack.goImportMeta(s.Domain))
}
