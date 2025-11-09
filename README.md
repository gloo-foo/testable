# testable

Testing utilities for gloo-foo commands. Test commands without touching actual files or I/O.

> **✨ Immutable Design**: All configuration methods return new `Runner` instances for safety and predictability.
> Remember to call `.Run()` to execute! (or use `Quick()`/`WithInput()` shortcuts)

## Quick Start

```go
import (
    "testing"
    "github.com/gloo-foo/testable/run"
    "github.com/gloo-foo/testable/assertion"
)

// Your command (normal gloo-foo command - no changes needed!)
func MyCommand(args ...any) gloo.Command {
    inputs := gloo.Initialize[gloo.File, Flags](args...)
    return command(inputs)
}

// Test it
func TestMyCommand(t *testing.T) {
    result := run.Command(MyCommand()).
        WithStdinLines("input1", "input2").
        Run()

    assertion.NoError(t, result.Err)
    assertion.Lines(t, result.Stdout, []string{
        "processed input1",
        "processed input2",
    })
}
```

**That's it!** Your command doesn't change - gloo-foo already abstracts I/O. This package just makes testing easy.

## Why This Works

Gloo-foo commands receive `(ctx, stdin, stdout, stderr)` from the framework. They're already testable! This package provides:

1. **`run`** - Execute commands with mock stdin, capture stdout/stderr
2. **`assertion`** - Helpful test assertions
3. **`splitter`** - Field splitting utilities (bonus for awk-style commands)

## Package: `run`

Execute commands and capture results:

```go
// Basic execution
result := run.Command(MyCommand()).WithStdin("test input").Run()

// Multiple lines
result := run.Command(MyCommand()).WithStdinLines("line1", "line2").Run()

// No input (Quick executes immediately)
result := run.Quick(MyCommand())

// With context
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
result := run.WithContext(ctx, MyCommand()).WithStdin("data").Run()

// Inject errors for testing
result := run.Command(MyCommand()).WithStdinError(errors.New("disk error")).Run()
result := run.Command(MyCommand()).WithStdoutError(errors.New("write failed")).Run()

// Result contains:
result.Stdout  // []string - captured output lines
result.Stderr  // []string - captured error lines
result.Err     // error - error returned by command
```

### Immutable Design

All `With*` methods return a **new** `Runner` with updated configuration. This makes the API:
- **Safe**: No mutation of shared state
- **Reusable**: Store base configurations and derive variations
- **Predictable**: Each call is independent

```go
// Base configuration
base := run.Command(MyCommand()).WithContext(ctx)

// Derive different test cases
test1 := base.WithStdinLines("case1").Run()
test2 := base.WithStdinLines("case2").Run()
// base is unchanged!
```

**Terminal Methods** (execute immediately):
- `.Run()` - Executes the configured command
- `Quick(cmd)` - Shorthand for `Command(cmd).Run()`
- `WithInput(cmd, "data")` - Shorthand for `Command(cmd).WithStdin("data").Run()`

## Package: `assertion`

Clear test assertions:

```go
// Compare lines
assertion.Lines(t, result.Stdout, []string{"expected", "output"})

// Check content
assertion.Contains(t, result.Stdout, "success")
assertion.NotContains(t, result.Stderr, "error")

// Counts
assertion.Count(t, result.Stdout, 5)
assertion.Empty(t, result.Stderr)

// Patterns
assertion.Prefix(t, result.Stdout, "INFO:")
assertion.Suffix(t, result.Stdout, " done")

// Errors
assertion.NoError(t, result.Err)
assertion.ErrorContains(t, result.Err, "file not found")

// Generic
assertion.Equal(t, actual, expected, "description")
assertion.True(t, condition, "should be true")
```

## Package: `splitter`

Field splitting utilities for awk-style commands:

```go
// Whitespace (collapses runs, like awk)
fields := splitter.Whitespace("a  b   c", " ")  // ["a", "b", "c"]

// Exact separator
fields := splitter.Exact("a,b,c", ",")  // ["a", "b", "c"]

// Fixed-width
split := splitter.Fixed(0, 10, 20)
fields := split("John      Smith     NYC       ", "")

// Character class
fields := splitter.CharacterClass("a,b;c:d", ",;:")  // ["a", "b", "c", "d"]
```

## Common Patterns

### Success Path
```go
func TestMyCommand_Success(t *testing.T) {
    result := run.Command(MyCommand()).
        WithStdinLines("input1", "input2")

    assertion.NoError(t, result.Err)
    assertion.Lines(t, result.Stdout, []string{
        "processed input1",
        "processed input2",
    })
}
```

### Empty Input
```go
func TestMyCommand_EmptyInput(t *testing.T) {
    result := run.Quick(MyCommand())

    assertion.NoError(t, result.Err)
    assertion.Empty(t, result.Stdout)
}
```

### Error Handling
```go
func TestMyCommand_InputError(t *testing.T) {
    result := run.Command(MyCommand()).
        WithStdinError(errors.New("read failed"))

    assertion.ErrorContains(t, result.Err, "read failed")
}
```

### With Files
```go
func TestMyCommand_WithFiles(t *testing.T) {
    // Command opens actual files (gloo-foo handles it)
    result := run.Quick(MyCommand("testdata/file1.txt", "testdata/file2.txt"))

    assertion.NoError(t, result.Err)
    // verify output
}
```

