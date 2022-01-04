package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLine(t *testing.T) {

	testCases := []struct {
		line string
	}{
		{line: "[1,2]"},
		{line: "[[1,2],3]"},
		{line: "[9,[8,7]]"},
		{line: "[[1,9],[8,5]]"},
		{line: "[[[[1,2],[3,4]],[[5,6],[7,8]]],9]"},
		{line: "[[[9,[3,8]],[[0,9],6]],[[[3,7],[4,9]],3]]"},
		{line: "[[[[1,3],[5,3]],[[1,3],[8,7]]],[[[4,9],[6,9]],[[8,2],[7,3]]]]"},
	}

	for _, test := range testCases {
		t.Run(test.line, func(t *testing.T) {
			parsed := ParseLine(test.line)

			require.Equal(t, test.line, parsed.String())
		})
	}
}
func TestParseArray(t *testing.T) {

	testCases := []struct {
		line   string
		output FlatNumber
	}{
		{line: "[1,2]", output: FlatNumber{Open, 1, 2, Close}},
		{line: "[[1,2],3]", output: FlatNumber{Open, Open, 1, 2, Close, 3, Close}},
		{line: "[9,[8,7]]", output: FlatNumber{Open, 9, Open, 8, 7, Close, Close}},
		{line: "[[1,9],[8,5]]", output: FlatNumber{Open, Open, 1, 9, Close, Open, 8, 5, Close, Close}},
	}

	for _, test := range testCases {
		t.Run(test.line, func(t *testing.T) {
			parsed := ParseArray(test.line)

			require.Equal(t, test.output, parsed)
		})
	}
}

func TestAddArray(t *testing.T) {
	tests := []struct {
		n1     string
		n2     string
		output string
	}{
		{
			n1:     "[[9,8],1]",
			n2:     "[[2,3],4]",
			output: "[[[9,8],1],[[2,3],4]]",
		},
	}
	for _, test := range tests {
		t.Run(test.output, func(t *testing.T) {
			n1 := ParseArray(test.n1)
			n2 := ParseArray(test.n2)

			require.Equal(t, test.output, AddArray(n1, n2).String())
		})
	}
}

func TestExplode(t *testing.T) {
	tests := []struct {
		number string
		output string
		same   bool
	}{
		{
			number: "[[[[0,9],2],3],4]",
			output: "[[[[0,9],2],3],4]",
			same:   true,
		},
		{
			number: "[[[[[9,8],1],2],3],4]",
			output: "[[[[0,9],2],3],4]",
		},
		{
			number: "[7,[6,[5,[4,[3,2]]]]]",
			output: "[7,[6,[5,[7,0]]]]",
		},
		{
			number: "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
			output: "[[3,[2,[8,0]]],[9,[5,[7,0]]]]",
		},
		{
			number: "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
			output: "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
		},
		{
			number: "[[6,[5,[4,[3,2]]]],1]",
			output: "[[6,[5,[7,0]]],3]",
		},
	}

	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			n := ParseArray(test.number)
			expected := ParseArray(test.output)

			exploded, isModified := Explode(n)
			require.Equal(t, !test.same, isModified)

			require.Equal(t, expected.String(), exploded.String())
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		number string
		output string
		same   bool
	}{
		{
			number: "[[[[0,7],4],[15,[0,13]]],[1,1]]",
			output: "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]",
		},
	}

	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			n := ParseArray(test.number)
			expected := ParseArray(test.output)

			exploded, isModified := Split(n)
			require.Equal(t, !test.same, isModified)

			require.Equal(t, expected.String(), exploded.String())
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		number string
		output string
	}{
		{
			number: "[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]",
			output: "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
	}

	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			n := ParseArray(test.number)
			expected := ParseArray(test.output)

			exploded := ReduceArray(n)

			require.Equal(t, expected.String(), exploded.String())
		})
	}
}

func TestMagnitude(t *testing.T) {
	tests := []struct {
		number string
		result int
	}{
		{
			number: "[9,1]",
			result: 29,
		},
		{
			number: "[1,9]",
			result: 21,
		},
		{
			number: "[[9,1],[1,9]]",
			result: 129,
		},
		{
			number: "[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
			result: 3488,
		},
	}

	for _, test := range tests {
		t.Run(test.number, func(t *testing.T) {
			n := ParseArray(test.number)

			require.Equal(t, test.result, n.Magnitude())
		})
	}
}
