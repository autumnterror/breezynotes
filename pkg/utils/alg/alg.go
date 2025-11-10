package alg

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func SliceCopy[T any](s []T, start, end int) []T {
	if start >= end {
		return nil
	}
	out := make([]T, end-start)
	copy(out, s[start:end])
	return out
}
