package testutils_test

import (
	"testing"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

func TestS(t *testing.T) {
	t.Run("integers", testSIntegers)
	t.Run("strings", testSStrings)
	t.Run("empty", testSEmpty)
	t.Run("structs", testSStructs)
}

func testSIntegers(t *testing.T) {
	got := testutils.S(1, 2, 3)
	want := []int{1, 2, 3}

	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("element %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func testSStrings(t *testing.T) {
	got := testutils.S("a", "b", "c")
	want := []string{"a", "b", "c"}

	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("element %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func testSEmpty(t *testing.T) {
	// Empty with type inference from usage
	var intSlice []int = testutils.S[int]()
	if len(intSlice) != 0 {
		t.Errorf("expected empty slice, got %v", intSlice)
	}

	// Empty with explicit type
	stringSlice := testutils.S[string]()
	if len(stringSlice) != 0 {
		t.Errorf("expected empty slice, got %v", stringSlice)
	}
}

func testSStructs(t *testing.T) {
	type testStruct struct {
		name  string
		value int
	}

	got := testutils.S(
		testStruct{name: "first", value: 1},
		testStruct{name: "second", value: 2},
	)

	if len(got) != 2 {
		t.Errorf("expected 2 elements, got %d", len(got))
	}
	if got[0].name != "first" || got[0].value != 1 {
		t.Errorf("first element incorrect: %+v", got[0])
	}
	if got[1].name != "second" || got[1].value != 2 {
		t.Errorf("second element incorrect: %+v", got[1])
	}
}
