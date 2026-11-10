package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var binPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "todo-test")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binPath = filepath.Join(dir, "todo")
	if out, err := exec.Command("go", "build", "-o", binPath, ".").CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n%s", err, out)
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func run(t *testing.T, dir string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = dir
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return out.String(), errOut.String(), code
}

func mustRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	stdout, stderr, code := run(t, dir, args...)
	if code != 0 {
		t.Fatalf("%v: exit %d: %s", args, code, stderr)
	}
	return stdout
}

func TestAddAndList(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "-add", "buy", "milk")
	mustRun(t, dir, "-add", "walk dog")

	got := mustRun(t, dir, "-list")
	want := "1. [ ] buy milk\n2. [ ] walk dog\n"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestListEmpty(t *testing.T) {
	if got := mustRun(t, t.TempDir(), "-list"); got != "" {
		t.Errorf("expected no output, got %q", got)
	}
}

func TestComplete(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "-add", "a")
	mustRun(t, dir, "-add", "b")
	mustRun(t, dir, "-complete", "2")

	got := mustRun(t, dir, "-list")
	want := "1. [ ] a\n2. [x] b\n"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "-add", "a")
	mustRun(t, dir, "-add", "b")
	mustRun(t, dir, "-delete", "1")

	got := mustRun(t, dir, "-list")
	want := "1. [ ] b\n"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestNotFound(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "-add", "a")

	for _, flag := range []string{"-complete", "-delete"} {
		stdout, stderr, code := run(t, dir, flag, "5")
		if code != 1 {
			t.Errorf("%s: expected exit 1, got %d", flag, code)
		}
		if stdout != "" {
			t.Errorf("%s: expected empty stdout, got %q", flag, stdout)
		}
		if stderr == "" {
			t.Errorf("%s: expected error on stderr", flag)
		}
	}
}

func TestAddWithoutTask(t *testing.T) {
	dir := t.TempDir()
	_, stderr, code := run(t, dir, "-add")
	if code != 2 {
		t.Errorf("expected exit 2, got %d", code)
	}
	if stderr == "" {
		t.Errorf("expected error on stderr")
	}
	if _, err := os.Stat(filepath.Join(dir, todoFileName)); err == nil {
		t.Errorf("todo file should not have been created")
	}
}

func TestNoFlags(t *testing.T) {
	_, stderr, code := run(t, t.TempDir())
	if code != 2 {
		t.Errorf("expected exit 2, got %d", code)
	}
	if stderr == "" {
		t.Errorf("expected usage error on stderr")
	}
}

func TestCorruptFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, todoFileName), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, stderr, code := run(t, dir, "-list")
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
	if stderr == "" {
		t.Errorf("expected error on stderr")
	}
}
