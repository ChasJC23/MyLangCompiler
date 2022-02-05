package main

type Token int

const (
	EOF           Token = 0
	INT_LITERAL   Token = 1
	FLOAT_LITERAL Token = 2
)

const (
	PREPROCESSOR        = '#'
	LINE_COMMENT        = ';'
	START_BLOCK_COMMENT = '['
	STOP_BLOCK_COMMENT  = ']'
	NEW_LINE            = '\n'
	RADIX               = '.'
	EOF_RUNE            = '\000'
)
