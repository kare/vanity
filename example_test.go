package vanity_test

import (
	"log"
	"net/http"
	"os"

	"kkn.fi/vanity"
)

func ExampleHandler() {
	errorLog := log.New(os.Stderr, "vanity: ", log.Ldate|log.Ltime|log.LUTC)
	srv := vanity.Handler(
		vanity.SetLogger(errorLog),
		vanity.VCSURL("https://github.com/kare"),
		vanity.VCS("git"),
	)
	http.Handle("/", srv)
	// Output:
}
