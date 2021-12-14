package main

import (
	"fmt"
	"strings"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	p := ProblemFromInput(lines)

	part1(p)
	part2(p)

	return nil
}

func part1(p Problem) {
	fmt.Println("Part1: ", solve(p, 10))
}

func part2(p Problem) {
	fmt.Println("Part2: ", solve(p, 40))
}

func solve(p Problem, nbSteps int) int {
	polymerPairs := pairsFromTemplate(p.polymerTemplate)

	for i := 0; i < nbSteps; i++ {
		polymerPairs = grow(polymerPairs, p.pairInsertions)
	}

	return countMaxMinusMin(p.polymerTemplate, polymerPairs)
}

func countMaxMinusMin(template string, pairs map[Pair]int) int {
	counts := make(map[byte]int)
	counts[template[0]] += 1
	counts[template[len(template)-1]] += 1
	for p, nb := range pairs {
		counts[p[0]] += nb
		counts[p[1]] += nb
	}

	allCounts := make([]int, 0, len(counts))
	for _, c := range counts {
		allCounts = append(allCounts, c/2)
	}

	min, max := utils.MinMax(allCounts...)

	return max - min
}

func pairsFromTemplate(template string) map[Pair]int {
	pairs := make(map[Pair]int)

	for i := 1; i < len(template); i++ {
		pairs[Pair(template[i-1:i+1])]++
	}

	return pairs
}

func grow(polymer map[Pair]int, insertions map[Pair]rune) map[Pair]int {
	res := make(map[Pair]int)

	for pair, nb := range polymer {
		ins := insertions[pair]
		res[Pair([]byte{pair[0], byte(ins)})] += nb
		res[Pair([]byte{byte(ins), pair[1]})] += nb
	}

	return res
}

type Problem struct {
	polymerTemplate string
	pairInsertions  map[Pair]rune
}

func (p Problem) String() string {
	var builder strings.Builder

	builder.WriteString(p.polymerTemplate)
	builder.WriteString("\n")

	for pair, ins := range p.pairInsertions {
		builder.WriteString(string(pair))
		builder.WriteString(" --> ")
		builder.WriteString(string(ins))
		builder.WriteString("\n")
	}

	return builder.String()
}

type Pair string

type PairInsertion struct {
	pair   Pair
	insert rune
}

func (pi PairInsertion) String() string {
	return fmt.Sprintf("%s -> %s", pi.pair, string(pi.insert))
}

func ProblemFromInput(lines []string) Problem {
	return Problem{polymerTemplate: lines[0], pairInsertions: insertionsFor(lines[2:])}
}

func insertionsFor(lines []string) map[Pair]rune {
	res := make(map[Pair]rune)

	for _, l := range lines {
		parts := strings.Split(l, " -> ")
		pair := Pair(parts[0])
		res[pair] = rune(parts[1][0])
	}

	return res
}
