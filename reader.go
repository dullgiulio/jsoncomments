// Copyright 2015 Giulio Iotti. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsoncomments

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

func isQuote(c rune) bool {
	if c == '"' || c == '\'' {
		return true
	}
	return false
}

// SkipCommentsReader implements a reader that discards comments.
//
// Comments are all contents after a # in non-quoted string context.
type SkipCommentsReader struct {
	scanner *bufio.Scanner
	buf     *bytes.Buffer
	eof     bool
	ooc     bool // Out-of-Context: cannot remove comments here.
}

// NewSkipCommentsReader creates a new SkipCommentsReader from underlying reader r.
func NewSkipCommentsReader(r io.Reader) *SkipCommentsReader {
	return &SkipCommentsReader{
		buf:     bytes.NewBufferString(""),
		scanner: bufio.NewScanner(r),
	}
}

// findRune finds rune r in a non-quoted context in s.
func (r *SkipCommentsReader) findRune(s string, ru rune) int {
	var esc bool // Next rune has been escaped
	p := -1
	for i, b := range s {
		// Skip the rune after an escape
		if esc {
			esc = false
			continue
		}
		// \ is used to escape the next rune in ooc mode
		if b == '\\' && r.ooc {
			esc = true
			continue
		}
		if isQuote(b) {
			r.ooc = !r.ooc
			continue
		}
		// Found a valid comment marker rune
		if b == ru && !r.ooc {
			p = i
			break
		}
	}
	return p
}

// load loads into an internal buffer to satisfy the next read.
// Comments are stripped before the resulting data is written to the buffer.
func (r *SkipCommentsReader) load(n int) error {
	for {
		if !r.scanner.Scan() {
			return r.scanner.Err()
		}
		line := r.scanner.Text()
		p := r.findRune(line, '#')
		if p >= 0 {
			line = line[:p]
		}
		r.buf.WriteString(line + "\n")
		if r.buf.Len() > n {
			return nil
		}
	}
}

// Read implements a io.Reader interface.
func (r *SkipCommentsReader) Read(p []byte) (n int, err error) {
	if r.eof {
		if r.buf.Len() > 0 {
			return r.buf.Read(p)
		}
		return 0, io.EOF
	}
	if r.buf.Len() < len(p) {
		r.load(len(p) - r.buf.Len())
	}
	return r.buf.Read(p)
}
