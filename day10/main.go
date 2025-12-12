package main

import (
	"container/heap"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/andydennehy/advent-of-code-2025/utils"
)

type queueItem struct {
	lights  string
	presses int
}

func main() {
	scanner, file := utils.ParseInput("input.txt")
	defer file.Close()

	var result1, result2 int
	for scanner.Scan() {
		var buttons [][]int
		var lights string
		var joltage []int
		fields := strings.Fields(scanner.Text())
		for _, field := range fields {
			switch field[0] {
			case '[':
				lights = field[1 : len(field)-1]
			case '(':
				var newButton []int
				for _, char := range strings.Split(field[1:len(field)-1], ",") {
					index, err := strconv.Atoi(char)
					utils.Check(err)
					newButton = append(newButton, index)
				}
				buttons = append(buttons, newButton)
			case '{':
				for _, char := range strings.Split(field[1:len(field)-1], ",") {
					index, err := strconv.Atoi(char)
					utils.Check(err)
					joltage = append(joltage, index)
				}
			}
		}
		result1 += part1(lights, buttons)
		result2 += part2(joltage, buttons)
	}

	fmt.Println(result1)
	fmt.Println(result2)
}

func toggle(s string, indices []int) string {
	b := []byte(s)
	for _, i := range indices {
		if b[i] == '.' {
			b[i] = '#'
		} else {
			b[i] = '.'
		}
	}
	return string(b)
}

func bfs(target string, buttons [][]int) int {
	initial := strings.Repeat(".", len(target)) // All lights off
	if initial == target {
		return 0
	}
	visited := make(map[string]bool)
	queue := []queueItem{{initial, 0}}

	for len(queue) > 0 {
		popped := queue[0]
		queue = queue[1:]

		if visited[popped.lights] {
			continue
		}

		if popped.lights == target {
			return popped.presses
		}

		for _, button := range buttons {
			next := toggle(popped.lights, button)
			queue = append(queue, queueItem{next, popped.presses + 1})
		}
		visited[popped.lights] = true
	}

	return -1
}

func part1(target string, buttons [][]int) int {
	return bfs(target, buttons)
}

