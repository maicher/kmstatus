package test

import (
	"log"
	"os"
)

func NewTempFile() *os.File {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func WriteLine(f *os.File, s string) {
	os.WriteFile(f.Name(), []byte(s+"\n"), 0644)
}
