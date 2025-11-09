package capture

import (
	"bytes"
	"errors"
	"testing"
)

// writeAll writes each line to buf, failing the test on any error.
func writeAll(t *testing.T, buf *Buffer, lines ...string) {
	t.Helper()
	for _, line := range lines {
		if err := buf.Write(line); err != nil {
			t.Fatalf("unexpected write error: %v", err)
		}
	}
}

func TestBuffer(t *testing.T) {
	buf := NewBuffer()

	writeAll(t, buf, "line1", "line2", "line3")

	lines := buf.Lines()
	if len(lines) != 3 {
		t.Errorf("got %d lines, want 3", len(lines))
	}

	expected := []string{"line1", "line2", "line3"}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("line[%d] = %q, want %q", i, line, expected[i])
		}
	}
}

func TestBuffer_String(t *testing.T) {
	buf := NewBuffer()
	writeAll(t, buf, "line1", "line2")

	result := buf.String()
	expected := "line1\nline2\n"
	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

func TestBuffer_Count(t *testing.T) {
	buf := NewBuffer()
	if buf.Count() != 0 {
		t.Errorf("initial count = %d, want 0", buf.Count())
	}

	writeAll(t, buf, "line1")
	if buf.Count() != 1 {
		t.Errorf("count after write = %d, want 1", buf.Count())
	}

	writeAll(t, buf, "line2")
	if buf.Count() != 2 {
		t.Errorf("count after 2 writes = %d, want 2", buf.Count())
	}
}

func TestBuffer_Clear(t *testing.T) {
	buf := NewBuffer()
	writeAll(t, buf, "line1", "line2")

	buf.Clear()

	if buf.Count() != 0 {
		t.Errorf("count after clear = %d, want 0", buf.Count())
	}

	lines := buf.Lines()
	if len(lines) != 0 {
		t.Errorf("lines after clear = %v, want empty", lines)
	}
}

func TestStdoutFunc(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := StdoutFunc(buf)

	err := writer("test output")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.String() != "test output\n" {
		t.Errorf("output = %q, want 'test output\\n'", buf.String())
	}
}

func TestStderrFunc(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := StderrFunc(buf)

	err := writer("error message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.String() != "error message\n" {
		t.Errorf("output = %q, want 'error message\\n'", buf.String())
	}
}

func TestErrorInjector(t *testing.T) {
	buf := NewBuffer()
	expectedErr := errors.New("write error")
	writer := NewErrorInjector(buf.Write, 2, expectedErr)

	// First two writes should succeed
	if err := writer.Write("line1"); err != nil {
		t.Fatalf("unexpected error on write 1: %v", err)
	}
	if err := writer.Write("line2"); err != nil {
		t.Fatalf("unexpected error on write 2: %v", err)
	}

	// Third write should fail
	if err := writer.Write("line3"); !errors.Is(err, expectedErr) {
		t.Errorf("expected injected error, got %v", err)
	}

	// After error, should continue working
	if err := writer.Write("line4"); err != nil {
		t.Fatalf("unexpected error after injected error: %v", err)
	}

	// Verify first two and fourth lines were written
	lines := buf.Lines()
	if len(lines) != 3 {
		t.Errorf("got %d lines, want 3 (excluding errored write)", len(lines))
	}
}

func TestErrorInjector_ImmediateError(t *testing.T) {
	buf := NewBuffer()
	expectedErr := errors.New("immediate error")
	writer := NewErrorInjector(buf.Write, 0, expectedErr)

	if err := writer.Write("line1"); !errors.Is(err, expectedErr) {
		t.Errorf("expected immediate error, got %v", err)
	}

	// Buffer should be empty since error happened before write
	if buf.Count() != 0 {
		t.Errorf("buffer count = %d, want 0", buf.Count())
	}
}

func TestTee(t *testing.T) {
	buf1 := NewBuffer()
	buf2 := NewBuffer()
	buf3 := NewBuffer()

	tee := NewTee(buf1.Write, buf2.Write, buf3.Write)

	err := tee.Write("test line")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All buffers should have the line
	for i, buf := range []*Buffer{buf1, buf2, buf3} {
		lines := buf.Lines()
		if len(lines) != 1 || lines[0] != "test line" {
			t.Errorf("buffer %d = %v, want ['test line']", i, lines)
		}
	}
}

func TestTee_WithError(t *testing.T) {
	buf1 := NewBuffer()
	expectedErr := errors.New("write error")
	errorWriter := NewErrorInjector(NewBuffer().Write, 0, expectedErr)
	buf3 := NewBuffer()

	tee := NewTee(buf1.Write, errorWriter.Write, buf3.Write)

	err := tee.Write("test line")
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error from tee, got %v", err)
	}

	// First and third buffers should still have the line
	if buf1.Count() != 1 {
		t.Errorf("buf1 count = %d, want 1", buf1.Count())
	}
	if buf3.Count() != 1 {
		t.Errorf("buf3 count = %d, want 1", buf3.Count())
	}
}
