package main

import (
	"fmt"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type Region struct {
	width, height int
	presents      [6]int
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var sizes [6]int
	currentCount := 0
	currentBox := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "x") { break }
		if strings.HasSuffix(line, ":") {
			continue
		}
		if line == "" {
			if currentCount > 0 {
				sizes[currentBox] = currentCount
				currentCount = 0
			}
			currentBox += 1
			continue
		}
		currentCount += strings.Count(line, "#")
	}

	var regions []Region
	for ok := true; ok; ok = (scanner.Scan()) {
		line := scanner.Text()

		var width, height int
		var presents [6]int
		fmt.Sscanf(line, "%dx%d: %d %d %d %d %d %d", 
			&width, &height, &presents[0], 
			&presents[1], &presents[2], &presents[3], 
			&presents[4], &presents[5])

		regions = append(regions, Region{
			width:    width,
			height:   height,
			presents: presents,
		})
	}

	fmt.Println("Sizes:", sizes)
	fmt.Println("Regions:", len(regions), regions[:3])
	fmt.Println(part1(regions, sizes))
}

func part1(regions []Region, sizes [6]int) int {
	// Input is actually simple
	var result int
	RegionLoop:
	for _, region := range regions {
		availableSpace := region.width * region.height
		for index, present := range region.presents {
			availableSpace -= present * sizes[index]
			if availableSpace < 0 { continue RegionLoop }
		}
		result += 1
	}
	return result
}
