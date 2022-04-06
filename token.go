package main

// standard tokens
const (
	INIT_TOKEN       = 1
	NIL_TOKEN        = 0
	EOF_TOKEN        = -1
	COMMENT_TOKEN    = -2
	CHAR_LITERAL     = -3
	STRING_LITERAL   = -4
	INT_LITERAL      = -5
	FLOAT_LITERAL    = -6
	TRUE_LITERAL     = -7
	FALSE_LITERAL    = -8
	IDENTIFIER_TOKEN = -9
	NEWLINE_TOKEN    = -10
	MAX_PARENS_TOKEN = -11
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
	OPEN_PARENS_MOD_2  = 0
	CLOSE_PARENS_MOD_2 = 1
)

const (
	RADIX        = '.'
	EOF_RUNE     = '\000'
	NEWLINE_RUNE = '\n'
)
