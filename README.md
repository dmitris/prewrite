prewrite
========

[![Build Status](https://travis-ci.org/dmitris/prewrite.svg?branch=master)](https://travis-ci.org/dmitris/prewrite)

Tool to rewrite import paths and [package import comments](https://golang.org/s/go14customimport) for vendoring by adding or removing a given path prefix.  The files are rewritten in-place with no backup (expectation is that version control is used), the output is gofmt'ed.

# Install
go get github.com/dmitris/prewrite

# Usage
prewrite -p prefix [-r] [-v] [path]

# Command-line arguments
* -from &lt;oldprefix&gt; -to &lt;newprefix&gt; -- rewrite import paths replacing oldprefix with newprefix
* -v        -- verbosely print the names of the changed files

Example: `prewrite -from go.oldcompany.com -new go.newcompany.com`.

If not provided, the path defaults to the current directory (will recursively traverse).  Either a single file or a directory (such as a root of a source tree) can be given.

# Examples

Change the prefix in all the imports (except the standard library) and package comment paths under the current directory:  
prewrite -from go.stealthy.com -to go.nextunicorn.com

Remove a prefix from all imports and package comment paths under the directory /tmp/foobar :  
prewrite -from go.stealthy.company.com/go.theunicorn.com -to go.nextunicorn.com /tmp/foobar


