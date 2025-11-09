# Test Template for Gloo-Foo Commands

This document explains how to write comprehensive tests for gloo-foo commands using the `testable` package. The `awk` command tests serve as the reference implementation.

## Test File Structure

```go
package mycommand_test  // Always use _test package

import (
    "testing"
    "github.com/gloo-foo/testable/run"
    "github.com/gloo-foo/testable/assertion"
    command "github.com/your/command"  // Alias to avoid conflicts
)
```

## Test Organization (Use Sections with Comments)

```go
// ==============================================================================
// Test Pure Functions
// ==============================================================================

// Test any pure functions (no I/O) first

// ==============================================================================
// Test Command Execution - Simple Cases
// ==============================================================================

// Basic success paths

// ==============================================================================
// Test Custom Behavior
// ==============================================================================

// Command-specific logic tests

// ==============================================================================
// Test Error Handling
// ==============================================================================

// All error paths

// ==============================================================================
// Test Edge Cases
// ==============================================================================

// Unicode, empty input, very long lines, etc.

// ==============================================================================
// Table-Driven Test Example
// ==============================================================================

// At least one table-driven test showing the pattern
```

## Pattern 1: Test Pure Functions First

Pure functions (no I/O) are easiest to test:

```go
func TestMyPureFunction(t *testing.T) {
    result := myPureFunction("input")
    assertion.Equal(t, result, "expected", "description")
}

// Table-driven for multiple cases
func TestMyPureFunction_Various(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"simple", "hello", "HELLO"},
        {"empty", "", ""},
        {"unicode", "日本語", "日本語"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := myPureFunction(tt.input)
            assertion.Equal(t, got, tt.want, "output")
        })
    }
}
```

## Pattern 2: Test Command Execution

### Success Path

```go
func TestMyCommand_Success(t *testing.T) {
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines("input1", "input2", "input3")

    assertion.NoError(t, result.Err)
    assertion.Lines(t, result.Stdout, []string{
        "expected output 1",
        "expected output 2",
        "expected output 3",
    })
}
```

### Empty Input

```go
func TestMyCommand_EmptyInput(t *testing.T) {
    result := run.Quick(mycommand.MyCommand())

    assertion.NoError(t, result.Err)
    assertion.Empty(t, result.Stdout)
}
```

### Single Line

```go
func TestMyCommand_SingleLine(t *testing.T) {
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines("single line")

    assertion.NoError(t, result.Err)
    assertion.Lines(t, result.Stdout, []string{"expected output"})
}
```

## Pattern 3: Test With Flags

```go
func TestMyCommand_WithFlags(t *testing.T) {
    result := run.Command(mycommand.MyCommand(
        mycommand.Flag1,
        mycommand.Flag2,
    )).WithStdinLines("input")

    assertion.NoError(t, result.Err)
    assertion.Lines(t, result.Stdout, []string{"expected with flags"})
}
```

## Pattern 4: Test With Files

```go
func TestMyCommand_WithFiles(t *testing.T) {
    // Command will open actual files
    result := run.Quick(mycommand.MyCommand(
        "testdata/file1.txt",
        "testdata/file2.txt",
    ))

    assertion.NoError(t, result.Err)
    // Verify output
}
```

## Pattern 5: Test Error Handling

### Input Errors

```go
func TestMyCommand_InputError(t *testing.T) {
    result := run.Command(mycommand.MyCommand()).
        WithStdinError(errors.New("read failed"))

    assertion.ErrorContains(t, result.Err, "read failed")
}
```

### Output Errors

```go
func TestMyCommand_OutputError(t *testing.T) {
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines("data").
        WithStdoutError(errors.New("write failed"))

    assertion.ErrorContains(t, result.Err, "write failed")
}
```

### Internal Errors

```go
func TestMyCommand_InternalError(t *testing.T) {
    // Test command's internal error handling
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines("invalid input that causes error")

    assertion.ErrorContains(t, result.Err, "expected error message")
}
```

## Pattern 6: Test Edge Cases

Always include these standard edge cases:

```go
func TestMyCommand_EmptyLines(t *testing.T) {
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines("", "", "")

    assertion.NoError(t, result.Err)
    // Verify behavior with empty lines
}

func TestMyCommand_VeryLongLine(t *testing.T) {
    longLine := strings.Repeat("a", 10000)
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines(longLine)

    assertion.NoError(t, result.Err)
    // Verify handling of long lines
}

func TestMyCommand_ManyLines(t *testing.T) {
    lines := make([]string, 1000)
    for i := range lines {
        lines[i] = fmt.Sprintf("line %d", i)
    }

    result := run.Command(mycommand.MyCommand()).
        WithStdinLines(lines...)

    assertion.NoError(t, result.Err)
    assertion.Count(t, result.Stdout, 1000)
}

func TestMyCommand_UnicodeHandling(t *testing.T) {
    result := run.Command(mycommand.MyCommand()).
        WithStdinLines("日本語", "Ελληνικά", "Русский")

    assertion.NoError(t, result.Err)
    // Verify Unicode handling
}
```

## Pattern 7: Table-Driven Tests

Always include at least one table-driven test:

