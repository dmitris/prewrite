prewrite
========

[![Build Status](https://travis-ci.org/dmitris/prewrite.svg?branch=master)](https://travis-ci.org/dmitris/prewrite)

Tool to rewrite import paths and [package import comments](https://golang.org/s/go14customimport) for vendoring by adding or removing a given path prefix.  The files are rewritten in-place with no backup (expectation is that version control is used), the output is gofmt'ed.

# Install
go get github.com/dmitris/prewrite

# Usage
prewrite -p prefix [-r] [-v] [path ...]

# Command-line arguments
* -p prefix -- prefix to add to imports and package import comments or remove (with -r) - required
* -r        -- remove the given prefix from import statements and package import comments
* -v        -- verbosely print the names of the changed files

If not provided, the path defaults to the current directory (will recursively traverse). Multiple targets can be given.

The last target parameter can be either a single file or a directory (such as a root of a source tree).

# Examples

Add a prefix to all imports (except the standard library) and package comment paths under the current directory:  
prewrite -p go.corp.company.com/x -v

Remove a prefix from all imports and package comment paths under the current directory:  
prewrite -p go.corp.company.com/x -r -v


