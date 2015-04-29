// Copyright 2015, Yahoo Inc. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
//
// Author: Dmitry Savintsev <dsavints@yahoo-inc.com>

package astmod

import (
	"io/ioutil"
	"os"
	"testing"
)

// TODO: change for OSS!
const pkgname = "go.corp.yahoo.com/dsavints/prewrite/astmod"
const prefix = `go.corp.example.com/x`

var input map[string]string
var helloworld, ext1, int1 string

// init reads testdata files and puts them in the input map
func init() {
	input = make(map[string]string)
	gopath := os.Getenv("GOPATH")
	base := gopath + "/src/" + pkgname + "/testdata/"
	var err error
	var b []byte
	b, err = ioutil.ReadFile(base + "helloworld.go")
	if err != nil {
		panic(err)
	}
	input["helloworld"] = string(b)

	b, err = ioutil.ReadFile(base + "int1.go")
	if err != nil {
		panic(err)
	}
	input["int1"] = string(b)

	b, err = ioutil.ReadFile(base + "ext1.go")
	if err != nil {
		panic(err)
	}
	input["ext1"] = string(b)
}

// tests table
var tests = []struct {
	in      string
	remove  bool
	wanted  string
	changed bool
	label   string
}{
	{in: "ext1",
		remove:  false,
		wanted:  "int1",
		changed: true,
		label:   "ext1",
	},
	{in: "int1",
		remove:  true,
		wanted:  "ext1",
		changed: true,
		label:   "int1",
	},
	// try to call rewrite on the file that has already been rewritten - expect a no-op
	{in: "int1",
		remove:  false,
		wanted:  "int1",
		changed: false,
		label:   "int1-noop",
	},
	{in: "helloworld",
		remove:  false,
		wanted:  "helloworld",
		changed: false,
		label:   "unmodified",
	},
}

func TestRewrite(t *testing.T) {
	for _, test := range tests {
		buf, err := Rewrite(test.label, input[test.in], prefix, test.remove)
		if err != nil {
			t.Error(err)
		}
		// nil buf means no changes - expect test.changed be false
		if buf == nil {
			if test.changed == true {
				t.Errorf("Error in %s - buf is nil but test.changed is true", test.label)
			}
			continue
		}
		if buf != nil && test.changed == false {
			t.Errorf("Error in %s - buf is non-nil but test.changed is false", test.label)
		}
		if buf.String() != input[test.wanted] {
			t.Errorf("Error in %s: Input:\n%s\n, Got:\n%s\nWanted:\n%s\nremove: %t\n",
				test.label, input[test.in], buf.String(), input[test.wanted], test.remove)
		}
	}
}
