package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"
	// "io/fs"
	// "os"
)

// var srcFilePath = flag.Arg(0)

func init() {
	flag.Parse()
	// if !fs.ValidPath(srcFilePath) {
	// 	fmt.Println("Path does not point to a valid source file.")
	// 	os.Exit(2)
	// }
}

func main() {
	// srcFile, err := os.Open(srcFilePath)
	// if err != nil {
	// 	panic(err)
	// }
	// defer srcFile.Close()

	ctx := NewOpContext()
	ops := []string{"+", "-", "*", "/", "=", "~", "!=", ">", "<", ">=", "<=", "&", "|", "^"}
	for _, v := range ops {
		ctx.AddOperator([]rune(v))
	}

	fmt.Println(ctx.opTree.ToString(false))

	reader := bufio.NewReader(strings.NewReader("; time to do a bit of maths \n [ my favourite number ] 5 + [ no you're a joke] 1 * -4 [inequality]>= 1.1;all done!"))
	tokeniser := NewTokeniser(reader, ctx)

	for tokeniser.currToken != EOF {
		fmt.Println(tokeniser.currToken)
		tokeniser.ReadToken()
	}

	// using runes
	// srcReader := bufio.NewReader(srcFile)
}