### Table-Driven Tests
```go
func TestMyCommand_Various(t *testing.T) {
    tests := []struct {
        name   string
        input  []string
        output []string
    }{
        {"simple", []string{"a", "b"}, []string{"A", "B"}},
        {"empty", []string{}, []string{}},
        {"unicode", []string{"日本語"}, []string{"日本語"}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := run.Command(MyCommand()).
                WithStdinLines(tt.input...)

            assertion.NoError(t, result.Err)
            assertion.Lines(t, result.Stdout, tt.output)
        })
    }
}
```

### Test Flags
```go
func TestMyCommand_WithFlags(t *testing.T) {
    result := run.Command(MyCommand(Flag1, Flag2)).
        WithStdin("input")

    assertion.NoError(t, result.Err)
    // verify output with flags
}
```

## Complete Example

```go
// mycommand.go
package mycommand

import (
    "bufio"
    "context"
    "io"
    "strings"
    gloo "github.com/gloo-foo/framework"
)

type Flags struct {
    Uppercase bool
}

type UppercaseFlag bool
const (
    Lowercase UppercaseFlag = false
    Uppercase UppercaseFlag = true
)
func (f UppercaseFlag) Configure(flags *Flags) {
    flags.Uppercase = bool(f)
}

type command gloo.Inputs[gloo.File, Flags]

func Transform(args ...any) gloo.Command {
    inputs := gloo.Initialize[gloo.File, Flags](args...)
    return command(inputs)
}

func (c command) Executor() gloo.CommandExecutor {
    inputs := gloo.Inputs[gloo.File, Flags](c)
    return inputs.Wrap(func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
        scanner := bufio.NewScanner(stdin)
        for scanner.Scan() {
            line := scanner.Text()
            if inputs.Flags.Uppercase {
                line = strings.ToUpper(line)
            }
            if _, err := stdout.Write([]byte(line + "\n")); err != nil {
                return err
            }
        }
        return scanner.Err()
    })
}
```

```go
// mycommand_test.go
package mycommand_test

import (
    "testing"
    "github.com/gloo-foo/testable/run"
    "github.com/gloo-foo/testable/assertion"
    "mycommand"
)

func TestTransform_Uppercase(t *testing.T) {
    result := run.Command(mycommand.Transform(mycommand.Uppercase)).
        WithStdinLines("hello", "world")

    assertion.NoError(t, result.Err)
    assertion.Lines(t, result.Stdout, []string{"HELLO", "WORLD"})
}

func TestTransform_EmptyInput(t *testing.T) {
    result := run.Quick(mycommand.Transform())

    assertion.NoError(t, result.Err)
    assertion.Empty(t, result.Stdout)
}

func TestTransform_WithFiles(t *testing.T) {
    result := run.Quick(mycommand.Transform("testdata/input.txt"))

    assertion.NoError(t, result.Err)
    // Check output
}
```

## Best Practices

### ✅ DO: Test Through gloo.Command

Commands already use gloo-foo's I/O abstraction - just test them directly:

```go
result := run.Command(MyCommand())
```

### ❌ DON'T: Rewrite Commands for Testing

Don't change your command structure. If gloo-foo handles I/O, you're already testable:

```go
// ❌ Don't do this
type MyExecutor struct {
    Reader CustomReader  // Unnecessary!
    Writer CustomWriter  // gloo-foo already handles this!
}

// ✅ Do this - use gloo-foo as designed
func (c command) Executor() gloo.CommandExecutor {
    return c.inputs.Wrap(func(ctx, stdin, stdout, stderr) error {
        // Your logic here
    })
}
```

### ✅ DO: Test Pure Functions Separately

Extract business logic into pure functions and test them directly:

```go
// Pure function - no I/O
func processLine(line string) string {
    return strings.ToUpper(line)
}

// Test without any mocking
func TestProcessLine(t *testing.T) {
    got := processLine("hello")
    assertion.Equal(t, got, "HELLO", "output")
}
```

### ✅ DO: Use Table-Driven Tests

```go
tests := []struct {
    name  string
    input string
    want  string
}{
    {"simple", "hello", "HELLO"},
    {"empty", "", ""},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := run.Command(Cmd()).WithStdin(tt.input)
        assertion.Lines(t, result.Stdout, []string{tt.want})
    })
}
```

## Coverage

```bash
# Generate coverage
go test -coverprofile=coverage.out

# View in terminal
go tool cover -func=coverage.out

# View in browser
go tool cover -html=coverage.out
```

### Coverage Checklist

- [ ] Success path
- [ ] Empty input
- [ ] All flags/options
- [ ] Error paths (stdin/stdout failures)
- [ ] Edge cases
- [ ] Multiple files (if applicable)
- [ ] All conditional branches

## How It Works

1. `run.Command()` takes your `gloo.Command`
2. Calls its `Executor()` method with mock stdin and capture buffers
3. Returns `Result` with captured stdout/stderr and any error

Your command doesn't change - it's already testable because gloo-foo abstracts I/O.

## Other Packages

This repo also contains:
- `stream` - Lower-level input abstractions (if you need custom readers)
- `capture` - Lower-level output abstractions (if you need custom writers)
- `scenario` - Test scenario builders (alternative to `run`)

These are useful if you're building custom test utilities, but **for testing gloo-foo commands, just use `run`**.

## License

MIT — see [LICENSE](LICENSE).
