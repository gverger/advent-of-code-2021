package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	ZERO = '0'
	ONE  = '1'
)

type rawData string

func RawData(data string) rawData {
	return rawData(data)
}

func DataFromHexString(hexString string) rawData {
	bytes, _ := hex.DecodeString(hexString)

	var dataBuilder strings.Builder
	for _, n := range bytes {
		dataBuilder.WriteString(fmt.Sprintf("%08b", n))
	}

	return rawData(dataBuilder.String())
}

func (d rawData) String() string {
	return string(d)
}

func (d rawData) parseBinary() int {
	res, err := strconv.ParseInt(d.String(), 2, 64)
	if err != nil {
		panic(err)
	}

	return int(res)
}

func (d rawData) DecodeBinary(length int) (int, rawData) {
	return d.FirstBits(length).parseBinary(), d.FromBit(length)
}

func (d rawData) At(i int) byte {
	return d[i]
}

func (d rawData) FromBit(start int) rawData {
	return d[start:]
}

func (d rawData) FirstBits(n int) rawData {
	return d[:n]
}

func (d rawData) Length() int {
	return len(d)
}
