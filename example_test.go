package vanity

import (
	"net/http"
)

func ExampleBasic() {
	http.Handle("/", GoDocRedirect("git", "kkn.fi", "github.com/kare"))
}
