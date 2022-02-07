package main

const (
	EOF              = ^0
	INT_LITERAL      = ^1
	FLOAT_LITERAL    = ^2
	OPEN_PARENS      = ^3
	CLOSE_PARENS     = ^4
	OPEN_CODE_BLOCK  = ^5
	CLOSE_CODE_BLOCK = ^6
	OPEN_COMMENT     = ^7
	CLOSE_COMMENT    = ^8
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
