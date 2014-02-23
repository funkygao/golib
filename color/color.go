package color

// Convert a string to string with color escape info that can output to console
func Colorize(colors []string, str string) string {
	buf.Reset()
	for _, color := range colors {
		buf.WriteString(color_table[color])
	}

	buf.WriteString(str)
	buf.WriteString(color_table[colorReset])
	return buf.String()
}

func Blue(str string) string {
	return colorStr("FgBlue", str)
}

func Red(str string) string {
	return colorStr("FgRed", str)
}

func Green(str string) string {
	return colorStr("FgGreen", str)
}

func Yellow(str string) string {
	return colorStr("FgYellow", str)
}

func Magenta(str string) string {
	return colorStr("FgMagenta", str)
}

func Cyan(str string) string {
	return colorStr("FgCyan", str)
}

func colorStr(color string, str string) string {
	return Colorize([]string{color}, str)
}
