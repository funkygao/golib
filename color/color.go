package color

import (
	"bytes"
)

// Convert a string to string with color escape info that can output to console
func Colorize(colors []string, str string) string {
	r := new(bytes.Buffer)
	for _, color := range colors {
		r.WriteString(color_table[color])
	}

	r.WriteString(str)
	r.WriteString(color_table[colorReset])
	return r.String()
}
