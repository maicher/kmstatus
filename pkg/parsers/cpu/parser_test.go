package cpu

import (
	"fmt"
	"testing"

	"github.io/maicher/stbar/pkg/test"
)

func TestParser_Parse_FileCanNotBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "cpx  1171962 591604 506805 67668597")

	parser := Parser{file: f}
	_, err := parser.Parse()

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func TestParser_Parse_FileCanBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "cpu  1171962 591604 506805 67668597")

	parser := Parser{file: f}
	parser.Parse()

	test.WriteLine(f, "cpu  1172013 591650 506843 67673329")
	cpu, err := parser.Parse()

	if load := fmt.Sprintf("%.1f", cpu.Load()); load != "2.8" {
		t.Fatalf("Load equals: %s, want: 2.8", load)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}
