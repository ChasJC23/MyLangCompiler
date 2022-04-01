package main

import (
	"bufio"
	"fmt"
	"os"
)

var srcFilePath string // = flag.Arg(0)
var opc *OpContext

func init() {
	//flag.Parse()
	//if !fs.ValidPath(srcFilePath) {
	srcFilePath = "example.idk"
	//}
}

func main() {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(srcFile)
	tokeniser := NewTokeniser(reader, opc)
	parser := NewParser(tokeniser)
	ast := parser.ParseSource()
	fmt.Println(ast)
	err = srcFile.Close()
	if err != nil {
		panic(err)
	}
}
