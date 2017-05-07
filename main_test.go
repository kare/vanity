package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	config := `/gist	git	https://github.com/kare/gist

/vanity	git	https://github.com/kare/vanity
`
	expected := map[string]*Package{
		"/gist":   NewPackage("/gist", "git", "https://github.com/kare/gist"),
		"/vanity": NewPackage("/vanity", "git", "https://github.com/kare/vanity"),
	}
	packages, err := readConfig(strings.NewReader(config))
	if err != nil {
		t.Fatalf("unexcepted configuration error: %v", err)
	}
	if len(packages) != 2 {
		t.Fatalf("expecting config for %v packages, but got %v", 4, len(packages))
	}
	for key, pack := range expected {
		if !reflect.DeepEqual(pack, packages[key]) {
			t.Fatalf("expected %v but got %v", pack, packages[key])
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
