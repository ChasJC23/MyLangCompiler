package main

import (
	"flag"
	"fmt"
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
	ctx.AddOperator([]rune("++"), 0)
	ctx.AddOperator([]rune("+"), 0)
	ctx.AddOperator([]rune("-"), 0)
	ctx.AddOperator([]rune("--"), 0)

	fmt.Println(ctx.PossibleCount([]rune("")))

	// using runes
	// srcReader := bufio.NewReader(srcFile)
}
