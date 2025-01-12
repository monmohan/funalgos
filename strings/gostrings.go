package main

import (
	"fmt"
	"sort"
)

func main() {
	//From go strings blog post
	/*const nihongo = "日本語"
	sl := []byte(nihongo)
	sl = append(sl[:3], []byte{0xa4, 0x9e}...)
	sl = append(sl, []byte(nihongo)[3:]...)
	fmt.Printf("%s\n", string(sl))
	for index, runeValue := range string(sl) {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}*/
	//UTF-8 Multi-byte
	gaccent := "\xc3\xa0"
	fmt.Printf("%s\n", gaccent)

	//Combined characters
	//a := "\x61"
	//acute := "\xcc\x81"
	//hat := "\xcc\x82"
	combination := []byte{0x61, 0xcc, 0x81}
	fmt.Printf("%s\n", string(combination))
	strings := []string{"c", gaccent, string(combination)}
	sort.Strings(strings)
	fmt.Printf("%q\n", strings)

}
