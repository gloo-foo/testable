package stream

import (
	"bufio"
	"io"
)

// LineReader abstracts line-by-line reading for testability.
// This interface allows commands to process input without depending on
// concrete io.Reader implementations, enabling easy mocking and testing.
type LineReader interface {
	// ReadLine returns the next line of input.
	// Returns io.EOF when no more lines are available.
	ReadLine() (string, error)
}

// BufferedReader wraps an io.Reader with buffered line reading.
// This is the standard implementation used in production code.
type BufferedReader struct {
	scanner *bufio.Scanner
}

// NewBufferedReader creates a LineReader that reads from the given io.Reader.
func NewBufferedReader(r io.Reader) *BufferedReader {
	return &BufferedReader{scanner: bufio.NewScanner(r)}
}

// ReadLine reads and returns the next line from the input.
func (b *BufferedReader) ReadLine() (string, error) {
	if b.scanner.Scan() {
		return b.scanner.Text(), nil
	}
	if err := b.scanner.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

// SliceReader reads lines from a pre-defined slice.
// Useful for testing with known input sequences.
type SliceReader struct {
	lines []string
	index int
}

// NewSliceReader creates a LineReader from a slice of strings.
func NewSliceReader(lines []string) *SliceReader {
	return &SliceReader{lines: lines, index: 0}
}

// ReadLine returns the next line from the slice.
func (s *SliceReader) ReadLine() (string, error) {
	if s.index >= len(s.lines) {
		return "", io.EOF
	}
	line := s.lines[s.index]
	s.index++
	return line, nil
}

// ErrorInjector wraps a LineReader and injects errors at specific points.
// Useful for testing error handling paths.
type ErrorInjector struct {
	reader    LineReader
	errorAt   int   // Line number to inject error (0-based)
	error     error // Error to inject
	callCount int   // Current call count
}

// NewErrorInjector creates a LineReader that injects an error after reading a specific number of lines.
func NewErrorInjector(reader LineReader, errorAt int, err error) *ErrorInjector {
	return &ErrorInjector{
		reader:  reader,
		errorAt: errorAt,
		error:   err,
	}
}

// ReadLine delegates to the wrapped reader until the error point is reached.
func (e *ErrorInjector) ReadLine() (string, error) {
	if e.callCount == e.errorAt && e.error != nil {
		e.callCount++
		return "", e.error
	}
	e.callCount++
	return e.reader.ReadLine()
}
