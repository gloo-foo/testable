package testable_test

import (
	"errors"
	"io"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/capture"
	"github.com/gloo-foo/testable/scenario"
	"github.com/gloo-foo/testable/splitter"
	"github.com/gloo-foo/testable/stream"
)

// Example command that processes lines
type ProcessCommand struct {
	Input  stream.LineReader
	Output capture.OutputFunc
	Prefix string
}

func (p *ProcessCommand) Execute() error {
	for {
		line, err := p.Input.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		result := p.Prefix + line
		if err := p.Output(result); err != nil {
			return err
		}
	}
	return nil
}

func TestProcessCommand_Success(t *testing.T) {
	test := scenario.Quick("input1", "input2", "input3")

	cmd := &ProcessCommand{
		Input:  test.Input,
		Output: test.Stdout,
		Prefix: "processed: ",
	}

	err := cmd.Execute()

	assertion.NoError(t, err)
	assertion.Lines(t, test.StdoutBuffer.Lines(), []string{
		"processed: input1",
		"processed: input2",
		"processed: input3",
	})
}

func TestProcessCommand_EmptyInput(t *testing.T) {
	test := scenario.Empty()

	cmd := &ProcessCommand{
		Input:  test.Input,
		Output: test.Stdout,
		Prefix: "prefix: ",
	}

	err := cmd.Execute()

	assertion.NoError(t, err)
	assertion.Empty(t, test.StdoutBuffer.Lines())
}

func TestProcessCommand_InputError(t *testing.T) {
	expectedErr := errors.New("read failed")
	test := scenario.New().
		WithInput("line1", "line2").
		WithInputError(1, expectedErr).
		Build()

	cmd := &ProcessCommand{
		Input:  test.Input,
		Output: test.Stdout,
	}

	err := cmd.Execute()

	assertion.ErrorContains(t, err, "read failed")
	assertion.Count(t, test.StdoutBuffer.Lines(), 1)
}

func TestProcessCommand_OutputError(t *testing.T) {
	expectedErr := errors.New("write failed")
	test := scenario.New().
		WithInput("line1", "line2").
		WithOutputError(1, expectedErr).
		Build()

	cmd := &ProcessCommand{
		Input:  test.Input,
		Output: test.Stdout,
	}

	err := cmd.Execute()

	assertion.ErrorContains(t, err, "write failed")
}

func TestSplitter_Integration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		splitter splitter.FieldFunc
		sep      string
		want     []string
	}{
		{
			name:     "csv",
			input:    "a,b,c,d",
			splitter: splitter.Whitespace,
			sep:      ",",
			want:     []string{"a", "b", "c", "d"},
		},
		{
			name:     "whitespace",
			input:    "one   two  three",
			splitter: splitter.Whitespace,
			sep:      " ",
			want:     []string{"one", "two", "three"},
		},
		{
			name:     "fixed width",
			input:    "abcdefghij",
			splitter: splitter.Fixed(0, 3, 6),
			sep:      "",
			want:     []string{"abc", "def", "ghij"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.splitter(tt.input, tt.sep)
			assertion.Equal(t, len(got), len(tt.want), "field count")
			for i := range got {
				assertion.Equal(t, got[i], tt.want[i], "field")
			}
		})
	}
}

func TestCapture_Tee(t *testing.T) {
	buf1 := capture.NewBuffer()
	buf2 := capture.NewBuffer()

	tee := capture.NewTee(buf1.Write, buf2.Write)

	if err := tee.Write("line1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tee.Write("line2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertion.Lines(t, buf1.Lines(), []string{"line1", "line2"})
	assertion.Lines(t, buf2.Lines(), []string{"line1", "line2"})
}

func TestAssertion_Examples(t *testing.T) {
	output := []string{"INFO: started", "INFO: processing", "INFO: done"}

	assertion.Count(t, output, 3)
	assertion.Prefix(t, output, "INFO:")
	assertion.Contains(t, output, "processing", "done")
	assertion.NotContains(t, output, "ERROR", "WARN")
}
