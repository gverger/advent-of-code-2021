package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/gverger/advent2021/utils"
	"github.com/gverger/advent2021/utils/maps"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	data := packetToBinary(lines[0])
	p, _ := decode(data)
	fmt.Println(p)

	fmt.Println("Sum of versions =", p.SumVersions())
	fmt.Println("Compute =", p.Compute())
	return nil
}

const (
	TypeLiteral = 4

	TypeSum     = 0
	TypeProduct = 1
	TypeMinimum = 2
	TypeMaximum = 3
	TypeGreater = 5
	TypeLess    = 6
	TypeEqual   = 7
)

func isLiteral(typeID int) bool {
	return typeID == TypeLiteral
}

func isOperator(typeID int) bool {
	return typeID != TypeLiteral
}

type Packet interface {
	Version() int
	TypeID() int

	SumVersions() int
	Compute() int
}

// Literal Packet
type LiteralPacket struct {
	Header

	value int
}

func NewLiteralPacket(h Header, v int) LiteralPacket {
	return LiteralPacket{Header: h, value: v}
}

func (p LiteralPacket) SumVersions() int {
	return p.Version()
}

func (p LiteralPacket) Compute() int {
	return p.value
}

// Operator Packet
type OperatorPacket struct {
	Header

	packets []Packet
}

func NewOperatorPacket(h Header, packets []Packet) OperatorPacket {
	return OperatorPacket{Header: h, packets: packets}
}

func (p OperatorPacket) SumVersions() int {
	res := p.Version()
	for _, sub := range p.packets {
		res += sub.SumVersions()
	}

	return res
}

func (p OperatorPacket) Compute() int {
	switch p.TypeID() {
	case TypeSum:
		res := 0
		for _, sub := range p.packets {
			res += sub.Compute()
		}
		return res
	case TypeProduct:
		res := 1
		for _, sub := range p.packets {
			res *= sub.Compute()
		}
		return res
	case TypeMinimum:
		res := p.packets[0].Compute()
		for _, sub := range p.packets {
			res = utils.Min(res, sub.Compute())
		}
		return res
	case TypeMaximum:
		res := p.packets[0].Compute()
		for _, sub := range p.packets {
			res = utils.Max(res, sub.Compute())
		}
		return res
	case TypeLess:
		a := p.packets[0].Compute()
		b := p.packets[1].Compute()
		if a < b {
			return 1
		}
		return 0
	case TypeGreater:
		a := p.packets[0].Compute()
		b := p.packets[1].Compute()
		if a > b {
			return 1
		}
		return 0
	case TypeEqual:
		a := p.packets[0].Compute()
		b := p.packets[1].Compute()
		if a == b {
			return 1
		}
		return 0
	}
	return -1
}

type Header struct {
	version int
	typeID  int
}

func (h Header) Version() int {
	return h.version
}

func (h Header) TypeID() int {
	return h.typeID
}

func decodeHeader(data string) (Header, string) {
	version := parseBinary(data[:3])
	typeID := parseBinary(data[3:6])

	return Header{version: version, typeID: typeID}, data[6:]
}

func decode(data string) (Packet, string) {
	h, data := decodeHeader(data)

	if isLiteral(h.typeID) {
		value, data := decodeLiteral(data)
		return NewLiteralPacket(h, value), data
	}

	packets, data := decodeOperator(data)
	return OperatorPacket{Header: h, packets: packets}, data
}

func decodeLiteral(data string) (int, string) {
	chunks := maps.String(data).ChunkEvery(5)

	nbChunks := 0
	var builder strings.Builder
	for _, c := range chunks {
		builder.WriteString(c[1:])
		nbChunks++
		if c[0] == '0' {
			break
		}
	}
	res, err := strconv.ParseInt(builder.String(), 2, 64)
	if err != nil {
		panic(err)
	}

	return int(res), data[nbChunks*5:]
}

func decodeOperator(data string) ([]Packet, string) {
	if data[0] == '0' {
		return decodeOperatorLength0(data[1:])
	}
	return decodeOperatorLength1(data[1:])
}

func decodeOperatorLength0(data string) ([]Packet, string) {
	l, data := decodeLength0(data)

	packets := make([]Packet, 0)
	d := data[:l]
	for len(d) > 0 {
		var p Packet
		p, d = decode(d)
		packets = append(packets, p)
	}

	return packets, data[l:]
}

func decodeOperatorLength1(data string) ([]Packet, string) {
	nbPackets, data := decodeLength1(data)

	packets := make([]Packet, 0)
	for i := 0; i < nbPackets; i++ {
		var p Packet
		p, data = decode(data)
		packets = append(packets, p)
	}

	return packets, data
}

func decodeLength0(data string) (int, string) {
	return parseBinary(data[:15]), data[15:]
}

func decodeLength1(data string) (int, string) {
	return parseBinary(data[:11]), data[11:]
}

func packetToBinary(packet string) string {
	bytes, _ := hex.DecodeString(packet)

	var dataBuilder strings.Builder
	for _, n := range bytes {
		dataBuilder.WriteString(fmt.Sprintf("%08b", n))
	}

	return dataBuilder.String()
}

func parseBinary(data string) int {
	res, err := strconv.ParseInt(data, 2, 64)
	if err != nil {
		panic(err)
	}

	return int(res)
}

func (p LiteralPacket) String() string {
	return fmt.Sprintf("LIT(%v, value=%d)", p.Header, p.value)
}

func (h Header) String() string {
	return fmt.Sprintf("H[v%d t%d]", h.Version(), h.TypeID())
}

func (p OperatorPacket) String() string {
	packets := make([]string, len(p.packets))

	for i, sub := range p.packets {
		packets[i] = fmt.Sprint(sub)
	}

	op := "OPE"
	switch p.TypeID() {
	case TypeSum:
		op = "SUM"
	case TypeProduct:
		op = "PROD"
	case TypeMinimum:
		op = "MIN"
	case TypeMaximum:
		op = "MAX"
	case TypeLess:
		op = "LESS"
	case TypeGreater:
		op = "GREATER"
	case TypeEqual:
		op = "EQ"
	}

	return fmt.Sprintf("%s(%v, packets=%v)", op, p.Header, packets)
}
