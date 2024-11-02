package main

import "fmt"

func main() {
	var input string
	fmt.Print("Enter string: ")
	fmt.Scanln(&input)

	length := len(input)

	if length%2 == 0 {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
