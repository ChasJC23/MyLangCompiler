package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
)

var srcFilePath = flag.Arg(0)

func init() {
	flag.Parse()
	if !fs.ValidPath(srcFilePath) {
		srcFilePath = "example.idk"
	}
}

func main() {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(srcFile)
	opc := NewOpContext()
	opc.opTree.AddOperator([]rune(";"), COMMENT_TOKEN)
	opc.opTree.AddOperator([]rune("<--"), OPEN_COMMENT_TOKEN)
	opc.opTree.AddOperator([]rune("-->"), CLOSE_COMMENT_TOKEN)
	fmt.Println(opc.opTree.ToString(true))
	tokeniser := NewTokeniser(reader, opc)
	for tokeniser.currToken != EOF_TOKEN {
		fmt.Println(tokeniser.comment)
		tokeniser.ReadToken()
	}
	defer srcFile.Close()
}
