package astmod

import (
	"go/parser"
	"go/token"
	"log"
	"sort"
	"strings"
)

// ImportSpec is a slimmed-down version of https://golang.org/pkg/go/ast/#ImportSpec.
type ImportSpec struct {
	Name, Path string
}

// String prints the import with the name field prepended (if not empty).
func (i *ImportSpec) String() string {
	if i == nil {
		return ""
	}
	if i.Name != "" {
		return i.Name + " " + i.Path
	}
	return i.Path
}

// ListImports returns a set of ImportSpec objects with imports
// matching the given prefix.
func ListImports(fname string, src interface{}, prefix string) ([]*ImportSpec, error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fname, src, parser.ImportsOnly)
	if err != nil {
		log.Printf("Error parsing file %s, source: [%s], error: %s", fname, src, err)
		return nil, err
	}
	var ret []*ImportSpec
	for _, imp := range f.Imports {
		path := strings.Replace(imp.Path.Value, `"`, "", -1)
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		spec := &ImportSpec{}
		if imp.Name != nil {
			spec.Name = imp.Name.Name
		}
		spec.Path = path
		ret = append(ret, spec)
	}
	return ret, nil
}

// ListImportPaths returns a slice of strings for the unique imports paths
// matching the given prefix.
func ListImportPaths(fname string, src interface{}, prefix string) ([]string, error) {
	imps, err := ListImports(fname, src, prefix)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, imp := range imps {
		ret = append(ret, imp.Path)
	}
	sort.Strings(ret)
	return ret, nil
}
