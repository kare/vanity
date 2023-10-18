package vanity_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kkn.fi/vanity"
)

func TestStaticDirIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "local directory not found",
			path:    "/not/found",
			wantErr: true,
		},
		{
			name:    "local directory found",
			path:    "/tmp",
			wantErr: false,
		},
	}
	for _, test := range tests {
		testCase := test
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "https://kkn.fi", nil)
			srv, err := vanity.NewHandlerWithOptions(
				vanity.StaticDir(testCase.path, "/files/"),
			)
			if (err != nil) != testCase.wantErr {
				t.Errorf("expecting error: %v, got '%v'", testCase.wantErr, err)
				return
			}
			if srv == nil {
				return
			}
			srv.ServeHTTP(rec, req)
			res := rec.Result()
			if res.StatusCode != http.StatusNotFound {
				t.Errorf("%v: expected response status 404, but got %v\n", testCase.name, res.StatusCode)
			}
		})
	}
}
