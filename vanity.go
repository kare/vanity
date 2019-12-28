package vanity

import (
	"bytes"
	"html/template"
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

type data struct {
	ImportRoot string
	VCS        string
	VCSRoot    string
}

var tmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.VCSRoot}}">
</head>
</html>
`))

// Redirect is a HTTP middleware that redirects browsers to pkg.go.dev or
// Go tool to VCS repository.
func Redirect(vcs, importPath, repoRoot string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		if r.FormValue("go-get") != "1" {
			url := "https://pkg.go.dev/" + r.Host + r.URL.Path
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

		d := &data{
			ImportRoot: r.Host + r.URL.Path,
			VCS:        vcs,
			VCSRoot:    vcsroot,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, d); err != nil {
			log.Printf("vanity: template execution error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=300")
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Printf("vanity: i/o error: %v", err)
		}
	})
}
