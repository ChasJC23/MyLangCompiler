package main

// standard tokens
const (
	INIT_TOKEN             = 1
	NIL_TOKEN              = 0
	EOF_TOKEN              = -1
	COMMENT_TOKEN          = -2
	CHAR_LITERAL           = -3
	STRING_LITERAL         = -4
	INT_LITERAL            = -5
	FLOAT_LITERAL          = -6
	IDENTIFIER_TOKEN       = -7
	OPEN_PARENS_TOKEN      = -8
	CLOSE_PARENS_TOKEN     = -9
	OPEN_CODE_BLOCK_TOKEN  = -10
	CLOSE_CODE_BLOCK_TOKEN = -11
	NEWLINE_TOKEN          = -12
	STATEMENT_ENDING_TOKEN = -13
	TRUE_LITERAL           = -14
	FALSE_LITERAL          = -15
)

// control flags
const (
	OPEN_COMMENT_FLAG  uint = 1 << 0
	CLOSE_COMMENT_FLAG uint = 1 << 1
	COMMENT_FLAG       uint = 1 << 2
	NEWLINE_FLAG       uint = 1 << 3
	OPEN_CHAR_FLAG     uint = 1 << 4
	CLOSE_CHAR_FLAG    uint = 1 << 5
	OPEN_STRING_FLAG   uint = 1 << 6
	CLOSE_STRING_FLAG  uint = 1 << 7
)

const (
	RADIX        = '.'
	EOF_RUNE     = '\000'
	NEWLINE_RUNE = '\n'
)
