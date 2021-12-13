package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	part1(NewCavernFromInput(lines))
	part2(NewCavernFromInput(lines))

	return nil
}

func part1(c Cavern) {
	nb := 0
	for i := 0; i < 100; i++ {
		nb += c.Step()
	}

	fmt.Println("Part1: nb of flashes in 100 steps:", nb)
}

func part2(c Cavern) {
	nb := 0
	step := 0
	for nb != 100 {
		step += 1
		nb = c.Step()
	}

	fmt.Println("Part2: synchronized at step", step)
}

type Cavern struct {
	energy [][]int
}

func NewCavernFromInput(lines []string) Cavern {
	energy := make([][]int, 10)
	c := Cavern{energy: energy}
	for i := range energy {
		energy[i] = make([]int, 10)
	}

	for row, l := range lines {
		for col, char := range l {
			energy, _ := strconv.Atoi(string(char))
			c.Set(col, row, energy)
		}
	}

	return c
}

func (c *Cavern) Set(x int, y int, energy int) {
	if x < 0 || y < 0 || x > 9 || y > 9 {
		return
	}
	c.energy[y][x] = energy
}

func (c *Cavern) Incr(x int, y int) int {
	if x < 0 || y < 0 || x > 9 || y > 9 {
		return -1
	}
	if c.energy[y][x] <= 9 {
		c.energy[y][x] = c.energy[y][x] + 1
	}

	return c.energy[y][x]
}

func (c *Cavern) Get(x int, y int) int {
	return c.energy[y][x]
}

func (c *Cavern) Step() int {
	nbFlashes := 0
	flashes := make([]Position, 0)

	incr := func(p Position) {
		c.Incr(p.x, p.y)
		if c.Get(p.x, p.y) > 9 {
			c.Set(p.x, p.y, 0)
			flashes = append(flashes, p)
			nbFlashes++
		}
	}

	for y, row := range c.energy {
		for x := range row {
			incr(Position{x: x, y: y})
		}
	}

	for len(flashes) > 0 {
		f := flashes[0]
		flashes = flashes[1:]

		for _, p := range f.Neighbours() {
			if c.Get(p.x, p.y) == 0 {
				continue
			}

			incr(p)
		}
	}

	return nbFlashes
}

type Position struct {
	x int
	y int
}

func (p Position) IsValid() bool {
	return p.x >= 0 && p.y >= 0 && p.x <= 9 && p.y <= 9
}

func (p Position) Neighbours() []Position {
	neighbours := make([]Position, 0)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			x := p.x + dx
			y := p.y + dy
			p := Position{x: x, y: y}
			if p.IsValid() {
				neighbours = append(neighbours, p)
			}
		}
	}

	return neighbours
}

func (c Cavern) String() string {
	res := make([]string, 0, len(c.energy))
	for _, row := range c.energy {
		sRow, _ := maps.Ints(row).ToStrings()
		for i, c := range sRow {
			if c == "10" {
				sRow[i] = "X"
			}
		}

		res = append(res, strings.Join(sRow, ""))
	}

	return strings.Join(res, "\n")
}
