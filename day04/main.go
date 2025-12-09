package main

import (
	"fmt"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var grid [][]rune
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]rune, 0)
		for _, char := range line {
			row = append(row, char)
		}
		grid = append(grid, row)
	}

	fmt.Println(part1(grid))
	fmt.Println(part2(grid))
}

func countNeighbors(grid [][]rune, row, col int) int {
	diffs := []int{-1, 0, 1}
	var result int
	for _, rowDiff := range diffs {
		searchRow := row + rowDiff
		if searchRow < 0 || searchRow > len(grid)-1 {
			continue
		}
		for _, colDiff := range diffs {
			searchCol := col + colDiff
			if searchCol < 0 || searchCol > len(grid[0])-1 {
				continue
			}
			if searchRow == row && searchCol == col {
				continue
			}
			if grid[searchRow][searchCol] == '@' {
				result++
			}
		}
	}
	return result
}

func part1(grid [][]rune) int {
	// Counts the number of items in the grid with fewer than four roll neighbors
	var result int
	rows := len(grid)
	cols := len(grid[0])
	for i := range rows {
		for j := range cols {
			if grid[i][j] == '@' && countNeighbors(grid, i, j) < 4 {
				result++
			}
		}
	}
	return result
}

func part2(grid [][]rune) int {
	rows := len(grid)
	cols := len(grid[0])
	var removedRolls int
	var stop bool = false
	for !stop {
		stop = true
		for i := range rows {
			for j := range cols {
				if grid[i][j] == '@' && countNeighbors(grid, i, j) < 4 {
					grid[i][j] = 'x'
					removedRolls++
				}
			}
		}
		for i := range rows {
			for j := range cols {
				if grid[i][j] == 'x' {
					grid[i][j] = '.'
					stop = false
				}
			}
		}
	}
	return removedRolls
}
