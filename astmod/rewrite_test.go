// Copyright 2015, Yahoo Inc. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
//
// Author: Dmitry Savintsev <dsavints@yahoo-inc.com>

package astmod

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const prefix = `go.corp.example.com/x`

// tests table
var tests = []struct {
	in      string
	from    string
	to      string
	wanted  string
	changed bool
	label   string
}{
	{in: "ext1.go",
		from:    "github.com",
		to:      "gitlab.com",
		changed: true,
		label:   "ext1",
	},
	{in: "helloworld.go",
		from:    "github.com",
		to:      "gitlab.com",
		changed: false,
		label:   "helloworld",
	},
	{in: "int1.go",
		from:    "go.corp.example.com",
		to:      "go.brandnewcorp.com",
		changed: true,
		label:   "int1",
	},
}

func TestRewrite(t *testing.T) {
	for _, tt := range tests {
		fname := filepath.Join("testdata", tt.in)
		inp, err := ioutil.ReadFile(fname)
		if err != nil {
			t.Fatalf("unable to read input file %s: %v", fname, err)
		}
		buf, err := Rewrite(tt.label, inp, tt.from, tt.to)
		if err != nil {
			t.Error(err)
		}
		// nil buf means no changes - expect test.changed be false
		if buf == nil {
			if tt.changed == true {
				t.Errorf("Error in %s - buf is nil but test.changed is true", tt.label)
			}
			continue
		}
		if buf != nil && tt.changed == false {
			t.Errorf("Error in %s - buf is non-nil but test.changed is false", tt.label)
		}
		b, err := ioutil.ReadFile(fname + ".golden")
		if err != nil {
			t.Fatalf("unable to read golden file for %s: %v", fname, err)
		}
		want := string(b)
		if buf.String() != want {
			t.Errorf("Error in %s: Input:\n%s\n, Got:\n%s\nWanted:\n%s",
				tt.label, string(inp), buf.String(), want)
		}
	}
}
