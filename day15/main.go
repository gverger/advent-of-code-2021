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
	m := NewMapFromInput(lines)

	fmt.Println("Part 1:", m.Risk(Position{x: 0, y: 0}, Position{x: m.Width() - 1, y: m.Height() - 1}))

	m = ExpandMap(m)
	fmt.Println("Part 2:", m.Risk(Position{x: 0, y: 0}, Position{x: m.Width() - 1, y: m.Height() - 1}))

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

		if len(pq.forRisk[r]) == 0 {
			delete(pq.forRisk, r)
		}

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

type Dijsktra struct {
	m       Map
	cost    Map
	pq      PQueue
	visited [][]bool
}

func (d Dijsktra) IsVisited(p Position) bool {
	return d.visited[p.y][p.x]
}

func NewDijsktra(m Map, from Position) Dijsktra {
	visited := make([][]bool, len(m))
	for i := range m {
		visited[i] = make([]bool, m.Width())
	}

	cost := make(Map, len(m))
	for i := range cost {
		cost[i] = make([]int, len(m[i]))
		for j := range cost[i] {
			cost[i][j] = m.Width() * m.Height()
		}

	}
	cost[from.y][from.x] = 0

	pq := NewPQueue()
	pq.Push(from, 0)

	return Dijsktra{
		m:       m,
		cost:    cost,
		pq:      pq,
		visited: visited,
	}
}

func (d *Dijsktra) nextNode() Position {
	next := d.pq.Pop()
	for d.IsVisited(next) {
		next = d.pq.Pop()
	}

	return next
}

func (d *Dijsktra) Step() Position {
	current := d.nextNode()

	d.visited[current.y][current.x] = true
	for _, p := range current.Neighbors() {
		if !d.m.IsValidPos(p) {
			continue
		}
		risk := d.cost.At(current) + d.m.At(p)
		if d.cost.At(p) <= risk {
			continue
		}

		d.cost[p.y][p.x] = risk

		d.pq.Push(p, risk)
	}

	return current
}

func (d Dijsktra) String() string {
	var builder strings.Builder
	for y, l := range d.m {
		for x, r := range l {
			if r < d.m.Width()*d.m.Height() {
				if d.IsVisited(Position{x: x, y: y}) {
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

	return builder.String()
}

func (m Map) Risk(from Position, to Position) int {
	d := NewDijsktra(m, from)
	rev := NewDijsktra(m, to)

	// fmt.Println(d)
	current := from

	for !rev.IsVisited(current) {
		current = d.Step()
		rev.Step()
	}

	// for current != to {
	// 	current = d.Step()
	// }

	// Uncomment to see visited nodes
	// dis := NewDijsktra(m, from)
	//
	// for p := range d.visited {
	// 	dis.visited[p] = true
	// }
	// for p := range rev.visited {
	// 	dis.visited[p] = true
	// }
	// fmt.Println(dis)

	// return d.cost.At(current) + m.At(current)
	return d.cost.At(current) + rev.cost.At(current) + m.At(to) - m.At(current)
}

func (m Map) At(p Position) int {
	return m[p.y][p.x]
}
