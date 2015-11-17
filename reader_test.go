// Copyright 2015 Giulio Iotti. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsoncomments

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// Notice: indents below are tabs.
var input string = `
{ # Some JSON
	"key": "value#",
	# A comment
	"#another": "#val\"#"
}
`
var output string = `
{ 
	"key": "value#",
	
	"#another": "#val\"#"
}
`

func TestStripComments(t *testing.T) {
	buf := bytes.NewBufferString("")
	r := NewReader(strings.NewReader(input))
	buf.ReadFrom(r)

	out := buf.String()
	if out != output {
		t.Errorf("Output not as expected: \n%#v\n-------------\n%#v", out, output)
	}
}

func TestSmallReader(t *testing.T) {
	var b [4]byte

	buf := bytes.NewBufferString("")
	r := NewReader(strings.NewReader(input))
	for {
		n, err := r.Read(b[:])
		if n > 0 {
			buf.Write(b[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Unexpected error: %s", err)
			break
		}
	}

	out := buf.String()
	if out != output {
		t.Errorf("Output not as expected: \n%#v\n-------------\n%#v", out, output)
	}
}
