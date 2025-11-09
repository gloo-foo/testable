package stream

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestBufferedReader(t *testing.T) {
	input := "line1\nline2\nline3"
	reader := NewBufferedReader(strings.NewReader(input))

	expected := []string{"line1", "line2", "line3"}
	actual := []string{}

	for {
		line, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		actual = append(actual, line)
	}

	if len(actual) != len(expected) {
		t.Errorf("got %d lines, want %d", len(actual), len(expected))
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("line[%d] = %q, want %q", i, actual[i], expected[i])
		}
	}
}

func TestBufferedReader_Empty(t *testing.T) {
	reader := NewBufferedReader(strings.NewReader(""))

	_, err := reader.ReadLine()
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF for empty input, got %v", err)
	}
}

// failingReader returns a non-EOF error on the first read so the scanner
// surfaces it through scanner.Err().
type failingReader struct{ err error }

func (f *failingReader) Read(_ []byte) (int, error) { return 0, f.err }

func TestBufferedReader_ScanError(t *testing.T) {
	wantErr := errors.New("read boom")
	reader := NewBufferedReader(&failingReader{err: wantErr})

	_, err := reader.ReadLine()
	if !errors.Is(err, wantErr) {
		t.Errorf("err = %v, want %v", err, wantErr)
	}
}

func TestSliceReader(t *testing.T) {
	lines := []string{"first", "second", "third"}
	reader := NewSliceReader(lines)

	for i, expected := range lines {
		line, err := reader.ReadLine()
		if err != nil {
			t.Fatalf("unexpected error at line %d: %v", i, err)
		}
		if line != expected {
			t.Errorf("line[%d] = %q, want %q", i, line, expected)
		}
	}

	// Should return EOF after all lines
	_, err := reader.ReadLine()
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF after all lines, got %v", err)
	}
}

func TestSliceReader_Empty(t *testing.T) {
	reader := NewSliceReader([]string{})

	_, err := reader.ReadLine()
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF for empty slice, got %v", err)
	}
}

func TestErrorInjector(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	base := NewSliceReader(lines)
	expectedErr := errors.New("injected error")
	reader := NewErrorInjector(base, 2, expectedErr)

	// First two lines should work
	line, err := reader.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error on line 1: %v", err)
	}
	if line != "line1" {
		t.Errorf("line 1 = %q, want 'line1'", line)
	}

	line, err = reader.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error on line 2: %v", err)
	}
	if line != "line2" {
		t.Errorf("line 2 = %q, want 'line2'", line)
	}

	// Third call should return error
	_, err = reader.ReadLine()
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected injected error, got %v", err)
	}

	// After error, should continue reading from wrapped reader
	line, err = reader.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error after injected error: %v", err)
	}
	if line != "line3" {
		t.Errorf("line after error = %q, want 'line3'", line)
	}
}

func TestErrorInjector_ImmediateError(t *testing.T) {
	base := NewSliceReader([]string{"line1"})
	expectedErr := errors.New("immediate error")
	reader := NewErrorInjector(base, 0, expectedErr)

	_, err := reader.ReadLine()
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected immediate error, got %v", err)
	}
}
