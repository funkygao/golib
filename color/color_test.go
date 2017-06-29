package color

import (
	"github.com/funkygao/assert"
	"testing"
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
	for c := range color_table {
		t.Logf("%15s %s", c, colored(c, c))
	}
}
