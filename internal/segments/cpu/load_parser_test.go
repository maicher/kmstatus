package cpu

import (
	"fmt"
	"testing"

	"github.com/maicher/kmstatus/internal/test"
)

func TestLoadParser_Parse_FileCanNotBeParsed(t *testing.T) {
	var load float64
	f := test.NewTempFile()
	test.WriteLine(f, "cpx  1171962 591604 506805 67668597")

	parser := LoadParser{statFile: f}
	err := parser.Parse(&load)

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func TestLoadParser_Parse_FileCanBeParsed(t *testing.T) {
	var load float64
	f := test.NewTempFile()
	test.WriteLine(f, "cpu  1171962 591604 506805 67668597")

	parser := LoadParser{statFile: f}
	parser.Parse(&load)

	test.WriteLine(f, "cpu  1172013 591650 506843 67673329")
	err := parser.Parse(&load)

	if l := fmt.Sprintf("%.1f", load); l != "2.8" {
		t.Fatalf("Load equals: %s, want: 2.8", l)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}

func TestLoadParser_Parse_IowaitIsIdle(t *testing.T) {
	var load float64
	f := test.NewTempFile()
	test.WriteLine(f, "cpu  0 0 0 0 0 0 0 0 0 0")

	parser := LoadParser{statFile: f}
	parser.Parse(&load)

	// 10 user, 10 irq, 40 idle, 40 iowait
	test.WriteLine(f, "cpu  10 0 0 40 40 10 0 0 0 0")
	err := parser.Parse(&load)

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
	if load != 20 {
		t.Fatalf("Load equals: %f, want: 20", load)
	}
}

func TestLoadParser_Parse_NoTimePassed(t *testing.T) {
	load := 5.0
	f := test.NewTempFile()
	test.WriteLine(f, "cpu  10 0 0 40 0 0 0 0 0 0")

	parser := LoadParser{statFile: f, stat: stat{active: 10, idle: 40}}
	parser.Parse(&load)

	if load != 5 {
		t.Fatalf("Load equals: %f, want: previous value 5", load)
	}
}
