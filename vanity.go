package vanity // import "kkn.fi/vanity"

import (
	"fmt"
	"net/http"
)

type (
	// Path is the path component of the HTTP request sent by Go tool or browser.
	Path string
	// VCS is the Version Control System used by the Go project.
	VCS struct {
		// System defines which version control system is used.
		// Usually git or hg.
		System string
		// VCSURL is the HTTPS URL for project's version control system.
		// Usually github.com or bitbucket.org address.
		URL string
	}
	// Package defines Go package that has vanity import defined by Path,
	// VCS system type and VCS URL.
	Package struct {
		// Path is path component of vanity url.
		Path *Path
		// VCS is version control system used by the project.
		VCS *VCS
	}
	// Server is the actual HTTP server for Go vanity domains.
	Server struct {
		// Domain is the vanity domain.
		Domain string
		// Config contains settings for vanity packages.
		Config map[Path]Package
	}
)

func NewPath(path string) *Path {
	p := Path(path)
	return &p
}

func NewVCS(system, url string) *VCS {
	v := &VCS{
		System: system,
		URL:    url,
	}
	return v
}

// NewPackage returns a new Package given a path and VCS.
func NewPackage(path *Path, vcs *VCS) *Package {
	p := &Package{
		Path: path,
		VCS:  vcs,
	}
	return p
}

// NewServer returns a new Vanity Server given domain name and vanity package configuration.
func NewServer(domain string, config map[Path]Package) *Server {
	s := &Server{
		Domain: domain,
		Config: config,
	}
	return s
}

// GoMetaContent creates a value from the <meta/> tag content attribute.
func (v VCS) GoMetaContent() string {
	return fmt.Sprintf("%v %v", v.System, v.URL)
}

// GoDocURL returns the HTTP URL to godoc.org.
func (p Package) GoDocURL(domain, path string) string {
	return fmt.Sprintf("https://godoc.org/%v%v", domain, path)
}

// GoImportLink creates the link used in HTML <meta/> tag
// where domain is the domain name of the server.
func (p Package) GoImportLink(domain string) string {
	return fmt.Sprintf("%v%v %v", domain, *p.Path, p.VCS.GoMetaContent())
}

// GoImportMeta creates the <meta/> HTML tag containing name and content attributes.
func (p Package) GoImportMeta(domain string) string {
	link := p.GoImportLink(domain)
	return fmt.Sprintf(`<meta name="go-import" content="%s">`, link)
}

// ServerHTTP is a HTTP Handler for Go vanity domain.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	pack, ok := s.Config[Path(r.URL.Path)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("go-get") != "1" {
		url := pack.GoDocURL(s.Domain, r.URL.Path)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprint(w, pack.GoImportMeta(s.Domain))
}
