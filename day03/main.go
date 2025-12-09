package main

import (
	"fmt"
	"strconv"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var result1, result2 int
	for scanner.Scan() {
		bank := scanner.Text()
		result1 += part1(bank)
		result2 += part2(bank)
	}

	fmt.Println(result1)
	fmt.Println(result2)
}

func part1(bank string) int {
	var first int
	for i := 1; i < len(bank)-1; i++ {
		if bank[i] > bank[first] {
			first = i
		}
	}

	second := first + 1
	for i := first + 2; i < len(bank); i++ {
		if bank[i] > bank[second] {
			second = i
		}
	}

	result, err := strconv.Atoi(string(bank[first]) + string(bank[second]))
	utils.Check(err)
	return result
}

func part2(bank string) int {
	lastMaxIndex := -1
	var digits []int = make([]int, 12)
	for i := range 12 {
		digits[i] = lastMaxIndex + 1 // Initialize to first valid index in range
		for j := lastMaxIndex + 1; j <= len(bank)-12+i; j++ {
			if bank[j] > bank[digits[i]] {
				digits[i] = j
			}
		}
		lastMaxIndex = digits[i]
	}
	var result int
	for _, digit := range digits {
		newDigit, err := strconv.Atoi(string(bank[digit]))
		utils.Check(err)
		result = result*10 + newDigit
	}
	return result
}
