package main

import (
	"bufio"
	"io"
	"unicode"
)

const (
	EOF = 0
)

const (
	PREPROCESSOR = "#"
)

type Tokeniser struct {
	reader    *bufio.Reader
	currRune  rune
	currToken int
	operators interface{}
}

func ReadToken(reader *bufio.Reader) int {
	currentRune, _, err := reader.ReadRune()
	// ignore whitespace
	for unicode.IsSpace(currentRune) {
		currentRune, _, err = reader.ReadRune()
	}

	// if we've reached the end of the file
	if err == io.EOF {
		return EOF
	} else if err != nil {
		// FUCK FUCK FUCK
		panic(err)
	}

	return 1
}
