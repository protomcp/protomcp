package testutils_test

import (
	"fmt"
	"strings"
	"testing"

	"google.golang.org/protobuf/types/pluginpb"

	"protomcp.org/protomcp/pkg/generator/testutils"
)

// mockT implements testutils.T for testing assertion functions
type mockT struct {
	errors       []string
	fatalCalls   []string
	helperCalled bool
}

func (m *mockT) Helper() {
	m.helperCalled = true
}

func (m *mockT) Errorf(format string, args ...any) {
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockT) Fatalf(format string, args ...any) {
	m.fatalCalls = append(m.fatalCalls, fmt.Sprintf(format, args...))
}

func TestAssertContains(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expected  string
		wantError bool
	}{
		{
			name:      "contains substring",
			content:   "Hello, World!",
			expected:  "World",
			wantError: false,
		},
		{
			name:      "does not contain substring",
			content:   "Hello, World!",
			expected:  "Goodbye",
			wantError: true,
		},
		{
			name:      "empty content",
			content:   "",
			expected:  "anything",
			wantError: true,
		},
		{
			name:      "empty expected",
			content:   "Hello, World!",
			expected:  "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{}
			testutils.AssertContains(mt, tt.content, tt.expected)

			if !mt.helperCalled {
				t.Error("Helper() was not called")
			}

			if tt.wantError && len(mt.errors) == 0 {
				t.Error("expected error but got none")
			}
			if !tt.wantError && len(mt.errors) > 0 {
				t.Errorf("unexpected error: %v", mt.errors)
			}
		})
	}
}

func TestAssertFileCount(t *testing.T) {
	tests := []struct {
		name      string
		fileCount int
		expected  int
		wantFatal bool
	}{
		{
			name:      "correct count",
			fileCount: 3,
			expected:  3,
			wantFatal: false,
		},
		{
			name:      "incorrect count",
			fileCount: 2,
			expected:  3,
			wantFatal: true,
		},
		{
			name:      "zero files",
			fileCount: 0,
			expected:  0,
			wantFatal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{}
			response := &pluginpb.CodeGeneratorResponse{
				File: make([]*pluginpb.CodeGeneratorResponse_File, tt.fileCount),
			}

			testutils.AssertFileCount(mt, response, tt.expected)

			if !mt.helperCalled {
				t.Error("Helper() was not called")
			}

			if tt.wantFatal && len(mt.fatalCalls) == 0 {
				t.Error("expected fatal but got none")
			}
			if !tt.wantFatal && len(mt.fatalCalls) > 0 {
				t.Errorf("unexpected fatal: %v", mt.fatalCalls)
			}
		})
	}
}

func TestAssertSliceEqual(t *testing.T) {
	tests := []struct {
		name      string
		got       []string
		want      []string
		wantError bool
	}{
		{
			name:      "equal slices",
			got:       testutils.S("a", "b", "c"),
			want:      testutils.S("a", "b", "c"),
			wantError: false,
		},
		{
			name:      "different slices",
			got:       testutils.S("a", "b"),
			want:      testutils.S("a", "b", "c"),
			wantError: true,
		},
		{
			name:      "empty slices",
			got:       testutils.S[string](),
			want:      testutils.S[string](),
			wantError: false,
		},
		{
			name:      "nil vs empty",
			got:       nil,
			want:      testutils.S[string](),
			wantError: false, // core.SliceEqual treats nil and empty as equal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{}
			testutils.AssertSliceEqual(mt, tt.got, tt.want, "test slice")

			if !mt.helperCalled {
				t.Error("Helper() was not called")
			}

			if tt.wantError && len(mt.errors) == 0 {
				t.Error("expected error but got none")
			}
			if !tt.wantError && len(mt.errors) > 0 {
				t.Errorf("unexpected error: %v", mt.errors)
			}
		})
	}

	// Test format string support
	t.Run("format string", func(t *testing.T) {
		mt := &mockT{}
		testutils.AssertSliceEqual(mt, testutils.S("a"), testutils.S("b"), "slice for %s", "test")
		if len(mt.errors) == 0 {
			t.Error("expected error")
		}
		if !strings.Contains(mt.errors[0], "slice for test") {
			t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
		}
	})
}

