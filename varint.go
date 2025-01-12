package main

import (
	"fmt"
)

// protobuf varint encoding
func main() {
	var num uint64 = 102235234234234234
	//print num in binary
	fmt.Printf("%b\n", num)
	//create a slice of bytes to hold the encoded number varint
	buf := make([]byte, 10)
	encodeVarint(buf, num)
	//print the slice of bytes in hex
	fmt.Printf("%02x\n", buf[0])
	fmt.Printf("%02x\n", buf[1])

}

// encode function
func encodeVarint(buf []byte, x uint64) {
	i := 0
	for {
		buf[i] = byte(x & 0x7f)
		x = x >> 7
		if x == 0 {
			break
		}
		buf[i] |= 0x80
		i++
	}

}
func decodeVarint(buf []byte) (x uint64) {
	msbmask := uint64(0x7f)
	var t uint64
	for i, b := range buf {
		//drop the high-order bit
		t = uint64(b) & msbmask
		//print t in hex
		fmt.Printf("%b\n", t)
		t = uint64(t << (7 * i))
		fmt.Printf("%b\n", t)
		x = x | uint64(t)
		//print x in binary
		fmt.Printf("%b\n", x)
	}
	return x
}

func decodeVarintCopilot(buf []byte) (x uint64) {
	for i, b := range buf {
		x |= uint64(b&0x7f) << (7 * i)
		if b&0x80 == 0 {
			break
		}
	}

	return x
}
