package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Tokeniser struct {
	reader        *bufio.Reader
	currRune      rune
	currToken     int
	opctx         *OpContext
	intLiteral    int64
	floatLiteral  float64
	charLiteral   rune
	stringLiteral string
	comment       string
	identifier    string
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

	// ignore whitespace
	tk.skipWhitespace()

	// just in case operator parsing fails partway through, we should at least try to treat it like a variable
	var identifierBuilder strings.Builder
	validIdentifier := true

	// operators, symbols, etc.
	possibleCount, branchDeducedOn := tk.opctx.opTree.PossibleCountRune(tk.currRune)
	if possibleCount > 0 {
		for possibleCount > 0 {
			if validIdentifier {
				if unicode.IsLetter(tk.currRune) || unicode.IsDigit(tk.currRune) {
					identifierBuilder.WriteRune(tk.currRune)
				} else {
					validIdentifier = false
				}
			}
			tk.readRune()
			possibleCount, branchDeducedOn = branchDeducedOn.PossibleCountRune(tk.currRune)
		}
		tk.currToken = branchDeducedOn.operatorToken
		controlOps := branchDeducedOn.controlOps
		// tokens with special meaning. Some need special care and attention.
		if controlOps != 0 {
			// skipping block comments
			if controlOps&OPEN_COMMENT_FLAG != 0 {
				tk.comment = tk.skipUntilControl(CLOSE_COMMENT_FLAG)
				tk.currToken = COMMENT_TOKEN
			} else
			// skipping line comments
			if controlOps&COMMENT_FLAG != 0 {
				tk.comment = tk.skipUntilControl(NEWLINE_FLAG)
			} else
			// parsing characters
			if controlOps&OPEN_CHAR_FLAG != 0 {
				// TODO: escaped code points? In any case, a character isn't always represented by itself in code (plus this could be improved anyways)
				charContent := tk.skipUntilControl(CLOSE_CHAR_FLAG)
				tk.charLiteral = []rune(charContent)[0]
				tk.currToken = CHAR_LITERAL
			} else
			// parsing strings
			if controlOps&OPEN_STRING_FLAG != 0 {
				tk.stringLiteral = tk.skipUntilControl(CLOSE_STRING_FLAG)
				tk.currToken = STRING_LITERAL
			}
		} else if tk.currToken != NIL_TOKEN {
			return
		}
	}

	// numeric literals
	if unicode.IsDigit(tk.currRune) || tk.currRune == RADIX {

		// set up with first character
		var litBuilder strings.Builder
		litBuilder.WriteRune(tk.currRune)
		hasRadix := tk.currRune == RADIX
		tk.readRune()

		// add following characters
		for unicode.IsDigit(tk.currRune) || !hasRadix && tk.currRune == RADIX {
			litBuilder.WriteRune(tk.currRune)
			hasRadix = hasRadix || tk.currRune == RADIX
			tk.readRune()
		}

		// check for subscript base notation
		var base int
		if IsSubscript(tk.currRune) {
			base = 0
			for IsSubscript(tk.currRune) {
				base *= 10
				base += NumericSubscriptValue(tk.currRune)
				tk.readRune()
			}
		} else {
			base = 10
		}

		// put things in the right places
		var err error
		if hasRadix {
			tk.floatLiteral, err = ParseFloat(litBuilder.String(), 64)
			tk.currToken = FLOAT_LITERAL
		} else {
			tk.intLiteral, err = strconv.ParseInt(litBuilder.String(), base, 64)
			tk.currToken = INT_LITERAL
		}
		if err != nil {
			panic("poorly formatted number")
		}
		return
	}

	// anything else has to be an identifier
	for unicode.IsLetter(tk.currRune) || unicode.IsDigit(tk.currRune) {
		identifierBuilder.WriteRune(tk.currRune)
		tk.readRune()
	}
	tk.currToken = IDENTIFIER_TOKEN
	tk.identifier = identifierBuilder.String()
}

func (tk *Tokeniser) skipUntilControl(controlBit uint) string {
	buff := make([]rune, 0)
	branch := tk.opctx.opTree.branches[tk.currRune]
	searching := true
	var builder strings.Builder
	depth := 1
	for searching {
		for branch == nil {
			length, _ := builder.WriteRune(tk.currRune)
			tk.readRune()
			branch = tk.opctx.opTree.branches[tk.currRune]
			depth = length
		}
		if (branch.controlOps & controlBit) != 0 {
			searching = false
		} else if (branch.childControlOps & controlBit) != 0 {
			length, _ := builder.WriteRune(tk.currRune)
			buff = append(buff, tk.currRune)
			branch = branch.branches[tk.currRune]
			depth += length
		} else {
			branch = tk.opctx.opTree.branches[tk.currRune]
			for len(buff) != 0 {
				buff = buff[1:]
				found := branch.GetBranch(buff)
				if (found.controlOps & controlBit) != 0 {
					searching = false
					break
				}
			}
			length, _ := builder.WriteRune(tk.currRune)
			depth = length
		}
		tk.readRune()
	}
	result := builder.String()
	return result[:len(result)-depth]
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
