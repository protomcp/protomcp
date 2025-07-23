package generator

import (
	"testing"
)

func TestLazyBuffer(t *testing.T) {
	t.Run("WriteString", testLazyBufferWriteString)
	t.Run("WriteString with empty strings", testLazyBufferWriteStringEmpty)
	t.Run("Printf", testLazyBufferPrintf)
	t.Run("WriteRunes", testLazyBufferWriteRunes)
	t.Run("method chaining", testLazyBufferChaining)
	t.Run("nil buffer", testLazyBufferNil)
}

func testLazyBufferWriteString(t *testing.T) {
	var buf LazyBuffer
	buf.WriteString("hello", " ", "world", "")
	if got := buf.String(); got != "hello world" {
		t.Errorf("WriteString() = %q, want %q", got, "hello world")
	}
}

func testLazyBufferWriteStringEmpty(t *testing.T) {
	var buf LazyBuffer
	buf.WriteString("", "hello", "", "", "world", "")
	if got := buf.String(); got != "helloworld" {
		t.Errorf("WriteString() = %q, want %q", got, "helloworld")
	}
}

func testLazyBufferPrintf(t *testing.T) {
	var buf LazyBuffer
	buf.Printf("Hello %s, number %d", "world", 42)
	if got := buf.String(); got != "Hello world, number 42" {
		t.Errorf("Printf() = %q, want %q", got, "Hello world, number 42")
	}
}

func testLazyBufferWriteRunes(t *testing.T) {
	var buf LazyBuffer
	buf.WriteRunes('H', 'e', 'l', 'l', 'o')
	buf.WriteRunes(' ', '世', '界') // Test with Unicode
	if got := buf.String(); got != "Hello 世界" {
		t.Errorf("WriteRunes() = %q, want %q", got, "Hello 世界")
	}
}

func testLazyBufferChaining(t *testing.T) {
	var buf LazyBuffer
	// Test method chaining
	result := buf.WriteString("Hello").
		WriteRunes(' ').
		WriteString("world").
		Printf(", %d", 42).
		WriteRunes('!').
		String()

	if result != "Hello world, 42!" {
		t.Errorf("chained methods = %q, want %q", result, "Hello world, 42!")
	}
}

func testLazyBufferNil(t *testing.T) {
	var buf *LazyBuffer
	// These should not panic
	buf.WriteString("test")
	buf.Printf("test %d", 1)
	buf.WriteRunes('a', 'b', 'c')
	if got := buf.String(); got != "" {
		t.Errorf("nil buffer String() = %q, want empty", got)
	}

	// Test that nil buffer returns nil for chaining
	if result := buf.WriteString("test"); result != nil {
		t.Errorf("nil buffer WriteString() = %v, want nil", result)
	}
}
