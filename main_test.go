package main

import (
	"reflect"
	"strings"
	"testing"

	"kkn.fi/cmd/vanity/internal"
)

func TestParseConfig(t *testing.T) {
	config := `/gist	git	https://github.com/kare/gist

/vanity	git	https://github.com/kare/vanity
`
	expected := []*internal.Package{
		internal.NewPackage("/gist", internal.NewVCS("git", "https://github.com/kare/gist")),
		internal.NewPackage("/vanity", internal.NewVCS("git", "https://github.com/kare/vanity")),
	}
	packages, err := readConfig(strings.NewReader(config))
	if err != nil {
		t.Fatalf("unexcepted configuration error: %v", err)
	}
	if len(packages) != 2 {
		t.Fatalf("expecting config for %v packages, but got %v", 4, len(packages))
	}
	for i, pack := range expected {
		if !reflect.DeepEqual(pack, packages[i]) {
			t.Fatalf("expected %v but got %v", pack, packages[i])
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
