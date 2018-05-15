package vanity

import (
	stdlog "log"
	"net/http"
	"os"
)

func ExampleRedirect() {
	http.Handle("/", Redirect("git", "kkn.fi", "github.com/kare"))
}

func ExampleSetLogger() {
	// stdlog is Go Standard Library's log package
	errorLog := stdlog.New(os.Stderr, "vanity: ", stdlog.Ldate|stdlog.Ltime|stdlog.LUTC)
	SetLogger(errorLog)
}
