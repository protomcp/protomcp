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
			got:       []string{"a", "b", "c"},
			want:      []string{"a", "b", "c"},
			wantError: false,
		},
		{
			name:      "different slices",
			got:       []string{"a", "b"},
			want:      []string{"a", "b", "c"},
			wantError: true,
		},
		{
			name:      "empty slices",
			got:       []string{},
			want:      []string{},
			wantError: false,
		},
		{
			name:      "nil vs empty",
			got:       nil,
			want:      []string{},
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
			got:       [][]string{{"a", "b"}, {"c", "d"}},
			want:      [][]string{{"a", "b"}, {"c", "d"}},
			wantError: false,
		},
		{
			name:      "different length",
			got:       [][]string{{"a", "b"}},
			want:      [][]string{{"a", "b"}, {"c", "d"}},
			wantError: true,
		},
		{
			name:      "different content",
			got:       [][]string{{"a", "b"}, {"c", "d"}},
			want:      [][]string{{"a", "b"}, {"e", "f"}},
			wantError: true,
		},
		{
			name:      "empty slices",
			got:       [][]string{},
			want:      [][]string{},
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
	// Test with strings
	t.Run("strings", func(t *testing.T) {
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
	})

	// Test with integers
	t.Run("integers", func(t *testing.T) {
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
	})

	// Test with booleans
	t.Run("booleans", func(t *testing.T) {
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
	})
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
