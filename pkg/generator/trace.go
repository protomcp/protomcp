// Package generator provides utilities for ProtoMCP code generation.
// It includes debugging and tracing facilities to help diagnose issues during
// development and testing.
package generator

import (
	"fmt"
	"io"
	"os"
	"path"

	"darvaza.org/core"
)

// TraceOutput is the output writer for trace messages. Defaults to os.Stderr.
// Can be overridden for testing purposes to capture trace output.
//
// Example:
//
//	var buf bytes.Buffer
//	generator.TraceOutput = &buf
//	generator.Trace("test message")
//	output := buf.String()
var TraceOutput io.Writer = os.Stderr

// Trace prints debug information if PROTOMCP_TRACE environment variable is set.
// The output includes the source location (file:line) of the caller, followed by
// the formatted message.
//
// The trace output format is:
//
//	[TRACE package/file.go:123] formatted message
//
// Environment:
//
//	PROTOMCP_TRACE - when set (any value), enables trace output
//
// Example:
//
//	// Enable tracing
//	os.Setenv("PROTOMCP_TRACE", "1")
//
//	// This will output: [TRACE example/file.go:42] Processing 5 items
//	generator.Trace("Processing %d items", itemCount)
//
// The function uses darvaza.org/core.StackFrame to capture the caller's location,
// skipping 1 frame to report the actual caller rather than this function.
func Trace(format string, args ...any) {
	if _, ok := os.LookupEnv("PROTOMCP_TRACE"); !ok {
		return
	}

	frame := core.StackFrame(1)
	s := makeTraceString(frame, format, args...)
	_, _ = io.WriteString(TraceOutput, s)
}

// getFilePath constructs a file path string from a stack frame.
// If the frame has a package name, it returns "package/filename".
// Otherwise, it returns the full file path.
func getFilePath(frame *core.Frame) string {
	var pkgName = frame.PkgName()
	if pkgName != "" {
		return fmt.Sprintf("%s/%s", pkgName, path.Base(frame.File()))
	}
	return frame.File()
}

// makeTraceString constructs a complete trace message string including location info.
// It handles both formatted messages (with args) and plain string messages.
// The frame parameter may be nil, in which case location info is omitted.
//
// Output format:
//
//	[TRACE file:line] message\n  (with location)
//	[TRACE] message\n          (without location)
//	[TRACE file:line]\n       (empty message with location)
//	[TRACE]\n                 (empty message without location)
func makeTraceString(frame *core.Frame, format string, args ...any) string {
	var buf LazyBuffer

	buf.WriteString("[TRACE")
	if frame != nil {
		// Get caller information
		var filePath = getFilePath(frame)
		var line = frame.Line()

		buf.Printf(" %s:%d", filePath, line)
	}

	if format == "" {
		// empty
		return buf.WriteString("]\n").String()
	}

	buf.WriteString("] ")
	if len(args) > 0 {
		// Format the trace message
		buf.Printf(format, args...)
	} else {
		// Raw message
		buf.WriteString(format)
	}

	return buf.WriteString("\n").String()
}
