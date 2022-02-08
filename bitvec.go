package main

import "math/bits"

func Getbit(bitvec []uint, bit int) uint {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> 6
	result := bitvec[index] & mask
	return result >> (bit & 0b111111)
}

func Setbit(bitvec []uint, bit int) {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> (5 + bits.UintSize>>6)
	bitvec[index] |= mask
}

func Resetbit(bitvec []uint, bit int) {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> (5 + bits.UintSize>>6)
	bitvec[index] &^= mask
}
