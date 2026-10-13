package testhelper

// testingT is the small testing interface used by this exercise.
type testingT interface {
	Helper()
	Errorf(format string, args ...any)
}
