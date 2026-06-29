package assertion

import (
	"errors"
	"testing"
)

// spy is a reporter that records calls instead of failing a test, so the
// failure branches of every assertion helper are reachable and verifiable.
type spy struct {
	errors      []string
	helperCalls int
}

func (s *spy) Helper() { s.helperCalls++ }

func (s *spy) Errorf(format string, args ...any) {
	s.errors = append(s.errors, Format(format, args...))
}

// failed reports whether the spy recorded any Errorf call.
func (s *spy) failed() bool { return len(s.errors) > 0 }

func TestLines_Equal(t *testing.T) {
	s := &spy{}
	Lines(s, []string{"a", "b"}, []string{"a", "b"})
	if s.failed() {
		t.Errorf("equal lines reported errors: %v", s.errors)
	}
	if s.helperCalls == 0 {
		t.Errorf("Helper was not called")
	}
}

func TestLines_CountMismatch(t *testing.T) {
	s := &spy{}
	Lines(s, []string{"a"}, []string{"a", "b"})
	if len(s.errors) != 3 {
		t.Errorf("got %d errors, want 3 (count + actual + expected)", len(s.errors))
	}
}

func TestLines_ContentMismatch(t *testing.T) {
	s := &spy{}
	Lines(s, []string{"a", "x"}, []string{"a", "b"})
	if len(s.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(s.errors))
	}
}

func TestContains(t *testing.T) {
	pass := &spy{}
	Contains(pass, []string{"hello world"}, "hello", "world")
	if pass.failed() {
		t.Errorf("present substrings reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Contains(fail, []string{"hello"}, "missing")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestNotContains(t *testing.T) {
	pass := &spy{}
	NotContains(pass, []string{"hello"}, "absent")
	if pass.failed() {
		t.Errorf("absent substring reported errors: %v", pass.errors)
	}

	fail := &spy{}
	NotContains(fail, []string{"hello"}, "ell")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestEmpty(t *testing.T) {
	pass := &spy{}
	Empty(pass, []string{})
	if pass.failed() {
		t.Errorf("empty slice reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Empty(fail, []string{"x"})
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestCount(t *testing.T) {
	pass := &spy{}
	Count(pass, []string{"a", "b"}, 2)
	if pass.failed() {
		t.Errorf("matching count reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Count(fail, []string{"a"}, 2)
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestPrefix(t *testing.T) {
	pass := &spy{}
	Prefix(pass, []string{"INFO: a", "INFO: b"}, "INFO:")
	if pass.failed() {
		t.Errorf("matching prefix reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Prefix(fail, []string{"INFO: a", "WARN: b"}, "INFO:")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestSuffix(t *testing.T) {
	pass := &spy{}
	Suffix(pass, []string{"a;", "b;"}, ";")
	if pass.failed() {
		t.Errorf("matching suffix reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Suffix(fail, []string{"a;", "b"}, ";")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestErrorContains(t *testing.T) {
	pass := &spy{}
	ErrorContains(pass, errors.New("boom failed"), "boom")
	if pass.failed() {
		t.Errorf("matching error reported errors: %v", pass.errors)
	}

	nilErr := &spy{}
	ErrorContains(nilErr, nil, "boom")
	if len(nilErr.errors) != 1 {
		t.Errorf("nil error: got %d errors, want 1", len(nilErr.errors))
	}

	mismatch := &spy{}
	ErrorContains(mismatch, errors.New("other"), "boom")
	if len(mismatch.errors) != 1 {
		t.Errorf("mismatch: got %d errors, want 1", len(mismatch.errors))
	}
}

func TestNoError(t *testing.T) {
	pass := &spy{}
	NoError(pass, nil)
	if pass.failed() {
		t.Errorf("nil error reported errors: %v", pass.errors)
	}

	fail := &spy{}
	NoError(fail, errors.New("boom"))
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestEqual(t *testing.T) {
	pass := &spy{}
	Equal(pass, 1, 1, "count")
	if pass.failed() {
		t.Errorf("equal values reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Equal(fail, 1, 2, "count")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestTrue(t *testing.T) {
	pass := &spy{}
	True(pass, true, "should hold")
	if pass.failed() {
		t.Errorf("true condition reported errors: %v", pass.errors)
	}

	fail := &spy{}
	True(fail, false, "should hold")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestFalse(t *testing.T) {
	pass := &spy{}
	False(pass, false, "should not hold")
	if pass.failed() {
		t.Errorf("false condition reported errors: %v", pass.errors)
	}

	fail := &spy{}
	False(fail, true, "should not hold")
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}

func TestFormat(t *testing.T) {
	got := Format("%s=%d", "x", 7)
	if got != "x=7" {
		t.Errorf("Format() = %q, want %q", got, "x=7")
	}
}

func TestError(t *testing.T) {
	pass := &spy{}
	Error(pass, errors.New("boom"))
	if pass.failed() {
		t.Errorf("non-nil error reported errors: %v", pass.errors)
	}

	fail := &spy{}
	Error(fail, nil)
	if len(fail.errors) != 1 {
		t.Errorf("got %d errors, want 1", len(fail.errors))
	}
}
