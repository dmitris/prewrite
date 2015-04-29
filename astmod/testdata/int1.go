// package doc comment for github.com/foo/bar
package bar // import "go.corp.example.com/x/github.com/foo/bar"

import _ "go.corp.example.com/x/github.com/abc/xyz"

func Bar() {
	// unrelated comment in func
	println("Hello, World! (from github.com/foo/bar)")
}
