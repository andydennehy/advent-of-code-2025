package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type Box struct {
	x, y, z int
}

type Pair struct {
	i, j   int
	distSq int
}

type UnionFind struct {
	parent []int
	rank   []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	size := make([]int, n)
	for i := range n {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent, rank, size}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])  // Path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX == rootY {
		return false  // Already in same circuit
	}
	// Union by rank
	if uf.rank[rootX] < uf.rank[rootY] {
		rootX, rootY = rootY, rootX
	}
	uf.parent[rootY] = rootX
	uf.size[rootX] += uf.size[rootY]
	if uf.rank[rootX] == uf.rank[rootY] {
		uf.rank[rootX]++
	}
	return true
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var boxes []Box
	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])
		boxes = append(boxes, Box{x, y, z})
	}

	fmt.Println(part1(boxes))
	fmt.Println(part2(boxes))
}

func getSortedPairs(boxes []Box) []Pair {
	n := len(boxes)
	pairs := make([]Pair, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			dx, dy, dz := boxes[i].x-boxes[j].x, boxes[i].y-boxes[j].y, boxes[i].z-boxes[j].z
			pairs = append(pairs, Pair{i, j, dx*dx + dy*dy + dz*dz})
		}
	}
	sort.Slice(pairs, func(a, b int) bool {
		return pairs[a].distSq < pairs[b].distSq
	})
	return pairs
}

func part1(boxes []Box) int {
	pairs := getSortedPairs(boxes)
	uf := NewUnionFind(len(boxes))

	for i := 0; i < 1000; i++ {
		uf.Union(pairs[i].i, pairs[i].j)
	}

	circuitSizes := make(map[int]int)
	for i := 0; i < len(boxes); i++ {
		root := uf.Find(i)
		circuitSizes[root] = uf.size[root]
	}

	sizes := make([]int, 0, len(circuitSizes))
	for _, size := range circuitSizes {
		sizes = append(sizes, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	return sizes[0] * sizes[1] * sizes[2]
}

func part2(boxes []Box) int {
	pairs := getSortedPairs(boxes)
	uf := NewUnionFind(len(boxes))

	var lastPair Pair
	unions := 0
	for _, p := range pairs {
		if uf.Union(p.i, p.j) {
			lastPair = p
			unions++
			if unions == len(boxes)-1 {
				break
			}
		}
	}

	return boxes[lastPair.i].x * boxes[lastPair.j].x
}
