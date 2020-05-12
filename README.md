prewrite
========

[![Build Status](https://travis-ci.org/dmitris/prewrite.svg?branch=master)](https://travis-ci.org/dmitris/prewrite)

The repository contains `plist` and `prewrite` tools under `cmd/` subdirectory.

# Install
go get github.com/dmitris/prewrite/...

# prewrite
## Usage:
```
prewrite -p prefix [-r] [-v] [path ...]
```

## Command-line arguments
* -from -- the old prefix to change in the imports statements (required)
* - to -- new prefix to add instead of the old one
* -v        -- verbosely print the names of the changed files

If not provided, the path defaults to the current directory (will recursively traverse).

The last target parameter can be either a single file or a directory (such as a root of a source tree).

## Examples

Change the prefix for all imports (except the standard library) under the current directory
from "github.com/*" to "gitlab.com/*"
`prewrite -from github.com -to gitlab.com`

# plist
## Usage
```
usage: plist -prefix [path ...]
  -all
    	list all imports (-p is ignored)
  -o string
    	output file to write the imports listing
  -p string
    	import prefix to match
  -v	verbose
```

