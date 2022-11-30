package cpu

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParser_Parse_FileDoesNotExist(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}

	writeStat(f, "cpu  1171962 591604 506805 67668597")

	parser := Parser{filePath: "some path"}
	_, err = parser.Parse()

	if err == nil {
		t.Fatalf("Empty error")
	}
}

func TestParser_Parse_FileCanNotBeParsed(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}

	writeStat(f, "cpx  1171962 591604 506805 67668597")

	parser := Parser{filePath: f.Name()}
	_, err = parser.Parse()

	if err == nil {
		t.Fatalf("Empty error")
	}
}

func TestParser_Parse_FileCanBeParsed(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}

	writeStat(f, "cpu  1171962 591604 506805 67668597")

	parser := Parser{filePath: f.Name()}
	parser.Parse()

	writeStat(f, "cpu  1172013 591650 506843 67673329")
	cpu, err := parser.Parse()

	if load := fmt.Sprintf("%.1f", cpu.Load()); load != "2.8" {
		t.Fatalf("Load equals: %s, want: 2.8", load)
	}
}

func writeStat(f *os.File, s string) {
	ioutil.WriteFile(f.Name(), []byte(s), 0644)
}
