package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	scanners := make([]*Scanner, 0)
	var currentScanner *Scanner

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		if strings.HasPrefix(l, "---") {
			currentScanner = newScanner(l[4 : len(l)-4])
			scanners = append(scanners, currentScanner)
			continue
		}

		currentScanner.beacons = append(currentScanner.beacons, parseBeacon(l))
	}

	done := make([]bool, len(scanners))
	todo := make([]int, 0)
	trBeacons := make([]func(Beacon) Beacon, len(scanners))

	trBeacons[0] = func(b Beacon) Beacon { return b }
	todo = append(todo, 0)
	done[0] = true

	for len(todo) > 0 {
		current := todo[0]
		todo = todo[1:]
		// done[current] = true
		scanner := scanners[current]

		for i, s := range scanners {
			if done[i] {
				continue
			}

			f, ok := rotationBetweenScanners(*scanner, *s)
			if ok {
				fmt.Printf("Close scanners: %s and %s\n", scanner.name, s.name)
				todo = append(todo, i)
				trBeacons[i] = func(b Beacon) Beacon { return trBeacons[current](f(b)) }
				done[i] = true
			}
		}
	}

	allBeacons := make(map[Beacon]bool)
	for i, s := range scanners {
		tr := trBeacons[i]
		for _, b := range s.beacons {
			allBeacons[tr(b)] = true
		}
	}

	fmt.Println("Count all =", len(allBeacons))

	scannerPos := make([]Beacon, len(scanners))

	for i := range scanners {
		scannerPos[i] = trBeacons[i](Beacon{})
	}

	oceanSize := 0
	for i := range scannerPos {
		for j := range scannerPos {
			oceanSize = utils.Max(oceanSize, int(scannerPos[i].manDistTo(scannerPos[j])))
		}
	}

	fmt.Println("Ocean Size =", oceanSize)

	return nil
}

func rotationBetweenScanners(s1, s2 Scanner) (func(Beacon) Beacon, bool) {
	d1 := distMaps(s1)
	d2 := distMaps(s2)

	same1 := make([]Beacon, 0)
	same2 := make([]Beacon, 0)
	for b1, dists := range d1 {
		for b2, dists2 := range d2 {
			sameDists := 0
			for d, nb := range dists {
				sameDists += utils.Min(nb, dists2[d])
			}
			if sameDists >= 12 {
				same1 = append(same1, b1)
				same2 = append(same2, b2)
				break
			}
		}
	}

	if len(same1) >= 12 {
		same1, same2 = reallySame(same1, same2)
		if len(same1) >= 12 {
			return rotationBetweenBeacons(same1, same2)
		}
	}
	return nil, false
}

func rotationBetweenBeacons(same1, same2 []Beacon) (func(Beacon) Beacon, bool) {
	if len(same1) == 0 {
		return func(b Beacon) Beacon { return b }, false
	}

	b1s := make([]Beacon, len(same1))
	o1 := same1[0]
	for i, b := range same1 {
		b1s[i] = Beacon{x: b.x - o1.x, y: b.y - o1.y, z: b.z - o1.z}
	}

	b2s := make([]Beacon, len(same2))
	o2 := same2[0]
	for i, b := range same2 {
		b2s[i] = Beacon{x: b.x - o2.x, y: b.y - o2.y, z: b.z - o2.z}
	}

	idx := 1
	for b2s[idx].x == b2s[idx].y || b2s[idx].y == b2s[idx].z || b2s[idx].x == b2s[idx].z {
		idx += 1
	}
	rx := findRotation(b1s[idx].x, b2s[idx])
	ry := findRotation(b1s[idx].y, b2s[idx])
	rz := findRotation(b1s[idx].z, b2s[idx])

	for i := range b2s {
		b := b2s[i]
		b2s[i] = Beacon{x: rx(b), y: ry(b), z: rz(b)}
	}

	nbOk := 0
	for i := range b2s {
		b1 := b1s[i]
		b2 := b2s[i]
		if b1 == b2 {
			nbOk++
		}
	}

	if nbOk < 12 {
		return nil, false
	}

	f := func(b Beacon) Beacon {
		b = Beacon{x: b.x - o2.x, y: b.y - o2.y, z: b.z - o2.z}
		b = Beacon{x: rx(b), y: ry(b), z: rz(b)}
		b = Beacon{x: b.x + o1.x, y: b.y + o1.y, z: b.z + o1.z}
		return b
	}

	return f, true
}

