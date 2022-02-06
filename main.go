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
	var r rune
	for r, _, err = reader.ReadRune(); err == nil; r, _, err = reader.ReadRune() {
		fmt.Println(r)
	}
	defer srcFile.Close()
}
