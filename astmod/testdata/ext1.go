// package doc comment for go.corp.company.com/abc/xyz
package bar // import "go.corp.company.com/abc/xyz"

import _ "go.corp.company.com/abc/xyz"

func Bar() {
	// unrelated comment in func
	println("Hello, World! (from go.corp.company.com/abc/xyz)")
}
