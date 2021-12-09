package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	hm := HeatMapFromInput(lines)

	fmt.Println("Part 1 =", part1Score(hm))
	fmt.Println("Part 1 =", part2Score(hm))
	return nil
}

func part1Score(hm HeatMap) int {
	points := hm.LowPoints()

	res := 0
	for _, p := range points {
		res += hm.At(p.x, p.y) + 1
	}

	return res
}

func part2Score(hm HeatMap) int {
	lowPoints := hm.LowPoints()
	sizes := make([]int, 0, len(lowPoints))

	for _, p := range lowPoints {
		sizes = append(sizes, len(hm.BasinFrom(p)))
	}

	sort.IntSlice(sizes).Sort()
	sizes = sizes[len(sizes)-3:]
	return sizes[0] * sizes[1] * sizes[2]
}

type HeatMap struct {
	data      [][]int
	maxHeight int
}

func HeatMapFromInput(lines []string) HeatMap {
	data := make([][]int, len(lines))

	maxHeight := 0
	for i, l := range lines {
		data[i] = make([]int, 0)
		for _, c := range l {
			height, _ := strconv.Atoi(string(c))
			if height > maxHeight {
				maxHeight = height
			}

			data[i] = append(data[i], height)
		}
	}

	return HeatMap{data: data, maxHeight: maxHeight}
}

func (hm HeatMap) At(x int, y int) int {
	if !hm.IsValidPos(x, y) {
		return hm.maxHeight + 1
	}
	return hm.data[y][x]
}

func (hm HeatMap) IsValidPos(x int, y int) bool {
	return x >= 0 && y >= 0 && x < hm.Width() && y < hm.Height()
}

func (hm HeatMap) Width() int {
	return len(hm.data[0])
}

func (hm HeatMap) Height() int {
	return len(hm.data)
}

type Position struct {
	x int
	y int
}

func (p Position) West() Position {
	return Position{x: p.x - 1, y: p.y}
}

func (p Position) East() Position {
	return Position{x: p.x + 1, y: p.y}
}

func (p Position) North() Position {
	return Position{x: p.x, y: p.y + 1}
}

func (p Position) South() Position {
	return Position{x: p.x, y: p.y - 1}
}

func (p Position) Neighbors() []Position {
	return []Position{p.West(), p.East(), p.North(), p.South()}
}

func (hm HeatMap) LowPoints() []Position {
	lowPoints := make([]Position, 0)

	for i := 0; i < hm.Width(); i++ {
		for j := 0; j < hm.Height(); j++ {
			minAround := utils.Min(hm.At(i-1, j), hm.At(i, j-1), hm.At(i+1, j), hm.At(i, j+1))
			if hm.At(i, j) < minAround {
				lowPoints = append(lowPoints, Position{x: i, y: j})
			}
		}
	}

	return lowPoints
}

func (hm HeatMap) BasinFrom(p Position) []Position {
	inBasin := make(map[Position]bool)

	toAdd := NewPositionSet()
	toAdd.Add(p)

	for len(toAdd) > 0 {
		p, _ := toAdd.PickOne()

		if inBasin[p] {
			continue
		}

		inBasin[p] = true

		for _, n := range p.Neighbors() {
			if hm.IsValidPos(n.x, n.y) && hm.At(n.x, n.y) > hm.At(p.x, p.y) && hm.At(n.x, n.y) < 9 {
				toAdd.Add(n)
			}
		}
	}

	res := make([]Position, 0, len(inBasin))
	for p := range inBasin {
		res = append(res, p)
	}

	return res
}

type PositionSet map[Position]bool

func NewPositionSet() PositionSet {
	return make(PositionSet)
}

func (ps PositionSet) Add(p Position) {
	ps[p] = true
}

func (ps PositionSet) PickOne() (Position, bool) {
	for p := range ps {
		delete(ps, p)
		return p, true
	}

	return Position{}, false
}
