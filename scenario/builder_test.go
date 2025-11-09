package scenario

import (
	"errors"
	"io"
	"testing"

	"github.com/gloo-foo/testable/stream"
)

func TestQuick(t *testing.T) {
	test := Quick("line1", "line2", "line3")

	// Read all lines
	lines := []string{}
	for {
		line, err := test.Input.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		lines = append(lines, line)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Errorf("got %d lines, want %d", len(lines), len(expected))
	}
	for i := range lines {
		if lines[i] != expected[i] {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], expected[i])
		}
	}
}

func TestEmpty(t *testing.T) {
	test := Empty()

	_, err := test.Input.ReadLine()
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF for empty test, got %v", err)
	}
}

func TestBuilder_WithInput(t *testing.T) {
	test := New().
		WithInput("a", "b", "c").
		Build()

	line, err := test.Input.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line != "a" {
		t.Errorf("first line = %q, want 'a'", line)
	}
}

func TestBuilder_CaptureOutput(t *testing.T) {
	test := Quick("input")

	// Write some output
	if err := test.Stdout("output line 1"); err != nil {
		t.Fatalf("unexpected stdout error: %v", err)
	}
	if err := test.Stdout("output line 2"); err != nil {
		t.Fatalf("unexpected stdout error: %v", err)
	}

	lines := test.StdoutBuffer.Lines()
	if len(lines) != 2 {
		t.Errorf("got %d output lines, want 2", len(lines))
	}
	if lines[0] != "output line 1" {
		t.Errorf("line 0 = %q, want 'output line 1'", lines[0])
	}
	if lines[1] != "output line 2" {
		t.Errorf("line 1 = %q, want 'output line 2'", lines[1])
	}
}

func TestBuilder_CaptureStderr(t *testing.T) {
	test := Quick()

	if err := test.Stderr("error 1"); err != nil {
		t.Fatalf("unexpected stderr error: %v", err)
	}
	if err := test.Stderr("error 2"); err != nil {
		t.Fatalf("unexpected stderr error: %v", err)
	}

	lines := test.StderrBuffer.Lines()
	if len(lines) != 2 {
		t.Errorf("got %d stderr lines, want 2", len(lines))
	}
	if lines[0] != "error 1" {
		t.Errorf("line 0 = %q, want 'error 1'", lines[0])
	}
}

func TestBuilder_WithInputError(t *testing.T) {
	expectedErr := errors.New("read error")
	test := New().
		WithInput("line1", "line2").
		WithInputError(1, expectedErr).
		Build()

	// First read should work
	line, err := test.Input.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error on first read: %v", err)
	}
	if line != "line1" {
		t.Errorf("first line = %q, want 'line1'", line)
	}

	// Second read should fail
	_, err = test.Input.ReadLine()
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected injected error, got %v", err)
	}
}

func TestBuilder_WithOutputError(t *testing.T) {
	expectedErr := errors.New("write error")
	test := New().
		WithOutputError(1, expectedErr).
		Build()

	// First write should work
	if err := test.Stdout("line1"); err != nil {
		t.Fatalf("unexpected error on first write: %v", err)
	}

	// Second write should fail
	if err := test.Stdout("line2"); !errors.Is(err, expectedErr) {
		t.Errorf("expected injected error, got %v", err)
	}
}

func TestBuilder_WithStderrError(t *testing.T) {
	expectedErr := errors.New("stderr error")
	test := New().
		WithStderrError(0, expectedErr).
		Build()

	// First write to stderr should fail
	if err := test.Stderr("error"); !errors.Is(err, expectedErr) {
		t.Errorf("expected injected error, got %v", err)
	}
}

func TestBuilder_WithInputReader(t *testing.T) {
	reader := stream.NewSliceReader([]string{"custom"})
	test := New().WithInputReader(reader).Build()

	line, err := test.Input.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line != "custom" {
		t.Errorf("line = %q, want 'custom'", line)
	}
}

func TestBuilder_NoInput(t *testing.T) {
	test := New().Build()

	// Should return EOF immediately
	_, err := test.Input.ReadLine()
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF with no input, got %v", err)
	}
}
