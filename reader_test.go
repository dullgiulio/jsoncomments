// Copyright 2015 Giulio Iotti. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsoncomments

import (
	"bytes"
	"strings"
	"testing"
)

func TestStripComments(t *testing.T) {
	// Notice: indents below are tabs.
	input := `
{ # Some JSON
	"key": "value#",
	# A comment
	"#another": "#val\"#"
}
`
	output := `
{ 
	"key": "value#",
	
	"#another": "#val\"#"
}
`
	buf := bytes.NewBufferString("")
	r := NewSkipCommentsReader(strings.NewReader(input))
	buf.ReadFrom(r)

	out := buf.String()
	if out != output {
		t.Errorf("Output not as expected: \n%#v\n-------------\n%#v", out, output)
	}
}
