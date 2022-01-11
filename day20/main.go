package main

import (
	"fmt"
	"strings"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	r := NewEnhanceRule(lines[0])

	g := Grid{outside: Dark}
	for j, line := range lines[2:] {
		for i, c := range line {
			if c == '#' {
				g.Set(Point{x: i, y: j})
			}
		}

	}

	fmt.Println(g)

	for i := 0; i < 50; i++ {
		g = g.Enhance(r)
	}

	fmt.Println(g)

	fmt.Println(len(g.data))
	return nil
}

type color int

const (
	Dark  color = 0
	Light color = 1
)

type Point struct {
	x int
	y int
}

type Grid struct {
	data    map[Point]color
	min     Point
	max     Point
	outside color
}

func (g *Grid) Set(p Point) {
	if g.data == nil {
		g.data = make(map[Point]color)
		g.min = p
		g.max = p
		g.data[p] = Light
		return
	}

	g.data[p] = Light
	if p.x < g.min.x {
		g.min.x = p.x
	}
	if p.y < g.min.y {
		g.min.y = p.y
	}
	if p.x > g.max.x {
		g.max.x = p.x
	}
	if p.y > g.max.y {
		g.max.y = p.y
	}
}

func (g Grid) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("GRID [%d,%d] --> [%d,%d]\n", g.min.x, g.min.y, g.max.x, g.max.y))

	for j := g.min.y; j <= g.max.y; j++ {
		for i := g.min.x; i <= g.max.x; i++ {
			c := "."
			if g.At(Point{x: i, y: j}) == Light {
				c = "#"
			}
			b.WriteString(c)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (g Grid) At(p Point) color {
	if p.x < g.min.x || p.y < g.min.y || p.x > g.max.x || p.y > g.max.y {
		return g.outside
	}

	if c, ok := g.data[p]; ok {
		return c
	}
	return Dark
}

func (g Grid) Around(p Point) int {
	res := 0

	for j := p.y - 1; j <= p.y+1; j++ {
		for i := p.x - 1; i <= p.x+1; i++ {
			value := int(g.At(Point{x: i, y: j}))
			res = res*2 + value
		}
	}

	return res
}

func (g Grid) Enhance(r EnhanceRule) Grid {
	eg := Grid{outside: r.For(g.Around(Point{x: g.min.x - 5, y: g.min.y - 5}))}

	for j := g.min.y - 3; j <= g.max.y+3; j++ {
		for i := g.min.x - 3; i <= g.max.x+3; i++ {
			p := Point{x: i, y: j}
			value := g.Around(p)
			if r.For(value) == Light {
				eg.Set(p)
			}
		}
	}
	return eg
}

type EnhanceRule string

func NewEnhanceRule(line string) EnhanceRule {
	return EnhanceRule(line)
}

func (e EnhanceRule) For(number int) color {
	if e[number] == '#' {
		return Light
	}
	return Dark
}
