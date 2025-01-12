package main

import "fmt"

func main() {
	//create a map of random numbers and their roman equivalents

	omanNumerals := map[int]string{
		42:   "XLII",
		197:  "CXCVII",
		3549: "MMMDXLIX",
		816:  "DCCCXVI",
		2024: "MMXXIV",
		91:   "XCI",
		444:  "CDXLIV",
		1888: "MDCCCLXXXVIII",
		73:   "LXXIII",
		505:  "DV",
	}
	//call IntToRoman function for each number in map and compare with the value in map
	for k, v := range omanNumerals {
		//compare the result with the value in map
		if IntToRoman(k) == v {
			fmt.Println("Test passed for ", k)
		} else {
			fmt.Println("Test failed for ", k)
		}

	}
	for k, v := range omanNumerals {
		//compare the result with the value in map
		if IntToRomanGreedy(k) == v {
			fmt.Println("Test passed for ", k)
		} else {
			fmt.Println("Test failed for ", k)
		}

	}

}

func IntToRoman(num int) string {
	//create a map of roman numerals
	romanMap := map[int]string{
		1:    "I",
		4:    "IV",
		5:    "V",
		9:    "IX",
		10:   "X",
		40:   "XL",
		50:   "L",
		90:   "XC",
		100:  "C",
		400:  "CD",
		500:  "D",
		900:  "CM",
		1000: "M",
	}
	result := ""
	//loop with decreasing powers of 10

	for num > 0 {
		if num > 2000 {
			quotient, remainder := num/1000, num%1000
			fmt.Println(num, quotient, remainder)
			for i := 0; i < quotient; i++ {
				result += romanMap[1000]
			}
			num = remainder
			continue
		}
		if num >= 1000 && num <= 2000 {
			result += romanMap[1000]
			num -= 1000
			continue

		}
		if num >= 900 && num <= 999 {
			result += romanMap[900]
			num -= 900
			continue
		}
		if num >= 500 && num <= 900 {
			result += romanMap[500]
			num -= 500
			continue
		}
		if num >= 400 && num <= 499 {
			result += romanMap[400]
			num -= 400
			continue
		}
		if num >= 100 && num <= 400 {
			quotient, remainder := num/100, num%100
			//fmt.Println(num, quotient, remainder)
			for i := 0; i < quotient; i++ {
				result += romanMap[100]
			}
			num = remainder
			continue
		}
		if num >= 90 && num <= 99 {
			result += romanMap[90]
			num -= 90
			continue
		}
		if num >= 50 && num <= 90 {
			result += romanMap[50]
			num -= 50
			continue
		}
		if num >= 40 && num <= 49 {
			result += romanMap[40]
			num -= 40
			continue
		}
		if num >= 10 && num <= 40 {
			quotient, remainder := num/10, num%10
			//fmt.Println(num, quotient, remainder)
			for i := 0; i < quotient; i++ {
				result += romanMap[10]
			}
			num = remainder
			continue
		}
		if num == 9 {
			result += romanMap[9]
			num -= 9
			continue
		}
		if num >= 5 && num <= 9 {
			result += romanMap[5]
			num -= 5
			continue
		}
		if num == 4 {
			result += romanMap[4]
			num -= 4
			continue
		}
		if num >= 1 && num <= 4 {
			quotient, remainder := num/1, num%1
			//fmt.Println(num, quotient, remainder)
			for i := 0; i < quotient; i++ {
				result += romanMap[1]
			}
			num = remainder
			continue

		}

	}
	return result

}

func IntToRomanGreedy(num int) string {
	//create a map of roman numerals
	//create the struct to order the roman numerals
	type roman struct {
		value int
		roman string
	}
	romanNumerals := []roman{
		{1, "I"},
		{4, "IV"},
		{5, "V"},
		{9, "IX"},
		{10, "X"},
		{40, "XL"},
		{50, "L"},
		{90, "XC"},
		{100, "C"},
		{400, "CD"},
		{500, "D"},
		{900, "CM"},
		{1000, "M"},
	}

	//loop in reverse of the roman numerals
	result := ""
	for i := len(romanNumerals) - 1; i >= 0; i-- {
		//fmt.Println(num, romanNumerals[i].value)
		//check if the number is greater than the value in the struct
		for num >= romanNumerals[i].value {
			//if yes, add the roman numeral to result and subtract the value from num
			result += romanNumerals[i].roman
			num -= romanNumerals[i].value
			//fmt.Println(result, num)
		}
	}
	return result

}
