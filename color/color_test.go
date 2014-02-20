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
