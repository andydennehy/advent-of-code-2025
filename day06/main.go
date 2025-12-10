package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	// Parse keeping raw lines for part2 (preserving spaces)
	var rawLines []string
	for scanner.Scan() {
		rawLines = append(rawLines, scanner.Text())
	}

	// Parse with Fields for part1 (original behavior)
	var lines [][]string
	for _, raw := range rawLines {
		lines = append(lines, strings.Fields(raw))
	}

	fmt.Println(part1(lines))
	fmt.Println(part2(rawLines))
}

func part1(lines [][]string) int {
	var result int
	for col := 0; col < len(lines[0]); col++ {
		operation := lines[len(lines)-1][col]
		var localResult int
		if operation == "*" {
			localResult = 1
		}

		for row := 0; row < len(lines)-1; row++ {
			operand, err := strconv.Atoi(lines[row][col])
			utils.Check(err)
			switch operation {
			case "+":
				localResult = localResult + operand
			case "*":
				localResult = localResult * operand
			}
		}
		result += localResult
	}
	return result
}

func part2(rawLines []string) int {
	maxLen := 0
	for _, line := range rawLines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Pad lines to equal length
	grid := make([]string, len(rawLines))
	for i, line := range rawLines {
		grid[i] = line + strings.Repeat(" ", maxLen-len(line))
	}

	operatorRow := len(grid) - 1
	var result int
	var problems [][]int

	var currentProblem []int
	for col := 0; col < maxLen; col++ {
		isSeparator := true
		for row := 0; row < operatorRow; row++ {
			if grid[row][col] != ' ' {
				isSeparator = false
				break
			}
		}

		if isSeparator {
			if len(currentProblem) > 0 {
				problems = append(problems, currentProblem)
				currentProblem = nil
			}
		} else {
			currentProblem = append(currentProblem, col)
		}
	}
	if len(currentProblem) > 0 {
		problems = append(problems, currentProblem)
	}

	for _, cols := range problems {
		var operation byte
		for _, col := range cols {
			if grid[operatorRow][col] != ' ' {
				operation = grid[operatorRow][col]
				break
			}
		}

		var localResult int
		if operation == '*' {
			localResult = 1
		}

		for i := len(cols) - 1; i >= 0; i-- {
			col := cols[i]
			var numStr string
			for row := 0; row < operatorRow; row++ {
				ch := grid[row][col]
				if ch != ' ' {
					numStr += string(ch)
				}
			}

			if numStr == "" { continue }

			num, err := strconv.Atoi(numStr)
			utils.Check(err)

			switch operation {
			case '+':
				localResult += num
			case '*':
				localResult *= num
			}
		}

		result += localResult
	}

	return result
}
