// Copyright 2015 Giulio Iotti. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsoncomments

import (
	"bufio"
	"bytes"
	"io"
)

func isQuote(c rune) bool {
	if c == '"' || c == '\'' {
		return true
	}
	return false
}

// Reader implements a reader that discards comments.
//
// Comments are all contents after a # in non-quoted string context.
type Reader struct {
	scanner *bufio.Scanner
	buf     *bytes.Buffer
	eof     bool
	ooc     bool // Out-of-Context: cannot remove comments here.
}

// NewReader creates a new Reader from underlying reader r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		buf:     bytes.NewBufferString(""),
		scanner: bufio.NewScanner(r),
	}
}

// findRune finds rune r in a non-quoted context in s.
func (r *Reader) findRune(s string, ru rune) int {
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
func (r *Reader) load(n int) error {
	var err error
	for {
		if !r.scanner.Scan() {
			if err := r.scanner.Err(); err != nil {
				return err
			}
			return io.EOF
		}
		line := r.scanner.Text()
		p := r.findRune(line, '#')
		if p >= 0 {
			line = line[:p]
		}
		r.buf.WriteString(line + "\n")
		if r.buf.Len() >= n {
			return err
		}
	}
}

// Read implements a io.Reader interface.
func (r *Reader) Read(p []byte) (n int, err error) {
	if r.buf.Len() < len(p) {
		if err := r.load(len(p) - r.buf.Len()); err != nil {
			n, _ = r.buf.Read(p)
			return n, err
		}
	}
	return r.buf.Read(p)
}
