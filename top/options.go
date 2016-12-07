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

func WithSumFooter() option {
	return func(t *top) {
		t.sumFooter = true
	}
}

func WithMaxFooter() option {
	return func(t *top) {
		t.maxFooter = true
	}
}
