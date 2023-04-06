package filesystem

import (
	"os"
	"path"
	"testing"
)

func Test_FSParser(t *testing.T) {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		panic("BASE_PATH not set")
	}

	parser := FSParser{path: path.Join(basePath, "pkg/test/df_test")}
	val, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	fs := val.(FS)

	if l := len(fs.Drives); l != 4 {
		t.Fatalf("Drives count: %d, want: 4", l)
	}

}
