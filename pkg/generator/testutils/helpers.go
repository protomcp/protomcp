package testutils

// S is a helper function for creating test slices in a more concise way.
// It takes variadic arguments and returns a slice of the same type.
// This is particularly useful in table-driven tests where many slice literals are used.
// The type is usually inferred from the context, making the code very clean.
//
// Example usage:
//
//	// Type inferred from field type
//	tests := []struct {
//	    input []string
//	}{
//	    {input: S("a", "b", "c")},  // Type inferred as []string
//	}
//
//	// Type inferred from function parameter
//	AssertSliceEqual(t, S("x", "y"), expected)  // Type inferred as []string
//
//	// Explicit type when needed
//	empty := S[int]()  // []int{}
func S[T any](v ...T) []T {
	if len(v) == 0 {
		return []T{}
	}
	return v
}
