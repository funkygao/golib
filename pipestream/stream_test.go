package pipestream

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestAll(t *testing.T) {
	s := New("/bin/ls")
	err := s.Open()
	assert.Equal(t, nil, err)
	r := s.Reader()
	line, _, err := r.ReadLine()
	s.Close()
	assert.Equal(t, nil, err)
	assert.Equal(t, "stream.go", string(line)) // current dir fist 'ls' output line is stream.go
}
