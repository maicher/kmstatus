package processes

import (
	"os"
	"path"
	"testing"
)

func Test_ProcessesParser(t *testing.T) {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		panic("BASE_PATH not set")
	}
	var d []data
	d = append(d, data{icon: "F", phrase: "firefox"})
	d = append(d, data{icon: "X", phrase: "xxx"})

	parser := Parser{command: path.Join(basePath, "internal/test/ps_test")}
	err := parser.Parse(d)
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	if d[0].active != true {
		t.Fatalf("firefox is not active, want active")
	}

	if d[1].active != false {
		t.Fatalf("xxx is active, want not active")
	}
}
