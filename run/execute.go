// Package run provides a small fluent harness for exercising a
// Command[[]byte, []byte] in tests: configure stdin, run, and inspect the
// captured output lines and error.
//
//	res := run.Command(myCmd).WithStdinLines("a", "b").Run()
//	// res.Stdout == []string{...}, res.Err == nil
//
// It complements the package-level testable.Test / testable.TestLines helpers
// with explicit input control (custom readers, injected read errors).
package run

import (
	"context"
	"io"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

// Result holds the captured output and error from a command execution.
type Result struct {
	Stdout []string // Output lines (the command's output stream)
	Err    error    // First error carried by the stream, if any
}

// Runner configures and executes a command for testing. Every With* method
// returns a NEW Runner, so a Runner is immutable and safe to reuse.
type Runner struct {
	cmd      gloo.Command[[]byte, []byte]
	ctx      context.Context
	stdin    io.Reader
	stdinErr error
}

// Command creates a Runner for cmd with empty stdin and a background context.
func Command(cmd gloo.Command[[]byte, []byte]) Runner {
	return Runner{cmd: cmd, ctx: context.Background(), stdin: strings.NewReader("")}
}

// WithContext (standalone) creates a Runner bound to ctx.
func WithContext(ctx context.Context, cmd gloo.Command[[]byte, []byte]) Runner {
	return Runner{cmd: cmd, ctx: ctx, stdin: strings.NewReader("")}
}

// WithContext returns a copy of the Runner bound to ctx.
func (r Runner) WithContext(ctx context.Context) Runner {
	r.ctx = ctx
	return r
}

// WithStdin returns a copy of the Runner whose input is input.
func (r Runner) WithStdin(input string) Runner {
	r.stdin = strings.NewReader(input)
	r.stdinErr = nil
	return r
}

// WithStdinLines returns a copy of the Runner whose input is lines joined by
// newlines (each line, including the last, terminated by '\n').
func (r Runner) WithStdinLines(lines ...string) Runner {
	if len(lines) == 0 {
		r.stdin = strings.NewReader("")
	} else {
		r.stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")
	}
	r.stdinErr = nil
	return r
}

// WithStdinReader returns a copy of the Runner reading input from reader.
func (r Runner) WithStdinReader(reader io.Reader) Runner {
	r.stdin = reader
	r.stdinErr = nil
	return r
}

// WithStdinError returns a copy of the Runner that fails the input stream with
// err, for exercising error propagation.
func (r Runner) WithStdinError(err error) Runner {
	r.stdinErr = err
	return r
}

// Run executes the configured command and returns the captured result. It is a
// terminal operation.
func (r Runner) Run() *Result {
	reader := r.stdin
	if r.stdinErr != nil {
		reader = &errorReader{err: r.stdinErr}
	}

	source := gloo.ByteReaderSource([]io.Reader{reader})
	output := r.cmd.Execute(r.ctx, source.Stream(r.ctx))
	items, err := gloo.Collect(r.ctx, output)

	lines := make([]string, len(items))
	for i, b := range items {
		lines[i] = string(b)
	}
	return &Result{Stdout: lines, Err: err}
}

// errorReader fails on the first read with its configured error.
type errorReader struct{ err error }

func (e *errorReader) Read(_ []byte) (int, error) { return 0, e.err }

// Quick executes cmd with empty stdin and returns the result.
func Quick(cmd gloo.Command[[]byte, []byte]) *Result {
	return Command(cmd).Run()
}

// WithInput executes cmd with the given stdin string and returns the result.
func WithInput(cmd gloo.Command[[]byte, []byte], stdin string) *Result {
	return Command(cmd).WithStdin(stdin).Run()
}
