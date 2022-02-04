package main

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
)

const (
	EOF           = 0
	INT_LITERAL   = 1
	FLOAT_LITERAL = 2
)

const (
	PREPROCESSOR        = '#'
	LINE_COMMENT        = ';'
	START_BLOCK_COMMENT = '['
	STOP_BLOCK_COMMENT  = ']'
	NEW_LINE            = '\n'
	RADIX               = '.'
)

type Tokeniser struct {
	reader       *bufio.Reader
	currRune     rune
	currToken    int
	operators    *OpContext
	intLiteral   int64
	floatLiteral float64
}

func NewTokeniser(reader *bufio.Reader, operators *OpContext) *Tokeniser {
	result := new(Tokeniser)
	result.reader = reader
	result.operators = operators
	result.readRune()
	result.ReadToken()
	return result
}

func (tokeniser *Tokeniser) ReadToken() {

	// ignore comments and surrounding whitespace if present
	tokeniser.skipComments()

	// detect EOF
	if tokeniser.currRune == '\000' {
		tokeniser.currToken = EOF
		return
	}

	// operators
	possibleCount, branchDeducedOn := tokeniser.operators.tree.PossibleCount_rune(tokeniser.currRune)
	if possibleCount > 0 {
		for ; possibleCount > 0; possibleCount, branchDeducedOn = branchDeducedOn.PossibleCount_rune(tokeniser.currRune) {
			tokeniser.readRune()
		}
		token := branchDeducedOn.operatorToken
		if token == -1 {
			panic("Invalid token")
		}
		tokeniser.currToken = token
		return
	}

	// numeric literals
	if unicode.IsNumber(tokeniser.currRune) || tokeniser.currRune == RADIX {
		builderSlice := make([]rune, 1, 32)
		builderSlice[0] = tokeniser.currRune
		hasRadix := tokeniser.currRune == RADIX
		tokeniser.readRune()
		for unicode.IsNumber(tokeniser.currRune) || !hasRadix && tokeniser.currRune == RADIX {
			builderSlice = append(builderSlice, tokeniser.currRune)
			hasRadix = hasRadix || tokeniser.currRune == RADIX
			tokeniser.readRune()
		}
		var err error
		if hasRadix {
			tokeniser.floatLiteral, err = strconv.ParseFloat(string(builderSlice), 64)
			tokeniser.currToken = FLOAT_LITERAL
		} else {
			tokeniser.intLiteral, err = strconv.ParseInt(string(builderSlice), 0, 64)
			tokeniser.currToken = INT_LITERAL
		}
		if err != nil {
			panic("poorly formatted number")
		}
	}
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
