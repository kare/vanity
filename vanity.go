package vanity // import "kkn.fi/vanity"

import (
	"fmt"
	"net/http"
)

type (
	// Path is the path component of the request sent by Go tool or browser.
	Path string
	// Package defines Go package that has vanity import defined.
	Package struct {
		// Path is path component of vanity url.
		Path string
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

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conf, ok := s.Config[Path(r.URL.Path)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("go-get") != "1" {
		url := fmt.Sprintf("https://godoc.org/%v%v", *s.Domain, r.URL.Path)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	i := fmt.Sprintf("%s%s %s %s", *s.Domain, conf.Path, conf.VCSSystem, conf.VCSURL)
	fmt.Fprintf(w, `<meta name="go-import" content="%s">`, i)
}
