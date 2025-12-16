package main

import (
	"fmt"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type queueItem struct {
	device, path string
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	adjList := make(map[string][]string)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ": ")
		adjList[parts[0]] = strings.Fields(parts[1])
	}

	fmt.Println(part1("you", "out", adjList))
	fmt.Println(part2("svr", "out", adjList))
}

func part1(start, end string, adjList map[string][]string) int {
	// BFS
	var result int
	queue := []queueItem{{start, start}}
	visited := make(map[string]bool)
	for len(queue) > 0 {
		popped := queue[0]
		queue = queue[1:]
		if visited[popped.path] { continue }
		if popped.device == end {
			result += 1
			continue
		}
		for _, neighbor := range adjList[popped.device] {
			nextPath := popped.path + neighbor
			if !visited[nextPath] {
				queue = append(queue, queueItem{neighbor, nextPath})
			}
		}
	}
	return result
}

func countPaths(start, end string, adjList map[string][]string, cache map[string]int) int {
	if start == end { return 1 }
	if paths, ok := cache[start + end]; ok { return paths }
	var localResult int
	for _, neighbor := range adjList[start] {
		localResult += countPaths(neighbor, end, adjList, cache)
	}
    cache[start+end] = localResult
	return localResult
}

func part2(start, end string, adjList map[string][]string) int {
	// DP with memo
	cache := make(map[string]int)
	return countPaths(start, "fft", adjList, cache) *
	    countPaths("fft", "dac", adjList, cache) *
        countPaths("dac", end, adjList, cache) +
        countPaths(start, "dac", adjList, cache) *
        countPaths("dac", "fft", adjList, cache) *
        countPaths("fft", end, adjList, cache)
}
