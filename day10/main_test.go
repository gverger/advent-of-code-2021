package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	s := NewStack()

	s.Push('a')
	s.Push('b')
	s.Push('c')

	pop := firstValue(s.Pop)

	require.Equal(t, 'c', pop())
	require.Equal(t, 'b', pop())
	require.Equal(t, 'a', pop())
}

func TestCorrupted(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		corrupted bool
		got       rune
		want      rune
		idx       int
	}{
		{
			name:      "empty",
			input:     "",
			corrupted: false,
		},
		{
			name:      "single",
			input:     "[]",
			corrupted: false,
		},
		{
			name:      "no left bracket",
			input:     "]",
			corrupted: true,
			idx:       0,
			got:       ']',
			want:      '?',
		},
		{
			name:      "no left bracket",
			input:     "[]]",
			corrupted: true,
			idx:       2,
			got:       ']',
			want:      '?',
		},
		{
			name:      "mixed ok",
			input:     "<[(){}[]]>",
			corrupted: false,
		},
		{
			name:      "mixed not ok",
			input:     "<[(){}([]]>",
			corrupted: true,
			idx:       9,
			got:       ']',
			want:      ')',
		},
		{
			name:      "test 1",
			input:     "{([(<{}[<>[]}>{[]{[(<()>",
			corrupted: true,
			idx:       12,
			got:       '}',
			want:      ']',
		},
		{
			name:      "test incomplete",
			input:     "[({(<(())[]>[[{[]{<()<>>",
			corrupted: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckSyntax(test.input)
			if !test.corrupted {
				require.Nil(t, err)
				return
			}

			require.Equal(t, string(test.got), string(err.got))
			require.Equal(t, string(test.want), string(err.want))
			require.Equal(t, test.idx, err.idx)
		})
	}
}

func TestMissingClosingCharacters(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		missing string
	}{
		{
			name:    "first",
			input:   "[({(<(())[]>[[{[]{<()<>>",
			missing: "}}]])})]",
		},
		{
			name:    "second",
			input:   "[(()[<>])]({[<{<<[]>>(",
			missing: ")}>]})",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			missing := MissingClosingCharacters(test.input)

			require.Equal(t, test.missing, missing)
		})
	}
}

func TestMissingCost(t *testing.T) {
	tests := []struct {
		name    string
		missing string
		cost    int
	}{
		{
			name:    "first",
			missing: "])}>",
			cost:    294,
		},
		{
			name:    "second",
			missing: "}}]])})]",
			cost:    288957,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cost := missingCost(test.missing)

			require.Equal(t, test.cost, cost)
		})
	}
}

func firstValue(f func() (rune, error)) func() rune {
	return func() rune {
		res, _ := f()
		return res
	}
}
