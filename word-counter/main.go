package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Mode int

const (
	Words Mode = iota
	Lines
	Bytes
)

func (m Mode) split() bufio.SplitFunc {
	switch m {
	case Lines:
		return bufio.ScanLines
	case Bytes:
		return bufio.ScanBytes
	default:
		return bufio.ScanWords
	}
}

func main() {
	lines := flag.Bool("l", false, "count lines")
	bytes := flag.Bool("b", false, "count bytes")
	flag.Parse()

	if *lines && *bytes {
		fmt.Fprintln(os.Stderr, "-l and -b are mutually exclusive")
		os.Exit(2)
	}

	mode := Words
	if *lines {
		mode = Lines
	} else if *bytes {
		mode = Bytes
	}

	n, err := Count(os.Stdin, mode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(n)
}

func Count(r io.Reader, mode Mode) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(mode.split())

	n := 0
	for scanner.Scan() {
		n++
	}
	return n, scanner.Err()
}
