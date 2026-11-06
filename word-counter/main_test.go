package main

import (
	"strings"
	"testing"
)

func TestCountWords(t *testing.T) {
	got, err := Count(strings.NewReader("Hello world! This is a test."), Words)
	if err != nil {
		t.Fatal(err)
	}
	if exp := 6; got != exp {
		t.Errorf("CountWords returned %d, expected %d", got, exp)
	}
}

func TestCountLines(t *testing.T) {
	got, err := Count(strings.NewReader("1 2 3 4 5 6 7\nline2\nline3 8"), Lines)
	if err != nil {
		t.Fatal(err)
	}
	if exp := 3; got != exp {
		t.Errorf("CountLines returned %d, expected %d", got, exp)
	}
}

func TestCountBytes(t *testing.T) {
	got, err := Count(strings.NewReader("héllo\nworld"), Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if exp := 12; got != exp {
		t.Errorf("CountBytes returned %d, expected %d", got, exp)
	}
}
