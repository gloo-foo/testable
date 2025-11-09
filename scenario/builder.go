package scenario

import (
	"github.com/gloo-foo/testable/capture"
	"github.com/gloo-foo/testable/stream"
)

// Test represents a complete test scenario with input, output, and error handling.
// Use Builder to construct test scenarios fluently.
type Test struct {
	Input  stream.LineReader
	Stdout capture.OutputFunc
	Stderr capture.OutputFunc

	// Captured outputs for verification
	StdoutBuffer *capture.Buffer
	StderrBuffer *capture.Buffer
}

// Builder constructs test scenarios fluently.
type Builder struct {
	input        stream.LineReader
	captureOut   *capture.Buffer
	captureErr   *capture.Buffer
	injectOutErr error
	injectErrErr error
	errAtLine    int
	errOnOutput  int
}

// New creates a new test scenario builder.
func New() *Builder {
	return &Builder{
		captureOut: capture.NewBuffer(),
		captureErr: capture.NewBuffer(),
	}
}

// WithInput sets the input source for the test.
func (b *Builder) WithInput(lines ...string) *Builder {
	b.input = stream.NewSliceReader(lines)
	return b
}

// WithInputReader sets a custom input reader.
func (b *Builder) WithInputReader(reader stream.LineReader) *Builder {
	b.input = reader
	return b
}

// WithInputError injects a read error after reading N lines.
func (b *Builder) WithInputError(afterLines int, err error) *Builder {
	b.errAtLine = afterLines
	if b.input != nil {
		b.input = stream.NewErrorInjector(b.input, afterLines, err)
	}
	return b
}

// WithOutputError injects a write error on stdout after N writes.
func (b *Builder) WithOutputError(afterWrites int, err error) *Builder {
	b.errOnOutput = afterWrites
	b.injectOutErr = err
	return b
}

// WithStderrError injects a write error on stderr after N writes.
func (b *Builder) WithStderrError(afterWrites int, err error) *Builder {
	b.injectErrErr = err
	return b
}

// Build constructs the final Test scenario.
func (b *Builder) Build() *Test {
	if b.input == nil {
		b.input = stream.NewSliceReader([]string{})
	}

	// Create output functions
	stdout := capture.OutputFunc(b.captureOut.Write)
	stderr := capture.OutputFunc(b.captureErr.Write)

	// Inject errors if specified
	if b.injectOutErr != nil {
		injector := capture.NewErrorInjector(stdout, b.errOnOutput, b.injectOutErr)
		stdout = injector.Write
	}
	if b.injectErrErr != nil {
		injector := capture.NewErrorInjector(stderr, 0, b.injectErrErr)
		stderr = injector.Write
	}

	return &Test{
		Input:        b.input,
		Stdout:       stdout,
		Stderr:       stderr,
		StdoutBuffer: b.captureOut,
		StderrBuffer: b.captureErr,
	}
}

// Quick creates a simple test scenario with the given input lines.
// Most common use case - just provide input and capture output.
func Quick(input ...string) *Test {
	return New().WithInput(input...).Build()
}

// Empty creates a test scenario with no input.
func Empty() *Test {
	return New().Build()
}
