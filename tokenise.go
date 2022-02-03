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
	PREPROCESSOR        = '#'
	LINE_COMMENT        = ';'
	START_BLOCK_COMMENT = '['
	STOP_BLOCK_COMMENT  = ']'
	NEW_LINE            = '\n'
)

type Tokeniser struct {
	reader    *bufio.Reader
	currRune  rune
	currToken int
	operators OpContext
}

func (tokeniser *Tokeniser) ReadToken() {

	// ignore comments and surrounding whitespace if present
	tokeniser.skipComments()

	// detect EOF
	if tokeniser.currRune == '\000' {
		tokeniser.currToken = EOF
		return
	}

	// TODO: literals

	// operators

}

func (tokeniser *Tokeniser) readRune() {
	currentRune, _, err := tokeniser.reader.ReadRune()
	if err == io.EOF {
		tokeniser.currRune = '\000'
	} else if err == nil {
		tokeniser.currRune = currentRune
	} else {
		// FUCK FUCK FUCK
		panic(err)
	}
}

func (tokeniser *Tokeniser) skipWhitespace() {
	for unicode.IsSpace(tokeniser.currRune) {
		tokeniser.readRune()
	}
}

func (tokeniser *Tokeniser) skipComments() {

	tokeniser.skipWhitespace()

	// while a comment exists
	for tokeniser.currRune == START_BLOCK_COMMENT || tokeniser.currRune == LINE_COMMENT {

		// remove block comments
		if tokeniser.currRune == START_BLOCK_COMMENT {
			for tokeniser.currRune != STOP_BLOCK_COMMENT {
				tokeniser.readRune()
			}
		} else
		// remove line comments
		if tokeniser.currRune == LINE_COMMENT {
			for tokeniser.currRune != NEW_LINE {
				tokeniser.readRune()
			}
		}

		// remove succeeding whitespace ready to check for more comments
		// or start reading source code instead
		tokeniser.skipWhitespace()
	}
}
