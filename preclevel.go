package main

import "math/bits"

type PrecedenceLevel struct {
	properties uint8
	// a bitvector representing which operators exist within this precedence level
	operators []uint
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

func getbit(bitvec []uint, bit int) uint {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> 6
	result := bitvec[index] & mask
	return result >> (bit & 0b111111)
}

func setbit(bitvec []uint, bit int) {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> (5 + bits.UintSize>>6)
	bitvec[index] |= mask
}

func resetbit(bitvec []uint, bit int) {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> (5 + bits.UintSize>>6)
	bitvec[index] &^= mask
}
