package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type Tile struct {
	row, col int
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var tiles []Tile
	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ",")
		col, _ := strconv.Atoi(coords[0])
		row, _ := strconv.Atoi(coords[1])
		tiles = append(tiles, Tile{row, col})
	}

	fmt.Println(part1(tiles))
	fmt.Println(part2(tiles))
}

func squareArea(t1, t2 Tile) int {
	rowDist := t1.row - t2.row
	if rowDist < 0 {
		rowDist = -rowDist
	}
	colDist := t1.col - t2.col
	if colDist < 0 {
		colDist = -colDist
	}
	return (rowDist + 1) * (colDist + 1)
}

func pointInPolygon(pr, pc float64, tiles []Tile) bool {
	n := len(tiles)
	inside := false

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		ri, ci := float64(tiles[i].row), float64(tiles[i].col)
		rj, cj := float64(tiles[j].row), float64(tiles[j].col)

		if (ri > pr) != (rj > pr) {
			crossCol := ci + (pr-ri)*(cj-ci)/(rj-ri)
			if pc < crossCol {
				inside = !inside
			}
		}
	}
	return inside
}

func part1(tiles []Tile) int {
	// O(n^2): bad
	var maxArea int
	for i, tile := range tiles {
		for _, otherTile := range tiles[i+1:] {
			area := squareArea(tile, otherTile)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

func segmentIntersectsRectInterior(r1, c1, r2, c2, minR, maxR, minC, maxC int) bool {
	dr := float64(r2 - r1)
	dc := float64(c2 - c1)

	var tRowLo, tRowHi, tColLo, tColHi float64

	if dr == 0 {
		if float64(r1) <= float64(minR) || float64(r1) >= float64(maxR) {
			return false
		}
		tRowLo, tRowHi = math.Inf(-1), math.Inf(1)
	} else {
		t1 := (float64(minR) - float64(r1)) / dr
		t2 := (float64(maxR) - float64(r1)) / dr
		if dr > 0 {
			tRowLo, tRowHi = t1, t2
		} else {
			tRowLo, tRowHi = t2, t1
		}
	}

	// Col constraints: minC < c1 + t*dc < maxC
	if dc == 0 {
		// Segment has constant col
		if float64(c1) <= float64(minC) || float64(c1) >= float64(maxC) {
			return false // Col is on or outside boundary
		}
		tColLo, tColHi = math.Inf(-1), math.Inf(1)
	} else {
		t1 := (float64(minC) - float64(c1)) / dc
		t2 := (float64(maxC) - float64(c1)) / dc
		if dc > 0 {
			tColLo, tColHi = t1, t2
		} else {
			tColLo, tColHi = t2, t1
		}
	}

	tLo := math.Max(math.Max(tRowLo, tColLo), 0.0)
	tHi := math.Min(math.Min(tRowHi, tColHi), 1.0)

	return tLo < tHi
}

func part2(tiles []Tile) int {
	n := len(tiles)

	var maxArea int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			t1, t2 := tiles[i], tiles[j]

			// Form rectangle from the two tiles
			minR := min(t1.row, t2.row)
			maxR := max(t1.row, t2.row)
			minC := min(t1.col, t2.col)
			maxC := max(t1.col, t2.col)

			// First check: the center of the rectangle must be inside the polygon
			centerR := float64(minR+maxR) / 2.0
			centerC := float64(minC+maxC) / 2.0
			if !pointInPolygon(centerR, centerC, tiles) {
				continue
			}

			// Second check: no polygon edge intersects the interior of this rectangle
			valid := true
			for k := range n {
				next := (k + 1) % n
				if segmentIntersectsRectInterior(
					tiles[k].row, tiles[k].col,
					tiles[next].row, tiles[next].col,
					minR, maxR, minC, maxC,
				) {
					valid = false
					break
				}
			}

			if valid {
				area := squareArea(t1, t2)
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}
	return maxArea
}
