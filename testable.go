// Package testable is the shared test harness for gloo-foo commands. It runs a
// Command[[]byte, []byte] against in-memory input and returns the captured
// output, so a command can be tested without touching real files or I/O.
//
// Test and TestLines are the convenience entry points used by the cmd-* modules.
// They delegate to the fn package — the production adapter that runs a command as
// an ordinary data function — so a command behaves identically under test and in
// use. The run subpackage offers a fluent Runner for finer control (custom
// readers, injected read errors, an explicit context).
package testable

import (
	"github.com/gloo-foo/fn"
	gloo "github.com/gloo-foo/framework"

	"github.com/gloo-foo/testable/run"
)

// Test runs cmd with the given input and returns its captured stdout as a single
// string, each output line terminated by '\n'. On command failure it returns the
// error and an empty string.
func Test(cmd gloo.Command[[]byte, []byte], input run.Input) (string, error) {
	return fn.Of(cmd).String(string(input))
}

// TestLines runs cmd with the given input and returns its captured stdout as
// lines. On command failure it returns the error and a nil slice.
func TestLines(cmd gloo.Command[[]byte, []byte], input run.Input) ([]string, error) {
	return fn.Of(cmd).Lines([]byte(input))
}
