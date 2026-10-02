package temperature

import (
	"testing"

	"github.com/maicher/kmstatus/internal/test"
)

func Test_TempParser_Parse_FileCanNotBeParsed(t *testing.T) {
	data := make([]data, 1)
	f := test.NewTempFile(t)
	test.WriteLine(t, f, "bla bla")

	parser := Parser{}
	parser.sensors = append(parser.sensors, sensor{file: f})
	err := parser.Parse(data)

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
}

func Test_TempParser_Parse_FileCanBeParsed(t *testing.T) {
	data := make([]data, 1)
	f := test.NewTempFile(t)
	test.WriteLine(t, f, "30000")

	parser := Parser{}
	parser.sensors = append(parser.sensors, sensor{file: f})
	err := parser.Parse(data)

	if val := data[0].Value; val != 30 {
		t.Fatalf("Temp equals: %d, want: 30", val)
	}

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
}

func Test_TempParser_Parse_OneSensorFails(t *testing.T) {
	data := make([]data, 2)
	broken := test.NewTempFile(t)
	test.WriteLine(t, broken, "bla bla")
	ok := test.NewTempFile(t)
	test.WriteLine(t, ok, "42000")

	parser := Parser{}
	parser.sensors = append(parser.sensors, sensor{file: broken}, sensor{file: ok})
	err := parser.Parse(data)

	if err == nil {
		t.Fatalf("Error nil, want: error")
	}
	if val := data[1].Value; val != 42 {
		t.Fatalf("Temp equals: %d, want: 42", val)
	}
}

func Test_SortNumerically(t *testing.T) {
	paths := []string{"/z/thermal_zone10/temp", "/z/thermal_zone2/temp", "/z/thermal_zone1/temp"}
	sortNumerically(paths)

	if paths[0] != "/z/thermal_zone1/temp" || paths[2] != "/z/thermal_zone10/temp" {
		t.Fatalf("Order: %v, want: zone1, zone2, zone10", paths)
	}
}
