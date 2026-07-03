package run_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	"github.com/gloo-foo/testable/run"
)

// UppercaseCommand is a tiny test command that upper-cases each input line.
func UppercaseCommand() gloo.Command[[]byte, []byte] {
	return patterns.Map(func(line []byte) ([]byte, error) {
		return []byte(strings.ToUpper(string(line))), nil
	})
}

func TestCommand_WithStdinLines(t *testing.T) {
	result := run.Command(UppercaseCommand()).
		WithStdinLines("hello", "world").
		Run()

	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 2 || result.Stdout[0] != "HELLO" || result.Stdout[1] != "WORLD" {
		t.Fatalf("got %q, want [HELLO WORLD]", result.Stdout)
	}
}

func TestCommand_WithStdin(t *testing.T) {
	result := run.Command(UppercaseCommand()).
		WithStdin("hello\nworld\n").
		Run()

	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 2 || result.Stdout[0] != "HELLO" {
		t.Fatalf("got %q", result.Stdout)
	}
}

func TestCommand_NoInput(t *testing.T) {
	result := run.Quick(UppercaseCommand())
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 0 {
		t.Errorf("got %d lines, want 0", len(result.Stdout))
	}
}

func TestCommand_WithStdinError(t *testing.T) {
	wantErr := errors.New("read error")
	result := run.Command(UppercaseCommand()).
		WithStdinError(wantErr).
		Run()

	if !errors.Is(result.Err, wantErr) {
		t.Errorf("err = %v, want it to wrap %v", result.Err, wantErr)
	}
}

func TestWithInput(t *testing.T) {
	result := run.WithInput(UppercaseCommand(), "test\n")
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 1 || result.Stdout[0] != "TEST" {
		t.Errorf("got %v, want [TEST]", result.Stdout)
	}
}

func TestCommand_EmptyLines(t *testing.T) {
	result := run.Command(UppercaseCommand()).WithStdinLines().Run()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 0 {
		t.Errorf("got %d lines, want 0", len(result.Stdout))
	}
}

func TestWithContext_Standalone(t *testing.T) {
	result := run.WithContext(context.Background(), UppercaseCommand()).
		WithStdin("hi\n").
		Run()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 1 || result.Stdout[0] != "HI" {
		t.Errorf("got %v, want [HI]", result.Stdout)
	}
}

func TestRunner_WithContext(t *testing.T) {
	result := run.Command(UppercaseCommand()).
		WithContext(context.Background()).
		WithStdin("hi\n").
		Run()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 1 || result.Stdout[0] != "HI" {
		t.Errorf("got %v, want [HI]", result.Stdout)
	}
}

func TestRunner_WithStdinReader(t *testing.T) {
	result := run.Command(UppercaseCommand()).
		WithStdinReader(strings.NewReader("reader\n")).
		Run()
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Stdout) != 1 || result.Stdout[0] != "READER" {
		t.Errorf("got %v, want [READER]", result.Stdout)
	}
}
