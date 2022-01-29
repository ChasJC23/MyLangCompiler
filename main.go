package main

import (
	"flag"
	"fmt"
)

// import "io"
// import "strings"

var nFlag = flag.Int("n", 1234, "just gimme a number")

func init() {
	flag.Parse()
}

func main() {
	fmt.Println(*nFlag)
}
