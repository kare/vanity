package vanity

import (
	"net/http"
)

func ExampleBasic() {
	http.Handle("/", Redirect("git", "kkn.fi", "github.com/kare"))
}
