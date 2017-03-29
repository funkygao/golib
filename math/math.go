package math

func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}


func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}
