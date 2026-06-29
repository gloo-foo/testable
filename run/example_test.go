package run_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	"github.com/gloo-foo/testable/run"
)

// upper upper-cases each input line; it stands in for a real cmd-* command.
func upper() gloo.Command[[]byte, []byte] {
	return patterns.Map(func(line []byte) ([]byte, error) {
		return []byte(strings.ToUpper(string(line))), nil
	})
}

func ExampleCommand() {
	res := run.Command(upper()).WithStdinLines("alpha", "beta").Run()
	fmt.Println(res.Err)
	fmt.Println(res.Stdout)
	// Output:
	// <nil>
	// [ALPHA BETA]
}

func ExampleQuick() {
	res := run.Quick(upper())
	fmt.Println(res.Err)
	fmt.Println(len(res.Stdout))
	// Output:
	// <nil>
	// 0
}
