package main

type PrecedenceLevel struct {
}

// might be able to use a bitmask...

/*

precedence level bitmask:
xxxxyyyyzzzzwwww
			ww00 prefix
			ww01 postfix
			ww10 infix left associative
			ww11 infix right associative
			w1ww repeatable
*/
