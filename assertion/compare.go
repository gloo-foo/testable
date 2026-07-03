package assertion

import (
	"fmt"
	"strings"
)

// reporter is the slice of *testing.T that the assertion helpers depend on.
// *testing.T satisfies it, so callers pass a *testing.T unchanged; tests inject
// a spy to exercise the failure branches without failing the host test.
type reporter interface {
	Helper()
	Errorf(format string, args ...any)
}

// Substring is a fragment of text expected (or forbidden) within output.
type Substring string

// LineCount is a number of output lines.
type LineCount int

// LinePrefix is text every output line must start with.
type LinePrefix string

// LineSuffix is text every output line must end with.
type LineSuffix string

// Label names the value under comparison in a failure message.
type Label string

// Condition is a boolean expectation checked by True and False.
type Condition bool

// Message describes an assertion in its failure report.
type Message string

// Template is a printf-style format template.
type Template string

// Lines verifies that actual output matches expected lines.
// Reports detailed differences using Errorf if they don't match.
func Lines(t reporter, actual, expected []string) {
	t.Helper()
	if len(actual) != len(expected) {
		reportLineCount(t, actual, expected)
		return
	}
	reportLineDiffs(t, actual, expected)
}

// reportLineCount reports a length mismatch between actual and expected output.
func reportLineCount(t reporter, actual, expected []string) {
	t.Errorf("line count mismatch:\n  got:  %d lines\n  want: %d lines", len(actual), len(expected))
	t.Errorf("actual output:\n%s", strings.Join(actual, "\n"))
	t.Errorf("expected output:\n%s", strings.Join(expected, "\n"))
}

// reportLineDiffs reports each line where actual and expected differ. The slices
// must already have equal length.
func reportLineDiffs(t reporter, actual, expected []string) {
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("line %d mismatch:\n  got:  %q\n  want: %q", i+1, actual[i], expected[i])
		}
	}
}

// Contains verifies that the actual output contains all expected strings.
func Contains(t reporter, actual []string, expected ...Substring) {
	t.Helper()
	full := strings.Join(actual, "\n")
	for _, exp := range expected {
		if !strings.Contains(full, string(exp)) {
			t.Errorf("output missing expected string:\n  want: %q\n  in:   %s", exp, full)
		}
	}
}

// NotContains verifies that the actual output does not contain any of the given strings.
func NotContains(t reporter, actual []string, unexpected ...Substring) {
	t.Helper()
	full := strings.Join(actual, "\n")
	for _, unexp := range unexpected {
		if strings.Contains(full, string(unexp)) {
			t.Errorf("output contains unexpected string:\n  found: %q\n  in:    %s", unexp, full)
		}
	}
}

// Empty verifies that the output is empty.
func Empty(t reporter, actual []string) {
	t.Helper()
	if len(actual) != 0 {
		t.Errorf("expected empty output, got %d lines:\n%s", len(actual), strings.Join(actual, "\n"))
	}
}

// Count verifies that the output has exactly the expected number of lines.
func Count(t reporter, actual []string, expected LineCount) {
	t.Helper()
	if LineCount(len(actual)) != expected {
		t.Errorf("line count mismatch:\n  got:  %d\n  want: %d", len(actual), expected)
	}
}

// Prefix verifies that each line starts with the expected prefix.
func Prefix(t reporter, actual []string, prefix LinePrefix) {
	t.Helper()
	for i, line := range actual {
		if !strings.HasPrefix(line, string(prefix)) {
			t.Errorf("line %d missing prefix:\n  line: %q\n  want prefix: %q", i+1, line, prefix)
		}
	}
}

// Suffix verifies that each line ends with the expected suffix.
func Suffix(t reporter, actual []string, suffix LineSuffix) {
	t.Helper()
	for i, line := range actual {
		if !strings.HasSuffix(line, string(suffix)) {
			t.Errorf("line %d missing suffix:\n  line: %q\n  want suffix: %q", i+1, line, suffix)
		}
	}
}

// ErrorContains verifies that an error contains the expected substring.
func ErrorContains(t reporter, err error, expected Substring) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", expected)
		return
	}
	if !strings.Contains(err.Error(), string(expected)) {
		t.Errorf("error message mismatch:\n  got:  %q\n  want substring: %q", err.Error(), expected)
	}
}

// NoError verifies that no error occurred.
func NoError(t reporter, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Equal compares two values and reports if they differ.
func Equal[T comparable](t reporter, actual, expected T, label Label) {
	t.Helper()
	if actual != expected {
		t.Errorf("%s mismatch:\n  got:  %v\n  want: %v", label, actual, expected)
	}
}

// True verifies that isCondition holds.
func True(t reporter, isCondition Condition, message Message) {
	t.Helper()
	if !isCondition {
		t.Errorf("condition failed: %s", message)
	}
}

// False verifies that isCondition does not hold.
func False(t reporter, isCondition Condition, message Message) {
	t.Helper()
	if isCondition {
		t.Errorf("condition should be false: %s", message)
	}
}

// Format is a helper for creating formatted error messages.
func Format(format Template, args ...any) string {
	return fmt.Sprintf(string(format), args...)
}

// Error verifies that an error occurred (err != nil).
func Error(t reporter, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}
