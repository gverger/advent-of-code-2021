package main

import (
	"strings"
	"testing"

	"github.com/gverger/advent2021/utils/maps"
	"github.com/stretchr/testify/require"
)

func TestStep(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		output []string
	}{
		{
			name: "no flash",
			input: []string{
				"5483143223",
				"2745854711",
				"5264556173",
				"6141336146",
				"6357385478",
				"4167524645",
				"2176841721",
				"6882881134",
				"4846848554",
				"5283751526",
			},
			output: []string{
				"6594254334",
				"3856965822",
				"6375667284",
				"7252447257",
				"7468496589",
				"5278635756",
				"3287952832",
				"7993992245",
				"5957959665",
				"6394862637",
			},
		},
		{
			name: "with flashes",
			input: []string{
				"6594254334",
				"3856965822",
				"6375667284",
				"7252447257",
				"7468496589",
				"5278635756",
				"3287952832",
				"7993992245",
				"5957959665",
				"6394862637",
			},
			output: []string{
				"8807476555",
				"5089087054",
				"8597889608",
				"8485769600",
				"8700908800",
				"6600088989",
				"6800005943",
				"0000007456",
				"9000000876",
				"8700006848",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewCavernFromInput(test.input)
			c.Step()
			require.Equal(t, test.output, cavernState(c))
		})
	}
}

func cavernState(c Cavern) []string {
	res := make([]string, len(c.energy))
	for i, row := range c.energy {
		s, _ := maps.Ints(row).ToStrings()
		res[i] = strings.Join(s, "")
	}

	return res
}
