package test

import (
	"os"
	"testing"
)

// NewTempFile creates a file which is closed and removed when the test finishes.
func NewTempFile(t testing.TB) *os.File {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })

	return f
}

// WriteLine replaces the content of the file with s followed by a newline.
func WriteLine(t testing.TB, f *os.File, s string) {
	t.Helper()

	err := os.WriteFile(f.Name(), []byte(s+"\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
