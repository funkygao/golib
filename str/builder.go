package str

import (
	"bytes"
)

type StringBuilder struct {
	*bytes.Buffer
}

func NewStringBuilder() *StringBuilder {
	this := new(StringBuilder)
	this.Buffer = new(bytes.Buffer)
	return this
}
