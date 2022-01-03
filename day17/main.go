package main

import (
	"fmt"
	"regexp"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	t := TargetFromInput(lines[0])

	fmt.Println("Height=", MaxHeight(t))
	fmt.Println("Nb Possibles=", len(Possibles(t)))

	return nil
}

type Direction struct {
	dx int
	dy int
}

func Possibles(t Target) []Direction {
	dxs := findDX(t)
	dys := findDY(t)

	nb := 0
	dirs := make([]Direction, 0)
	for _, dx := range dxs {
		for _, dy := range dys {
			p := Probe{dx: dx, dy: dy}

			for p.x < t.minX {
				p.Step()
			}

			for p.x <= t.maxX && p.y >= t.minY {
				if p.y <= t.maxY {
					dirs = append(dirs, Direction{dx: dx, dy: dy})
					nb++
					break
				}
				p.Step()
			}
		}
	}

	return dirs
}

func MaxHeight(t Target) int {
	dy := utils.Max(findDY(t)...)
	return dy * (dy + 1) / 2
}

func findDX(t Target) []int {
	minX := t.minX
	maxX := t.maxX

	res := make([]int, 0)

	for ix := 0; ix <= maxX; ix++ {
		dx := ix + 1
		for x := 0; x <= maxX; x += dx {
			if x >= minX {
				res = append(res, ix)
				break
			}
			if dx > 0 {
				dx--
			} else {
				break
			}
		}
	}

	return res
}

func findDY(t Target) []int {
	maxY := -t.minY
	minY := -t.maxY

	res := make([]int, 0)

	for iy := 0; iy <= maxY; iy++ {
		dy := iy - 1
		for y := 0; y <= maxY; y += dy {
			if y >= minY {
				res = append(res, -iy)
				if iy > 1 {
					res = append(res, iy-1)
				}
				break
			}
			dy++
		}
	}

	return res
}

type Target struct {
	minX int
	minY int
	maxX int
	maxY int
}

func TargetFromInput(line string) Target {
	re := regexp.MustCompile("-?[0-9]+")
	matches := re.FindAllString(line, -1)
	numbers, _ := maps.Strings(matches).ToInts()
	return Target{
		minX: numbers[0],
		maxX: numbers[1],
		minY: numbers[2],
		maxY: numbers[3],
	}
}

type Probe struct {
	x  int
	y  int
	dx int
	dy int
}

func (p *Probe) Step() {
	p.x += p.dx
	p.y += p.dy

	if p.dx < 0 {
		p.dx++
	} else if p.dx > 0 {
		p.dx--
	}

	p.dy--
}