func TestAssertSliceOfSlicesEqual(t *testing.T) {
	tests := []struct {
		name      string
		got       [][]string
		want      [][]string
		wantError bool
	}{
		{
			name:      "equal nested slices",
			got:       testutils.S(testutils.S("a", "b"), testutils.S("c", "d")),
			want:      testutils.S(testutils.S("a", "b"), testutils.S("c", "d")),
			wantError: false,
		},
		{
			name:      "different length",
			got:       testutils.S(testutils.S("a", "b")),
			want:      testutils.S(testutils.S("a", "b"), testutils.S("c", "d")),
			wantError: true,
		},
		{
			name:      "different content",
			got:       testutils.S(testutils.S("a", "b"), testutils.S("c", "d")),
			want:      testutils.S(testutils.S("a", "b"), testutils.S("e", "f")),
			wantError: true,
		},
		{
			name:      "empty slices",
			got:       testutils.S[[]string](),
			want:      testutils.S[[]string](),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &mockT{}
			testutils.AssertSliceOfSlicesEqual(mt, tt.got, tt.want, "test nested slice")

			if !mt.helperCalled {
				t.Error("Helper() was not called")
			}

			if tt.wantError && len(mt.errors) == 0 {
				t.Error("expected error but got none")
			}
			if !tt.wantError && len(mt.errors) > 0 {
				t.Errorf("unexpected error: %v", mt.errors)
			}

			// Check that we got error for each mismatch
			if tt.name == "different content" && len(mt.errors) != 1 {
				t.Errorf("expected 1 error for content mismatch, got %d", len(mt.errors))
			}
		})
	}
}

func TestAssertEqual(t *testing.T) {
	t.Run("strings", testAssertEqualStrings)
	t.Run("integers", testAssertEqualIntegers)
	t.Run("booleans", testAssertEqualBooleans)
	t.Run("format string", testAssertEqualFormatString)
}

func testAssertEqualStrings(t *testing.T) {
	mt := &mockT{}
	testutils.AssertEqual(mt, "hello", "hello", "greeting")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for equal strings, got: %v", mt.errors)
	}

	mt = &mockT{}
	testutils.AssertEqual(mt, "hello", "world", "greeting")
	if len(mt.errors) == 0 {
		t.Error("expected error for different strings")
	}
	if !strings.Contains(mt.errors[0], "greeting") {
		t.Errorf("error should mention field name 'greeting', got: %s", mt.errors[0])
	}
}

func testAssertEqualIntegers(t *testing.T) {
	mt := &mockT{}
	testutils.AssertEqual(mt, 42, 42, "answer")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for equal integers, got: %v", mt.errors)
	}

	mt = &mockT{}
	testutils.AssertEqual(mt, 42, 24, "answer")
	if len(mt.errors) == 0 {
		t.Error("expected error for different integers")
	}
}

func testAssertEqualBooleans(t *testing.T) {
	mt := &mockT{}
	testutils.AssertEqual(mt, true, true, "flag")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for equal booleans, got: %v", mt.errors)
	}

	mt = &mockT{}
	testutils.AssertEqual(mt, true, false, "flag")
	if len(mt.errors) == 0 {
		t.Error("expected error for different booleans")
	}
}

func testAssertEqualFormatString(t *testing.T) {
	mt := &mockT{}
	testutils.AssertEqual(mt, "foo", "bar", "field[%d]", 42)
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "field[42]") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}

	// Test multiple format args
	mt = &mockT{}
	testutils.AssertEqual(mt, 1, 2, "%s.%s", "obj", "prop")
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "obj.prop") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

func TestAssertTypeIs(t *testing.T) {
	t.Run("successful type assertion", testAssertTypeIsSuccess)
	t.Run("failed type assertion", testAssertTypeIsFailed)
	t.Run("with format string", testAssertTypeIsFormatString)
	t.Run("interface types", testAssertTypeIsInterface)
}

