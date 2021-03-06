// Copyright 2015, Yahoo Inc. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
//
// Author: Dmitry Savintsev <dsavints@yahoo-inc.com>

// prewrite tool rewrites import paths and package import comments for vendoring
// by adding or removing a given path prefix. The files are rewritten
// in-place with no backup (expectation is that version control is used),
// the output is gofmt'ed.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmitris/prewrite/astmod"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: prewrite [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	fromOpt := flag.String("from", "", "old prefix to be changed")
	toOpt := flag.String("to", "", "new prefix to add instead of the old onw")
	verboseOpt := flag.Bool("v", false, "verbose")
	flag.Usage = usage
	flag.Parse()
	fromPrefix := *fromOpt
	toPrefix := *toOpt
	verbose := *verboseOpt

	var root string
	var err error
	if len(flag.Args()) < 1 {
		root, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		root = flag.Arg(0)
	}

	log.Printf("from %s to %s", fromPrefix, toPrefix)
	if fromPrefix == "" || toPrefix == "" {
		usage()
		os.Exit(1)
	}

	// add trailing slash if not already there
	if fromPrefix[len(fromPrefix)-1] != '/' {
		fromPrefix += "/"
	}
	if toPrefix[len(toPrefix)-1] != '/' {
		toPrefix += "/"
	}

	processor := makeVisitor(fromPrefix, toPrefix, verbose)
	_, err = os.Stat(root)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("Error - the traversal root %s does not exist, please double-check\n", root)
	}
	err = filepath.Walk(root, processor)
	if err != nil {
		log.Fatalf("Error processing %s: %s\n", flag.Arg(0), err)
	}

}

// makeVisitor returns a rewriting function with parameters bound with a closure
func makeVisitor(from, to string, verbose bool) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			return nil
		}
		// special cases
		if skipFile(path) {
			return nil
		}
		src, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Fatal error reading file %s\n", path)
		}
		buf, err := astmod.Rewrite(path, src, from, to)
		if err != nil {
			log.Fatalf("Fatal error rewriting AST, file %s - error: %s\n", path, err)
		}
		// check if there were any mods done for the file, return if non
		if buf == nil {
			return nil
		}
		err = ioutil.WriteFile(path, buf.Bytes(), f.Mode())
		if err != nil {
			log.Fatalf("Fatal error - unable to write to file %s: %s\n", path, err)
		}
		if verbose {
			fmt.Println(path)
		}
		return nil
	}
}

func skipFile(fname string) bool {
	// known special cases
	skip := [...]string{
		"golang.org/x/tools/go/loader/testdata/badpkgdecl.go",
	}
	for _, s := range skip {
		if strings.HasSuffix(fname, s) {
			return true
		}
	}
	return false
}
