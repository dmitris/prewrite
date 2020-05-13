// Copyright Verizon Meida, Inc. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
//
// Author: Dmitry Savintsev <dsavints@verizonmedia.com>

// Plist list imports matching the given prefix.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dmitris/prewrite/astmod"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: plist -prefix [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	allOpt := flag.Bool("all", false, "list all imports (-p is ignored)")
	prefixOpt := flag.String("p", "", "import prefix to match")
	outputOpt := flag.String("o", "", "output file to write the imports listing")
	verboseOpt := flag.Bool("v", false, "verbose")
	flag.Usage = usage
	flag.Parse()
	prefix := *prefixOpt
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

	if prefix == "" && !*allOpt {
		log.Fatalf("invalid flags: no prefix given with -p and -all not set")
	}
	if *allOpt {
		prefix = ""
	}

	var (
		impts []string
		mu    sync.Mutex
		w     io.Writer
	)

	if *outputOpt == "" {
		w = os.Stdout
	} else {
		f, err := os.Open(*outputOpt)
		if err != nil {
			log.Fatalf("unable to open output file %s: %v", *outputOpt, err)
		}
		w = f
		defer f.Close()
	}
	lister := makeLister(w, prefix, &impts, &mu, verbose)
	_, err = os.Stat(root)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("Error - the traversal root %s does not exist, please double-check\n", root)
	}
	if err := filepath.Walk(root, lister); err != nil {
		log.Fatalf("Error listing imports %s: %s\n", flag.Arg(0), err)
	}
	uniq := map[string]bool{}
	var allImports []string
	for _, impt := range impts {
		if _, ok := uniq[impt]; !ok {
			uniq[impt] = true
			allImports = append(allImports, impt)
		}
	}
	sort.Strings(allImports)
	for _, impt := range allImports {
		fmt.Fprintln(w, impt)
	}
	return
}

// makeLister returns an imports listing function with parameters bound with a closure
func makeLister(w io.Writer,
	prefix string,
	imptsOut *[]string,
	mu *sync.Mutex,
	verbose bool) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			return nil
		}
		src, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Fatal error reading file %s\n", path)
		}
		impts, err := astmod.ListImportPaths(path, src, prefix)
		if err != nil {
			log.Fatalf("Fatal error rewriting AST, file %s - error: %s\n", path, err)
		}

		mu.Lock()
		for _, imp := range impts {
			*imptsOut = append(*imptsOut, imp)
		}
		mu.Unlock()

		return nil
	}
}
