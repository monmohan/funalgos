package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Printf("%s\n", string([]byte{0xf0, 0x9f, 0x98, 0x8d}))
	//open utf8-trucate/cases file
	cases, err := os.OpenFile("utf8-truncate/cases", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cases)
	op, err := os.OpenFile("output", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	for scanner.Scan() {
		//read line
		buf := bytes.NewBuffer(scanner.Bytes())
		tLen := uint8(buf.Next(1)[0])
		fmt.Printf("tLen: %d\n", tLen)
		readChar(buf, int(tLen), op)
		op.Write([]byte("\n"))

	}
	defer cases.Close()
	defer op.Close()

}
func readChar(buf *bytes.Buffer, tLen int, op *os.File) {

	charsLeft := tLen

	for buf.Len() > 0 {
		mv := 0
		next := buf.Next(1)[0]
		//create a byte slice
		char := make([]byte, 0)
		char = append(char, next)
		switch {
		case next&0x80 == 0:
			//ascii
			mv = 0
		case next&0xf0 == 0xf0:
			//4 byte utf8
			mv = 3
		case next&0xe0 == 0xe0:
			//3 byte utf8
			mv = 2
		case next&0xc0 == 0xc0:
			//2 byte utf8
			mv = 1

		}

		charsLeft -= (mv + 1)
		if charsLeft < 0 {
			//truncate
			break

		}
		//append the rest of the bytes
		char = append(char, buf.Next(mv)...)
		op.Write(char)
		//fmt.Printf("% x        %q\n", char, char)

	}
}
