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

// Rewrite modifies the AST to rewrite import statements.
// src should be compatible with go/parser/#ParseFile:
// (The type of the argument for the src parameter must be string, []byte, or io.Reader.)
//
// return of nil, nil (no result, no error) means no changes are needed
func Rewrite(fname string, src interface{}, fromPrefix, toPrefix string) (buf *bytes.Buffer, err error) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, fname, src, parser.ParseComments)
	if err != nil {
		log.Printf("Error parsing file %s, source: [%s], error: %s", fname, src, err)
		return nil, err
	}

	changed, err := RewriteImports(f, fromPrefix, toPrefix)
	if err != nil {
		log.Printf("Error rewriting imports in the AST: file %s - %s", fname, err)
		return nil, err
	}
	if !changed {
		return nil, nil
	}
	buf = &bytes.Buffer{}
	err = format.Node(buf, fset, f)
	return buf, err
}

// RewriteImports rewrites imports in the passed AST (in-place).
// It returns bool changed set to true if any changes were made
// and non-nil err on error
func RewriteImports(f *ast.File, fromPrefix, toPrefix string) (changed bool, err error) {
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
		if strings.HasPrefix(imp, fromPrefix) {
			changed = true
			impNode.Path.Value = strconv.Quote(toPrefix + imp[len(fromPrefix):])
		}
	}
	return
}
