package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	m := NewMapFromInput(lines)
	m = ExpandMap(m)

	fmt.Println("MAX RISK =", m.MaxEstimatedRisk())
	fmt.Println(m.Risk(Position{x: 0, y: 0}, Position{x: m.Width() - 1, y: m.Height() - 1}))

	return nil
}

func ExpandMap(m Map) Map {
	res := make(Map, len(m)*5)
	for i := range res {
		offset := i / m.Height()
		res[i] = make([]int, len(m[i%len(m)])*5)
		for j, r := range m[i%len(m)] {
			res[i][j] = (r+offset-1)%9 + 1
			res[i][j+m.Width()] = (r+offset)%9 + 1
			res[i][j+2*m.Width()] = (r+offset+1)%9 + 1
			res[i][j+3*m.Width()] = (r+offset+2)%9 + 1
			res[i][j+4*m.Width()] = (r+offset+3)%9 + 1
		}
	}

	return res
}

type Map [][]int

func (m Map) String() string {
	var builder strings.Builder
	for _, l := range m {
		s, _ := maps.Ints(l).ToStrings()
		builder.WriteString(strings.Join(s, ""))
		builder.WriteString("\n")
	}

	return builder.String()
}

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

func (m Map) MaxEstimatedRisk() int {
	risk := 0
	for i := 1; i < m.Height(); i++ {
		risk += m.At(Position{x: 0, y: i})
	}

	for i := 1; i < m.Width(); i++ {
		risk += m.At(Position{x: i, y: m.Height() - 1})
	}

	return risk
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

type PQueue struct {
	forRisk map[int][]Position
	maxRisk int
	minRisk int
}

func NewPQueue() PQueue {
	return PQueue{
		forRisk: make(map[int][]Position),
		maxRisk: 0,
		minRisk: 0,
	}
}

func (pq *PQueue) Pop() Position {
	for r := pq.minRisk; r <= pq.maxRisk; r++ {
		if len(pq.forRisk[r]) == 0 {
			continue
		}

		pq.minRisk = r

		l := len(pq.forRisk[r]) - 1
		p := pq.forRisk[r][l]
		pq.forRisk[r] = pq.forRisk[r][:l]

		return p
	}

	panic(fmt.Sprintf("%v", pq.forRisk))
}

func (pq *PQueue) Push(p Position, risk int) {
	if risk < pq.minRisk {
		pq.minRisk = risk
	}

	if risk > pq.maxRisk {
		pq.maxRisk = risk
	}

	pq.forRisk[risk] = append(pq.forRisk[risk], p)
}

func (m Map) Risk(from Position, to Position) int {
	cost := make(Map, len(m))
	for i := range cost {
		cost[i] = make([]int, len(m[i]))
		for j := range cost[i] {
			cost[i][j] = m.Width() * m.Height()
		}

	}
	cost[from.y][from.x] = 0

	pq := NewPQueue()
	visited := make(map[Position]bool)

	pq.Push(from, 0)

	current := pq.Pop()
	nbVisits := 1
	nbDuplicated := 0

	for current != to {
		nbVisits++
		visited[current] = true
		for _, p := range current.Neighbors() {
			if !m.IsValidPos(p) {
				continue
			}
			risk := cost.At(current) + m.At(p)
			if cost.At(p) <= risk {
				nbDuplicated++
				continue
			}

			cost[p.y][p.x] = risk

			pq.Push(p, risk)
		}
		for visited[current] {
			current = pq.Pop()
		}
	}

	fmt.Println("Nb Visits =", nbVisits)
	fmt.Println("Nb Duplicated =", nbDuplicated)
	return cost.At(to)
}

func printRiskMap(m Map, seen map[Position]bool) {
	var builder strings.Builder
	for y, l := range m {
		for x, r := range l {
			if r < m.Width()*m.Height() {
				if seen[Position{x: x, y: y}] {
					builder.WriteString("X")
				} else {
					builder.WriteString(".")
				}
			} else {
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	fmt.Print(builder.String())

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
