package main

type Token int

const (
	EOF           Token = -1
	INT_LITERAL   Token = -2
	FLOAT_LITERAL Token = -3
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
