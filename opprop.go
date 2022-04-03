package main

type OpProp struct {
	subsequentSymbols  []int
	codeBlockArguments uint
	requireParens      uint
	argumentCount      int
	initSymbol         string // for debugging purposes
}
