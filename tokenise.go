package main

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
)

type Tokeniser struct {
	reader       *bufio.Reader
	currRune     rune
	currToken    Token
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

func (tk *Tokeniser) ReadToken() {

	// ignore comments and surrounding whitespace if present
	tk.skipComments()

	// detect EOF
	if tk.currRune == '\000' {
		tk.currToken = EOF
		return
	}

	// operators
	possibleCount, branchDeducedOn := tk.operators.opTree.PossibleCount_rune(tk.currRune)
	if possibleCount > 0 {
		for possibleCount > 0 {
			tk.readRune()
			possibleCount, branchDeducedOn = branchDeducedOn.PossibleCount_rune(tk.currRune)
		}
		tk.currToken = branchDeducedOn.operatorToken
		if tk.currToken == -1 {
			panic("Invalid token")
		}
		return
	}

	// numeric literals
	if unicode.IsNumber(tk.currRune) || tk.currRune == RADIX {

		// set up with first character
		builderSlice := make([]rune, 1, 32)
		builderSlice[0] = tk.currRune
		hasRadix := tk.currRune == RADIX
		tk.readRune()

		// add following characters
		for unicode.IsNumber(tk.currRune) || !hasRadix && tk.currRune == RADIX {
			builderSlice = append(builderSlice, tk.currRune)
			hasRadix = hasRadix || tk.currRune == RADIX
			tk.readRune()
		}

		// put things in the right places
		var err error
		if hasRadix {
			tk.floatLiteral, err = strconv.ParseFloat(string(builderSlice), 64)
			tk.currToken = FLOAT_LITERAL
		} else {
			tk.intLiteral, err = strconv.ParseInt(string(builderSlice), 0, 64)
			tk.currToken = INT_LITERAL
		}
		if err != nil {
			panic("poorly formatted number")
		}
	}
}

func (tk *Tokeniser) readRune() {
	currentRune, _, err := tk.reader.ReadRune()
	if err == io.EOF {
		tk.currRune = '\000'
	} else if err == nil {
		tk.currRune = currentRune
	} else {
		panic(err)
	}
}

func (tk *Tokeniser) skipWhitespace() {
	for unicode.IsSpace(tk.currRune) {
		tk.readRune()
	}
}

func (tk *Tokeniser) skipComments() {

	tk.skipWhitespace()

	// while a comment exists
	for tk.currRune == START_BLOCK_COMMENT || tk.currRune == LINE_COMMENT {

		// remove block comments
		if tk.currRune == START_BLOCK_COMMENT {
			for tk.currRune != STOP_BLOCK_COMMENT && tk.currRune != EOF_RUNE {
				tk.readRune()
			}
			tk.readRune()
		} else
		// remove line comments
		if tk.currRune == LINE_COMMENT {
			for tk.currRune != NEW_LINE && tk.currRune != EOF_RUNE {
				tk.readRune()
			}
		}

		// remove succeeding whitespace ready to check for more comments
		// or start reading source code instead
		tk.skipWhitespace()
	}
}
