package generator

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
	// Save original values
	origOutput := DebugOutput
	origEnv := os.Getenv("PROTOMCP_DEBUG")

	// Restore after test
	t.Cleanup(func() {
		DebugOutput = origOutput
		if origEnv != "" {
			_ = os.Setenv("PROTOMCP_DEBUG", origEnv)
		} else {
			_ = os.Unsetenv("PROTOMCP_DEBUG")
		}
	})

	t.Run("debug disabled", testDebugDisabled)
	t.Run("debug enabled with simple message", testDebugEnabledSimple)
	t.Run("debug enabled with formatted message", testDebugEnabledFormatted)
	t.Run("debug enabled with empty string", testDebugEnabledEmpty)
}

func testDebugDisabled(t *testing.T) {
	_ = os.Unsetenv("PROTOMCP_DEBUG")
	var buf bytes.Buffer
	DebugOutput = &buf

	Debug("This should not appear")

	if buf.Len() > 0 {
		t.Errorf("Debug() wrote output when disabled: %q", buf.String())
	}
}

func testDebugEnabledSimple(t *testing.T) {
	_ = os.Setenv("PROTOMCP_DEBUG", "1")
	var buf bytes.Buffer
	DebugOutput = &buf

	Debug("Test message")

	output := buf.String()
	expected := "[DEBUG] Test message\n"
	if output != expected {
		t.Errorf("Debug() = %q, want %q", output, expected)
	}
}

func testDebugEnabledFormatted(t *testing.T) {
	_ = os.Setenv("PROTOMCP_DEBUG", "1")
	var buf bytes.Buffer
	DebugOutput = &buf

	Debug("Processing %s: %d items", "test", 42)

	output := buf.String()
	expected := "[DEBUG] Processing test: 42 items\n"
	if output != expected {
		t.Errorf("Debug() = %q, want %q", output, expected)
	}
}

func testDebugEnabledEmpty(t *testing.T) {
	_ = os.Setenv("PROTOMCP_DEBUG", "1")
	var buf bytes.Buffer
	DebugOutput = &buf

	Debug("")

	output := buf.String()
	expected := "[DEBUG]\n"
	if output != expected {
		t.Errorf("Debug() = %q, want %q", output, expected)
	}
}

func TestMakeDebugString(t *testing.T) {
	t.Run("empty format", func(t *testing.T) {
		got := makeDebugString("")
		expected := "[DEBUG]\n"
		if got != expected {
			t.Errorf("makeDebugString() = %q, want %q", got, expected)
		}
	})

	t.Run("simple message", func(t *testing.T) {
		got := makeDebugString("Hello world")
		expected := "[DEBUG] Hello world\n"
		if got != expected {
			t.Errorf("makeDebugString() = %q, want %q", got, expected)
		}
	})

	t.Run("formatted message", func(t *testing.T) {
		got := makeDebugString("Processing %s with value %d", "test", 123)
		expected := "[DEBUG] Processing test with value 123\n"
		if got != expected {
			t.Errorf("makeDebugString() = %q, want %q", got, expected)
		}
	})
}

// saveAndRestoreEnv saves environment variables and restores them on cleanup
func saveAndRestoreEnv(t *testing.T, variables ...string) {
	origVals := make(map[string]string)
	for _, v := range variables {
		origVals[v] = os.Getenv(v)
	}

	t.Cleanup(func() {
		for v, val := range origVals {
			if val != "" {
				_ = os.Setenv(v, val)
			} else {
				_ = os.Unsetenv(v)
			}
		}
	})
}

// TestDebugVsTrace compares the output of Debug and Trace
func TestDebugVsTrace(t *testing.T) {
	// Save original values
	origTraceOutput := TraceOutput
	origDebugOutput := DebugOutput

	t.Cleanup(func() {
		TraceOutput = origTraceOutput
		DebugOutput = origDebugOutput
	})

	saveAndRestoreEnv(t, "PROTOMCP_TRACE", "PROTOMCP_DEBUG")

	_ = os.Setenv("PROTOMCP_TRACE", "1")
	_ = os.Setenv("PROTOMCP_DEBUG", "1")

	var traceBuf bytes.Buffer
	var debugBuf bytes.Buffer
	TraceOutput = &traceBuf
	DebugOutput = &debugBuf

	// Same message to both
	msg := "Processing task %s"
	arg := "example"
	Trace(msg, arg)
	Debug(msg, arg)

	traceOut := traceBuf.String()
	debugOut := debugBuf.String()

	// Verify trace output
	verifyTraceOutput(t, traceOut)

	// Verify debug output
	expected := "[DEBUG] Processing task example\n"
	if debugOut != expected {
		t.Errorf("Debug output = %q, want %q", debugOut, expected)
	}
}

// verifyTraceOutput checks that trace output has expected format
func verifyTraceOutput(t *testing.T, output string) {
	t.Helper()

	if !strings.Contains(output, "[TRACE ") {
		t.Error("Trace output missing [TRACE prefix with location")
	}
	if !strings.Contains(output, ".go:") {
		t.Error("Trace output missing file location")
	}
	if !strings.Contains(output, "Processing task example") {
		t.Error("Trace output missing message content")
	}
}

// TestDebugIntegration tests the debug functionality in a more realistic scenario
func TestDebugIntegration(t *testing.T) {
	// Setup test environment
	origOutput := DebugOutput
	t.Cleanup(func() {
		DebugOutput = origOutput
	})
	saveAndRestoreEnv(t, "PROTOMCP_DEBUG")

	_ = os.Setenv("PROTOMCP_DEBUG", "1")
	var buf bytes.Buffer
	DebugOutput = &buf

	// Run test scenario
	simulateWork("test-task", 3)

	// Verify output
	verifyDebugIntegrationOutput(t, buf.String())
}

// simulateWork simulates a function that uses debugging
func simulateWork(name string, count int) {
	Debug("Starting work on %s", name)
	for i := range count {
		Debug("Processing item %d", i)
	}
	Debug("Completed %s", name)
}

// verifyDebugIntegrationOutput checks the debug output from integration test
func verifyDebugIntegrationOutput(t *testing.T, output string) {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Check line count
	if len(lines) != 5 {
		t.Fatalf("Expected 5 debug lines, got %d", len(lines))
	}

	// Check format of all lines
	for i, line := range lines {
		verifyDebugLineFormat(t, i, line)
	}

	// Check specific content
	expectedFirst := "[DEBUG] Starting work on test-task"
	expectedLast := "[DEBUG] Completed test-task"

	if lines[0] != expectedFirst {
		t.Errorf("First line = %q, want %q", lines[0], expectedFirst)
	}
	if lines[4] != expectedLast {
		t.Errorf("Last line = %q, want %q", lines[4], expectedLast)
	}
}

// verifyDebugLineFormat checks that a debug line has correct format
func verifyDebugLineFormat(t *testing.T, lineNum int, line string) {
	t.Helper()

	if !strings.HasPrefix(line, "[DEBUG] ") {
		t.Errorf("Line %d missing [DEBUG] prefix: %q", lineNum, line)
	}
	if strings.Contains(line, ".go:") {
		t.Errorf("Line %d should not contain file info: %q", lineNum, line)
	}
}
