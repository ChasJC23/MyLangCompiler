package main

type PrecedenceLevel struct {
	properties uint8
	// pointers to what each operator means in this precedence level
	operators   map[int]*OpProp
	defaultNext *PrecedenceLevel
}

/*
precedence level bitmask:

------00 prefix
------01 postfix
------10 infix left associative
------11 infix right associative
-----X-- implied operation
*/

const (
	PREFIX                  = 0b000
	POSTFIX                 = 0b001
	INFIX_LEFT_ASSOCIATIVE  = 0b010
	INFIX_RIGHT_ASSOCIATIVE = 0b011
	IMPLIED_OPERATION       = 0b100
	DELIMITER               = 0b101
)

func (pl *PrecedenceLevel) OperatorExists(op *OpProp) bool {
	for _, v := range pl.operators {
		if v == op {
			return true
		}
	}
	return false
}
