package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type Range struct {
	start int
	end   int
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	scanner.Scan()
	rangeStrings := strings.Split(scanner.Text(), ",")
	var ranges []Range
	for _, rangeString := range rangeStrings {
		rangeValues := strings.Split(rangeString, "-")
		rangeStart, err := strconv.Atoi(rangeValues[0])
		utils.Check(err)
		rangeEnd, err := strconv.Atoi(rangeValues[1])
		utils.Check(err)
		ranges = append(ranges, Range{rangeStart, rangeEnd})
	}

	fmt.Println(part1(ranges))
	fmt.Println(part2(ranges))
}

func isRepeated(number int, parts int) bool {
	numberStr := strconv.Itoa(number)
	if len(numberStr)%parts != 0 {
		return false
	}
	substrLength := len(numberStr) / parts
	pattern := numberStr[:substrLength]
	for i := 1; i < parts; i++ {
		if numberStr[i*substrLength:(i+1)*substrLength] != pattern {
			return false
		}
	}
	return true
}

func part1(ranges []Range) int {
	// Brute force works for the first part
	// O(n)
	var result int
	for _, idRange := range ranges {
		for current := idRange.start; current < idRange.end; current++ {
			if isRepeated(current, 2) {
				result += current
			}
		}
	}
	return result
}

func part2(ranges []Range) int {
	// For the second part, try all substring lengths that divide the string length
	// This is O(m^2*n) worst case, where m is the maximum number string length a
	// and n the amount of integers contained in all ranges
	var result int
	for _, idRange := range ranges {
		for current := idRange.start; current < idRange.end; current++ {
			currLength := len(strconv.Itoa(current))
			for i := 2; i <= currLength; i++ {
				if isRepeated(current, i) {
					result += current
					break
				}
			}
		}
	}
	return result
}
