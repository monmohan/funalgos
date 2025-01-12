package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

func byteToFloat32(b []byte) float32 {
	//convert byte array to float32
	return float32(math.Float32frombits(binary.BigEndian.Uint32(b)))

}

func main() {
	nanBits := []byte{0x7f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	nan := math.Float64frombits(binary.BigEndian.Uint64(nanBits))
	fmt.Printf("%v\n", nan)

}

//hand crafted float32
/*func testByteToFloat() {
	// create a byte slice to test byteToFloat32
	b := []byte{0x41, 0x1c, 0x00, 0x00}
	b1 := []byte{0x3f, 0x60, 0x00, 0x00}
	fmt.Printf("%v\n", byteToFloat32(b))
	fmt.Printf("%v\n", byteToFloat32(b1))
	fmt.Printf("%b\n", math.Float32bits(9.75))
}*/
//conceal method of the exercise
func encryptinNaN(s string) (float64, error) {
	//create a byte slice equivalent of NaN
	nanBits := []byte{0x7f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//nan := math.Float64frombits(binary.BigEndian.Uint64(nanBits))
	//create a byte slice equivalent of s
	sBits := []byte(s)
	if len(sBits) > 6 {
		return 0, fmt.Errorf("string too long")
	}
	sLen := byte(len(sBits))
	fmt.Printf("%v\n", sLen)
	//otherwise, append 0x00 to sBits until its length is 6
	for len(sBits) < 6 {
		sBits = append(sBits, 0x00)
	}
	nanBits[1] = nanBits[1] | sLen

	//conceal, append first two nanBits byte to SBits
	sBits = append(nanBits[:2], sBits...)
	return math.Float64frombits(binary.BigEndian.Uint64(sBits)), nil

}

func decryptinNaN(f float64) (string, error) {
	//create a byte slice equivalent of f
	fBits := math.Float64bits(f)
	//check if f is NaN
	if f != f { //simple NaN check
		//if so, decrypt
		//get the length of the string
		sLen := (fBits & (0x0007000000000000)) >> 48
		//print sLen
		fmt.Printf("%v\n", sLen)
		//get the string
		sBits := (fBits & (0x0000ffffffffffff))
		//read last SLen bytes
		b := make([]byte, sLen)
		for i := 0; i < int(sLen); i++ {

			b[i] = byte(sBits >> (8 * (8 - (3 + i))))
		}

		fmt.Printf("%s\n", b)

		return string(b), nil
	}
	return "", fmt.Errorf("not a NaN")

}
