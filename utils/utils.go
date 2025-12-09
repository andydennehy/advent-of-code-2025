package utils

import (
	"bufio"
	"os"
)

// Check panics if the error is not nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseInput(filename string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(filename)
	Check(err)	
	scanner := bufio.NewScanner(file)
	return scanner, file
}
