package io

import (
	"bufio"
)

// ReadLine is a helper func for bufio's ReadLine that
// read the complete line no matter how long it is.
// Note: EOL is stripped.
func ReadLine(bio *bufio.Reader) ([]byte, error) {
	line, isPrefix, err := bio.ReadLine()
	if !isPrefix {
		return line, err
	}

	// line is too long, read till eol
	buf := append([]byte(nil), line...)
	for isPrefix && err == nil {
		line, isPrefix, err = bio.ReadLine()
		buf = append(buf, line...)
	}
	return buf, err
}
