package vanity // import "kkn.fi/vanity"

import (
	"fmt"
	"net/http"
)

type (
	// Path is the path component of the HTTP request sent by Go tool or browser.
	Path string
	// Package defines Go package that has vanity import defined by Path,
	// VCS system type and VCS URL.
	Package struct {
		// Path is path component of vanity url.
		Path Path
		// VCSSystem is version control system used by the project.
		// Usually 'git' or 'hg'.
		VCSSystem string
		// VCSURL is the HTTP URL for project's version control system.
		VCSURL string
	}
	// Server is the actual HTTP server for Go vanity domains.
	Server struct {
		// Domain is the vanity domain.
		Domain *string
		// Config contains settings for vanity packages.
		Config map[Path]Package
	}
)

// NewPackage returns a new Package given a path, VCS system and VCS URL.
func NewPackage(path, vcssystem, vcsurl string) *Package {
	p := &Package{
		Path:      Path(path),
		VCSSystem: vcssystem,
		VCSURL:    vcsurl,
	}
	return p
}

// NewServer returns a new Vanity Server given domain name and vanity package configuration.
func NewServer(domain string, config map[Path]Package) *Server {
	s := &Server{
		Domain: &domain,
		Config: config,
	}
	return s
}

// GoDocURL returns the HTTP URL to godoc.org.
func (p Package) GoDocURL(domain, path string) string {
	return fmt.Sprintf("https://godoc.org/%v%v", domain, path)
}

// GoImportLink creates the link used in HTML <meta/> tag
// where domain is the domain name of the server.
func (p Package) GoImportLink(domain string) string {
	return fmt.Sprintf("%s%s %s %s", domain, p.Path, p.VCSSystem, p.VCSURL)
}

// GoImportMeta creates the <meta/> HTML tag containing name and content attributes.
func (p Package) GoImportMeta(domain string) string {
	link := p.GoImportLink(domain)
	return fmt.Sprintf(`<meta name="go-import" content="%s">`, link)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	pack, ok := s.Config[Path(r.URL.Path)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("go-get") != "1" {
		url := pack.GoDocURL(*s.Domain, r.URL.Path)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprint(w, pack.GoImportMeta(*s.Domain))
}
