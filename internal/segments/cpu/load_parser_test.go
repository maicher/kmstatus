package cpu

import (
	"fmt"
	"testing"

	"github.com/maicher/kmst/internal/test"
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
