package testable

import (
	"errors"
	"strings"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// upper is a tiny command that upper-cases each input line.
func upper() gloo.Command[[]byte, []byte] {
	return patterns.Map(func(line []byte) ([]byte, error) {
		return []byte(strings.ToUpper(string(line))), nil
	})
}

// failing is a command that fails on the first line.
func failing() gloo.Command[[]byte, []byte] {
	return patterns.Map(func([]byte) ([]byte, error) {
		return nil, errors.New("boom")
	})
}

// logSpy records Helper/Fatalf/Log calls so Assert's branches are observable
// without aborting the host test the way a real *testing.T.Fatalf would.
type logSpy struct {
	helperCalls int
	fatalfMsgs  []string
	logged      []any
}

func (s *logSpy) Helper() { s.helperCalls++ }

func (s *logSpy) Fatalf(format string, _ ...any) {
	s.fatalfMsgs = append(s.fatalfMsgs, format)
}

func (s *logSpy) Log(args ...any) { s.logged = append(s.logged, args...) }

func TestTest_OutputRetainsTrailingNewline(t *testing.T) {
	out, err := Test(upper(), "hello\nworld\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "HELLO\nWORLD\n" {
		t.Errorf("Test() = %q, want %q", out, "HELLO\nWORLD\n")
	}
}

func TestTest_PropagatesError(t *testing.T) {
	out, err := Test(failing(), "x\n")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if out != "" {
		t.Errorf("Test() output = %q, want empty on error", out)
	}
}

func TestTestLines_SplitsAndStrips(t *testing.T) {
	lines, err := TestLines(upper(), "a\nb\nc\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"A", "B", "C"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %v", len(lines), len(want), lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestTestLines_Empty(t *testing.T) {
	lines, err := TestLines(upper(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("got %d lines, want 0", len(lines))
	}
}

func TestTestLines_PropagatesError(t *testing.T) {
	lines, err := TestLines(failing(), "x\n")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if lines != nil {
		t.Errorf("lines = %v, want nil on error", lines)
	}
}

func TestLogger_AssertSuccessLogsOutput(t *testing.T) {
	s := &logSpy{}
	Logger(s).Assert("the output", nil)
	if len(s.fatalfMsgs) != 0 {
		t.Errorf("Fatalf called on success: %v", s.fatalfMsgs)
	}
	if len(s.logged) != 1 || s.logged[0] != "the output" {
		t.Errorf("logged = %v, want [the output]", s.logged)
	}
}

func TestLogger_AssertFailureCallsFatalf(t *testing.T) {
	s := &logSpy{}
	Logger(s).Assert("ignored", errors.New("boom"))
	if len(s.fatalfMsgs) != 1 {
		t.Errorf("got %d Fatalf calls, want 1", len(s.fatalfMsgs))
	}
	if len(s.logged) != 0 {
		t.Errorf("logged on failure = %v, want none", s.logged)
	}
}
