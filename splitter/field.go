package splitter

import "strings"

// FieldFunc splits a line into fields based on a separator.
// This abstraction enables custom field splitting logic and testing.
type FieldFunc func(line, separator string) []string

// Whitespace splits on whitespace (spaces, tabs, etc.) when separator is " ".
// For any other separator, it splits on that exact string.
// This mimics awk's default field splitting behavior.
func Whitespace(line, separator string) []string {
	if separator == " " {
		// Default: split on runs of whitespace
		return strings.Fields(line)
	}
	// Custom separator: split on exact match
	return strings.Split(line, separator)
}

// Exact always splits on the exact separator string.
// Unlike Whitespace, it treats " " as a literal space, not as "any whitespace".
func Exact(line, separator string) []string {
	return strings.Split(line, separator)
}

// Pattern creates a FieldFunc that uses a custom splitting function.
// This allows for regex-based or other complex splitting logic.
func Pattern(splitFunc func(string) []string) FieldFunc {
	return func(line, _ string) []string {
		return splitFunc(line)
	}
}

// fieldEnd returns the exclusive end index of the field starting at positions[i],
// clamped to the end of line.
func fieldEnd(positions []int, i, lineLen int) int {
	if i+1 >= len(positions) {
		return lineLen
	}
	return min(positions[i+1], lineLen)
}

// fixedFields splits line at the given start positions, clamping each field to
// the end of the line and stopping once a start position is past the line.
func fixedFields(positions []int, line string) []string {
	if len(positions) == 0 {
		return []string{line}
	}
	fields := make([]string, 0, len(positions))
	for i, pos := range positions {
		if pos >= len(line) {
			break
		}
		fields = append(fields, line[pos:fieldEnd(positions, i, len(line))])
	}
	return fields
}

// Fixed creates a FieldFunc that splits at fixed-width positions.
// Positions define where each field starts (0-based).
// The last position extends to the end of the line.
func Fixed(positions ...int) FieldFunc {
	return func(line, _ string) []string {
		return fixedFields(positions, line)
	}
}

// CharacterClass splits on any character from the separator string.
// Similar to splitting on a regex character class [chars].
func CharacterClass(line, separator string) []string {
	if separator == "" {
		return []string{line}
	}

	return strings.FieldsFunc(line, func(r rune) bool {
		return strings.ContainsRune(separator, r)
	})
}
