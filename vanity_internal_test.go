package vanity

import "testing"

func TestValidateGoPackage(t *testing.T) {
	tests := []struct {
		name    string
		pkg     string
		wantErr bool
	}{
		{
			name:    "all allowed characters",
			pkg:     "/foo/bar/123/foo_abc",
			wantErr: false,
		},
		{
			name:    "contains query",
			pkg:     "/foobar?go-get=1",
			wantErr: false,
		},
		{
			name:    "contains dot",
			pkg:     "/foo.bar",
			wantErr: false,
		},
		{
			name:    "contains minus",
			pkg:     "/foo-bar",
			wantErr: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if err := validateGoPkgPath(test.pkg); (err != nil) != test.wantErr {
				t.Errorf("expecting error, none reported. error: %v", err)
			}
		})
	}
}
