package limits

// SetLimits is a no-op on Windows since it's not required there.
func SetLimits() error {
	return nil
}
