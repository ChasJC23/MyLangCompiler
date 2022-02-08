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
-----X-- repeatable
----X--- implied operation
XXXX---- argument count - 1
*/
