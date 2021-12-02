package part1

import m "github.com/gverger/advent2021/day2/moves"

type Pos struct {
	X int
	Y int
}

func NewPos() Pos {
	return Pos{}
}

func (p Pos) Apply(moves []m.Move) Pos {
	newPos := Pos{X: p.X, Y: p.Y}

	for _, current := range moves {
		switch current.Dir {
		case m.Up:
			newPos.Y -= current.Units
		case m.Down:
			newPos.Y += current.Units
		case m.Forward:
			newPos.X += current.Units
		}
	}

	return newPos
}
