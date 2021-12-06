package main

import (
	"fmt"
	"strings"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	vents := ventsFromInput(lines)
	// vents = part1Filter(vents)

	minX, maxX := minmax(vents[0].startX, vents[0].endX)
	minY, maxY := minmax(vents[0].startY, vents[0].endY)

	for _, v := range vents[1:] {
		minX, maxX = minmax(v.startX, v.endX, minX, maxX)
		minY, maxY = minmax(v.startY, v.endY, minY, maxY)
	}

	dangerous := 0
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			crossingVents := 0
			for _, v := range vents {
				if v.IsIn(x, y) {
					crossingVents += 1
					if crossingVents >= 2 {
						dangerous += 1
						break
					}
				}
			}
		}
	}

	fmt.Println("Result:", dangerous)

	return nil
}

func part1Filter(vents []Vent) []Vent {
	filtered := make([]Vent, 0)
	for _, v := range vents {
		if !v.IsHorizontal() && !v.IsVertical() {
			continue
		}
		filtered = append(filtered, v)
	}

	return filtered
}

type Vent struct {
	startX int
	startY int
	endX   int
	endY   int
}

func ventsFromInput(input []string) []Vent {
	vents := make([]Vent, 0, len(input))
	for _, line := range input {
		vents = append(vents, ventFromInput(line))
	}

	return vents
}

// ventFromInput takes an input like `A,B -> C,D`
func ventFromInput(input string) Vent {
	parts := strings.Fields(input)

	start, _ := maps.Strings(strings.Split(parts[0], ",")).ToInts()
	end, _ := maps.Strings(strings.Split(parts[2], ",")).ToInts()

	v := Vent{}
	v.startX, v.endX = start[0], end[0]
	v.startY, v.endY = start[1], end[1]

	return v
}

func (v Vent) String() string {
	return fmt.Sprintf("%d,%d -> %d,%d", v.startX, v.startY, v.endX, v.endY)
}

func (v Vent) IsIn(x int, y int) bool {
	minx, maxx := minmax(v.startX, v.endX)

	if x < minx || x > maxx {
		return false
	}
	if v.IsHorizontal() {
		return y == v.startY
	}

	miny, maxy := minmax(v.startY, v.endY)

	if y < miny || y > maxy {
		return false
	}
	if v.IsVertical() {
		return x == v.startX
	}

	slope := (v.endY - v.startY) / (v.endX - v.startX)
	b := v.endY - slope*v.endX

	return y == slope*x+b
}

func (v Vent) IsHorizontal() bool {
	return v.startY == v.endY
}

func (v Vent) IsVertical() bool {
	return v.startX == v.endX
}

func minmax(values ...int) (int, int) {
	min := values[0]
	max := values[0]

	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return min, max
}
