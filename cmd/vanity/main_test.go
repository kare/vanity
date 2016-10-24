package main

import (
	"strings"
	"testing"

	"kkn.fi/vanity"
)

func TestParseConfig(t *testing.T) {
	config := `/gist	git	https://github.com/kare/gist
/vanity	git	https://github.com/kare/vanity
/vanity/cmd	git	https://github.com/kare/vanity

/vanity/cmd/vanity	git	https://github.com/kare/vanity`

	expected := map[vanity.Path]vanity.Package{
		"/gist":              *vanity.NewPackage("/gist", "git", "https://github.com/kare/gist"),
		"/vanity":            *vanity.NewPackage("/vanity", "git", "https://github.com/kare/vanity"),
		"/vanity/cmd":        *vanity.NewPackage("/vanity", "git", "https://github.com/kare/vanity"),
		"/vanity/cmd/vanity": *vanity.NewPackage("/vanity", "git", "https://github.com/kare/vanity"),
	}
	conf, err := readConfig(strings.NewReader(config))
	if err != nil {
		t.Fatalf("unexcepted configuration error: %v", err)
	}
	if len(conf) != 4 {
		t.Fatalf("expecting config for %v packages, but got %v", 4, len(conf))
	}
	for p, c := range expected {
		if c != conf[p] {
			t.Fatalf("expected %v but got %v", c, conf[p])
		}
	}
}

func TestParseBrokenConfig(t *testing.T) {
	config := "/gist git"
	_, err := readConfig(strings.NewReader(config))
	if err == nil {
		t.Fatal("broken configuration did not return error")
	}
}
