package astmod

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestListImports(t *testing.T) {
	// tests table
	tests := []struct {
		in      string
		prefix  string
		imports []string
		label   string
	}{
		{in: "ext1.go",
			prefix:  "github.com",
			imports: []string{"github.com/abc/xyz"},
			label:   "ext1 correct prefix",
		},
		{in: "ext1.go",
			prefix:  "gitlab.com",
			imports: nil,
			label:   "ext1 non-existent prefix",
		},
		{in: "helloworld.go",
			prefix:  "",
			imports: []string{"fmt"},
			label:   "helloworld",
		},
		{in: "int1.go",
			prefix:  "go.corp.example.com",
			imports: []string{"go.corp.example.com/abc/xyz"},
			label:   "int1",
		},
	}
	for _, tt := range tests {
		fname := filepath.Join("testdata", tt.in)
		imports, err := ListImports(fname, nil, tt.prefix)
		if err != nil {
			t.Fatal(err)
		}
		var importPaths []string
		for _, imp := range imports {
			importPaths = append(importPaths, imp.Path)
		}
		if !reflect.DeepEqual(importPaths, tt.imports) {
			t.Errorf("%s: bad result %v, want %v", tt.label, importPaths, tt.imports)
		}
	}
}
