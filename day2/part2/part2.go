package part2

import m "github.com/gverger/advent2021/day2/moves"

type Pos struct {
	X   int
	Y   int
	Aim int
}

func NewPos() Pos {
	return Pos{}
}

func (p Pos) Apply(moves []m.Move) Pos {
	newPos := Pos{X: p.X, Y: p.Y, Aim: p.Aim}

	for _, current := range moves {
		switch current.Dir {
		case m.Up:
			newPos.Aim -= current.Units
		case m.Down:
			newPos.Aim += current.Units
		case m.Forward:
			newPos.X += current.Units
			newPos.Y += newPos.Aim * current.Units
		}
	}

	return newPos
}
