package color

import (
	"strings"
	"testing"

	"github.com/funkygao/assert"
)

func TestColorize(t *testing.T) {
	assert.Equal(t, "\x1b[30m\x1b[41mhello\x1b[0m",
		Colorize([]string{"FgBlack", "BgRed"}, "hello"))
}

func TestRed(t *testing.T) {
	assert.Equal(t, "\x1b[31mhello\x1b[0m", Red("hello"))
}

func TestBlueWithFormat(t *testing.T) {
	assert.Equal(t, "\x1b[34mhello world\x1b[0m", Blue("hello %s", "world"))
}

func TestColorTable(t *testing.T) {
	// fg
	for c := range color_table {
		if strings.HasPrefix(c, "Fg") {
			t.Logf("%15s %s", c, colored(c, c))
		}
	}

	t.Log()

	// bg
	for c := range color_table {
		if strings.HasPrefix(c, "Bg") {
			t.Logf("%15s %s", c, colored(c, c))
		}
	}

	t.Log()

	// misc
	for c := range color_table {
		if !(strings.HasPrefix(c, "Fg") || strings.HasPrefix(c, "Bg")) {
			t.Logf("%15s %s", c, colored(c, c))
		}
	}

	t.Log()
	t.Log(Colorize([]string{FgCyan, Underscore, Blink, Bright, BgYellow}, "hello"))
}
