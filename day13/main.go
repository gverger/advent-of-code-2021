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
	dotLines := make([]string, 0)
	for len(lines[0]) > 0 {
		dotLines = append(dotLines, lines[0])
		lines = lines[1:]
	}

	lines = lines[1:]
	foldLines := lines

	s := NewSheetFromInput(dotLines)

	fold(&s, []string{foldLines[0]})
	fmt.Println("Nb Dots after one fold =", s.NbDots())
	fold(&s, foldLines[1:])
	fmt.Println(s)

	return nil
}

func fold(s *Sheet, lines []string) {
	for _, l := range lines {
		parts := strings.Split(l, "=")
		value, _ := strconv.Atoi(parts[1])
		if parts[0][len(parts[0])-1] == 'x' {
			s.FoldX(value)
		} else {
			s.FoldY(value)
		}
	}
}

type Sheet struct {
	dots map[Position]bool
	maxX int
	maxY int
}

func NewSheetFromInput(lines []string) Sheet {
	s := Sheet{dots: make(map[Position]bool)}
	for _, l := range lines {
		sCoords := strings.Split(l, ",")

		coords, _ := maps.Strings(sCoords).ToInts()
		s.AddDot(coords[0], coords[1])
	}

	return s
}

func (s Sheet) NbDots() int {
	return len(s.dots)
}

func (s *Sheet) AddDot(x int, y int) {
	s.dots[Position{x: x, y: y}] = true
	s.maxX = utils.Max(x, s.maxX)
	s.maxY = utils.Max(y, s.maxY)
}

func (s *Sheet) FoldX(x int) {
	for p := range s.dots {
		if p.x <= x {
			continue
		}

		delete(s.dots, p)
		s.AddDot(2*x-p.x, p.y)
	}
	s.maxX = x - 1
}

func (s *Sheet) FoldY(y int) {
	for p := range s.dots {
		if p.y <= y {
			continue
		}

		delete(s.dots, p)
		s.AddDot(p.x, 2*y-p.y)
	}
	s.maxY = y - 1
}

func (s Sheet) String() string {
	var builder strings.Builder
	builder.WriteString("  ")
	for x := 0; x <= s.maxX; x++ {
		builder.WriteString(strconv.Itoa(x % 10))
	}
	builder.WriteString("\n ")
	for x := 0; x <= s.maxX+2; x++ {
		builder.WriteString("-")
	}
	builder.WriteString("\n")
	for y := 0; y <= s.maxY; y++ {
		builder.WriteString(strconv.Itoa(y % 10))
		builder.WriteString("|")
		for x := 0; x <= s.maxX; x++ {
			if _, ok := s.dots[Position{x: x, y: y}]; ok {
				builder.WriteString("#")
			} else {
				builder.WriteString(" ")
			}
		}
		builder.WriteString("|\n")
	}
	builder.WriteString(" ")
	for x := 0; x <= s.maxX+2; x++ {
		builder.WriteString("-")
	}

	return builder.String()
}

type Position struct {
	x int
	y int
}
