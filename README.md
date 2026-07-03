# testable

[![CI](https://github.com/gloo-foo/testable/actions/workflows/ci.yml/badge.svg)](https://github.com/gloo-foo/testable/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/gloo-foo/testable)](https://goreportcard.com/report/github.com/gloo-foo/testable) [![Go Reference](https://pkg.go.dev/badge/github.com/gloo-foo/testable.svg)](https://pkg.go.dev/github.com/gloo-foo/testable) [![License](https://img.shields.io/github/license/gloo-foo/testable)](https://github.com/gloo-foo/testable/blob/main/LICENSE)

Shared test harness for [gloo-foo](https://github.com/gloo-foo) commands. It runs a `Command[[]byte, []byte]` against in-memory input and returns the captured output, so a command is tested without touching real files or I/O.

Runnable, compiler-checked usage lives in the package examples (`go test`-verified, rendered by `go doc`); this README describes what each package is for. Run `go doc github.com/gloo-foo/testable` (and the subpackages) for the example source.

## Model

A gloo-foo command is a `Command[[]byte, []byte]`: it consumes a stream of input lines and produces a stream of output lines. The harness feeds it input, collects the output stream, and hands back the lines plus the first error the stream carried. There is a single output stream — the command's output — so the harness captures one stream, not a separate stdout/stderr pair.

## Packages

### `testable` — the entry points

The convenience functions every `cmd-*` module uses. Both run the command and collect its output:

- `Test(cmd, input) (string, error)` — output as one string, each line terminated by `\n`.
- `TestLines(cmd, input) ([]string, error)` — output as a slice of lines.

On command failure each returns the error (and, respectively, an empty string or a nil slice). Both delegate execution to `run`, so behaviour is identical to the fluent API below.

### `run` — the fluent runner

Finer control over execution. `run.Command(cmd)` returns an immutable `Runner`; every `With*` method returns a **new** `Runner`, so a base configuration can be shared and derived without mutation. Configuration is applied, then `Run()` executes and returns a `*Result` (`Stdout []string`, `Err error`).

Input configuration: `WithStdin` (a string), `WithStdinLines` (lines joined and newline-terminated), `WithStdinReader` (any `io.Reader`), `WithStdinError` (fail the input stream with a given error, to exercise error propagation), and `WithContext` (bind a `context.Context`). Shortcuts that execute immediately: `Quick(cmd)` (empty stdin) and `WithInput(cmd, stdin)`.

### `assertion` — test assertions

Assertion helpers that take a `*testing.T` (or any value with `Helper`/`Errorf`) and report differences with readable diffs:

| Function | Checks |
|---|---|
| `Lines(t, actual, expected)` | line-by-line equality, reporting the first differing line |
| `Contains(t, actual, want...)` / `NotContains(t, actual, unexpected...)` | substring presence / absence across the joined output |
| `Count(t, actual, n)` / `Empty(t, actual)` | exact line count / no output |
| `Prefix(t, actual, prefix)` / `Suffix(t, actual, suffix)` | every line starts / ends with a string |
| `Equal[T](t, actual, expected, label)` | equality of two comparable values |
| `True(t, cond, msg)` / `False(t, cond, msg)` | boolean conditions |
| `NoError(t, err)` / `Error(t, err)` / `ErrorContains(t, err, want)` | error presence / absence / message substring |

### `splitter` — field splitting

Field-splitting utilities for awk-style commands: `Whitespace` (collapses runs of whitespace when the separator is `" "`, otherwise splits on the exact separator), `Exact` (always the literal separator), `CharacterClass` (split on any character in a set), and `Fixed` / `Pattern` (fixed-width positions, or a caller-supplied split function).

## Quality

Every package holds 100% statement coverage and passes the shared gomatic gate (`gofumpt`, `go vet`, `staticcheck`, `golangci-lint`, `govulncheck`, `goreleaser check`). Run it with `make check`.

## License

MIT — see [LICENSE](LICENSE).
