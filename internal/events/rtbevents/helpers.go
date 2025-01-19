package rtbevents

//go:inline
func positiveNumber(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

//go:inline
func b2u(b bool) uint {
	if b {
		return 1
	}
	return 0
}
