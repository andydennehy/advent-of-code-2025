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

func ParseInput(filename string) []string {
	file, err := os.Open(filename)
	Check(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
