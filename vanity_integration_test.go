package vanity_test

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kkn.fi/vanity"
)

func tmpdir(t *testing.T) string {
	t.Helper()
	dir, err := os.CreateTemp("/tmp", "permission-")
	if err != nil {
		t.Fatalf("error while creating test temp dir: %v", err)
	}
	return dir.Name()
}

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}

func TestStaticDirPermissionsIntegration(t *testing.T) {
	integrationTest(t)
	tests := []struct {
		name        string
		permissions fs.FileMode
		wantErr     bool
	}{
		{
			name:        "directory doesn't have any permissions",
			permissions: 0000,
			wantErr:     true,
		},
		{
			name:        "directory has all exec permissions",
			permissions: 0111,
			wantErr:     true,
		},
		{
			name:        "directory has all write permissions",
			permissions: 0222,
			wantErr:     true,
		},
		{
			name:        "directory has all read permissions",
			permissions: 0444,
			wantErr:     true,
		},
		{
			name:        "directory has other read permissions",
			permissions: 0004,
			wantErr:     true,
		},
		{
			name:        "directory has group read permissions",
			permissions: 0040,
			wantErr:     true,
		},
		{
			name:        "directory has user read permissions",
			permissions: 0400,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		testCase := test
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "https://kkn.fi", nil)
			path := tmpdir(t)
			if err := os.Chmod(path, testCase.permissions); err != nil {
				t.Fatalf("error while setting tmp dir permissions: %v", err)
			}

			srv, err := vanity.NewHandlerWithOptions(
				vanity.StaticDir(path, "/files/"),
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

func TestStaticDirIntegration(t *testing.T) {
	integrationTest(t)
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
		{
			name:    "directory is a file",
			path:    "/etc/passwd",
			wantErr: true,
		},
		{
			name:    "directory doesn't have read permissions",
			path:    "",
			wantErr: true,
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
