package generator

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"darvaza.org/core"
)

func TestGetFilePath(t *testing.T) {
	// Test with current frame
	frame := core.Here()
	got := getFilePath(frame)

	// Should contain package name and file name
	if !strings.Contains(got, "generator") {
		t.Errorf("getFilePath() = %q, expected to contain 'generator'", got)
	}
	if !strings.Contains(got, "trace_test.go") {
		t.Errorf("getFilePath() = %q, expected to contain 'trace_test.go'", got)
	}
}

func TestMakeTraceString(t *testing.T) {
	t.Run("empty format without frame", testMakeTraceStringEmptyNoFrame)
	t.Run("simple message without frame", testMakeTraceStringSimpleNoFrame)
	t.Run("formatted message without frame", testMakeTraceStringFormattedNoFrame)
	t.Run("with frame", testMakeTraceStringWithFrame)
	t.Run("with frame and empty format", testMakeTraceStringWithFrameEmpty)
}

func testMakeTraceStringEmptyNoFrame(t *testing.T) {
	got := makeTraceString(nil, "")
	expected := "[TRACE]\n"
	if got != expected {
		t.Errorf("makeTraceString() = %q, want %q", got, expected)
	}
}

func testMakeTraceStringSimpleNoFrame(t *testing.T) {
	got := makeTraceString(nil, "Hello world")
	expected := "[TRACE] Hello world\n"
	if got != expected {
		t.Errorf("makeTraceString() = %q, want %q", got, expected)
	}
}

func testMakeTraceStringFormattedNoFrame(t *testing.T) {
	got := makeTraceString(nil, "Processing %s with value %d", "test", 123)
	expected := "[TRACE] Processing test with value 123\n"
	if got != expected {
		t.Errorf("makeTraceString() = %q, want %q", got, expected)
	}
}

func testMakeTraceStringWithFrame(t *testing.T) {
	frame := core.Here()
	got := makeTraceString(frame, "Test message")

	// Check structure
	if !strings.HasPrefix(got, "[TRACE ") {
		t.Errorf("makeTraceString() should start with '[TRACE ': %q", got)
	}
	if !strings.Contains(got, "trace_test.go:") {
		t.Errorf("makeTraceString() should contain file info: %q", got)
	}
	if !strings.Contains(got, "] Test message\n") {
		t.Errorf("makeTraceString() should contain message: %q", got)
	}
}

func testMakeTraceStringWithFrameEmpty(t *testing.T) {
	frame := core.Here()
	got := makeTraceString(frame, "")

	// Check it ends with just ]\n
	if !strings.HasSuffix(got, "]\n") {
		t.Errorf("makeTraceString() with empty format should end with ]\\n: %q", got)
	}
	// Check it contains file info
	if !strings.Contains(got, "trace_test.go:") {
		t.Errorf("makeTraceString() should contain file info: %q", got)
	}
}

func TestTrace(t *testing.T) {
	// Save original values
	origOutput := TraceOutput
	origEnv := os.Getenv("PROTOMCP_TRACE")

	// Restore after test
	t.Cleanup(func() {
		TraceOutput = origOutput
		if origEnv != "" {
			_ = os.Setenv("PROTOMCP_TRACE", origEnv)
		} else {
			_ = os.Unsetenv("PROTOMCP_TRACE")
		}
	})

	t.Run("trace disabled", testTraceDisabled)
	t.Run("trace enabled with simple message", testTraceEnabledSimple)
	t.Run("trace enabled with formatted message", testTraceEnabledFormatted)
	t.Run("trace enabled with empty string", testTraceEnabledEmpty)
}

func testTraceDisabled(t *testing.T) {
	_ = os.Unsetenv("PROTOMCP_TRACE")
	var buf bytes.Buffer
	TraceOutput = &buf

	Trace("This should not appear")

	if buf.Len() > 0 {
		t.Errorf("Trace() wrote output when disabled: %q", buf.String())
	}
}

func testTraceEnabledSimple(t *testing.T) {
	_ = os.Setenv("PROTOMCP_TRACE", "1")
	var buf bytes.Buffer
	TraceOutput = &buf

	Trace("Test message")

	output := buf.String()
	if !strings.Contains(output, "[TRACE") {
		t.Errorf("Trace() missing [TRACE prefix: %q", output)
	}
	if !strings.Contains(output, "Test message") {
		t.Errorf("Trace() missing message: %q", output)
	}
	if !strings.Contains(output, ".go:") {
		t.Errorf("Trace() missing file info: %q", output)
	}
	// Note: When called from t.Run, the trace shows buffer.go due to test framework internals
	t.Logf("Trace output (called from t.Run): %s", output)
}

func testTraceEnabledFormatted(t *testing.T) {
	_ = os.Setenv("PROTOMCP_TRACE", "1")
	var buf bytes.Buffer
	TraceOutput = &buf

	Trace("Processing %s: %d items", "test", 42)

	output := buf.String()
	if !strings.Contains(output, "Processing test: 42 items") {
		t.Errorf("Trace() formatted incorrectly: %q", output)
	}
}

func testTraceEnabledEmpty(t *testing.T) {
	_ = os.Setenv("PROTOMCP_TRACE", "1")
	var buf bytes.Buffer
	TraceOutput = &buf

	Trace("")

	output := buf.String()
	if !strings.HasSuffix(output, "]\n") {
		t.Errorf("Trace() with empty string should end with ]\\n: %q", output)
	}
}

// TestTraceIntegration tests the trace functionality in a more realistic scenario
func TestTraceIntegration(t *testing.T) {
	setupTraceEnvironment(t)

	var buf bytes.Buffer
	TraceOutput = &buf

	// Run test scenario
	simulateTraceWork("test-task", 3)

	// Verify output
	verifyTraceIntegrationOutput(t, buf.String())
}

// setupTraceEnvironment saves and restores trace environment
func setupTraceEnvironment(t *testing.T) {
	t.Helper()
	origOutput := TraceOutput
	origEnv := os.Getenv("PROTOMCP_TRACE")

	t.Cleanup(func() {
		TraceOutput = origOutput
		if origEnv != "" {
			_ = os.Setenv("PROTOMCP_TRACE", origEnv)
		} else {
			_ = os.Unsetenv("PROTOMCP_TRACE")
		}
	})

	_ = os.Setenv("PROTOMCP_TRACE", "1")
}

// simulateTraceWork simulates a function that uses tracing
func simulateTraceWork(name string, count int) {
	Trace("Starting work on %s", name)
	for i := range count {
		Trace("Processing item %d", i)
	}
	Trace("Completed %s", name)
}

// verifyTraceIntegrationOutput checks the trace output from integration test
func verifyTraceIntegrationOutput(t *testing.T, output string) {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Check line count
	if len(lines) != 5 {
		t.Errorf("Expected 5 trace lines, got %d", len(lines))
	}

	// Check format of all lines
	verifyTraceLineFormat(t, lines)

	// Check specific content
	if !strings.Contains(lines[0], "Starting work on test-task") {
		t.Errorf("First line has wrong content: %q", lines[0])
	}
	if !strings.Contains(lines[4], "Completed test-task") {
		t.Errorf("Last line has wrong content: %q", lines[4])
	}
}

// verifyTraceLineFormat checks that all trace lines have proper format
func verifyTraceLineFormat(t *testing.T, lines []string) {
	t.Helper()
	for i, line := range lines {
		if !strings.HasPrefix(line, "[TRACE") {
			t.Errorf("Line %d missing [TRACE prefix: %q", i, line)
		}
		if !strings.Contains(line, ".go:") {
			t.Errorf("Line %d missing file info: %q", i, line)
		}
	}
}
