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

	// operators, symbols, etc.
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
				tk.comment = tk.skipUntilControl(CLOSE_COMMENT_TOKEN)
				tk.currToken = COMMENT_TOKEN
			} else
			// skipping line comments
			if tk.currToken == COMMENT_TOKEN {
				tk.comment = tk.skipUntilControl(NEWLINE_TOKEN)
			} else
			// parsing characters
			if tk.currToken == OPEN_CHAR_LITERAL {
				// TODO: escaped code points? In any case, a character isn't always represented by itself in code
				// (plus this could be improved anyways)
				charContent := tk.skipUntilControl(CLOSE_CHAR_LITERAL)
				tk.charLiteral = []rune(charContent)[0]
				tk.currToken = CHAR_LITERAL
			} else
			// parsing strings
			if tk.currToken == OPEN_STRING_LITERAL {
				tk.stringLiteral = tk.skipUntilControl(CLOSE_STRING_LITERAL)
				tk.currToken = STRING_LITERAL
			}
		}
		return
	}

	// numeric literals
	if unicode.IsNumber(tk.currRune) || tk.currRune == RADIX {

		// set up with first character
		var litBuilder strings.Builder
		litBuilder.WriteRune(tk.currRune)
		hasRadix := tk.currRune == RADIX
		tk.readRune()

		// add following characters
		for unicode.IsNumber(tk.currRune) || !hasRadix && tk.currRune == RADIX {
			litBuilder.WriteRune(tk.currRune)
			hasRadix = hasRadix || tk.currRune == RADIX
			tk.readRune()
		}

		// put things in the right places
		var err error
		if hasRadix {
			tk.floatLiteral, err = strconv.ParseFloat(litBuilder.String(), 64)
			tk.currToken = FLOAT_LITERAL
		} else {
			tk.intLiteral, err = strconv.ParseInt(litBuilder.String(), 0, 64)
			tk.currToken = INT_LITERAL
		}
		if err != nil {
			panic("poorly formatted number")
		}
		return
	}

	// anything else has to be an identifier
	var idenBuilder strings.Builder
	for !unicode.IsSpace(tk.currRune) {
		idenBuilder.WriteRune(tk.currRune)
	}
	tk.currToken = IDENTIFIER_TOKEN
	tk.identifier = idenBuilder.String()
}

func (tk *Tokeniser) skipUntilControl(token int) string {
	controlBit := uint32(1 << ^token)
	buff := make([]rune, 0)
	branch := tk.opctx.opTree.branches[tk.currRune]
	searching := true
	var builder strings.Builder
	depth := 1
	for searching {
		for branch == nil {
			len, _ := builder.WriteRune(tk.currRune)
			tk.readRune()
			branch = tk.opctx.opTree.branches[tk.currRune]
			depth = len
		}
		if branch.operatorToken == token {
			searching = false
		} else if (branch.controlOps & controlBit) != 0 {
			len, _ := builder.WriteRune(tk.currRune)
			buff = append(buff, tk.currRune)
			branch = branch.branches[tk.currRune]
			depth += len
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
			len, _ := builder.WriteRune(tk.currRune)
			depth = len
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
