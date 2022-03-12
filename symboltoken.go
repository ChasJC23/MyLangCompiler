package main

import (
	"strconv"
	"strings"
)

func IsSubscript(r rune) bool {
	return '\u2080' <= r && r < '\u209a'
}

func NumericSubscriptValue(r rune) int {
	return int(r - '\u2080')
}

func separateFloat(s string) (iml, fml, exl int, hasRadix, hasExp, negMan, negExp, manSign, expSign bool) {
	index := 0
	if s[index] == '-' {
		negMan = true
		manSign = true
		index++
	} else if s[index] == '+' {
		manSign = true
		index++
	}
	for index < len(s) {
		if s[index] == '.' || s[index] == 'e' {
			hasRadix = hasRadix || s[index] == '.'
			hasExp = hasExp || s[index] == 'e'
			break
		}
		iml++
		index++
	}
	if index >= len(s) {
		return
	}
	if s[index] == '.' {
		index++
	}
	for index < len(s) {
		if s[index] == 'e' {
			hasExp = true
			break
		}
		fml++
		index++
	}
	index++
	if index >= len(s) {
		return
	}
	if s[index] == '-' {
		negExp = true
		expSign = true
		index++
	} else if s[index] == '+' {
		expSign = true
		index++
	}
	for index < len(s) {
		exl++
		index++
	}
	return
}

func squareBase(s string, base int) string {
	var numBuilder strings.Builder
	index := 0
	// put in the sign if it exists
	if s[index] == '-' || s[index] == '+' {
		numBuilder.WriteByte(s[index])
		index++
	}
	digits := len(s) - index
	// if the length is odd, write in our first character
	if digits%2 == 1 {
		numBuilder.WriteByte(s[index])
		index++
	}
	byteBase := uint8(base)
	for index < digits {
		ms := s[index]
		index++
		ls := s[index]
		index++
		next := (ms-'0')*byteBase + ls
		if next > '9' {
			next = next - '9' + 'A' - 1
		}
		numBuilder.WriteByte(next)
	}
	return numBuilder.String()
}

func ParseFloat(s string, base int) (float64, error) {
	if base == 10 {
		return strconv.ParseFloat(s, 64)
	} else if base == 16 {
		return strconv.ParseFloat("0x"+s, 64)
	} else {
		// TODO: get rid of unused return arguments, separateFloat is only ever used here
		iml, fml, _, hasradix, hasexp, _, _, mansign, _ := separateFloat(s)
		if base == 4 {
			var numBuilder strings.Builder
			var intMantissa string
			numBuilder.WriteString("0x")
			index := iml
			if mansign {
				numBuilder.WriteByte(s[index])
				index++
				intMantissa = s[1 : 1+iml]
			} else {
				intMantissa = s[0:iml]
			}
			numBuilder.WriteString(squareBase(intMantissa, base))
			if hasradix {
				numBuilder.WriteByte(s[index])
				index++
			}
			fracMantissa := s[index : index+fml]
			index += fml
			// TODO: the fractional part of the mantissa needs to be treated differently, sort it out in squareBase?
			numBuilder.WriteString(squareBase(fracMantissa, base))
			if hasexp {
				numBuilder.WriteByte(s[index])
				index++
			}
			exponent := s[index:]
			numBuilder.WriteString(squareBase(exponent, base))
			return strconv.ParseFloat(numBuilder.String(), 64)
		}
		return 0, nil
	}
}
