// package doc comment for github.com/foo/bar
package bar // import "github.com/foo/bar"

import _ "github.com/abc/xyz"

func Bar() {
	// unrelated comment in func
	println("Hello, World! (from github.com/foo/bar)")
}
