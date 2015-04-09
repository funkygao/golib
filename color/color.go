package color

import (
	"bytes"
	"fmt"
)

// Convert a string to string with color escape info that can output to console.
func Colorize(colors []string, format string, a ...interface{}) (s string) {
	buf := new(bytes.Buffer)
	for _, color := range colors {
		buf.WriteString(color_table[color])
	}

	buf.WriteString(fmt.Sprintf(format, a...))
	buf.WriteString(color_table[Reset])
	s = buf.String()
	return
}

func Blue(format string, a ...interface{}) string {
	return colored(FgBlue, format, a...)
}

func Red(format string, a ...interface{}) string {
	return colored(FgRed, format, a...)
}

func Green(format string, a ...interface{}) string {
	return colored(FgGreen, format, a...)
}

func Yellow(format string, a ...interface{}) string {
	return colored(FgYellow, format, a...)
}

func Magenta(format string, a ...interface{}) string {
	return colored(FgMagenta, format, a...)
}

func Cyan(format string, a ...interface{}) string {
	return colored(FgCyan, format, a...)
}

func colored(color string, format string, a ...interface{}) string {
	return Colorize([]string{color}, fmt.Sprintf(format, a...))
}
