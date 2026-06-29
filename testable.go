// Package testable is the shared test harness for gloo-foo commands. It runs a
// Command[[]byte, []byte] against in-memory input and returns the captured
// output, so a command can be tested without touching real files or I/O.
//
// Test and TestLines are the convenience entry points used by the cmd-* modules;
// both delegate execution to the run subpackage, which offers a fluent Runner
// for finer control (custom readers, injected read errors, an explicit context).
package testable

import (
	"strings"

	gloo "github.com/gloo-foo/framework"

	"github.com/gloo-foo/testable/run"
)

// Test runs cmd with the given input and returns its captured stdout as a single
// string, each output line terminated by '\n'. On command failure it returns the
// error and an empty string.
func Test(cmd gloo.Command[[]byte, []byte], input string) (string, error) {
	res := run.WithInput(cmd, input)
	if res.Err != nil {
		return "", res.Err
	}
	var b strings.Builder
	for _, line := range res.Stdout {
		_, _ = b.WriteString(line)
		_ = b.WriteByte('\n')
	}
	return b.String(), nil
}

// TestLines runs cmd with the given input and returns its captured stdout as
// lines. On command failure it returns the error and a nil slice.
func TestLines(cmd gloo.Command[[]byte, []byte], input string) ([]string, error) {
	res := run.WithInput(cmd, input)
	if res.Err != nil {
		return nil, res.Err
	}
	return res.Stdout, nil
}
