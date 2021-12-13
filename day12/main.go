package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {

	fmt.Println("Nb paths (part 1):", part1(NewGraphFromInput(lines)))
	fmt.Println("Nb paths (part 2):", part2(NewGraphFromInput(lines)))

	return nil
}

func part1(g Graph) int {
	nb := 0
	fn := func(_ Path) {
		nb++
	}

	g.Part1Paths(fn)

	return nb
}

func part2(g Graph) int {
	nb := 0
	fn := func(_ Path) {
		nb++
	}

	g.Part2Paths(fn)

	return nb
}

type Graph struct {
	edge map[string][]string
}

func NewGraphFromInput(lines []string) Graph {
	g := Graph{edge: make(map[string][]string)}

	for _, l := range lines {
		ends := strings.Split(l, "-")
		g.AddEdge(ends[0], ends[1])
	}

	return g
}

func (g Graph) Part1Paths(exec func(path Path)) {
	p := NewPathPart1()
	p.Append("start")
	g.pathsFrom("start", &p, exec)
}

func (g Graph) Part2Paths(exec func(path Path)) {
	p := NewPathPart2()
	p.Append("start")
	g.pathsFrom("start", &p, exec)
}

func (g Graph) pathsFrom(start string, done Path, exec func(path Path)) {
	if start == "end" {
		exec(done)
	}

	for _, to := range g.edge[start] {
		if done.Visited(to) {
			continue
		}
		done.Append(to)
		g.pathsFrom(to, done, exec)
		done.DeleteLast()
	}
}

func (g *Graph) AddEdge(from string, to string) {
	g.edge[from] = append(g.edge[from], to)
	g.edge[to] = append(g.edge[to], from)
}

func (g Graph) String() string {
	var builder strings.Builder
	for from, tos := range g.edge {
		builder.WriteString(fmt.Sprintf("%s --> %v\n", from, tos))
	}

	return builder.String()
}

type Path interface {
	Append(string)
	Visited(string) bool
	DeleteLast()
	String() string
}

type PathPart1 struct {
	visited map[string]struct{}
	path    []string
}

func NewPathPart1() PathPart1 {
	return PathPart1{visited: make(map[string]struct{}), path: make([]string, 0)}
}

func (p PathPart1) Visited(node string) bool {
	if unicode.IsUpper(rune(node[0])) {
		return false
	}
	_, ok := p.visited[node]

	return ok
}

func (p *PathPart1) Append(node string) {
	p.path = append(p.path, node)
	p.visited[node] = struct{}{}
}

func (p *PathPart1) DeleteLast() {
	last := p.path[len(p.path)-1]
	delete(p.visited, last)

	p.path = p.path[:len(p.path)-1]
}

func (p PathPart1) String() string {
	return strings.Join(p.path, " --> ")
}

type PathPart2 struct {
	visited map[string]int
	path    []string
}

func NewPathPart2() PathPart2 {
	return PathPart2{visited: make(map[string]int), path: make([]string, 0)}
}

func (p PathPart2) Visited(node string) bool {
	if unicode.IsUpper(rune(node[0])) {
		return false
	}
	nb := p.visited[node]

	if nb == 0 {
		return false
	}

	for _, n := range p.path {
		if p.visited[n] == 2 && unicode.IsLower(rune(n[0])) {
			return true
		}
	}

	if node == "start" || node == "end" {
		return true
	}

	return false
}

func (p *PathPart2) Append(node string) {
	p.path = append(p.path, node)
	p.visited[node]++
}

func (p *PathPart2) DeleteLast() {
	last := p.path[len(p.path)-1]
	p.visited[last]--

	p.path = p.path[:len(p.path)-1]
}

func (p PathPart2) String() string {
	return strings.Join(p.path, ",")
}
