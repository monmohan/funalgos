package main

import (
	"fmt"
	"testing"
)

// test encyrption and decryption
func TestEncryptDecrypt(t *testing.T) {
	//create some strings and call encryptinNaN and decryptinNaN and check if they are equal
	//create a slice of strings
	//for each string in the slice
	//call encryptinNaN
	//call decryptinNaN
	//check if the decrypted string is equal to the original string
	//if not, call t.Errorf
	strings := []string{"a", "ab", "abc", "abcd", "abcde", "abcdef"}
	//add some multi-byte strings less than 6 bytes
	strings = append(strings, "日本")
	strings = append(strings, string([]byte{0x61, 0xcc, 0x81}))

	for _, s := range strings {
		//print s
		fmt.Printf("%v\n", s)
		encVal, err := encryptinNaN(s)
		fmt.Printf("Encryoted Value = %v\n", encVal)

		if err != nil {
			t.Errorf("Error encrypting %v", err)
		}
		decryptedVal, err := decryptinNaN(encVal)
		fmt.Printf("Decrypted %v\n", decryptedVal)
		if err != nil {
			t.Errorf("Error decrypting %v", err)
		}
		if decryptedVal != s {
			t.Errorf("Expected %v, got %v", s, decryptedVal)
		}
	}

}
