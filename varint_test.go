package main

import (
	"fmt"
	"math"
	"testing"
)

// write test for encodeVarint and decodeVarint together
func TestVarint(t *testing.T) {
	//from 1 to max int 64
	//for num := uint64(0); num < math.MaxUint16; num++ {
	for num := uint64(math.MaxUint16); num < math.MaxUint32; num++ {
		// print num in binary
		fmt.Printf("%b\n", num)
		// create a slice of bytes to hold the encoded number varint
		buf := make([]byte, 10)
		encodeVarint(buf, num)
		// print the slice of bytes in hex
		fmt.Printf("%02x\n", buf[0])
		fmt.Printf("%02x\n", buf[1])
		// decode the varint
		x := decodeVarint(buf)
		//x := decodeVarintCopilot(buf)

		fmt.Printf("%v\n", x)
		if x != num {
			t.Errorf("Expected %v, got %v", num, x)
		}
	}

}
