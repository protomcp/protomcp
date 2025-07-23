package generator

import (
	"fmt"
	"strings"
)

// LazyBuffer is a wrapper around strings.Builder that provides convenience methods
// for building output strings without requiring error handling. Unlike strings.Builder,
// LazyBuffer methods don't return errors, making it "lazy" - you can chain calls
// without `_, _ =` error handling. It also handles nil receivers gracefully.
type LazyBuffer strings.Builder

// WriteString appends one or more strings to the buffer, ignoring empty strings.
// Safe to call on nil receiver. Returns the buffer for method chaining.
func (b *LazyBuffer) WriteString(ss ...string) *LazyBuffer {
	if b != nil {
		sb := (*strings.Builder)(b)
		for _, s := range ss {
			if s != "" {
				_, _ = sb.WriteString(s)
			}
		}
	}
	return b
}

// WriteRunes appends one or more runes to the buffer.
// Safe to call on nil receiver. Returns the buffer for method chaining.
func (b *LazyBuffer) WriteRunes(rr ...rune) *LazyBuffer {
	if b != nil {
		sb := (*strings.Builder)(b)
		for _, r := range rr {
			_, _ = sb.WriteRune(r)
		}
	}
	return b
}

// Printf appends a formatted string to the buffer using fmt.Fprintf.
// Safe to call on nil receiver. Returns the buffer for method chaining.
func (b *LazyBuffer) Printf(format string, args ...any) *LazyBuffer {
	if b != nil {
		sb := (*strings.Builder)(b)
		_, _ = fmt.Fprintf(sb, format, args...)
	}
	return b
}

// String returns the accumulated string. Returns empty string if receiver is nil.
func (b *LazyBuffer) String() string {
	if b != nil {
		return (*strings.Builder)(b).String()
	}
	return ""
}
