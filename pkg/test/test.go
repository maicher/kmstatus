package test

import (
	"io/ioutil"
	"log"
	"os"
)

func NewTempFile() *os.File {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func WriteLine(f *os.File, s string) {
	ioutil.WriteFile(f.Name(), []byte(s+"\n"), 0644)
}
