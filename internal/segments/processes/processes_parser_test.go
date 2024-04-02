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
	var data []Data
	data = append(data, Data{icon: "F", phrase: "firefox"})
	data = append(data, Data{icon: "X", phrase: "xxx"})

	parser := ProcessesParser{command: path.Join(basePath, "internal/test/ps_test")}
	err := parser.Parse(data)
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	if data[0].active != true {
		t.Fatalf("firefox is not active, want active")
	}

	if data[1].active != false {
		t.Fatalf("xxx is active, want not active")
	}
}
