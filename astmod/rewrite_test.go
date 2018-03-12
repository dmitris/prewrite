// Copyright 2015, Yahoo Inc. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
//
// Author: Dmitry Savintsev <dsavints@yahoo-inc.com>

package astmod

import (
	"io/ioutil"
	"testing"
)

const (
	from = `go.corp.company.com`
	to   = `go.newcompany.com`
)

// tests table
var tests = []struct {
	file    string
	from    string
	to      string
	want    string
	changed bool
	label   string
}{
	{
		file:    "testdata/ext1.go",
		from:    "go.corp.example.com",
		to:      "go.newcompany.com",
		want:    "testdata/int1.go",
		changed: true,
		label:   "ext1",
	},
	{
		file:    "testdata/int1.go",
		from:    "go.corp.example.com",
		to:      "go.newcompany.com",
		want:    "testdata/int1.go",
		changed: false,
		label:   "int1",
	},
	{
		file:    "testdata/helloworld.go",
		from:    "go.corp.example.com",
		to:      "go.newcompany.com",
		want:    "testdata/helloworld.go",
		changed: false,
		label:   "unmodified",
	},
}

func TestRewrite(t *testing.T) {
	for _, tt := range tests {
		in, err := ioutil.ReadFile(tt.file)
		if err != nil {
			t.Fatal(err)
		}
		wantBytes, err := ioutil.ReadFile(tt.want)
		if err != nil {
			t.Fatal(err)
		}
		want := string(wantBytes)
		buf, err := Rewrite(tt.label, in, from, to)
		if err != nil {
			t.Error(err)
			continue
		}
		if buf == nil {
			if tt.changed == false {
				continue // OK - no changes were made, as expected
			}
			t.Errorf("test '%s': changes expected but none made (buf=nil)", tt.label)
			continue
		}
		got := buf.String()
		if got != want {
			t.Errorf("%s case: got\n%s\nwant contents of %s:\n%s", tt.label, got, tt.want, want)
		}
	}
}
