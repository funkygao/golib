package top

import (
	"github.com/nsf/termbox-go"
)

type option func(*top)

func WithHeaderColor(color termbox.Attribute) option {
	return func(t *top) {
		t.headerColor = color
	}
}
