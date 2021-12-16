package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeHeader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Header
	}{
		{
			name:  "simple packet",
			input: "D2FE28",
			want:  Header{version: 6, typeID: 4},
		},
		{
			name:  "operator packet length ID 0",
			input: "38006F45291200",
			want:  Header{version: 1, typeID: 6},
		},
		{
			name:  "operator packet length ID 1",
			input: "EE00D40C823060",
			want:  Header{version: 7, typeID: 3},
		},
		{
			name:  "operator packet sum 16",
			input: "8A004A801A8002F478",
			want:  Header{version: 4, typeID: 2},
		},
		{
			name:  "operator packet sum 12",
			input: "620080001611562C8802118E34",
			want:  Header{version: 3, typeID: 0},
		},
		{
			name:  "operator packet sum 23",
			input: "C0015000016115A2E0802F182340",
			want:  Header{version: 6, typeID: 0},
		},
		{
			name:  "operator packet sum 31",
			input: "A0016C880162017C3686B18A3D4780",
			want:  Header{version: 5, typeID: 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := packetToBinary(test.input)
			p, _ := decode(data)
			assert.Equal(t, test.want.version, p.Version(), "version")
			assert.Equal(t, test.want.typeID, p.TypeID(), "type ID")
		})
	}
}

func TestDecodeLiteral(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple packet",
			input: "D2FE28",
			want:  2021,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := packetToBinary(test.input)
			p, _ := decode(data)

			require.IsType(t, LiteralPacket{}, p)
			assert.Equal(t, test.want, p.(LiteralPacket).value, "value")
		})
	}
}

func TestDecodeLength(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "operator packet length ID 0",
			input: "38006F45291200",
			want:  27,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := packetToBinary(test.input)
			_, data = decodeHeader(data)
			value, data := decodeLength0(data[1:])
			assert.Equal(t, test.want, value, "length")
		})
	}
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple packet",
			input: "D2FE28",
			want:  6,
		},
		{
			name:  "operator packet length ID 0",
			input: "38006F45291200",
			want:  9,
		},
		{
			name:  "operator packet length ID 1",
			input: "EE00D40C823060",
			want:  14,
		},
		{
			name:  "operator packet sum 16",
			input: "8A004A801A8002F478",
			want:  16,
		},
		{
			name:  "operator packet sum 12",
			input: "620080001611562C8802118E34",
			want:  12,
		},
		{
			name:  "operator packet sum 23",
			input: "C0015000016115A2E0802F182340",
			want:  23,
		},
		{
			name:  "operator packet sum 31",
			input: "A0016C880162017C3686B18A3D4780",
			want:  31,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := packetToBinary(test.input)
			p, _ := decode(data)
			assert.Equal(t, test.want, p.SumVersions(), "sum versions")
		})
	}
}

func TestCompute(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		// C200B40A82 finds the sum of 1 and 2, resulting in the value 3.
		{
			name:  "sum",
			input: "C200B40A82",
			want:  3,
		},
		// 04005AC33890 finds the product of 6 and 9, resulting in the value 54.
		{
			name:  "product",
			input: "04005AC33890",
			want:  54,
		},
		// 880086C3E88112 finds the minimum of 7, 8, and 9, resulting in the value 7.
		{
			name:  "minimum",
			input: "880086C3E88112",
			want:  7,
		},
		// CE00C43D881120 finds the maximum of 7, 8, and 9, resulting in the value 9.
		{
			name:  "maximum",
			input: "CE00C43D881120",
			want:  9,
		},
		// D8005AC2A8F0 produces 1, because 5 is less than 15.
		{
			name:  "less than",
			input: "D8005AC2A8F0",
			want:  1,
		},
		// F600BC2D8F produces 0, because 5 is not greater than 15.
		{
			name:  "greater than",
			input: "F600BC2D8F",
			want:  0,
		},
		// 9C005AC2F8F0 produces 0, because 5 is not equal to 15.
		{
			name:  "not equal",
			input: "9C005AC2F8F0",
			want:  0,
		},
		// 9C0141080250320F1802104A08 produces 1, because 1 + 3 = 2 * 2.
		{
			name:  "equal",
			input: "9C0141080250320F1802104A08",
			want:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := packetToBinary(test.input)
			p, _ := decode(data)

			require.Equal(t, test.want, p.Compute())
		})
	}
}
