package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()
	w.Close()
	os.Stdout = oldOut

	var buf bytes.Buffer
	buf.ReadFrom(r)

	if buf.String() != "Hello, World!\n" {
		t.Error("Not Hello, World! ", buf.String())
	}

}
