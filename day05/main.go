package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type Range struct {
	start, end int
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var freshRanges []Range
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
		line := strings.Split(scanner.Text(), "-")
		start, err := strconv.Atoi(line[0])
		utils.Check(err)
		end, err := strconv.Atoi(line[1])
		utils.Check(err)
		freshRanges = append(freshRanges, Range{start: start, end: end})
	}

	scanner.Scan()
	var ingredientIds []int
	for scanner.Scan() {
		id, err := strconv.Atoi(scanner.Text())
		utils.Check(err)
		ingredientIds = append(ingredientIds, id)
	}

	fmt.Println(part1(freshRanges, ingredientIds))
	fmt.Println(part2(freshRanges))
}

func part1(freshRanges []Range, ingredientIds []int) int {
	// Brute force
	var result int
	for _, id := range ingredientIds {
		for _, idRange := range freshRanges {
			if id >= idRange.start && id <= idRange.end {
				result++
				break
			}
		}
	}
	return result
}

func sortAndMerge(ranges []Range) []Range {
	if len(ranges) == 0 {
		return ranges
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	merged := []Range{ranges[0]}

	for _, r := range ranges[1:] {
		last := &merged[len(merged)-1]
		if r.start <= last.end {
			if r.end > last.end {
				last.end = r.end
			}
		} else {
			merged = append(merged, r)
		}
	}

	return merged
}

func part2(freshRanges []Range) int {
	sortedRanges := sortAndMerge(freshRanges)
	var result int
	for _, idRange := range sortedRanges {
		result += idRange.end - idRange.start + 1
	}
	return result
}