```go
func TestMyCommand_TableDriven(t *testing.T) {
    tests := []struct {
        name   string
        input  []string
        flags  []any
        output []string
    }{
        {
            name:   "simple case",
            input:  []string{"a", "b"},
            flags:  []any{},
            output: []string{"A", "B"},
        },
        {
            name:   "with flag",
            input:  []string{"a", "b"},
            flags:  []any{mycommand.SomeFlag},
            output: []string{"A+", "B+"},
        },
        {
            name:   "empty input",
            input:  []string{},
            flags:  []any{},
            output: []string{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := run.Command(mycommand.MyCommand(tt.flags...)).
                WithStdinLines(tt.input...)

            assertion.NoError(t, result.Err)
            assertion.Lines(t, result.Stdout, tt.output)
        })
    }
}
```

## Pattern 8: Test Custom Types/Programs

If your command uses custom types (like awk's Program interface):

```go
// Define test implementations
type TestProgram struct {
    // Embed default implementation if available
    command.DefaultProgram
}

func (p TestProgram) CustomMethod() string {
    return "test"
}

func TestMyCommand_WithCustomType(t *testing.T) {
    result := run.Command(mycommand.MyCommand(TestProgram{})).
        WithStdinLines("input")

    assertion.NoError(t, result.Err)
    // Verify custom behavior
}
```

## Assertion Patterns

### Use Specific Assertions

```go
// ✅ Good - specific assertion
assertion.Lines(t, result.Stdout, []string{"expected"})

// ❌ Bad - manual comparison
if len(result.Stdout) != 1 || result.Stdout[0] != "expected" {
    t.Error("mismatch")
}
```

### Check Errors Properly

```go
// ✅ Good - check error presence and content
assertion.ErrorContains(t, result.Err, "expected message")

// ❌ Bad - just check presence
if result.Err == nil {
    t.Error("expected error")
}
```

### Use Helper Assertions

```go
// Check for content
assertion.Contains(t, result.Stdout, "success")
assertion.NotContains(t, result.Stderr, "error")

// Check patterns
assertion.Prefix(t, result.Stdout, "INFO:")
assertion.Suffix(t, result.Stdout, " complete")

// Check counts
assertion.Count(t, result.Stdout, 5)
assertion.Empty(t, result.Stderr)
```

## Test Naming Conventions

```go
// Format: Test<Type>_<Scenario>
TestMyCommand_Success
TestMyCommand_EmptyInput
TestMyCommand_InputError
TestMyCommand_WithFlags

// For methods: Test<Type>_<Method>
TestContext_Field
TestContext_SetField

// For variations: Test<Type>_<Method>_<Variation>
TestContext_Field_OutOfBounds
TestMyCommand_WithFlags_CaseSensitive
```

## Coverage Checklist

- [ ] All exported functions tested
- [ ] All exported methods tested
- [ ] Success path tested
- [ ] Empty input tested
- [ ] All flags/options tested
- [ ] All error paths tested (Begin, Action, End, I/O)
- [ ] Edge cases tested (empty, long, many, unicode)
- [ ] Multiple files tested (if applicable)
- [ ] All conditional branches tested
- [ ] At least one table-driven test

## Running Tests

```bash
# Run all tests
go test -v

# Check coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestMyCommand_Success
```

## Complete Example

See `awk/command_test.go` for a complete reference implementation that achieves 100% coverage with:
- 25+ test functions
- Pure function tests
- Command execution tests
- Custom program tests
- Error handling tests
- Edge case tests
- Table-driven tests

## Tips for 100% Coverage

1. **Test pure functions separately** - They're easy to test and don't need `run` package
2. **Test all error paths** - Use `WithStdinError`, `WithStdoutError`, custom error scenarios
3. **Test all flags** - Each flag should have at least one test
4. **Test edge cases** - Empty, single, many, long, unicode
5. **Use table-driven tests** - Catches more cases with less code
6. **Check coverage regularly** - `go test -cover` shows what's missing
7. **Test interfaces/contracts** - If your command uses interfaces, test all methods

## Common Mistakes to Avoid

### ❌ Don't Mix Test and Production Code

```go
// Bad - don't do this
type MyCommand struct {
    Reader TestableReader  // Don't add test hooks to production code
}
```

### ❌ Don't Create Test-Specific Command Variants

```go
// Bad - don't do this
func MyCommandForTesting(args ...any) gloo.Command {
    // Special test version
}
```

### ✅ Do Use the Command As-Is

```go
// Good - test the actual command
result := run.Command(mycommand.MyCommand())
```

### ❌ Don't Forget to Test Errors

```go
// Bad - only tests success
func TestMyCommand(t *testing.T) {
    result := run.Command(cmd).WithStdin("data")
    assertion.NoError(t, result.Err)
}
```

### ✅ Do Test All Paths

```go
// Good - tests success and errors
func TestMyCommand_Success(t *testing.T) { /* ... */ }
func TestMyCommand_InputError(t *testing.T) { /* ... */ }
func TestMyCommand_OutputError(t *testing.T) { /* ... */ }
```

## Summary

1. Organize tests into clear sections
2. Test pure functions first
3. Test command execution with `run` package
4. Test all flags and options
5. Test all error paths
6. Test edge cases (empty, long, many, unicode)
7. Include table-driven tests
8. Use specific assertions
9. Aim for 100% coverage
10. Follow the `awk/command_test.go` template

This pattern works for all gloo-foo commands because they all implement the same `gloo.Command` interface and use the same I/O abstraction.

