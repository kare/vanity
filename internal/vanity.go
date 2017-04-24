package internal // import "kkn.fi/cmd/vanity/internal"

import (
	"fmt"
	"net/http"
	"strings"
)

type (
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
		// Name is the name of the Go package.
		Name string
		// VCS is version control system used by the project.
		VCS *VCS
	}
	// Server is the actual HTTP server for Go vanity domains.
	Server struct {
		// Domain is the vanity domain.
		Domain string
		// Packages contains settings for vanity packages.
		Packages []*Package
	}
)

func (p Package) name() string {
	path := p.Name
	c := strings.Index(path, "/")
	if c == -1 {
		return path
	}
	path = path[c+1:]
	c = strings.Index(path, "/")
	if c == -1 {
		return path
	}
	return path[:c]

}

// NewVCS creates a new VCS base on system and url.
func NewVCS(system, url string) *VCS {
	v := &VCS{
		System: system,
		URL:    url,
	}
	return v
}

func (v VCS) String() string {
	return fmt.Sprintf("VCS{%v, %v}", v.System, v.URL)
}

// NewPackage returns a new Package given a path and VCS.
func NewPackage(path string, vcs *VCS) *Package {
	p := &Package{
		Name: path,
		VCS:  vcs,
	}
	return p
}

func (p Package) String() string {
	return fmt.Sprintf("Package{%v, %v}", p.Name, p.VCS)
}

// NewServer returns a new Vanity Server given domain name and
// vanity package configuration.
func NewServer(domain string, config []*Package) *Server {
	s := &Server{
		Domain:   domain,
		Packages: config,
	}
	return s
}

// goMetaContent creates a value from the <meta/> tag content attribute.
func (v VCS) goMetaContent() string {
	return fmt.Sprintf("%v %v", v.System, v.URL)
}

// goDocURL returns the HTTP URL to godoc.org.
func (p Package) goDocURL() string {
	return fmt.Sprintf("https://godoc.org/%v", p.Name)
}

// goImportLink creates the link used in HTML <meta/> tag
// where domain is the domain name of the server.
func (p Package) goImportLink(domain string) string {
	path := p.name()
	return fmt.Sprintf("%v/%v %v", domain, path, p.VCS.goMetaContent())
}

// goImportMeta creates the <meta/> HTML tag containing name and content attributes.
func (p Package) goImportMeta(domain string) string {
	link := p.goImportLink(domain)
	return fmt.Sprintf(`<meta name="go-import" content="%s">`, link)
}

func (s Server) find(path string) *Package {
	for _, p := range s.Packages {
		if p.Name == s.Domain+path {
			return p
		}
	}
	return nil
}

// ServeHTTP is an HTTP Handler for Go vanity domain.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	pack := s.find(r.URL.Path)
	if pack == nil {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("go-get") != "1" {
		url := pack.goDocURL()
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprint(w, pack.goImportMeta(s.Domain))
}
