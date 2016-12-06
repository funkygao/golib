package top

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type top struct {
	header string
	rowsCh chan []string
	rowFmt string

	headerColor termbox.Attribute

	w, h  int
	close chan struct{}
}

func New(header, rowFmt string, options ...option) *top {
	t := &top{
		header:      header,
		rowFmt:      rowFmt,
		headerColor: termbox.ColorBlue,
		close:       make(chan struct{}),
		rowsCh:      make(chan []string),
	}
	for _, opt := range options {
		opt(t)
	}
	return t
}

func (t *top) Start() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()

	t.w, t.h = termbox.Size()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	go func() {
		for {
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventMouse:
				// not handled

			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					close(t.close)
				}

				switch ev.Ch {
				case 'q':
					close(t.close)
				}
			}
		}
	}()

	for {
		select {
		case rows := <-t.rowsCh:
			t.render(rows)

		case <-t.close:
			return nil
		}
	}
}

func (t *top) Refresh(rows []string) {
	t.rowsCh <- rows
}

func (t *top) Stop() {
	close(t.close)
}

func (t *top) drawRow(row string, y int, fg, bg termbox.Attribute) {
	tuples := strings.Split(row, "|")
	a := make([]interface{}, 0, len(tuples))
	for _, elem := range tuples {
		a = append(a, elem)
	}
	row = fmt.Sprintf(t.rowFmt, a...)
	x := 0
	for _, ch := range row {
		termbox.SetCell(x, y, ch, fg, bg)

		// wide string must be considered
		w := runewidth.RuneWidth(ch)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(ch)) {
			w = 1
		}
		x += w
	}

}

func (t *top) render(rows []string) {
	t.drawRow(t.header, 0, termbox.ColorDefault, t.headerColor)

	rowN := t.h - 1 // header deducted
	if rowN > len(rows) {
		rowN = len(rows)
	}
	for i := 0; i < rowN; i++ {
		t.drawRow(rows[i], i+1, termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.Flush()
}