func part2(joltage []int, buttons [][]int) int {
	// This is basically a system of linear equations. This
	// performs Gaussian elimination
	rows := len(joltage)
	cols := len(buttons)

	buttonMasks := make([]uint32, cols)
	for i, btn := range buttons {
		for _, idx := range btn {
			buttonMasks[i] |= 1 << idx
		}
	}

	grid := make([]*big.Rat, rows*cols)
	for i := range grid {
		grid[i] = new(big.Rat)
	}
	rhs := make([]*big.Rat, rows)
	for i := range rhs {
		rhs[i] = big.NewRat(int64(joltage[i]), 1)
	}

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			if buttonMasks[col]&(1<<row) != 0 {
				grid[row*cols+col].SetInt64(1)
			}
		}
	}

	for i := 0; i < rows && i < cols; i++ {
		pivotRow, pivotCol := -1, -1
		for c := i; c < cols && pivotRow < 0; c++ {
			for r := i; r < rows; r++ {
				if grid[r*cols+c].Sign() != 0 {
					pivotRow, pivotCol = r, c
					break
				}
			}
		}
		if pivotRow < 0 {
			break
		}

		if pivotRow != i {
			for c := range cols {
				grid[i*cols+c], grid[pivotRow*cols+c] = grid[pivotRow*cols+c], grid[i*cols+c]
			}
			rhs[i], rhs[pivotRow] = rhs[pivotRow], rhs[i]
		}

		if pivotCol != i {
			for r := 0; r < rows; r++ {
				grid[r*cols+i], grid[r*cols+pivotCol] = grid[r*cols+pivotCol], grid[r*cols+i]
			}
			buttonMasks[i], buttonMasks[pivotCol] = buttonMasks[pivotCol], buttonMasks[i]
		}

		pivot := new(big.Rat).Set(grid[i*cols+i])
		if pivot.Cmp(big.NewRat(1, 1)) != 0 {
			for c := i; c < cols; c++ {
				grid[i*cols+c].Quo(grid[i*cols+c], pivot)
			}
			rhs[i].Quo(rhs[i], pivot)
		}

		for r := 0; r < rows; r++ {
			if r != i && grid[r*cols+i].Sign() != 0 {
				factor := new(big.Rat).Set(grid[r*cols+i])
				for c := i; c < cols; c++ {
					sub := new(big.Rat).Mul(factor, grid[i*cols+c])
					grid[r*cols+c].Sub(grid[r*cols+c], sub)
				}
				sub := new(big.Rat).Mul(factor, rhs[i])
				rhs[r].Sub(rhs[r], sub)
			}
		}
	}

	numNonzeroRows := 0
	for r := rows - 1; r >= 0; r-- {
		for c := range cols {
			if grid[r*cols+c].Sign() != 0 {
				numNonzeroRows = r + 1
				break
			}
		}
		if numNonzeroRows > 0 {
			break
		}
	}

	rows = numNonzeroRows
	numFreeVars := cols - rows

	if numFreeVars == 0 {
		sum := 0
		for i := 0; i < rows; i++ {
			if !rhs[i].IsInt() || rhs[i].Sign() < 0 {
				return -1
			}
			sum += int(rhs[i].Num().Int64())
		}
		return sum
	}

	// Setup for priority queue search over free variables
	maxPresses := make([]int, numFreeVars)
	pressDiff := make([]*big.Rat, numFreeVars)
	factors := make([]int, numFreeVars)

	for i := 0; i < numFreeVars; i++ {
		freeCol := rows + i
		maxPresses[i] = 1 << 30
		for j := 0; j < len(joltage); j++ {
			if buttonMasks[freeCol]&(1<<j) != 0 && joltage[j] < maxPresses[i] {
				maxPresses[i] = joltage[j]
			}
		}
		pressDiff[i] = big.NewRat(1, 1)
		for r := 0; r < rows; r++ {
			pressDiff[i].Sub(pressDiff[i], grid[r*cols+freeCol])
		}
	}

	factor := 1
	for i := 0; i < numFreeVars; i++ {
		factors[i] = factor
		factor *= maxPresses[i] + 1
	}

	start := 0
	for i := 0; i < numFreeVars; i++ {
		if pressDiff[i].Sign() < 0 {
			start += maxPresses[i] * factors[i]
		}
	}

	pq := &ratHeap{}
	heap.Init(pq)
	heap.Push(pq, &heapItem{addedPresses: new(big.Rat), state: start})
	checked := make(map[int]bool)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*heapItem)
		if checked[item.state] {
			continue
		}
		checked[item.state] = true

		trailingPresses := make([]int, numFreeVars)
		for i := range numFreeVars {
			trailingPresses[i] = (item.state / factors[i]) % (maxPresses[i] + 1)
		}

		presses := make([]*big.Rat, rows)
		valid := true
		for r := 0; r < rows; r++ {
			presses[r] = new(big.Rat).Set(rhs[r])
			for i := range numFreeVars {
				sub := new(big.Rat).Mul(grid[r*cols+(rows+i)], big.NewRat(int64(trailingPresses[i]), 1))
				presses[r].Sub(presses[r], sub)
			}
			if presses[r].Sign() < 0 || !presses[r].IsInt() {
				valid = false
				break
			}
		}

		if valid {
			total := 0
			for r := 0; r < rows; r++ {
				total += int(presses[r].Num().Int64())
			}
			for i := range numFreeVars {
				total += trailingPresses[i]
			}
			return total
		}

		for i := range numFreeVars {
			if pressDiff[i].Sign() < 0 && trailingPresses[i] > 0 {
				newState := item.state - factors[i]
				if !checked[newState] {
					heap.Push(pq, &heapItem{new(big.Rat).Sub(item.addedPresses, pressDiff[i]), newState})
				}
			}
			if trailingPresses[i] < maxPresses[i] {
				newState := item.state + factors[i]
				if !checked[newState] {
					heap.Push(pq, &heapItem{new(big.Rat).Add(item.addedPresses, pressDiff[i]), newState})
				}
			}
		}
	}
	return -1
}

type heapItem struct {
	addedPresses *big.Rat
	state        int
}

type ratHeap []*heapItem

func (h ratHeap) Len() int           { return len(h) }
func (h ratHeap) Less(i, j int) bool { return h[i].addedPresses.Cmp(h[j].addedPresses) < 0 }
func (h ratHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *ratHeap) Push(x any)        { *h = append(*h, x.(*heapItem)) }
func (h *ratHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
