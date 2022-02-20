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
	// let's try parsing a C program!
	opc.AddFixedTokenOperator([]rune("//"), COMMENT_TOKEN, COMMENT_FLAG)
	opc.AddControlOperator([]rune("/*"), OPEN_COMMENT_FLAG)
	opc.AddControlOperator([]rune("*/"), CLOSE_COMMENT_FLAG)
	opc.AddFixedTokenOperator([]rune("{"), OPEN_CODE_BLOCK_TOKEN, 0)
	opc.AddFixedTokenOperator([]rune("}"), CLOSE_CODE_BLOCK_TOKEN, 0)
	opc.AddFixedTokenOperator([]rune("("), OPEN_PARENS_TOKEN, 0)
	opc.AddFixedTokenOperator([]rune(")"), OPEN_PARENS_TOKEN, 0)
	opc.AddFixedTokenOperator([]rune(";"), STATEMENT_ENDING_TOKEN, 0)
	fmt.Println(opc.opTree.String(true))
	tokeniser := NewTokeniser(reader, opc)
	for tokeniser.currToken != EOF_TOKEN {
		fmt.Println(tokeniser.comment)
		tokeniser.ReadToken()
	}
	defer srcFile.Close()
}
