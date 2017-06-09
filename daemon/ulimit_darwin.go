// +build darwin

package daemon

// EnsureServerUlimit ensures OS settings satisfy a damon requirement.
// If not satisfied, will panic.
func EnsureServerUlimit() {
	mustMaxOpenFile(10000)
}
