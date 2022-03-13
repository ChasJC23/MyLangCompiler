package main

import "math/bits"

func GetBit(bitVec []uint, bit int) uint {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> 6
	result := bitVec[index] & mask
	return result >> (bit & 0b111111)
}

func SetBit(bitVec []uint, bit int) {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> (5 + bits.UintSize>>6)
	bitVec[index] |= mask
}

func ResetBit(bitVec []uint, bit int) {
	var mask uint = 1 << (bit & 0b111111)
	index := bit >> (5 + bits.UintSize>>6)
	bitVec[index] &^= mask
}
