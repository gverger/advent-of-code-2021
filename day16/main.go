package main

import (
	"fmt"

	"github.com/gverger/advent2021/utils"
)

func main() {
	utils.Main(run)
}

func run(lines []string) error {
	data := DataFromHexString(lines[0])
	p, _ := decode(data)
	fmt.Println(p)

	fmt.Println("Sum of versions =", p.SumVersions())
	fmt.Println("Compute =", p.Compute())
	return nil
}

const (
	TypeLiteral = 4
)

func isLiteral(typeID int) bool {
	return typeID == TypeLiteral
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

	packets   []Packet
	operation operation
}

func NewOperatorPacket(h Header, packets []Packet) OperatorPacket {
	return OperatorPacket{Header: h, packets: packets, operation: operationForType[h.TypeID()]}
}

func (p OperatorPacket) SumVersions() int {
	res := p.Version()
	for _, sub := range p.packets {
		res += sub.SumVersions()
	}

	return res
}

func (p OperatorPacket) Compute() int {
	return p.operation(p.packets)
}

func decodeHeader(data rawData) (Header, rawData) {
	version, data := data.DecodeBinary(3)
	typeID, data := data.DecodeBinary(3)

	return Header{version: version, typeID: typeID}, data
}

func decode(data rawData) (Packet, rawData) {
	h, data := decodeHeader(data)

	if isLiteral(h.typeID) {
		value, data := decodeLiteralContent(data)
		return NewLiteralPacket(h, value), data
	}

	packets, data := decodeOperatorContent(data)
	return NewOperatorPacket(h, packets), data
}

func decodeLiteralContent(data rawData) (int, rawData) {
	value := 0
	for {
		chunk := data.FirstBits(5)
		value = value*16 + chunk.FromBit(1).parseBinary()
		data = data.FromBit(5)

		if chunk.At(0) == ZERO {
			break
		}
	}
	return value, data
}

func decodeOperatorContent(data rawData) ([]Packet, rawData) {
	if data.At(0) == ZERO {
		return decodeOperatorLength0(data.FromBit(1))
	}

	return decodeOperatorLength1(data.FromBit(1))
}

func decodeOperatorLength0(data rawData) ([]Packet, rawData) {
	l, data := data.DecodeBinary(15)

	packets := make([]Packet, 0)
	d := data.FirstBits(l)
	for d.Length() > 0 {
		var p Packet
		p, d = decode(d)
		packets = append(packets, p)
	}

	return packets, data.FromBit(l)
}

func decodeOperatorLength1(data rawData) ([]Packet, rawData) {
	nbPackets, data := data.DecodeBinary(11)

	packets := make([]Packet, 0)
	for i := 0; i < nbPackets; i++ {
		var p Packet
		p, data = decode(data)
		packets = append(packets, p)
	}

	return packets, data
}
