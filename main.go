package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"unicode"
)

const (
	EOF = 0
)

type Tokeniser interface {
	ReadToken()
}

func (reader buffio.Reader) ReadToken() int {
	r, size, err := reader.ReadRune()

	// ignore whitespace
	for unicode.IsSpace(r) {
		r, size, err = reader.ReadRune()
	}

	// if we've reached the end of the file
	if r == io.EOF {
		return EOF
	}
}

var srcFilePath = flag.Arg(0)
var fileSystem = os.DirFS(os.Getwd())

func init() {
	flag.Parse()
	if !fs.ValidPath(srcFlag) {
		fmt.Println("Path does not point to a valid source file.")
		os.Exit(2)
	}
}

func main() {
	srcFile := os.Open(srcFilePath)
	defer srcFile.Close()

	// using runes
	srcReader := bufio.NewReader(srcFile)

	r, size, err := srcReader.ReadRune()
}