func testAssertTypeIsSuccess(t *testing.T) {
	mt := &mockT{}
	var val any = "hello"

	result, ok := testutils.AssertTypeIs[string](mt, val, "value")

	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
	if len(mt.errors) > 0 {
		t.Errorf("unexpected error: %v", mt.errors)
	}
	if !ok {
		t.Error("expected type assertion to succeed")
	}
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func testAssertTypeIsFailed(t *testing.T) {
	mt := &mockT{}
	var val any = 123

	result, ok := testutils.AssertTypeIs[string](mt, val, "value")

	if len(mt.errors) == 0 {
		t.Error("expected error but got none")
	}
	if !strings.Contains(mt.errors[0], "expected type string") {
		t.Errorf("error should mention expected type: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "got int") {
		t.Errorf("error should mention actual type: %s", mt.errors[0])
	}
	if ok {
		t.Error("expected type assertion to fail")
	}
	if result != "" {
		t.Errorf("expected empty string (zero value), got %q", result)
	}
}

func testAssertTypeIsFormatString(t *testing.T) {
	mt := &mockT{}
	var val any = 42.5

	_, ok := testutils.AssertTypeIs[string](mt, val, "item[%d].%s", 3, "value")

	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "item[3].value") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
	if ok {
		t.Error("expected type assertion to fail")
	}
}

func testAssertTypeIsInterface(t *testing.T) {
	mt := &mockT{}

	// Test with error interface
	var val any = fmt.Errorf("test error")
	err, ok := testutils.AssertTypeIs[error](mt, val, "error value")

	if len(mt.errors) > 0 {
		t.Errorf("unexpected error: %v", mt.errors)
	}
	if !ok {
		t.Error("expected type assertion to succeed")
	}
	if err == nil || err.Error() != "test error" {
		t.Errorf("expected 'test error', got %v", err)
	}

	// Test with custom interface
	type Stringer interface {
		String() string
	}
	mt = &mockT{}
	var strVal any = "hello"
	result, ok := testutils.AssertTypeIs[Stringer](mt, strVal, "stringer")

	if len(mt.errors) == 0 {
		t.Error("expected error for non-Stringer type")
	}
	if ok {
		t.Error("expected type assertion to fail")
	}
	if result != nil {
		t.Errorf("expected nil (zero value for interface), got %v", result)
	}
}

func TestMockTBehaviour(t *testing.T) {
	t.Run("Helper is called", func(t *testing.T) {
		mt := &mockT{}
		testutils.AssertContains(mt, "test", "test")
		if !mt.helperCalled {
			t.Error("Helper() should have been called")
		}
	})

	t.Run("Errorf collects errors", func(t *testing.T) {
		mt := &mockT{}
		mt.Errorf("error %d: %s", 1, "test error")
		mt.Errorf("error %d: %s", 2, "another error")

		if len(mt.errors) != 2 {
			t.Errorf("expected 2 errors, got %d", len(mt.errors))
		}
		if mt.errors[0] != "error 1: test error" {
			t.Errorf("unexpected error format: %s", mt.errors[0])
		}
	})

	t.Run("Fatalf collects fatal calls", func(t *testing.T) {
		mt := &mockT{}
		mt.Fatalf("fatal error: %s", "critical")

		if len(mt.fatalCalls) != 1 {
			t.Errorf("expected 1 fatal, got %d", len(mt.fatalCalls))
		}
		if mt.fatalCalls[0] != "fatal error: critical" {
			t.Errorf("unexpected fatal format: %s", mt.fatalCalls[0])
		}
	})
}

func TestAssertNotEqual(t *testing.T) {
	t.Run("different values", testAssertNotEqualDifferent)
	t.Run("equal values", testAssertNotEqualSame)
	t.Run("format string", testAssertNotEqualFormatString)
}

func testAssertNotEqualDifferent(t *testing.T) {
	mt := &mockT{}
	testutils.AssertNotEqual(mt, "hello", "world", "strings")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for different values, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertNotEqualSame(t *testing.T) {
	mt := &mockT{}
	testutils.AssertNotEqual(mt, 42, 42, "answer")
	if len(mt.errors) == 0 {
		t.Error("expected error for equal values")
	}
	if !strings.Contains(mt.errors[0], "answer") {
		t.Errorf("error should mention field name 'answer', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want not equal") {
		t.Errorf("error should mention 'want not equal', got: %s", mt.errors[0])
	}
}

func testAssertNotEqualFormatString(t *testing.T) {
	mt := &mockT{}
	testutils.AssertNotEqual(mt, true, true, "flag[%d]", 5)
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "flag[5]") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

func TestAssertNil(t *testing.T) {
	t.Run("nil value", testAssertNilSuccess)
	t.Run("non-nil value", testAssertNilFailed)
	t.Run("format string", testAssertNilFormatString)
}

func testAssertNilSuccess(t *testing.T) {
	mt := &mockT{}
	var ptr *string
	testutils.AssertNil(mt, ptr, "pointer")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for nil value, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertNilFailed(t *testing.T) {
	mt := &mockT{}
	value := "not nil"
	testutils.AssertNil(mt, &value, "pointer")
	if len(mt.errors) == 0 {
		t.Error("expected error for non-nil value")
	}
	if !strings.Contains(mt.errors[0], "pointer") {
		t.Errorf("error should mention field name 'pointer', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want nil") {
		t.Errorf("error should mention 'want nil', got: %s", mt.errors[0])
	}
}

func testAssertNilFormatString(t *testing.T) {
	mt := &mockT{}
	testutils.AssertNil(mt, 123, "value[%s]", "test")
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "value[test]") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

func TestAssertNotNil(t *testing.T) {
	t.Run("non-nil value", testAssertNotNilSuccess)
	t.Run("nil value", testAssertNotNilFailed)
	t.Run("format string", testAssertNotNilFormatString)
}

func testAssertNotNilSuccess(t *testing.T) {
	mt := &mockT{}
	value := "not nil"
	testutils.AssertNotNil(mt, &value, "pointer")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for non-nil value, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertNotNilFailed(t *testing.T) {
	mt := &mockT{}
	var ptr *string
	testutils.AssertNotNil(mt, ptr, "pointer")
	if len(mt.errors) == 0 {
		t.Error("expected error for nil value")
	}
	if !strings.Contains(mt.errors[0], "pointer") {
		t.Errorf("error should mention field name 'pointer', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want not nil") {
		t.Errorf("error should mention 'want not nil', got: %s", mt.errors[0])
	}
}

func testAssertNotNilFormatString(t *testing.T) {
	mt := &mockT{}
	var ptr *int
	testutils.AssertNotNil(mt, ptr, "field.%s", "value")
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "field.value") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

type isNilCoverageTest struct {
	value     any
	name      string
	wantError bool
}

func newIsNilCoverageTest(name string, value any, wantError bool) isNilCoverageTest {
	return isNilCoverageTest{
		name:      name,
		value:     value,
		wantError: wantError,
	}
}

func (tc isNilCoverageTest) test(t *testing.T) {
	mt := &mockT{}
	testutils.AssertNil(mt, tc.value, "test value")

	if tc.wantError && len(mt.errors) == 0 {
		t.Error("expected error but got none")
	}
	if !tc.wantError && len(mt.errors) > 0 {
		t.Errorf("unexpected error: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func TestIsNilCoverage(t *testing.T) {
	tests := testutils.S(
		// Fast path pointer types
		newIsNilCoverageTest("nil *int fast path", (*int)(nil), false),
		newIsNilCoverageTest("nil *bool fast path", (*bool)(nil), false),
		newIsNilCoverageTest("non-nil *int", func() *int { v := 42; return &v }(), true),
		newIsNilCoverageTest("non-nil *bool", func() *bool { v := true; return &v }(), true),
		// Reflection path types
		newIsNilCoverageTest("nil slice", ([]string)(nil), false),
		newIsNilCoverageTest("nil map", (map[string]int)(nil), false),
		newIsNilCoverageTest("nil channel", (chan int)(nil), false),
		newIsNilCoverageTest("nil function", (func())(nil), false),
		newIsNilCoverageTest("nil interface", (any)(nil), false),
		newIsNilCoverageTest("non-nil slice", []string{"test"}, true),
		newIsNilCoverageTest("non-nil map", make(map[string]int), true),
	)

	for _, tc := range tests {
		t.Run(tc.name, tc.test)
	}
}

func TestAssertTrue(t *testing.T) {
	t.Run("true condition", testAssertTrueSuccess)
	t.Run("false condition", testAssertTrueFailed)
	t.Run("format string", testAssertTrueFormatString)
}

func testAssertTrueSuccess(t *testing.T) {
	mt := &mockT{}
	testutils.AssertTrue(mt, true, "condition")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for true condition, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertTrueFailed(t *testing.T) {
	mt := &mockT{}
	testutils.AssertTrue(mt, false, "condition")
	if len(mt.errors) == 0 {
		t.Error("expected error for false condition")
	}
	if !strings.Contains(mt.errors[0], "condition") {
		t.Errorf("error should mention field name 'condition', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want true") {
		t.Errorf("error should mention 'want true', got: %s", mt.errors[0])
	}
}

func testAssertTrueFormatString(t *testing.T) {
	mt := &mockT{}
	testutils.AssertTrue(mt, false, "check[%d]", 42)
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "check[42]") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

func TestAssertFalse(t *testing.T) {
	t.Run("false condition", testAssertFalseSuccess)
	t.Run("true condition", testAssertFalseFailed)
	t.Run("format string", testAssertFalseFormatString)
}

func testAssertFalseSuccess(t *testing.T) {
	mt := &mockT{}
	testutils.AssertFalse(mt, false, "condition")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors for false condition, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertFalseFailed(t *testing.T) {
	mt := &mockT{}
	testutils.AssertFalse(mt, true, "condition")
	if len(mt.errors) == 0 {
		t.Error("expected error for true condition")
	}
	if !strings.Contains(mt.errors[0], "condition") {
		t.Errorf("error should mention field name 'condition', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want false") {
		t.Errorf("error should mention 'want false', got: %s", mt.errors[0])
	}
}

func testAssertFalseFormatString(t *testing.T) {
	mt := &mockT{}
	testutils.AssertFalse(mt, true, "flag.%s", "enabled")
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "flag.enabled") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

func TestAssertError(t *testing.T) {
	t.Run("error present", testAssertErrorSuccess)
	t.Run("no error", testAssertErrorFailed)
	t.Run("format string", testAssertErrorFormatString)
}

func testAssertErrorSuccess(t *testing.T) {
	mt := &mockT{}
	err := fmt.Errorf("test error")
	testutils.AssertError(mt, err, "operation")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors when error is present, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertErrorFailed(t *testing.T) {
	mt := &mockT{}
	testutils.AssertError(mt, nil, "operation")
	if len(mt.errors) == 0 {
		t.Error("expected error when error is nil")
	}
	if !strings.Contains(mt.errors[0], "operation") {
		t.Errorf("error should mention field name 'operation', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want error") {
		t.Errorf("error should mention 'want error', got: %s", mt.errors[0])
	}
}

func testAssertErrorFormatString(t *testing.T) {
	mt := &mockT{}
	testutils.AssertError(mt, nil, "call[%d]", 3)
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "call[3]") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}

func TestAssertNoError(t *testing.T) {
	t.Run("no error", testAssertNoErrorSuccess)
	t.Run("error present", testAssertNoErrorFailed)
	t.Run("format string", testAssertNoErrorFormatString)
}

func testAssertNoErrorSuccess(t *testing.T) {
	mt := &mockT{}
	testutils.AssertNoError(mt, nil, "operation")
	if len(mt.errors) > 0 {
		t.Errorf("expected no errors when error is nil, got: %v", mt.errors)
	}
	if !mt.helperCalled {
		t.Error("Helper() was not called")
	}
}

func testAssertNoErrorFailed(t *testing.T) {
	mt := &mockT{}
	err := fmt.Errorf("unexpected error")
	testutils.AssertNoError(mt, err, "operation")
	if len(mt.errors) == 0 {
		t.Error("expected error when error is not nil")
	}
	if !strings.Contains(mt.errors[0], "operation") {
		t.Errorf("error should mention field name 'operation', got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "unexpected error") {
		t.Errorf("error should contain the error message, got: %s", mt.errors[0])
	}
	if !strings.Contains(mt.errors[0], "want no error") {
		t.Errorf("error should mention 'want no error', got: %s", mt.errors[0])
	}
}

func testAssertNoErrorFormatString(t *testing.T) {
	mt := &mockT{}
	err := fmt.Errorf("test failure")
	testutils.AssertNoError(mt, err, "func.%s", "Init")
	if len(mt.errors) == 0 {
		t.Error("expected error")
	}
	if !strings.Contains(mt.errors[0], "func.Init") {
		t.Errorf("expected formatted name in error, got: %s", mt.errors[0])
	}
}
