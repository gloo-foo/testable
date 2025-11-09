package splitter

import (
	"reflect"
	"testing"
)

func TestWhitespace(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		separator string
		want      []string
	}{
		{
			name:      "default whitespace",
			line:      "one  two   three",
			separator: " ",
			want:      []string{"one", "two", "three"},
		},
		{
			name:      "tabs and spaces",
			line:      "a\t\tb  c",
			separator: " ",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "comma separator",
			line:      "a,b,c",
			separator: ",",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "colon separator",
			line:      "first:second:third",
			separator: ":",
			want:      []string{"first", "second", "third"},
		},
		{
			name:      "empty line whitespace",
			line:      "",
			separator: " ",
			want:      []string{},
		},
		{
			name:      "empty line comma",
			line:      "",
			separator: ",",
			want:      []string{""},
		},
		{
			name:      "multi-char separator",
			line:      "a::b::c",
			separator: "::",
			want:      []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Whitespace(tt.line, tt.separator)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Whitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExact(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		separator string
		want      []string
	}{
		{
			name:      "single space treated literally",
			line:      "a b  c",
			separator: " ",
			want:      []string{"a", "b", "", "c"},
		},
		{
			name:      "comma",
			line:      "a,b,c",
			separator: ",",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "empty fields",
			line:      "a,,c",
			separator: ",",
			want:      []string{"a", "", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exact(tt.line, tt.separator)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPattern(t *testing.T) {
	// Custom splitter that reverses the line
	customSplit := func(line string) []string {
		runes := []rune(line)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return []string{string(runes)}
	}

	splitter := Pattern(customSplit)
	got := splitter("hello", "")
	want := []string{"olleh"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Pattern() = %v, want %v", got, want)
	}
}

func TestFixed(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		positions []int
		want      []string
	}{
		{
			name:      "basic fixed width",
			line:      "abcdefghij",
			positions: []int{0, 3, 6},
			want:      []string{"abc", "def", "ghij"},
		},
		{
			name:      "uneven fields",
			line:      "12345678",
			positions: []int{0, 2, 5},
			want:      []string{"12", "345", "678"},
		},
		{
			name:      "line shorter than positions",
			line:      "short",
			positions: []int{0, 3, 10},
			want:      []string{"sho", "rt"},
		},
		{
			name:      "single position",
			line:      "test line",
			positions: []int{5},
			want:      []string{"line"},
		},
		{
			name:      "no positions",
			line:      "test",
			positions: []int{},
			want:      []string{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitter := Fixed(tt.positions...)
			got := splitter(tt.line, "")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fixed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacterClass(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		separator string
		want      []string
	}{
		{
			name:      "multiple delimiters",
			line:      "a,b;c:d",
			separator: ",;:",
			want:      []string{"a", "b", "c", "d"},
		},
		{
			name:      "whitespace class",
			line:      "a b\tc\nd",
			separator: " \t\n",
			want:      []string{"a", "b", "c", "d"},
		},
		{
			name:      "empty separator",
			line:      "test",
			separator: "",
			want:      []string{"test"},
		},
		{
			name:      "no matches",
			line:      "test",
			separator: ",;:",
			want:      []string{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CharacterClass(tt.line, tt.separator)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CharacterClass() = %v, want %v", got, tt.want)
			}
		})
	}
}
