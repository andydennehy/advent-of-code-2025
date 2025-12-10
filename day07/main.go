package main

import (
	"fmt"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type Pos struct {
	row int
	col int
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var start Pos
	splitters := make(map[Pos]bool)
	var row int
	for scanner.Scan() {
		for col, char := range scanner.Text() {
			if char == 'S' {
				start = Pos{row, col}
			} // Assuming well behaved input
			if char == '^' {
				splitters[Pos{row, col}] = true
			}
		}
		row++
	}

	fmt.Println(part1(start, splitters, row))
	fmt.Println(part2(start, splitters, row))
}

func part1(start Pos, splitters map[Pos]bool, height int) int {
	var splits int

	stack := []Pos{start}
	visited := make(map[Pos]bool) // Marks places where beams have been
	for len(stack) > 0 {
		popped := stack[0]
		currentPos := Pos{popped.row, popped.col}
		stack = stack[1:]

		if visited[currentPos] {
			continue
		}

		if splitters[currentPos] {
			left := Pos{currentPos.row, currentPos.col - 1}
			right := Pos{currentPos.row, currentPos.col + 1}
			splits++
			if !visited[left] { stack = append(stack, left) }
			if !visited[right] { stack = append(stack, right) }
		} else {
			nextPos := Pos{currentPos.row + 1, currentPos.col}
			if nextPos.row < height-1 && !visited[nextPos] {
				stack = append(stack, nextPos)
			}
		}
		visited[currentPos] = true
	}

	return splits
}

func countTimelines(currentPos Pos, splitters map[Pos]bool, memo map[Pos]int, height int) int {
	if val, ok := memo[currentPos]; ok {
		return val
	}

	nextPos := Pos{currentPos.row + 1, currentPos.col}
	if nextPos.row >= height {
		return 1
	} // Base case

	var result int
	if splitters[currentPos] {
		left := Pos{currentPos.row, currentPos.col - 1}
		right := Pos{currentPos.row, currentPos.col + 1}
		result = countTimelines(left, splitters, memo, height) + countTimelines(right, splitters, memo, height)
	} else {
		result = countTimelines(nextPos, splitters, memo, height)
	}

	memo[currentPos] = result
	return result
}

func part2(start Pos, splitters map[Pos]bool, height int) int {
	memo := make(map[Pos]int)
	return countTimelines(start, splitters, memo, height)
}
