package vanity

import (
	"net/http"
)

func ExampleRedirect() {
	http.Handle("/", Redirect("git", "kkn.fi", "github.com/kare"))
}
