package generator

import (
	"io"
	"os"
)

// DebugOutput is the output writer for debug messages. Defaults to os.Stderr.
// Can be overridden for testing purposes to capture debug output.
//
// Example:
//
//	var buf bytes.Buffer
//	generator.DebugOutput = &buf
//	generator.Debug("test message")
//	output := buf.String()
var DebugOutput io.Writer = os.Stderr

// Debug prints debug information if PROTOMCP_DEBUG environment variable is set.
// Unlike Trace, Debug does not include source location information.
//
// The debug output format is:
//
//	[DEBUG] formatted message
//
// Environment:
//
//	PROTOMCP_DEBUG - when set (any value), enables debug output
//
// Example:
//
//	// Enable debugging
//	os.Setenv("PROTOMCP_DEBUG", "1")
//
//	// This will output: [DEBUG] Processing 5 items
//	generator.Debug("Processing %d items", itemCount)
//
// Debug is useful for general debugging output where source location is not
// needed, while Trace is better for understanding code flow and execution paths.
func Debug(format string, args ...any) {
	if _, ok := os.LookupEnv("PROTOMCP_DEBUG"); !ok {
		return
	}

	s := makeDebugString(format, args...)
	_, _ = io.WriteString(DebugOutput, s)
}

// makeDebugString constructs a debug message string without location info.
// It handles both formatted messages (with args) and plain string messages.
//
// Output format:
//
//	[DEBUG] message\n  (with message)
//	[DEBUG]\n          (without message)
func makeDebugString(format string, args ...any) string {
	var buf LazyBuffer

	buf.WriteString("[DEBUG]")

	if format == "" {
		return buf.WriteString("\n").String()
	}

	buf.WriteString(" ")
	if len(args) > 0 {
		// Format the debug message
		buf.Printf(format, args...)
	} else {
		// Raw message
		buf.WriteString(format)
	}

	return buf.WriteString("\n").String()
}
