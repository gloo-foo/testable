package capture

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// OutputFunc is a function type that writes output.
// This abstraction enables testing by allowing output to be captured
// or redirected without depending on concrete io.Writer implementations.
type OutputFunc func(output string) error

// StdoutFunc creates an OutputFunc that writes to an io.Writer with newlines.
// This is the standard implementation for production use.
func StdoutFunc(w io.Writer) OutputFunc {
	return func(output string) error {
		_, err := fmt.Fprintln(w, output)
		return err
	}
}

// StderrFunc creates an OutputFunc that writes to an io.Writer with newlines.
// Identical to StdoutFunc but semantically distinct for stderr.
func StderrFunc(w io.Writer) OutputFunc {
	return func(output string) error {
		_, err := fmt.Fprintln(w, output)
		return err
	}
}

// Buffer captures output into a thread-safe buffer for testing.
// Use this to verify command output in tests.
type Buffer struct {
	lines []string
	mu    sync.Mutex
}

// NewBuffer creates a new output buffer.
func NewBuffer() *Buffer {
	return &Buffer{lines: make([]string, 0)}
}

// Write captures a single line of output.
func (b *Buffer) Write(output string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, output)
	return nil
}

// Lines returns all captured output lines.
func (b *Buffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	result := make([]string, len(b.lines))
	copy(result, b.lines)
	return result
}

// String returns all captured output as a single string with newlines.
func (b *Buffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	buf := &bytes.Buffer{}
	for _, line := range b.lines {
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.String()
}

// Count returns the number of lines captured.
func (b *Buffer) Count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.lines)
}

// Clear removes all captured output.
func (b *Buffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = make([]string, 0)
}

// ErrorInjector returns an error after a specified number of writes.
// Useful for testing error handling in output paths.
type ErrorInjector struct {
	wrapped   OutputFunc
	errorAt   int   // Write number to inject error (0-based)
	error     error // Error to inject
	callCount int   // Current call count
}

// NewErrorInjector creates an OutputFunc that injects an error at a specific call.
func NewErrorInjector(wrapped OutputFunc, errorAt int, err error) *ErrorInjector {
	return &ErrorInjector{
		wrapped: wrapped,
		errorAt: errorAt,
		error:   err,
	}
}

// Write delegates to the wrapped function until the error point is reached.
func (e *ErrorInjector) Write(output string) error {
	if e.callCount == e.errorAt && e.error != nil {
		e.callCount++
		return e.error
	}
	e.callCount++
	return e.wrapped(output)
}

// Tee writes output to multiple OutputFuncs simultaneously.
// Similar to Unix tee command, useful for capturing output while still displaying it.
type Tee struct {
	outputs []OutputFunc
}

// NewTee creates an OutputFunc that writes to multiple destinations.
func NewTee(outputs ...OutputFunc) *Tee {
	return &Tee{outputs: outputs}
}

// Write sends output to all wrapped OutputFuncs.
// Returns the first error encountered, but attempts all writes.
func (t *Tee) Write(output string) error {
	var firstErr error
	for _, out := range t.outputs {
		if err := out(output); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
