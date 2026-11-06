package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("Hello world! This is a test.")
	count := Count(b, false)
	exp := 6
	if count != exp {
		t.Errorf("CountWords returned %d, expected %d", count, exp)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("1 2 3 4 5 6 7\nline2\nline3 8")
	exp := 3
	res := Count(b, true)
	if res != exp {
		t.Errorf("CountLines returned %d, expected %d", res, exp)
	}
}
