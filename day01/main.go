package main

import (
	"fmt"
	"strconv"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

func main() {
	lines := utils.ParseInput("input.txt")
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}

func part1(lines []string) int {
	var dial int = 50
	var zeros int
	for _, line := range lines {
		direction := line[0]
		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}
		// Modulo logic
		if direction == 'L' {
			dial = (dial - steps%100 + 100) % 100
		} else {
			dial = (dial + steps) % 100
		}
		if dial == 0 {
			zeros++
		}
	}
	return zeros
}

func part2(lines []string) int {
	var dial int = 50
	var zeros int
	for _, line := range lines {
		direction := line[0]
		steps, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}
		if direction == 'L' {
			if dial == 0 {
				zeros += steps / 100
			} else {
				zeros += (steps + 100 - dial) / 100
			}
			dial = (dial - steps%100 + 100) % 100
		} else {
			zeros += (dial + steps) / 100
			dial = (dial + steps) % 100
		}
	}
	return zeros
}
