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
	currToken    int
	opctx        *OpContext
	intLiteral   int64
	floatLiteral float64
}

func NewTokeniser(reader *bufio.Reader, operators *OpContext) *Tokeniser {
	result := new(Tokeniser)
	result.reader = reader
	result.opctx = operators
	result.readRune()
	result.ReadToken()
	return result
}

func (tk *Tokeniser) ReadToken() {

	// ignore comments and surrounding whitespace if present
	tk.skipWhitespace()

	// operators, symbols, almost everything
	possibleCount, branchDeducedOn := tk.opctx.opTree.PossibleCount_rune(tk.currRune)
	if possibleCount > 0 {
		for possibleCount > 0 {
			tk.readRune()
			possibleCount, branchDeducedOn = branchDeducedOn.PossibleCount_rune(tk.currRune)
		}
		tk.currToken = branchDeducedOn.operatorToken
		if tk.currToken == NIL_TOKEN {
			panic("Invalid token")
		} else
		// tokens with special meaning. Some need special care and attention.
		if tk.currToken < NIL_TOKEN {
			// skipping block comments
			if tk.currToken == OPEN_COMMENT_TOKEN {
				tk.skipUntilControl(CLOSE_COMMENT_TOKEN)
				tk.currToken = COMMENT_TOKEN
			} else
			// skipping line comments
			if tk.currToken == COMMENT_TOKEN {
				tk.skipUntilControl(NEWLINE_TOKEN)
				tk.currToken = COMMENT_TOKEN
			}
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

func (tk *Tokeniser) skipUntilControl(token int) {
	controlBit := uint16(1 << ^token)
	buff := make([]rune, 0)
	branch := tk.opctx.opTree.branches[tk.currRune]
	searching := true
	for searching {
		for branch == nil {
			tk.readRune()
			branch = tk.opctx.opTree.branches[tk.currRune]
		}
		if branch.operatorToken == token {
			searching = false
		} else if (branch.controlOps & controlBit) != 0 {
			buff = append(buff, tk.currRune)
			branch = branch.branches[tk.currRune]
		} else {
			branch = tk.opctx.opTree.branches[tk.currRune]
			for len(buff) != 0 {
				buff = buff[1:]
				found := branch.GetToken(buff)
				if token == found {
					searching = false
					break
				}
			}
			buff = make([]rune, 0)
		}
		tk.readRune()
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
