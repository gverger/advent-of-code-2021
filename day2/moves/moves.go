package moves

import (
	"fmt"
	"strconv"
	"strings"
)

type Move struct {
	Dir   Direction
	Units int
}

type Direction string

const (
	Forward Direction = "forward"
	Up      Direction = "up"
	Down    Direction = "down"
)

func (dir Direction) IsValid() bool {
	return dir == Forward || dir == Up || dir == Down
}

func FromInput(lines []string) ([]Move, error) {
	moves := make([]Move, 0, len(lines))

	for _, l := range lines {
		m, err := NewMove(l)
		if err != nil {
			return moves, err
		}
		moves = append(moves, m)
	}

	return moves, nil
}

func NewMove(line string) (Move, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return Move{}, fmt.Errorf("not a valid input: %q", line)
	}

	dir := Direction(parts[0])
	if !dir.IsValid() {
		return Move{}, fmt.Errorf("not a valid direction: %q", dir)
	}

	units, err := strconv.Atoi(parts[1])
	if err != nil {
		return Move{}, err
	}

	return Move{Dir: dir, Units: units}, nil
}
