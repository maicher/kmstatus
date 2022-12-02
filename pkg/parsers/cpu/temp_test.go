package cpu

import (
	"os"
	"testing"

	"github.io/maicher/stbar/pkg/test"
)

func Test_TempParser_Parse_FileCanNotBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "bla bla")

	parser := TempParser{files: []*os.File{f}}
	_, err := parser.Parse()

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func Test_TempParser_Parse_FileCanBeParsed(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, "30000")

	parser := TempParser{files: []*os.File{f}}
	temp, err := parser.Parse()

	if val := temp[0]; val != 30 {
		t.Fatalf("Temp equals: %d, want: 30", val)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}
