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

func separateFloat(s string, base int) (ml int, hasExp, manSign, expSign bool) {
	index := 0
	if s[index] == '-' || s[index] == '+' {
		manSign = true
		index++
	}
	for index < len(s) {
		if s[index] == 'e' && base <= 10 || s[index] == 'p' {
			hasExp = true
			break
		}
		ml++
		index++
	}
	index++
	if index >= len(s) {
		return
	}
	if s[index] == '-' || s[index] == '+' {
		expSign = true
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
	// before anything, check if there's a radix in here anywhere
	radixIndex := len(s)
	for i := index; i < len(s); i++ {
		if s[i] == RADIX {
			radixIndex = i
			break
		}
	}
	// this is the number of digits before the radix point
	digits := radixIndex - index
	// and this is the total number of characters ignoring leading sign
	totalCount := len(s) - index
	// if the length is odd, write in our first character
	if digits%2 == 1 {
		numBuilder.WriteByte(s[index])
		index++
	}
	byteBase := uint8(base)
	for index < totalCount {
		// the radix is here!
		if index == digits {
			numBuilder.WriteByte(s[index])
			index++
		}
		if index == len(s) {
			break
		}
		ms := s[index]
		index++
		var ls byte
		if index == len(s) {
			ls = '0'
		} else {
			ls = s[index]
		}
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
		ml, hasExp, manSign, _ := separateFloat(s, base)
		if base == 4 || base == 2 {
			var numBuilder strings.Builder
			var mantissa string
			numBuilder.WriteString("0x")
			index := ml
			if manSign {
				numBuilder.WriteByte(s[index])
				index++
				mantissa = s[1 : 1+ml]
			} else {
				mantissa = s[0:ml]
			}
			if base == 2 {
				mantissa = squareBase(mantissa, 2)
			}
			numBuilder.WriteString(squareBase(mantissa, 4))
			if hasExp {
				// we might have 'e' at this position instead of 'p'
				numBuilder.WriteByte('p')
				index++
			}
			// TODO: the exponent needs to be base 10, so either I need to convert to base 10 to use the built in atof or just write my own...
			exponent := s[index:]
			numBuilder.WriteString(exponent)
			return strconv.ParseFloat(numBuilder.String(), 64)
		}
		return 0, nil
	}
}
