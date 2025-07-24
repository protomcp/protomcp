package testutils

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/pluginpb"

	"darvaza.org/core"
)

// AssertContains checks if the content contains the expected string
func AssertContains(t T, content, expected string) {
	t.Helper()
	if !strings.Contains(content, expected) {
		t.Errorf("generated content missing %q", expected)
	}
}

// AssertFileCount checks if the expected number of files were generated
func AssertFileCount(t T, response *pluginpb.CodeGeneratorResponse, expected int) {
	t.Helper()
	if len(response.File) != expected {
		t.Fatalf("expected %d generated file(s), got %d", expected, len(response.File))
	}
}

// AssertSliceEqual checks if two string slices are equal.
// The name parameter can be a simple string or a format string with args.
func AssertSliceEqual(t T, got, want []string, name string, args ...any) {
	t.Helper()
	if !core.SliceEqual(got, want) {
		label := name
		if len(args) > 0 {
			label = fmt.Sprintf(name, args...)
		}
		t.Errorf("%s = %v, want %v", label, got, want)
	}
}

// AssertSliceOfSlicesEqual checks if two [][]string are equal.
// The name parameter can be a simple string or a format string with args.
func AssertSliceOfSlicesEqual(t T, got, want [][]string, name string, args ...any) {
	t.Helper()
	label := name
	if len(args) > 0 {
		label = fmt.Sprintf(name, args...)
	}
	if len(got) != len(want) {
		t.Errorf("%s: got %d groups, want %d groups", label, len(got), len(want))
		return
	}
	for i := range got {
		if !core.SliceEqual(got[i], want[i]) {
			t.Errorf("%s[%d] = %v, want %v", label, i, got[i], want[i])
		}
	}
}

// AssertEqual checks if two values are equal.
// The name parameter can be a simple string or a format string with args.
func AssertEqual[V comparable](t T, got, want V, name string, args ...any) {
	t.Helper()
	if got != want {
		label := name
		if len(args) > 0 {
			label = fmt.Sprintf(name, args...)
		}
		t.Errorf("%s = %v, want %v", label, got, want)
	}
}

// AssertTypeIs checks if a value is of the expected type and returns the typed value with success indicator.
// It fails the test if the type assertion fails but still returns the boolean to avoid silent errors.
// The name parameter can be a simple string or a format string with args.
//
// Example usage:
//
//	var val any = "hello"
//	str, ok := AssertTypeIs[string](t, val, "greeting")
//	if ok {
//		// str is now typed as string and can be used directly
//	}
//
//	// With format string
//	result, ok := AssertTypeIs[*MyType](t, obj, "response[%d]", index)
func AssertTypeIs[V any](t T, value any, name string, args ...any) (V, bool) {
	t.Helper()
	result, ok := value.(V)
	if !ok {
		label := name
		if len(args) > 0 {
			label = fmt.Sprintf(name, args...)
		}
		var zero V
		t.Errorf("%s: expected type %T, got %T", label, zero, value)
		return zero, false
	}
	return result, true
}
