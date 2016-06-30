package io

import (
	"bufio"
	"os"
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

// ReadLines reads a whole file into memory and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
