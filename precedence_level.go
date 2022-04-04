package main

type PrecedenceLevel struct {
	properties uint8
	// pointers to what each operator means in this precedence level
	operators map[int]*OpProp
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
	PREFIX                  = 0b00
	POSTFIX                 = 0b01
	INFIX_LEFT_ASSOCIATIVE  = 0b10
	INFIX_RIGHT_ASSOCIATIVE = 0b11
	IMPLIED_OPERATION       = 0b100
	DELIMITER               = 0b100
)

func (pl *PrecedenceLevel) OperatorExists(op *OpProp) bool {
	for _, v := range pl.operators {
		if v == op {
			return true
		}
	}
	return false
}
