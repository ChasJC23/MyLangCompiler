package main

type OpProp struct {
	subsequentSymbols   []int
	argumentPrecedences []*PrecedenceLevel
	requireParens       uint
	argumentCount       int
	initSymbol          string // for debugging purposes
}
