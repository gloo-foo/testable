package testable_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	"github.com/gloo-foo/testable"
)

// upper upper-cases each input line; it stands in for a real cmd-* command.
func upper() gloo.Command[[]byte, []byte] {
	return patterns.Map(func(line []byte) ([]byte, error) {
		return []byte(strings.ToUpper(string(line))), nil
	})
}

func ExampleTestLines() {
	lines, err := testable.TestLines(upper(), "hello\nworld\n")
	fmt.Println(err)
	fmt.Println(lines)
	// Output:
	// <nil>
	// [HELLO WORLD]
}

func ExampleTest() {
	out, err := testable.Test(upper(), "hello\nworld\n")
	fmt.Printf("%q\n", out)
	fmt.Println(err)
	// Output:
	// "HELLO\nWORLD\n"
	// <nil>
}
