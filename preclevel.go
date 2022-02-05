package main

type PrecedenceLevel uint8

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
