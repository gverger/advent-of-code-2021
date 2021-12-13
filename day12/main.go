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

	part1(NewGraphFromInput(lines))

	return nil
}

func part1(g Graph) {
	nb := 0
	fn := func(_ Path) {
		nb++
	}

	g.Paths(fn)

	fmt.Println("Nb Paths =", nb)
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

func (g Graph) Paths(exec func(path Path)) {
	p := NewPath()
	p.Append("start")
	g.pathsFrom("start", p, exec)
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

type Path struct {
	visited map[string]struct{}
	path    []string
}

func NewPath() Path {
	return Path{visited: make(map[string]struct{}), path: make([]string, 0)}
}

func (p Path) Visited(node string) bool {
	if unicode.IsUpper(rune(node[0])) {
		return false
	}
	_, ok := p.visited[node]

	return ok
}

func (p *Path) Append(node string) {
	p.path = append(p.path, node)
	p.visited[node] = struct{}{}
}

func (p *Path) DeleteLast() {
	last := p.path[len(p.path)-1]
	delete(p.visited, last)

	p.path = p.path[:len(p.path)-1]
}
