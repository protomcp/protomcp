package testutils

import (
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

// AssertSliceEqual checks if two string slices are equal
func AssertSliceEqual(t T, got, want []string, name string) {
	t.Helper()
	if !core.SliceEqual(got, want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

// AssertSliceOfSlicesEqual checks if two [][]string are equal
func AssertSliceOfSlicesEqual(t T, got, want [][]string, name string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: got %d groups, want %d groups", name, len(got), len(want))
		return
	}
	for i := range got {
		if !core.SliceEqual(got[i], want[i]) {
			t.Errorf("%s[%d] = %v, want %v", name, i, got[i], want[i])
		}
	}
}

// AssertEqual checks if two values are equal
func AssertEqual[V comparable](t T, got, want V, name string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
