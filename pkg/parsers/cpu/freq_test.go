package cpu

import (
	"os"
	"testing"

	"github.io/maicher/stbar/pkg/test"
)

func Test_FreqParser_Parse_FileCanNotBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "")

	parser := FreqParser{files: []*os.File{f}}
	_, err := parser.Parse()

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func Test_FreqParser_Parse_FileCanBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "1200000")

	parser := FreqParser{files: []*os.File{f}}
	v, err := parser.Parse()
	freq := v.(Freq)

	if freq != 1200000 {
		t.Fatalf("Freq equals: %d, want: 1200000", freq)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}
