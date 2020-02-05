package vanity_test

import (
	"log"
	"net/http"
	"os"
	"strings"

	"kkn.fi/vanity"
)

func ExampleHandler() {
	errorLog := log.New(os.Stderr, "vanity: ", log.Ldate|log.Ltime|log.LUTC)
	indexPage := strings.NewReader("<html><h1>Vanity</h1></html>")
	srv := vanity.Handler(
		vanity.ModuleServerURL("https://pkg.go.dev"),
		vanity.SetLogger(errorLog),
		vanity.VCSURL("https://github.com/kare"),
		vanity.VCS("git"),
		vanity.IndexPage(indexPage),
	)
	http.Handle("/", srv)
	// Output:
}
