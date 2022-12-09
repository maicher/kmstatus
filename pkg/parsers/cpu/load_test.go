package cpu

import (
	"fmt"
	"testing"

	"github.io/maicher/kmstatus/pkg/test"
)

func TestLoadParser_Parse_FileCanNotBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "cpx  1171962 591604 506805 67668597")

	parser := LoadParser{statFile: f}
	_, err := parser.Parse()

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func TestLoadParser_Parse_FileCanBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "cpu  1171962 591604 506805 67668597")

	parser := LoadParser{statFile: f}
	parser.Parse()

	test.WriteLine(f, "cpu  1172013 591650 506843 67673329")
	v, err := parser.Parse()
	load := v.(Load)

	if l := fmt.Sprintf("%.1f", load); l != "2.8" {
		t.Fatalf("Load equals: %s, want: 2.8", l)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}
