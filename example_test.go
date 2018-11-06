package vanity_test

import (
	"log"
	"net/http"
	"os"

	"kkn.fi/vanity"
)

func ExampleRedirect() {
	http.Handle("/", vanity.Redirect("git", "kkn.fi", "github.com/kare"))
	// Output:
}

func ExampleSetLogger() {
	errorLog := log.New(os.Stderr, "vanity: ", log.Ldate|log.Ltime|log.LUTC)
	vanity.SetLogger(errorLog)
	// Output:
}
