package temperature

import (
	"testing"

	"github.com/maicher/kmstatus/internal/test"
)

func Test_TempParser_Parse_FileCanNotBeParsed(t *testing.T) {
	data := make([]data, 1)
	f := test.NewTempFile()
	test.WriteLine(f, "bla bla")

	parser := TemperatureParser{}
	parser.sensors = append(parser.sensors, sensor{file: f})
	err := parser.Parse(data)

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func Test_TempParser_Parse_FileCanBeParsed(t *testing.T) {
	data := make([]data, 1)
	f := test.NewTempFile()
	test.WriteLine(f, "30000")

	parser := TemperatureParser{}
	parser.sensors = append(parser.sensors, sensor{file: f})
	err := parser.Parse(data)

	if val := data[0].Value; val != 30 {
		t.Fatalf("Temp equals: %d, want: 30", val)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}
