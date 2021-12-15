package main

import (
	"errors"
	"fmt"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	fmt.Println(NewMapFromInput(lines).Risk(Position{x: 0, y: 0}, Position{x: len(lines) - 1, y: len(lines) - 1}))

	return nil
}

type Map [][]int

func NewMapFromInput(lines []string) Map {
	res := make(Map, len(lines))
	for i, l := range lines {
		res[i], _ = maps.String(l).ToInts()
	}

	return res
}

func (m Map) Width() int {
	return len(m[0])
}

func (m Map) Height() int {
	return len(m)
}
func (m Map) IsValidPos(p Position) bool {
	return p.x >= 0 && p.y >= 0 && p.x < m.Width() && p.y < m.Height()
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

func (m Map) Risk(from Position, to Position) int {
	visited := make([][]int, len(m))
	for i := range visited {
		visited[i] = make([]int, len(m[i]))
		for j := range visited[i] {
			visited[i][j] = m.Width() * m.Height()
		}

	}
	visited[from.y][from.x] = 0

	return m.risk(from, to, visited)
}

func (m Map) risk(from Position, to Position, visited Map) int {

	for _, p := range from.Neighbors() {
		if !m.IsValidPos(p) {
			continue
		}
		risk := visited.At(from) + m.At(p)
		if visited.At(p) <= risk {
			continue
		}
		visited[p.y][p.x] = risk
		m.risk(p, to, visited)
	}

	return visited.At(to)
}

func (m Map) At(p Position) int {
	return m[p.y][p.x]
}

type Stack struct {
	data []Position
}

func NewStack() Stack {
	return Stack{data: make([]Position, 0)}
}

func (s *Stack) Push(value Position) {
	s.data = append(s.data, value)
}

func (s Stack) IsEmpty() bool {
	return len(s.data) == 0
}

var errorEmptyStack = errors.New("stack is empty")

func (s *Stack) Pop() (Position, error) {
	if len(s.data) == 0 {
		return Position{}, errorEmptyStack
	}
	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return res, nil
}
