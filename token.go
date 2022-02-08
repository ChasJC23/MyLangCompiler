package main

const (
	INIT_TOKEN          = 1
	NIL_TOKEN           = 0
	EOF_TOKEN           = -1
	INT_LITERAL         = -2
	FLOAT_LITERAL       = -3
	OPEN_PARENS         = -4
	CLOSE_PARENS        = -5
	OPEN_CODE_BLOCK     = -6
	CLOSE_CODE_BLOCK    = -7
	OPEN_COMMENT_TOKEN  = -8
	CLOSE_COMMENT_TOKEN = -9
	COMMENT_TOKEN       = -10
	NEWLINE_TOKEN       = -11
)

const (
	RADIX        = '.'
	EOF_RUNE     = '\000'
	NEWLINE_RUNE = '\n'
)
