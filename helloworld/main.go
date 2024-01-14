package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	//counts := make(map[string]int)

	// countLines(os.Stdin, counts)

	// for line, n := range counts {
	// 	if n > 1 {
	// 		fmt.Printf("%d\t%s\n", n, line)
	// 	}
	// }

	ReadFile("inputfile.txt")
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
}

func ReadFile(filename string) {
	readFile, err := os.Open(filename)
	check(err)
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	counts := make(map[string]int)
	bigInt := 0
	for fileScanner.Scan() {
		counts[fileScanner.Text()] = findFirstDigit(fileScanner.Text())*10 + findLastDigit(fileScanner.Text())
		bigInt += counts[fileScanner.Text()]
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}

	fmt.Printf("\n\nTotal:%d", bigInt)
}

func findFirstDigit(str string) int {
	for _, char := range str {
		if unicode.IsDigit(char) {
			return int(char - '0')
		}
	}
	return 0
}

func findLastDigit(str string) int {
	last := len(str) - 1
	for i := range str {
		char := rune(str[last-i])
		if unicode.IsDigit(char) {
			return int(char - '0')
		}
	}
	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
