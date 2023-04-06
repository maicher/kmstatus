package processes

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

	parser := ProcessesParser{path: path.Join(basePath, "pkg/test/ps_test")}
	val, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	ps := val.(PS)

	if l := ps.Uniq(); l != 14 {
		t.Fatalf("Processes count: %d, want: 14", l)
	}

	if c := ps.Find("firefox"); c != 1 {
		t.Fatalf("Processes firefox count: %d, want: 1", c)
	}

	if c := ps.FindByPrefix("chrome"); c != 2 {
		t.Fatalf("Processes chrome count: %d, want: 2", c)
	}
}
