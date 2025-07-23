package testutils

// T is a minimal interface for testing assertions.
// It includes only the methods our assertion functions actually use.
type T interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}
