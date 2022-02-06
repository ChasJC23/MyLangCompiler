package main

import (
	"flag"
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

	bitvec := []uint64{0, 0, 0}

	setbit(bitvec, 4)
	setbit(bitvec, 2)
	setbit(bitvec, 2)
	resetbit(bitvec, 2)

	// using runes
	// srcReader := bufio.NewReader(srcFile)
}
