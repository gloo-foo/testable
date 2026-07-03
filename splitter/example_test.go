package splitter_test

import (
	"fmt"

	"github.com/gloo-foo/testable/splitter"
)

func ExampleWhitespace() {
	fmt.Println(splitter.Whitespace("a  b   c", " "))
	// Output: [a b c]
}

func ExampleExact() {
	fmt.Println(splitter.Exact("a,b,,c", ","))
	// Output: [a b  c]
}

func ExampleFixed() {
	split := splitter.Fixed(0, 3, 6)
	fmt.Println(split("abcdefghij", ""))
	// Output: [abc def ghij]
}

func ExampleCharacterClass() {
	fmt.Println(splitter.CharacterClass("a,b;c:d", ",;:"))
	// Output: [a b c d]
}
