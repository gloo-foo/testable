package splitter

import "strings"

// Line is a single line of input text to be split into fields.
type Line string

// Separator selects how a Line is split into fields.
type Separator string

// Position is a 0-based starting column of a fixed-width field.
type Position int

// FieldFunc splits a line into fields based on a separator.
// This abstraction enables custom field splitting logic and testing.
type FieldFunc func(line Line, separator Separator) []string

// Whitespace splits on whitespace (spaces, tabs, etc.) when separator is " ".
// For any other separator, it splits on that exact string.
// This mimics awk's default field splitting behavior.
func Whitespace(line Line, separator Separator) []string {
	if separator == " " {
		// Default: split on runs of whitespace
		return strings.Fields(string(line))
	}
	// Custom separator: split on exact match
	return strings.Split(string(line), string(separator))
}

// Exact always splits on the exact separator string.
// Unlike Whitespace, it treats " " as a literal space, not as "any whitespace".
func Exact(line Line, separator Separator) []string {
	return strings.Split(string(line), string(separator))
}

// Pattern creates a FieldFunc that uses a custom splitting function.
// This allows for regex-based or other complex splitting logic.
func Pattern(splitFunc func(string) []string) FieldFunc {
	return func(line Line, _ Separator) []string {
		return splitFunc(string(line))
	}
}

// fieldEnd returns the exclusive end index of the field bounded by the
// remaining start positions, clamped to the end of line.
func fieldEnd(rest []Position, line Line) int {
	if len(rest) == 0 {
		return len(line)
	}
	return min(int(rest[0]), len(line))
}

// fixedFields splits line at the given start positions, clamping each field to
// the end of the line and stopping once a start position is past the line.
func fixedFields(positions []Position, line Line) []string {
	if len(positions) == 0 {
		return []string{string(line)}
	}
	fields := make([]string, 0, len(positions))
	for i, pos := range positions {
		if int(pos) >= len(line) {
			break
		}
		fields = append(fields, string(line[pos:fieldEnd(positions[i+1:], line)]))
	}
	return fields
}

// Fixed creates a FieldFunc that splits at fixed-width positions.
// Positions define where each field starts (0-based).
// The last position extends to the end of the line.
func Fixed(positions ...Position) FieldFunc {
	return func(line Line, _ Separator) []string {
		return fixedFields(positions, line)
	}
}

// CharacterClass splits on any character from the separator string.
// Similar to splitting on a regex character class [chars].
func CharacterClass(line Line, separator Separator) []string {
	if separator == "" {
		return []string{string(line)}
	}

	return strings.FieldsFunc(string(line), func(r rune) bool {
		return strings.ContainsRune(string(separator), r)
	})
}
