package testutils

import (
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/protobuf/types/pluginpb"

	"darvaza.org/core"
)

// formatLabel formats a label with optional args, or returns the label as-is
func formatLabel(name string, args []any) string {
	if len(args) > 0 {
		return fmt.Sprintf(name, args...)
	}
	return name
}

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
		label := formatLabel(name, args)
		t.Errorf("%s = %v, want %v", label, got, want)
	}
}

// AssertSliceOfSlicesEqual checks if two [][]string are equal.
// The name parameter can be a simple string or a format string with args.
func AssertSliceOfSlicesEqual(t T, got, want [][]string, name string, args ...any) {
	t.Helper()
	label := formatLabel(name, args)
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
		label := formatLabel(name, args)
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
		label := formatLabel(name, args)
		var zero V
		t.Errorf("%s: expected type %T, got %T", label, zero, value)
		return zero, false
	}
	return result, true
}

// AssertNotEqual checks if two values are not equal.
// The name parameter can be a simple string or a format string with args.
func AssertNotEqual[V comparable](t T, got, want V, name string, args ...any) {
	t.Helper()
	if got == want {
		label := formatLabel(name, args)
		t.Errorf("%s = %v, want not equal to %v", label, got, want)
	}
}

// AssertNil checks if a value is nil.
// The name parameter can be a simple string or a format string with args.
func AssertNil(t T, value any, name string, args ...any) {
	t.Helper()
	if !isNil(value) {
		label := formatLabel(name, args)
		t.Errorf("%s = %v, want nil", label, value)
	}
}

// AssertNotNil checks if a value is not nil.
// The name parameter can be a simple string or a format string with args.
func AssertNotNil(t T, value any, name string, args ...any) {
	t.Helper()
	if isNil(value) {
		label := formatLabel(name, args)
		t.Errorf("%s = nil, want not nil", label)
	}
}

// isNil checks if a value is nil, handling interface cases correctly
func isNil(value any) bool {
	if value == nil {
		return true
	}

	// Fast path for common pointer types - avoid reflection for these
	switch v := value.(type) {
	case *string:
		return v == nil
	case *int:
		return v == nil
	case *bool:
		return v == nil
	}

	// Reflection fallback for other types
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return true
	}
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// AssertTrue checks if a condition is true.
// The name parameter can be a simple string or a format string with args.
func AssertTrue(t T, condition bool, name string, args ...any) {
	t.Helper()
	AssertEqual(t, condition, true, name, args...)
}

// AssertFalse checks if a condition is false.
// The name parameter can be a simple string or a format string with args.
func AssertFalse(t T, condition bool, name string, args ...any) {
	t.Helper()
	AssertEqual(t, condition, false, name, args...)
}

// AssertError checks if an error is not nil.
// The name parameter can be a simple string or a format string with args.
func AssertError(t T, err error, name string, args ...any) {
	t.Helper()
	if err == nil {
		label := formatLabel(name, args)
		t.Errorf("%s = nil, want error", label)
	}
}

// AssertNoError checks if an error is nil.
// The name parameter can be a simple string or a format string with args.
func AssertNoError(t T, err error, name string, args ...any) {
	t.Helper()
	if err != nil {
		label := formatLabel(name, args)
		t.Errorf("%s = %v, want no error", label, err)
	}
}
