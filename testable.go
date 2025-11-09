package testable

import (
	"bytes"
	"context"
	"io"
	"strings"

	gloo "github.com/gloo-foo/framework"
)

// Test runs a Command[[]byte, []byte] with string input and returns captured output.
func Test(cmd gloo.Command[[]byte, []byte], input string) (string, error) {
	ctx := context.Background()
	reader := strings.NewReader(input)
	source := gloo.ByteReaderSource([]io.Reader{reader})
	output := cmd.Execute(ctx, source.Stream(ctx))
	lines, err := gloo.Collect(ctx, output)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, line := range lines {
		buf.Write(line)
		buf.WriteByte('\n')
	}
	return buf.String(), nil
}

// reporter is the slice of *testing.T that TestLogger depends on. *testing.T
// satisfies it, so callers pass a *testing.T unchanged; tests inject a spy to
// exercise the failure branch without aborting the host test.
type reporter interface {
	Helper()
	Fatalf(format string, args ...any)
	Log(args ...any)
}

// TestLogger provides test assertion helpers.
type TestLogger struct {
	t reporter
}

// Logger creates a TestLogger for the given test.
func Logger(t reporter) *TestLogger {
	return &TestLogger{t: t}
}

// Assert verifies the command executed successfully and logs the output.
func (l *TestLogger) Assert(output string, err error) {
	l.t.Helper()
	if err != nil {
		l.t.Fatalf("command failed: %v", err)
		return
	}
	l.t.Log(output)
}

// TestLines runs a Command[[]byte, []byte] with string input and returns output as lines.
func TestLines(cmd gloo.Command[[]byte, []byte], input string) ([]string, error) {
	output, err := Test(cmd, input)
	if err != nil {
		return nil, err
	}
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(strings.TrimRight(output, "\n"), "\n"), nil
}