type Beacon struct {
	x int
	y int
	z int
}

func parseBeacon(line string) Beacon {
	coords, err := maps.Strings(strings.Split(line, ",")).ToInts()
	if err != nil {
		panic(err)
	}
	return Beacon{x: coords[0], y: coords[1], z: coords[2]}
}

func (b Beacon) distTo(other Beacon) float64 {
	return math.Sqrt(math.Pow(float64(b.x-other.x), 2) + math.Pow(float64(b.y-other.y), 2) + math.Pow(float64(b.z-other.z), 2))
}

func (b Beacon) manDistTo(other Beacon) int {
	return utils.Abs(b.x-other.x) + utils.Abs(b.y-other.y) + utils.Abs(b.z-other.z)
}

type Scanner struct {
	name    string
	beacons []Beacon
}

func newScanner(name string) *Scanner {
	return &Scanner{name: name, beacons: make([]Beacon, 0)}
}

func toInt(number string) int {
	n, _ := strconv.Atoi(number)
	return n
}

func (s Scanner) Distances() map[Beacon][]float64 {
	res := make(map[Beacon][]float64)
	for _, b := range s.beacons {
		res[b] = make([]float64, 0, len(s.beacons))
		for _, b2 := range s.beacons {
			res[b] = append(res[b], b.distTo(b2))
		}
	}

	return res
}

func reallySame(same1, same2 []Beacon) ([]Beacon, []Beacon) {
	if len(same1) == 0 {
		return nil, nil
	}

	b1s := make([]Beacon, len(same1))
	o1 := same1[0]
	for i, b := range same1 {
		b1s[i] = Beacon{x: b.x - o1.x, y: b.y - o1.y, z: b.z - o1.z}
	}

	b2s := make([]Beacon, len(same2))
	o2 := same2[0]
	for i, b := range same2 {
		b2s[i] = Beacon{x: b.x - o2.x, y: b.y - o2.y, z: b.z - o2.z}
	}

	idx := 1
	for b2s[idx].x == b2s[idx].y || b2s[idx].y == b2s[idx].z || b2s[idx].x == b2s[idx].z {
		idx += 1
	}
	rx := findRotation(b1s[idx].x, b2s[idx])
	ry := findRotation(b1s[idx].y, b2s[idx])
	rz := findRotation(b1s[idx].z, b2s[idx])

	for i := range b2s {
		b := b2s[i]
		b2s[i] = Beacon{x: rx(b), y: ry(b), z: rz(b)}
	}

	res1 := make([]Beacon, 0)
	res2 := make([]Beacon, 0)

	for i := range b2s {
		b1 := b1s[i]
		b2 := b2s[i]
		if b1 == b2 {
			res1 = append(res1, same1[i])
			res2 = append(res2, same2[i])
		} else {
			fmt.Println("ERROR")
		}
	}

	return res1, res2
}

func findRotation(goal int, b Beacon) func(Beacon) int {
	for _, r := range Rotations {
		if r(b) == goal {
			return r
		}
	}

	return nil
}

func X(b Beacon) int      { return b.x }
func MinusX(b Beacon) int { return -b.x }
func Y(b Beacon) int      { return b.y }
func MinusY(b Beacon) int { return -b.y }
func Z(b Beacon) int      { return b.z }
func MinusZ(b Beacon) int { return -b.z }

var Rotations = []func(Beacon) int{X, MinusX, Y, MinusY, Z, MinusZ}

func same12Dists(s1, s2 Scanner) bool {
	d1 := distMaps(s1)
	d2 := distMaps(s2)

	nbSame := 0
	for _, dists := range d1 {
		for _, dists2 := range d2 {
			sameDists := 0
			for d, nb := range dists {
				sameDists += utils.Min(nb, dists2[d])
			}
			if sameDists >= 12 {
				nbSame++
				break
			}
		}
	}

	return nbSame >= 12
}

func distMaps(s Scanner) map[Beacon]map[float64]int {
	res := make(map[Beacon]map[float64]int)
	for b, dists := range s.Distances() {
		res[b] = make(map[float64]int)
		for _, d := range dists {
			res[b][d] += 1
		}
	}

	return res
}
