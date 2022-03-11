package main

import "strconv"

func IsSubscript(r rune) bool {
	return '\u2080' <= r && r < '\u209a'
}

func NumericSubscriptValue(r rune) int {
	return int(r - '\u2080')
}

func SeparateFloat(s string) (int, int, int) {
	// we first group together pairs of symbols
	intMantissaLen := 0
	fracMantissaLen := 0
	exponentLen := 0
	index := 0
	for index < len(s) {
		if s[index] == '.' || s[index] == 'e' {
			break
		}
		intMantissaLen++
		index++
	}
	if index < len(s) {
		if s[index] != 'e' {
			index++
		}
		for index < len(s) {
			if s[index] == 'e' {
				break
			}
			fracMantissaLen++
			index++
		}
		index++
		for index < len(s) {
			exponentLen++
			index++
		}
	}
	return intMantissaLen, fracMantissaLen, exponentLen
}

func ParseFloat(s string, base int) (float64, error) {
	if base == 10 {
		return strconv.ParseFloat(s, 64)
	} else if base == 16 {
		return strconv.ParseFloat("0x"+s, 64)
	} else {
		// we first group together pairs of symbols
		intMantissaLen := 0
		index := 0
		for s[index] != '.' && index < len(s) {
			intMantissaLen++
			index++
		}
		fracMantissaLen := 0
		index++
		for s[index] != 'e' && index < len(s) {
			fracMantissaLen++
			index++
		}
		exponentLen := 0
		index++
		for index < len(s) {
			exponentLen++
			index++
		}
		if base == 8 {

		}
		return 0, nil
	}
}
