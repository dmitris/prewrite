// Copyright 2015, Yahoo Inc. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
//
// Author: Dmitry Savintsev <dsavints@yahoo-inc.com>

package astmod

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"
)

// Rewrite modifies the AST to rewrite import statements and package import comments.
// src should be compatible with go/parser/#ParseFile:
// (The type of the argument for the src parameter must be string, []byte, or io.Reader.)
//
// return of nil, nil (no result, no error) means no changes are needed
func Rewrite(fname string, src interface{}, from, to string) (buf *bytes.Buffer, err error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fname, src, parser.ParseComments)
	if err != nil {
		log.Printf("Error parsing file %s, source: [%s], error: %s", fname, src, err)
		return nil, err
	}

	changed, err := RewriteImports(f, from, to)
	if err != nil {
		log.Printf("Error rewriting imports in the AST: file %s - %s", fname, err)
		return nil, err
	}
	changed2, err := RewriteImportComments(f, fset, from, to)
	if err != nil {
		log.Printf("Error rewriting import comments in the AST: file %s - %s", fname, err)
		return nil, err
	}
	if !changed && !changed2 {
		return nil, nil
	}
	buf = &bytes.Buffer{}
	err = format.Node(buf, fset, f)
	return buf, err
}

// RewriteImports rewrites imports in the passed AST (in-place).
// It returns bool changed set to true if any changes were made
// and non-nil err on error.
func RewriteImports(f *ast.File, from, to string) (changed bool, err error) {
	for _, impNode := range f.Imports {
		imp, err := strconv.Unquote(impNode.Path.Value)
		if err != nil {
			log.Printf("Error unquoting import value %v - %s\n", impNode.Path.Value, err)
			return false, err
		}
		// skip standard library imports and relative references
		if !strings.Contains(imp, ".") || strings.HasPrefix(imp, ".") {
			continue
		}

		if strings.HasPrefix(imp, from) {
			changed = true
			newimp := strings.Replace(impNode.Path.Value, from, to, 1)
			impNode.Path.Value = newimp
		}
	}
	return
}

// RewriteImportComments rewrites package import comments (https://golang.org/s/go14customimport)
func RewriteImportComments(f *ast.File, fset *token.FileSet, from, to string) (changed bool, err error) {
	pkgpos := fset.Position(f.Package)
	// Print the AST.
	// ast.Print(fset, f)
	newcommentgroups := make([]*ast.CommentGroup, 0)
	for _, c := range f.Comments {
		commentpos := fset.Position(c.Pos())
		// keep the comment if we are not on the "package <X>" line
		// or the comment after the package statement does not look like import comment
		if commentpos.Line != pkgpos.Line ||
			!strings.HasPrefix(c.Text(), `import "`) {
			newcommentgroups = append(newcommentgroups, c)
			continue
		}
		parts := strings.Split(strings.Trim(c.Text(), "\n\r\t "), " ")
		oldimp, err := strconv.Unquote(parts[1])
		if err != nil {
			log.Fatalf("Error unquoting import value [%v] - %s\n", parts[1], err)
		}
		// if the prefix is not there = nothing to remove, keep the comment
		if !strings.HasPrefix(oldimp, from) {
			newcommentgroups = append(newcommentgroups, c)
			continue
		}
		newimp := strings.Replace(oldimp, from, to, 1)
		changed = true
		c2 := ast.Comment{Slash: c.Pos(), Text: `// import ` + strconv.Quote(newimp)}
		cg := ast.CommentGroup{List: []*ast.Comment{&c2}}
		newcommentgroups = append(newcommentgroups, &cg)
	}
	// change the AST only if there are pending mods
	if changed {
		f.Comments = newcommentgroups
	}
	return changed, nil
}
